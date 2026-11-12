package card_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

// La franja de cifras de una lista, como la .b-strip del proto: número en
// negritas, etiqueta, tono y un enlace opcional al final.
func TestStrip_CifrasConTonoYExtra(t *testing.T) {
	body := render(t, context.Background(), card.Strip([]card.StripItem{
		{Value: "7", Label: "activas"},
		{Value: "2", Label: "bloqueos", Tone: "warn"},
		{Value: "1", Label: "fuera de tus dominios", Extra: templ.Raw(`<a href="/auth">¿por qué?</a>`)},
	}))

	assert.Contains(t, body, `<b class="font-semibold text-kiban-ink">7</b> activas`)
	assert.Contains(t, body, "text-[#B26B00]", "el tono warn pinta el número")
	assert.Contains(t, body, `<a href="/auth">¿por qué?</a>`)
	assert.Contains(t, body, "bg-[#F4F7FF]")
}

func TestStrip_SinCifrasNoPintaNada(t *testing.T) {
	assert.Empty(t, render(t, context.Background(), card.Strip(nil)))
}
