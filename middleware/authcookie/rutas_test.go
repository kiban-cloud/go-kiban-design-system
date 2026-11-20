package authcookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/authcookie"
	"github.com/kiban-cloud/go-kiban-design-system/middleware/rutas"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

// porRutas sirve una pantalla de workfloo detrás de authcookie, con o sin el
// envoltorio de middleware/rutas, y devuelve la respuesta.
func porRutas(t *testing.T, method, target string, flamboyant, envuelto bool, prepare func(*http.Request)) (*httptest.ResponseRecorder, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	reached := false
	engine := gin.New()
	m := authcookie.New(fakeAuthorize{features: []string{authcookie.FeatureShellFlamboyant}}, "/login")
	engine.Handle(method, "/workfloo/workfloos/:id", m.Middleware(), func(c *gin.Context) {
		reached = true
		c.Status(http.StatusOK)
	})
	var h http.Handler = engine.Handler()
	if envuelto {
		h = rutas.AceptarNombresNuevos(h)
	}

	req := httptest.NewRequest(method, target, nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	if flamboyant {
		req.AddCookie(&http.Cookie{Name: layout.UICookie, Value: layout.UIFlamboyant})
	}
	if prepare != nil {
		prepare(req)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec, reached
}

func TestMiddleware_FlamboyantConNombreViejoVaAlNuevo(t *testing.T) {
	rec, reached := porRutas(t, http.MethodGet, "/workfloo/workfloos/wf1?tab=2", true, true, nil)

	assert.False(t, reached)
	assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
	assert.Equal(t, "/solicitudes/workfloos/wf1?tab=2", rec.Header().Get("Location"))
}

func TestMiddleware_FlamboyantConNombreNuevoSeQueda(t *testing.T) {
	rec, reached := porRutas(t, http.MethodGet, "/solicitudes/workfloos/wf1", true, true, nil)

	assert.True(t, reached)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_SinFlamboyantConNombreNuevoVuelveAlViejo(t *testing.T) {
	rec, reached := porRutas(t, http.MethodGet, "/solicitudes/workfloos/wf1?tab=2", false, true, nil)

	assert.False(t, reached)
	assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
	assert.Equal(t, "/workfloo/workfloos/wf1?tab=2", rec.Header().Get("Location"))
}

func TestMiddleware_SinFlamboyantConNombreViejoSeQueda(t *testing.T) {
	rec, reached := porRutas(t, http.MethodGet, "/workfloo/workfloos/wf1", false, true, nil)

	assert.True(t, reached, "el diseño de siempre no cambia")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_ElRedirectConservaElMetodo(t *testing.T) {
	rec, _ := porRutas(t, http.MethodPost, "/workfloo/workfloos/wf1", true, true, func(r *http.Request) {
		r.Header.Set("HX-Request", "true")
	})

	assert.Equal(t, http.StatusPermanentRedirect, rec.Code, "308: el POST de htmx llega igual al nombre nuevo")
	assert.Equal(t, "/solicitudes/workfloos/wf1", rec.Header().Get("Location"))
}

func TestMiddleware_SinElEnvoltorioNoRedirige(t *testing.T) {
	rec, reached := porRutas(t, http.MethodGet, "/workfloo/workfloos/wf1", true, false, nil)

	assert.True(t, reached, "un módulo que no sirve el nombre nuevo no se manda a él")
	assert.Equal(t, http.StatusOK, rec.Code)
}
