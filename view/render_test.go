package view_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	"github.com/kiban-cloud/go-kiban-design-system/view"
)

func renderRecorder(t *testing.T, htmx bool) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest("GET", "/anything?page=3", nil)
	if htmx {
		c.Request.Header.Set("HX-Request", "true")
	}
	body := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p>hi</p>")
		return err
	})
	view.Render(c, 200, body)
	return rec
}

// Cada pantalla responde la MISMA url con dos cuerpos distintos según el header
// de htmx (página completa vs fragmento). Sin `Vary: HX-Request` el navegador
// no sabe que son recursos distintos y puede servir un fragmento cacheado como
// documento entero — la página vuelve sin <head>, o sea sin estilos ni scripts.
// Pasa al volver atrás a una url que htmx empujó (paginación, filtros).
func TestRender_AlwaysVariesOnHxRequest(t *testing.T) {
	for _, tc := range []struct {
		name string
		htmx bool
	}{
		{"navegación normal", false},
		{"petición htmx", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := renderRecorder(t, tc.htmx)
			if got := rec.Header().Get("Vary"); got != "HX-Request" {
				t.Errorf("Vary = %q, se esperaba HX-Request", got)
			}
		})
	}
}

// El header no debe pisar un Vary que el handler haya puesto antes (por eso se
// agrega, no se asigna).
func TestRender_AppendsToExistingVary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest("GET", "/anything", nil)
	c.Writer.Header().Set("Vary", "Accept-Language")

	view.Render(c, 200, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "x")
		return err
	}))

	got := rec.Header().Values("Vary")
	if len(got) != 2 || got[0] != "Accept-Language" || got[1] != "HX-Request" {
		t.Errorf("Vary = %v, se esperaba conservar Accept-Language y sumar HX-Request", got)
	}
}
