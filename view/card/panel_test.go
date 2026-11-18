package card_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

func TestPanelSection_ConYSinTitulo(t *testing.T) {
	con := render(t, context.Background(), card.PanelSection("Datos del pago"))
	assert.Contains(t, con, "data-kiban-panel-section")
	assert.Contains(t, con, ">Datos del pago</h3>")
	assert.Contains(t, con, "<p>cuerpo</p>")

	sin := render(t, context.Background(), card.PanelSection(""))
	assert.NotContains(t, sin, "<h3", "sin título sólo la tarjeta")
	assert.Contains(t, sin, "<p>cuerpo</p>")
}

// Note: «warning» en ámbar; cualquier otro tono, informativo en azul.
func TestNote_Tonos(t *testing.T) {
	aviso := render(t, context.Background(), card.Note("warning", "Se retira todo el saldo"))
	assert.Contains(t, aviso, `data-kiban-note="warning"`)
	assert.Contains(t, aviso, "bg-[#FFF4E5] text-[#7A3E00]")
	assert.Contains(t, aviso, ">Se retira todo el saldo</div>")

	for _, tone := range []string{"info", "", "otro"} {
		body := render(t, context.Background(), card.Note(tone, "x"))
		assert.Contains(t, body, `data-kiban-note="info"`, "tono %q", tone)
		assert.Contains(t, body, "bg-[#EAF0FF] text-[#1A3FA6]", "tono %q", tone)
	}
}

func TestCallout_TonosYDesconocidoEsInfo(t *testing.T) {
	cases := map[string]struct{ tone, class string }{
		"ok":   {"ok", "bg-[#E3FBF5] text-[#007A63]"},
		"warn": {"warn", "bg-[#FFF6E5] text-[#8A5A00]"},
		"bad":  {"bad", "bg-[#FDECEC] text-[#B42318]"},
		"info": {"info", "bg-[#EAF0FF] text-[#0036C4]"},
		"":     {"info", "bg-[#EAF0FF] text-[#0036C4]"},
		"raro": {"info", "bg-[#EAF0FF] text-[#0036C4]"},
	}
	for in, want := range cases {
		body := render(t, context.Background(), card.Callout(card.CalloutOptions{Tone: in, Text: "Timbrado"}))
		assert.Contains(t, body, `data-kiban-callout="`+want.tone+`"`, "tono %q", in)
		assert.Contains(t, body, want.class, "tono %q", in)
	}
}

// Check cambia el ícono de alerta por una palomita; Attrs van al contenedor.
func TestCallout_PalomitaYAtributos(t *testing.T) {
	alerta := render(t, context.Background(), card.Callout(card.CalloutOptions{Text: "El timbrado falló"}))
	assert.Contains(t, alerta, `<circle cx="12" cy="12" r="10">`)
	assert.NotContains(t, alerta, `points="20 6 9 17 4 12"`)
	assert.Contains(t, alerta, `<span class="min-w-0">El timbrado falló</span>`)

	ok := render(t, context.Background(), card.Callout(card.CalloutOptions{Tone: "ok", Text: "Llegó", Check: true, Attrs: templ.Attributes{"role": "status"}}))
	assert.Contains(t, ok, `points="20 6 9 17 4 12"`)
	assert.NotContains(t, ok, `<circle`)
	assert.Contains(t, ok, `role="status"`)
}
