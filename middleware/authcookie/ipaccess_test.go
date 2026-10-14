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

// fakeConsole hace lo que el caso de uso de go-kiban: arma la sesión con la
// lista Consola y aplica la misma regla (ConsoleIpAllowed) con la IP que le
// pasa el middleware. Si no pasa, devuelve ErrorIpNotAllowed con la
// autorización completa.
type fakeConsole struct {
	features []string
	rules    *domain_core_authorization_model.IpAccessRules
	status   string
}

func (f fakeConsole) AuthorizeWithSessionUseCase(_, spaceID, token string, sandbox bool, ip string) (domain_core_authorization_model.Authorization, error) {
	status := f.status
	if status == "" {
		status = "Cuenta Activa"
	}
	auth := domain_core_authorization_model.Authorization{
		SpaceId: spaceID, Token: token, Sandbox: sandbox, Ip: ip,
		BillingAccountStatus: status,
		Features:             f.features,
		ConsoleIpAccess:      f.rules,
	}
	if !auth.ConsoleIpAllowed(infrastructure_core_env.Get(infrastructure_core_env.TOOL)) {
		return auth, domain_core_authorization_model.ErrorIpNotAllowed
	}
	return auth, nil
}

var oficina = &domain_core_authorization_model.IpAccessRules{Mode: "list", Entries: []domain_core_authorization_model.IpAccessRule{
	{Cidr: "189.203.44.0/24", Tools: []string{"WORKFLOO"}},
}}

var conFlag = []string{domain_core_authorization_model.FEATURE_IP_ACCESS_LIST}

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
	res := serveIp(t, "WORKFLOO", "/workfloo", "189.203.44.17", fakeConsole{features: conFlag, rules: oficina}, false)
	assert.True(t, res.reachedApp)
	assert.Equal(t, http.StatusOK, res.rec.Code)
}

func TestIpAccess_IpFueraVeLaPantallaDeBloqueo(t *testing.T) {
	res := serveIp(t, "WORKFLOO", "/workfloo", "203.0.113.9", fakeConsole{features: conFlag, rules: oficina}, false)
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
	res := serveIp(t, "WORKFLOO", "/workfloo/executions", "203.0.113.9", fakeConsole{features: conFlag, rules: oficina}, true)
	assert.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	assert.Equal(t, "true", res.rec.Header().Get("HX-Refresh"))
	assert.Empty(t, res.rec.Body.String())
}

// Una herramienta que ninguna IP de la lista cubre queda libre.
func TestIpAccess_HerramientaNoCubiertaQuedaLibre(t *testing.T) {
	res := serveIp(t, "KLIN", "/klin", "203.0.113.9", fakeConsole{features: conFlag, rules: oficina}, false)
	assert.True(t, res.reachedApp)
}

// Cambiar de espacio y salir son la salida de la pantalla: nunca se bloquean.
func TestIpAccess_CambiarDeEspacioYSalirNoSeBloquean(t *testing.T) {
	for _, path := range []string{"/workfloo/spaces/switch", "/workfloo/logout"} {
		res := serveIp(t, "WORKFLOO", path, "203.0.113.9", fakeConsole{features: conFlag, rules: oficina}, false)
		assert.True(t, res.reachedApp, path)
	}
}

func TestIpAccess_SinFlagNoSeFiltra(t *testing.T) {
	res := serveIp(t, "WORKFLOO", "/workfloo", "203.0.113.9", fakeConsole{rules: oficina}, false)
	assert.True(t, res.reachedApp)
}

// El firewall va primero: con la IP fuera de la lista se ve la pantalla de IP,
// aunque la cuenta además esté suspendida.
func TestIpAccess_ElFirewallVaAntesQueLaCuenta(t *testing.T) {
	f := fakeConsole{features: conFlag, rules: oficina, status: domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED}
	res := serveIp(t, "WORKFLOO", "/workfloo", "203.0.113.9", f, false)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	assert.Contains(t, res.rec.Body.String(), `data-ip-blocked="WORKFLOO"`)
	assert.NotContains(t, res.rec.Body.String(), "data-account-blocked")
}
