package flash_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/flash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

// Los banners de siempre no cambian: siguen trayendo la × y su temporizador.
func TestBanner_StaysDismissible(t *testing.T) {
	body := render(t, flash.Success("Guardado."))

	assert.Contains(t, body, "data-kiban-flash-dismiss")
	assert.Contains(t, body, `data-kiban-flash-autodismiss-ms="5000"`)
	assert.Contains(t, body, "Guardado.")
}

// Un banner que explica por qué algo no se puede hacer no es una notificación:
// cerrarlo esconde la única explicación de la pantalla.
func TestAlert_IsNotDismissibleByDefault(t *testing.T) {
	body := render(t, flash.Alert(flash.Options{Variant: "warning", Message: "No hay integración activa."}))

	assert.NotContains(t, body, "data-kiban-flash-dismiss")
	assert.NotContains(t, body, "data-kiban-flash-autodismiss-ms")
	assert.Contains(t, body, "No hay integración activa.")
}

// Y si no se puede cerrar a mano, tampoco puede irse solo.
func TestAlert_DismissibleDrivesTheTimer(t *testing.T) {
	sticky := render(t, flash.Alert(flash.Options{Variant: "info", Message: "Aviso."}))
	assert.NotContains(t, sticky, "data-kiban-flash-autodismiss-ms")

	ephemeral := render(t, flash.Alert(flash.Options{Variant: "info", Message: "Aviso.", Dismissible: true}))
	assert.Contains(t, ephemeral, `data-kiban-flash-autodismiss-ms="6000"`)
	assert.Contains(t, ephemeral, "data-kiban-flash-dismiss")
}

// La acción va DENTRO del banner: debajo se lee como un párrafo aparte y
// pierde la relación con lo que el banner dice.
func TestAlert_RendersTheActionInsideTheBox(t *testing.T) {
	body := render(t, flash.Alert(flash.Options{
		Variant:     "warning",
		Message:     "No hay integración activa.",
		ActionLabel: "Configurar una integración",
		ActionHref:  "/rekon/facturacion/configuracion?tab=integraciones",
	}))

	assert.Contains(t, body, `href="/rekon/facturacion/configuracion?tab=integraciones"`)
	assert.Contains(t, body, "Configurar una integración")
	// Dentro del mismo div del banner, no después de cerrarlo.
	assert.Greater(t, len(body), 0)
	assert.NotContains(t, body, "</div><a ")
}

// Con atributos y sin href sale un <button>: una acción htmx no necesita un
// enlace que no lleva a ninguna parte.
func TestAlert_ActionWithoutHrefIsAButton(t *testing.T) {
	body := render(t, flash.Alert(flash.Options{
		Variant:     "error",
		Message:     "No pudimos timbrar.",
		ActionLabel: "Reintentar",
		ActionAttrs: templ.Attributes{"hx-post": "/rekon/facturacion/1/reintentar"},
	}))

	assert.Contains(t, body, `hx-post="/rekon/facturacion/1/reintentar"`)
	assert.Contains(t, body, "<button")
	assert.NotContains(t, body, "<a\n")
}

// Sin label no hay acción, aunque vengan atributos: un botón sin texto es un
// control invisible.
func TestAlert_WithoutLabelRendersNoAction(t *testing.T) {
	body := render(t, flash.Alert(flash.Options{
		Message:     "Sólo texto.",
		ActionHref:  "/algun/lado",
		ActionAttrs: templ.Attributes{"hx-post": "/algun/lado"},
	}))

	assert.NotContains(t, body, "/algun/lado")
}
