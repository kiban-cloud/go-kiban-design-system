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

// redirige sirve un handler de workfloo que responde h, con flamboyant o sin
// él, detrás del envoltorio de rutas.
func redirige(t *testing.T, flamboyant bool, h gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	m := authcookie.New(fakeAuthorize{features: []string{authcookie.FeatureShellFlamboyant}}, "/login")
	engine.POST("/workfloo/workfloos/:id", m.Middleware(), h)

	req := httptest.NewRequest(http.MethodPost, "/solicitudes/workfloos/wf1", nil)
	if !flamboyant {
		req = httptest.NewRequest(http.MethodPost, "/workfloo/workfloos/wf1", nil)
	}
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	if flamboyant {
		req.AddCookie(&http.Cookie{Name: layout.UICookie, Value: layout.UIFlamboyant})
	}
	rec := httptest.NewRecorder()
	rutas.AceptarNombresNuevos(engine.Handler()).ServeHTTP(rec, req)
	return rec
}

func TestMiddleware_FlamboyantRedirigeConElNombreNuevo(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/workfloo/workfloos/wf2?ok=1")
	})

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/solicitudes/workfloos/wf2?ok=1", rec.Header().Get("Location"))
}

func TestMiddleware_FlamboyantHXRedirectSinCuerpoConElNombreNuevo(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.Header("HX-Redirect", "/workfloo/workfloos/wf2")
		c.Status(http.StatusOK)
	})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "/solicitudes/workfloos/wf2", rec.Header().Get("HX-Redirect"))
}

func TestMiddleware_FlamboyantPushUrlPuestoDespuesDelStatus(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.Header("HX-Push-Url", "/workfloo/workfloos?page=2")
		c.String(http.StatusOK, "tabla")
	})

	assert.Equal(t, "/solicitudes/workfloos?page=2", rec.Header().Get("HX-Push-Url"))
	assert.Equal(t, "tabla", rec.Body.String())
}

func TestMiddleware_FlamboyantNoTocaOtrosModulosNiJSON(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.Header("HX-Location", `{"path":"/workfloo/x"}`)
		c.Header("HX-Push-Url", "/crm/customers")
		c.Status(http.StatusOK)
	})

	assert.Equal(t, `{"path":"/workfloo/x"}`, rec.Header().Get("HX-Location"))
	assert.Equal(t, "/crm/customers", rec.Header().Get("HX-Push-Url"))
}

func TestMiddleware_SinFlamboyantLosRedirectsNoCambian(t *testing.T) {
	rec := redirige(t, false, func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/workfloo/workfloos/wf2")
	})

	assert.Equal(t, "/workfloo/workfloos/wf2", rec.Header().Get("Location"))
}

const paginaWorkfloo = `<a href="/workfloo/workfloos?page=2&amp;q=x">Flujos</a>` +
	`<button hx-get="/workfloo/workfloos/wf1/edit" hx-push-url="/workfloo/workfloos/wf1">e</button>` +
	`<form action="/klin-internal/deliveries/9" data-kiban-panel-url="/rekon/payments/p1"></form>` +
	`<img src="/klin/static/img/sic/bc.png">` +
	`<a href="/crm/customers">c</a><a href="/workfloo-templates">t</a><a href="/portal/workfloo/1">p</a>` +
	`<a href="https://cloud.kiban.com/workfloo/x">abs</a>` +
	`<script>var u = '/workfloo/workfloos/' + id;</script>` +
	`<p>Ve a /workfloo/workfloos en el menú</p>`

func TestMiddleware_FlamboyantReescribeLasURLsDelHTML(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(paginaWorkfloo))
	})

	body := rec.Body.String()
	for _, want := range []string{
		`href="/solicitudes/workfloos?page=2&amp;q=x"`,
		`hx-get="/solicitudes/workfloos/wf1/edit"`, `hx-push-url="/solicitudes/workfloos/wf1"`,
		`action="/cumplimiento-interno/deliveries/9"`, `data-kiban-panel-url="/pagos/payments/p1"`,
		`src="/cumplimiento/static/img/sic/bc.png"`,
	} {
		assert.Contains(t, body, want)
	}
	for _, igual := range []string{
		`href="/crm/customers"`, `href="/workfloo-templates"`, `href="/portal/workfloo/1"`,
		`href="https://cloud.kiban.com/workfloo/x"`,
		`var u = '/workfloo/workfloos/' + id;`, `Ve a /workfloo/workfloos en el menú`,
	} {
		assert.Contains(t, body, igual, "esto no se toca")
	}
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_FlamboyantConservaElStatusDelHTML(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.Data(http.StatusUnprocessableEntity, "text/html; charset=utf-8", []byte(`<form hx-post="/workfloo/x"></form>`))
	})

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t, `<form hx-post="/solicitudes/x"></form>`, rec.Body.String())
}

func TestMiddleware_FlamboyantNoTocaLoQueNoEsHTML(t *testing.T) {
	rec := redirige(t, true, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"href": "/workfloo/x"})
	})

	assert.JSONEq(t, `{"href":"/workfloo/x"}`, rec.Body.String())
}

func TestMiddleware_SinFlamboyantElHTMLNoCambia(t *testing.T) {
	rec := redirige(t, false, func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(paginaWorkfloo))
	})

	assert.Equal(t, paginaWorkfloo, rec.Body.String())
}
