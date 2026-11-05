package card_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

var flam = feature.Enable(context.Background(), flamboyant.Flag)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(templ.WithChildren(ctx, templ.Raw(`<p>cuerpo</p>`)), &buf))
	return buf.String()
}

func TestStep_ClasicoSoloTitulo(t *testing.T) {
	body := render(t, context.Background(), card.Step(card.StepOptions{Title: "Datos generales"}))

	assert.Contains(t, body, `<div class="font-heading font-medium text-base">Datos generales</div>`)
	assert.Contains(t, body, "rounded-md p-6", "es un Card de siempre")
	assert.Contains(t, body, "<p>cuerpo</p>")
	assert.NotContains(t, body, "data-kiban-step")
}

func TestStep_ClasicoConAyuda(t *testing.T) {
	body := render(t, context.Background(), card.Step(card.StepOptions{Title: "Fiscales", HintText: "Para facturar"}))

	assert.Contains(t, body, `<div><div class="font-heading font-medium text-base">Fiscales</div><div class="text-xs text-kiban-ink3 mt-1">Para facturar</div></div>`)
}

func TestStep_ClasicoConAccion(t *testing.T) {
	body := render(t, context.Background(), card.Step(card.StepOptions{Title: "Contactos", Right: templ.Raw(`<button>Agregar</button>`)}))

	assert.Contains(t, body, `<div class="flex items-center justify-between"><div class="font-heading font-medium text-base">Contactos</div><button>Agregar</button></div>`)
}

func TestStep_FlamboyantNumeradoConRequeridos(t *testing.T) {
	body := render(t, flam, card.Step(card.StepOptions{Key: "fiscal", Title: "Fiscales", HintText: "Para facturar", Required: []string{"rfc", "regimen"}}))

	assert.Contains(t, body, `data-kiban-step="fiscal"`)
	assert.Contains(t, body, `data-kiban-step-required="rfc regimen"`)
	assert.Contains(t, body, "data-kiban-step-n")
	assert.Contains(t, body, "Para facturar")
	assert.Contains(t, body, "<p>cuerpo</p>")
	assert.Contains(t, body, "rounded-[10px]")
}

func TestCard_FlamboyantCambiaBordeYRadio(t *testing.T) {
	assert.Contains(t, render(t, context.Background(), card.Card("")), "bg-white border rounded-md p-6 space-y-5 border-kiban-border")
	body := render(t, flam, card.Card(""))
	assert.Contains(t, body, "rounded-[10px]")
	assert.Contains(t, body, "border-[#CDD3DF]")
}
