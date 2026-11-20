package rutas_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/rutas"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// sesionPropia es un middleware de sesión que no es authcookie (el de staff
// de klin): prende las banderas a su modo y termina con rutas.Aplicar.
func sesionPropia(on bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if on {
			c.Request = c.Request.WithContext(feature.Enable(c.Request.Context(), flamboyant.Flag))
		}
		rutas.Aplicar(c)
	}
}

func staff(t *testing.T, on bool, target string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/klin-internal/deliveries/:id", sesionPropia(on), func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<a href="/klin-internal/deliveries">volver</a>`))
	})
	rec := httptest.NewRecorder()
	rutas.AceptarNombresNuevos(engine.Handler()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestAplicar_OtroMiddlewareDeSesion(t *testing.T) {
	t.Run("flamboyant con el nombre viejo va al nuevo", func(t *testing.T) {
		rec := staff(t, true, "/klin-internal/deliveries/9?tab=1")
		assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
		assert.Equal(t, "/cumplimiento-interno/deliveries/9?tab=1", rec.Header().Get("Location"))
	})
	t.Run("flamboyant con el nombre nuevo reescribe la respuesta", func(t *testing.T) {
		rec := staff(t, true, "/cumplimiento-interno/deliveries/9")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `<a href="/cumplimiento-interno/deliveries">volver</a>`, rec.Body.String())
	})
	t.Run("sin flamboyant con el nombre nuevo vuelve al viejo", func(t *testing.T) {
		rec := staff(t, false, "/cumplimiento-interno/deliveries/9")
		assert.Equal(t, http.StatusPermanentRedirect, rec.Code)
		assert.Equal(t, "/klin-internal/deliveries/9", rec.Header().Get("Location"))
	})
	t.Run("sin flamboyant con el nombre viejo no cambia nada", func(t *testing.T) {
		rec := staff(t, false, "/klin-internal/deliveries/9")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `<a href="/klin-internal/deliveries">volver</a>`, rec.Body.String())
	})
}
