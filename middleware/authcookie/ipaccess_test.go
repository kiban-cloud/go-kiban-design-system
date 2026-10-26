package authcookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	domain_core_authorization_model "github.com/kiban-cloud/go-kiban/domain/authorization/model"
	infrastructure_core_env "github.com/kiban-cloud/go-kiban/infrastructure_core/env"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/authcookie"
)

// fakeConsole hace lo que el caso de uso de go-kiban con la respuesta de
// kiban-cloud: la sesión, y ErrorIpNotAllowed (con la autorización completa)
// si kiban-cloud rechazó esa IP para esta herramienta. bloqueadas es lo que
// kiban-cloud rechazaría, por herramienta.
type fakeConsole struct {
	bloqueadas map[string][]string
	status     string
}

func (f fakeConsole) AuthorizeWithSessionUseCase(_, spaceID, token string, sandbox bool, ip string) (domain_core_authorization_model.Authorization, error) {
	status := f.status
	if status == "" {
		status = "Cuenta Activa"
	}
	auth := domain_core_authorization_model.Authorization{
		SpaceId: spaceID, Token: token, Sandbox: sandbox, Ip: ip,
		BillingAccountStatus: status,
	}
	for _, b := range f.bloqueadas[infrastructure_core_env.Get(infrastructure_core_env.TOOL)] {
		if b == ip {
			return auth, domain_core_authorization_model.ErrorIpNotAllowed
		}
	}
	return auth, nil
}

// oficina: kiban-cloud rechaza 203.0.113.9 en workfloo; klin queda libre.
var oficina = map[string][]string{"WORKFLOO": {"203.0.113.9"}}

// serveIp corre el middleware en una herramienta (TOOL) con la IP del cliente
// tal como la deja el balanceador: "<cliente>, <balanceador>".
func serveIp(t *testing.T, tool, path, ipCliente string, f fakeConsole, htmxReq bool) result {
	t.Helper()
	gin.SetMode(gin.TestMode)
	infrastructure_core_env.Set(infrastructure_core_env.TOOL, tool)
	t.Cleanup(func() { infrastructure_core_env.Set(infrastructure_core_env.TOOL, "") })

	res := result{rec: httptest.NewRecorder()}
	engine := gin.New()
	m := authcookie.New(f, "/login")
	engine.Any("/*any", m.Middleware(), func(c *gin.Context) {
		res.reachedApp = true
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	req.Header.Set("X-Forwarded-For", ipCliente+", 34.36.10.1")
	if htmxReq {
		req.Header.Set("HX-Request", "true")
	}
	engine.ServeHTTP(res.rec, req)
	return res
}

func TestIpAccess_IpDeLaListaEntra(t *testing.T) {
	res := serveIp(t, "WORKFLOO", "/workfloo", "189.203.44.17", fakeConsole{bloqueadas: oficina}, false)
	assert.True(t, res.reachedApp)
	assert.Equal(t, http.StatusOK, res.rec.Code)
}

func TestIpAccess_IpFueraVeLaPantallaDeBloqueo(t *testing.T) {
	res := serveIp(t, "WORKFLOO", "/workfloo", "203.0.113.9", fakeConsole{bloqueadas: oficina}, false)
	assert.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	body := res.rec.Body.String()
	assert.Contains(t, body, `data-ip-blocked="WORKFLOO"`)
	assert.Contains(t, body, "203.0.113.9")
	assert.Contains(t, body, "Acceso no permitido desde esta IP")
	assert.Contains(t, body, "workfloo")
}

// htmx no puede tomar la página como swap: recarga y la página completa
// muestra la pantalla, igual que con la cuenta suspendida.
func TestIpAccess_HtmxRecargaEnVezDeSwap(t *testing.T) {
	res := serveIp(t, "WORKFLOO", "/workfloo/executions", "203.0.113.9", fakeConsole{bloqueadas: oficina}, true)
	assert.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	assert.Equal(t, "true", res.rec.Header().Get("HX-Refresh"))
	assert.Empty(t, res.rec.Body.String())
}

// Una herramienta que ninguna IP de la lista cubre queda libre.
func TestIpAccess_HerramientaNoCubiertaQuedaLibre(t *testing.T) {
	res := serveIp(t, "KLIN", "/klin", "203.0.113.9", fakeConsole{bloqueadas: oficina}, false)
	assert.True(t, res.reachedApp)
}

// Cambiar de espacio y salir son la salida de la pantalla: nunca se bloquean.
func TestIpAccess_CambiarDeEspacioYSalirNoSeBloquean(t *testing.T) {
	for _, path := range []string{"/workfloo/spaces/switch", "/workfloo/logout"} {
		res := serveIp(t, "WORKFLOO", path, "203.0.113.9", fakeConsole{bloqueadas: oficina}, false)
		assert.True(t, res.reachedApp, path)
	}
}

// Lo que kiban-cloud deja pasar, pasa: el middleware no decide por su cuenta.
func TestIpAccess_SinRechazoDeKibanCloudPasa(t *testing.T) {
	res := serveIp(t, "WORKFLOO", "/workfloo", "203.0.113.9", fakeConsole{}, false)
	assert.True(t, res.reachedApp)
}

// El firewall va primero: con la IP fuera de la lista se ve la pantalla de IP,
// aunque la cuenta además esté suspendida.
func TestIpAccess_ElFirewallVaAntesQueLaCuenta(t *testing.T) {
	f := fakeConsole{bloqueadas: oficina, status: domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED}
	res := serveIp(t, "WORKFLOO", "/workfloo", "203.0.113.9", f, false)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	assert.Contains(t, res.rec.Body.String(), `data-ip-blocked="WORKFLOO"`)
	assert.NotContains(t, res.rec.Body.String(), "data-account-blocked")
}

// fakeIpVista registra con qué IP se pidió la autorización.
type fakeIpVista struct{ ip *string }

func (f fakeIpVista) AuthorizeWithSessionUseCase(_, spaceID, token string, sandbox bool, ip string) (domain_core_authorization_model.Authorization, error) {
	*f.ip = ip
	return domain_core_authorization_model.Authorization{SpaceId: spaceID, Token: token, Ip: ip, BillingAccountStatus: "Cuenta Activa"}, nil
}

// Otra herramienta que llama con la sesión del usuario sale por la IP de su
// servidor: la lista Consola se evalúa con la IP que reenvía, y solo si trae la
// api key interna de esta herramienta.
func TestIpAccess_IpReenviadaPorOtraHerramienta(t *testing.T) {
	casos := []struct {
		nombre   string
		apiKey   string
		esperada string
	}{
		{"con la api key interna se evalúa la IP del usuario", "llave-interna", "189.245.28.84"},
		{"sin la api key interna se evalúa la IP de origen", "otra", "34.96.55.51"},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			t.Setenv(infrastructure_core_env.INTERNAL_API_KEY, "llave-interna")
			var vista string
			engine := gin.New()
			engine.Any("/*any", authcookie.New(fakeIpVista{ip: &vista}, "/login").Middleware(), func(c *gin.Context) { c.Status(http.StatusOK) })

			req := httptest.NewRequest(http.MethodGet, "/workfloo", nil)
			req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
			req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
			req.Header.Set("X-Forwarded-For", "34.96.55.51, 8.233.84.187")
			req.Header.Set("X-Kiban-Client-Ip", "189.245.28.84")
			req.Header.Set("x-api-key", caso.apiKey)
			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, caso.esperada, vista)
		})
	}
}
