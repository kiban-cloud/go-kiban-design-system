package authcookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"
	domain_core_authorization_model "github.com/kiban-cloud/go-kiban/domain/authorization/model"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/authcookie"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

type fakeAuthorize struct {
	status   string
	err      error
	features []string
}

func (f fakeAuthorize) AuthorizeWithSessionUseCase(_, spaceID, token string, sandbox bool, _ string) (domain_core_authorization_model.Authorization, error) {
	if f.err != nil {
		return domain_core_authorization_model.Authorization{}, f.err
	}
	return domain_core_authorization_model.Authorization{
		SpaceId:              spaceID,
		Token:                token,
		Sandbox:              sandbox,
		BillingAccountStatus: f.status,
		Features:             f.features,
	}, nil
}

type result struct {
	rec        *httptest.ResponseRecorder
	reachedApp bool
}

func serve(t *testing.T, status string, prepare func(*http.Request)) result {
	t.Helper()
	gin.SetMode(gin.TestMode)

	res := result{rec: httptest.NewRecorder()}
	engine := gin.New()
	m := authcookie.New(fakeAuthorize{status: status}, "/login")
	engine.GET("/tool", m.Middleware(), func(c *gin.Context) {
		res.reachedApp = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/tool", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	if prepare != nil {
		prepare(req)
	}
	engine.ServeHTTP(res.rec, req)
	return res
}

func TestMiddleware_ActiveAccountReachesTheTool(t *testing.T) {
	res := serve(t, "Cuenta Activa", nil)

	assert.True(t, res.reachedApp)
	assert.Equal(t, http.StatusOK, res.rec.Code)
}

func TestMiddleware_EmptyStatusIsNotBlocked(t *testing.T) {
	// kiban-cloud deploys older than the billingAccountStatus field send
	// nothing; that must not lock every user out.
	res := serve(t, "", nil)

	assert.True(t, res.reachedApp)
}

func TestMiddleware_SuspendedAccountGetsTheFullPage(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, nil)

	require.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	body := res.rec.Body.String()
	assert.Contains(t, body, `data-account-blocked="suspended"`)
	assert.Contains(t, body, "Tu cuenta está suspendida")
	assert.Contains(t, body, `href="/kiban-cloud/billing"`)
	assert.Contains(t, body, "<!doctype html>")
}

func TestMiddleware_InactiveAccountPointsToSupport(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_INACTIVE, nil)

	require.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	body := res.rec.Body.String()
	assert.Contains(t, body, `data-account-blocked="inactive"`)
	assert.Contains(t, body, `href="mailto:help@kiban.com"`)
	assert.NotContains(t, body, `href="/kiban-cloud/billing"`)
}

func TestMiddleware_BlockedHtmxRequestRefreshesInsteadOfSwapping(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, func(r *http.Request) {
		r.Header.Set("HX-Request", "true")
	})

	require.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	assert.Equal(t, "true", res.rec.Header().Get("HX-Refresh"))
	assert.Empty(t, res.rec.Body.String())
}

func TestMiddleware_BlockedPageFollowsTheLanguageCookie(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: fsi18n.CookieLang, Value: "en"})
	})

	assert.Contains(t, res.rec.Body.String(), "Your account is suspended")
}

func TestMiddleware_BlockedPageFallsBackToAcceptLanguage(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, func(r *http.Request) {
		r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	})

	assert.Contains(t, res.rec.Body.String(), "Your account is suspended")
}

// El contrato Plan2025Anual vencido lo decide kiban-cloud y llega como error de
// la sesión: la herramienta pinta la pantalla de cuenta cortada, no manda al
// login como haría con una sesión inválida.
func TestMiddleware_ContractExpiredGetsTheFullPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	reached := false
	engine := gin.New()
	m := authcookie.New(fakeAuthorize{err: domain_core_authorization_model.ContractExpiredError()}, "/login")
	engine.GET("/tool", m.Middleware(), func(c *gin.Context) {
		reached = true
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/tool", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	engine.ServeHTTP(rec, req)

	require.False(t, reached)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `data-account-blocked="contract_expired"`)
	assert.Contains(t, body, "Tu contrato anual terminó")
	assert.Contains(t, body, `href="mailto:help@kiban.com"`)
	assert.NotContains(t, body, `href="/kiban-cloud/billing"`)
	assert.Empty(t, rec.Header().Get("Location"), "no redirige al login")
}

// shellFlags corre el middleware con los feature flags de la compañía y la
// cookie kiban_ui dados, y devuelve qué banderas del shell vio la herramienta.
func shellFlags(t *testing.T, features []string, ui string) (shell, sw bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	m := authcookie.New(fakeAuthorize{features: features}, "/login")
	engine.GET("/tool", m.Middleware(), func(c *gin.Context) {
		ctx := c.Request.Context()
		shell = feature.Enabled(ctx, layout.FlagShellFlamboyant)
		sw = feature.Enabled(ctx, layout.FlagShellFlamboyantSwitch)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/tool", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	if ui != "" {
		req.AddCookie(&http.Cookie{Name: layout.UICookie, Value: ui})
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	return shell, sw
}

func TestMiddleware_SinFeatureFlagNoHayShellNiSwitch(t *testing.T) {
	shell, sw := shellFlags(t, nil, "")

	assert.False(t, shell)
	assert.False(t, sw)
}

func TestMiddleware_ConFeatureFlagYSinCookieSeOfreceElSwitch(t *testing.T) {
	shell, sw := shellFlags(t, []string{authcookie.FeatureShellFlamboyant}, "")

	assert.False(t, shell, "sin elección del usuario se ve el diseño actual")
	assert.True(t, sw)
}

func TestMiddleware_ConFeatureFlagYCookieSePintaFlamboyant(t *testing.T) {
	shell, sw := shellFlags(t, []string{authcookie.FeatureShellFlamboyant}, layout.UIFlamboyant)

	assert.True(t, shell)
	assert.True(t, sw)
}

func TestMiddleware_CookieClasicoDejaElDisenoActual(t *testing.T) {
	shell, sw := shellFlags(t, []string{authcookie.FeatureShellFlamboyant}, layout.UIClasico)

	assert.False(t, shell)
	assert.True(t, sw)
}

func TestMiddleware_CookieSinFeatureFlagNoPrendeNada(t *testing.T) {
	shell, sw := shellFlags(t, []string{"OTRO_FLAG"}, layout.UIFlamboyant)

	assert.False(t, shell, "la cookie sola no alcanza: la compañía necesita el flag")
	assert.False(t, sw)
}

func TestMiddleware_FeatureLocalPrendeElSwitchSinFlagDeLaCompania(t *testing.T) {
	t.Setenv(authcookie.EnvLocalFeatures, "OTRO, "+authcookie.FeatureShellFlamboyant)

	shell, sw := shellFlags(t, nil, layout.UIFlamboyant)

	assert.True(t, shell)
	assert.True(t, sw)
}

func TestMiddleware_FeatureLocalSeIgnoraEnCloudRun(t *testing.T) {
	t.Setenv(authcookie.EnvLocalFeatures, authcookie.FeatureShellFlamboyant)
	t.Setenv("K_REVISION", "crm-00042")

	shell, sw := shellFlags(t, nil, layout.UIFlamboyant)

	assert.False(t, shell, "desplegado, la variable no puede prender nada")
	assert.False(t, sw)
}
