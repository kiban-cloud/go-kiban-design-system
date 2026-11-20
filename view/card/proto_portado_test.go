package card_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// Piezas portadas de kiban-proto: apagado flamboyant, lo de siempre;
// prendido, lo del proto.

var encendido = feature.Enable(context.Background(), flamboyant.Flag)

func pintar(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

var _ = strings.Contains

func TestBento_SinKpis_ApagadoDejaLaRejillaPrendidoNo(t *testing.T) {
	o := card.BentoOptions{Big: true, Title: "B", Trend: []float64{1, 2}}
	const kpis = "md:flex-col md:flex-nowrap"
	assert.Contains(t, pintar(t, context.Background(), card.Bento(o)), kpis, "apagado, como siempre")
	assert.NotContains(t, pintar(t, encendido, card.Bento(o)), kpis)
}

func TestBento_BodyEnLugarDeLaTendenciaYDosFilas(t *testing.T) {
	o := card.BentoOptions{Big: true, Rows: 2, Title: "B", Trend: []float64{1, 2}, Body: templ.Raw(`<div data-cuerpo></div>`)}
	body := pintar(t, encendido, card.Bento(o))
	assert.Contains(t, body, "data-cuerpo")
	assert.Contains(t, body, "xl:row-span-2")
	assert.NotContains(t, body, "<polyline", "Body reemplaza la tendencia")
	assert.NotContains(t, pintar(t, encendido, card.Bento(card.BentoOptions{Title: "B"})), "xl:row-span-2")
}
