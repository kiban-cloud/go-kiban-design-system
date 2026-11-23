package authcookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/authcookie"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

// novedadesFlags corre el middleware con los feature flags dados y devuelve qué
// vio la herramienta: la bandera novedades y la del shell flamboyant.
func novedadesFlags(t *testing.T, features []string) (novedades, shell bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	m := authcookie.New(fakeAuthorize{features: features}, "/login")
	engine.GET("/tool", m.Middleware(), func(c *gin.Context) {
		novedades = feature.Enabled(c.Request.Context(), layout.FlagNovedades)
		shell = feature.Enabled(c.Request.Context(), layout.FlagShellFlamboyantSwitch)
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/tool", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	return novedades, shell
}

func TestMiddleware_NovedadesConElFlagDeLaCompania(t *testing.T) {
	novedades, shell := novedadesFlags(t, []string{authcookie.FeatureNovedades})
	assert.True(t, novedades)
	assert.False(t, shell, "NOVEDADES no prende flamboyant")
}

func TestMiddleware_NovedadesConFlamboyant(t *testing.T) {
	novedades, shell := novedadesFlags(t, []string{authcookie.FeatureShellFlamboyant, authcookie.FeatureNovedades})
	assert.True(t, novedades)
	assert.True(t, shell)
}

func TestMiddleware_SinNovedadesSinBandera(t *testing.T) {
	novedades, _ := novedadesFlags(t, []string{authcookie.FeatureShellFlamboyant})
	assert.False(t, novedades)
}

func TestMiddleware_NovedadesLocal(t *testing.T) {
	t.Setenv(authcookie.EnvLocalFeatures, authcookie.FeatureNovedades)
	novedades, _ := novedadesFlags(t, nil)
	assert.True(t, novedades)
}
