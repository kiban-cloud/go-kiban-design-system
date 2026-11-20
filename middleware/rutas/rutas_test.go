package rutas_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/rutas"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

type visto struct {
	path, query      string
	acepta, porNuevo bool
}

func servir(t *testing.T, target string) (*httptest.ResponseRecorder, visto) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	var v visto
	engine := gin.New()
	engine.GET("/workfloo/workfloos/:id", func(c *gin.Context) {
		ctx := c.Request.Context()
		v = visto{c.Request.URL.Path, c.Request.URL.RawQuery, flamboyant.AceptaNombresNuevos(ctx), flamboyant.PorNombreNuevo(ctx)}
		c.String(http.StatusOK, c.Param("id"))
	})
	rec := httptest.NewRecorder()
	rutas.AceptarNombresNuevos(engine.Handler()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec, v
}

func TestAceptarNombresNuevos_ElNombreNuevoLlegaALaRutaDeSiempre(t *testing.T) {
	rec, v := servir(t, "/solicitudes/workfloos/wf1?tab=2")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "wf1", rec.Body.String())
	assert.Equal(t, visto{"/workfloo/workfloos/wf1", "tab=2", true, true}, v)
}

func TestAceptarNombresNuevos_ElNombreViejoPasaTalCualYQuedaMarcado(t *testing.T) {
	rec, v := servir(t, "/workfloo/workfloos/wf1")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, visto{"/workfloo/workfloos/wf1", "", true, false}, v)
}

func TestAceptarNombresNuevos_OtroPrefijoNoSeTraduce(t *testing.T) {
	rec, _ := servir(t, "/pagosx/workfloos/wf1")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAceptarNombresNuevos_TraduceElPathCrudo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var raw string
	h := rutas.AceptarNombresNuevos(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		raw = r.URL.RawPath
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/solicitudes/a%2Fb", nil))

	assert.Equal(t, "/workfloo/a%2Fb", raw)
}
