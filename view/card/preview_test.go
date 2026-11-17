package card_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

func TestPreview_CajaConAtributos(t *testing.T) {
	body := render(t, context.Background(), card.Preview(templ.Attributes{"id": "resumen", "hx-post": "/rekon/orders/preview"}))

	assert.Contains(t, body, "data-kiban-preview")
	assert.Contains(t, body, `aria-live="polite"`)
	assert.Contains(t, body, `id="resumen"`)
	assert.Contains(t, body, `hx-post="/rekon/orders/preview"`)
	assert.Contains(t, body, "<p>cuerpo</p>")
}

func TestPreviewBody_MontoMonedaYRenglones(t *testing.T) {
	body := render(t, context.Background(), card.PreviewBody(card.PreviewOptions{
		Eyebrow: "Tienda Ana", Amount: "$1,500.00", Currency: "MXN",
		Rows: []card.PreviewRow{{Label: "Concepto", Value: "Renta"}, {Label: "CLABE", Value: "012180001234567891", Mono: true}},
	}))

	assert.Contains(t, body, ">Tienda Ana</span>")
	assert.Contains(t, body, "$1,500.00")
	assert.Contains(t, body, ">MXN</span>")
	assert.Contains(t, body, ">Concepto</span>")
	assert.Contains(t, body, ">Renta</b>")
	assert.Equal(t, 1, strings.Count(body, "font-mono"), "sólo la CLABE en monoespaciado")
}

func TestPreviewBody_SinAntetituloNiMoneda(t *testing.T) {
	body := render(t, context.Background(), card.PreviewBody(card.PreviewOptions{Amount: "$0.00"}))

	assert.NotContains(t, body, "uppercase", "sin Eyebrow no hay antetítulo")
	assert.NotContains(t, body, "text-sm font-normal", "sin Currency no hay moneda")
	assert.Contains(t, body, "$0.00")
}
