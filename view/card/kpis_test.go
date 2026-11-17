package card_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

// Las cifras de arriba de rekon: etiqueta con ícono, valor con unidad y
// detalle; la destacada va en azul kiban.
func TestKpis_CifrasConUnidadDetalleYDestacada(t *testing.T) {
	body := render(t, context.Background(), card.Kpis([]card.KpiOptions{
		{Label: "Por cobrar", Icon: templ.Raw(`<svg id="i"></svg>`), Value: "$11,109.16", Unit: "MXN", Detail: "3 órdenes pendientes", Featured: true},
		{Label: "Vencido sin pagar", Value: "$0.00"},
	}))

	assert.Contains(t, body, "data-kiban-kpis")
	assert.Equal(t, 2, strings.Count(body, "data-kiban-kpi>"), "una tarjeta por cifra")
	assert.Contains(t, body, `<svg id="i"></svg>`)
	assert.Contains(t, body, "Por cobrar")
	assert.Contains(t, body, "$11,109.16")
	assert.Contains(t, body, ">MXN</small>")
	assert.Contains(t, body, ">3 órdenes pendientes</span>")
	assert.Equal(t, 1, strings.Count(body, "border-kiban-primary bg-kiban-primary"), "sólo la destacada va en azul")
	assert.Equal(t, 1, strings.Count(body, "<small"), "sin unidad no hay <small>")
}

func TestKpis_SinCifrasNoPintaNada(t *testing.T) {
	assert.Empty(t, render(t, context.Background(), card.Kpis(nil)))
}
