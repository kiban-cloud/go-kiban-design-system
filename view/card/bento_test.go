package card_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

func bento() card.BentoOptions {
	return card.BentoOptions{
		Color: "#8000FF", Soft: "#F3E8FF", Icon: templ.Raw(`<svg id="pagos"></svg>`),
		Title: "Pagos", Text: "Lo que entra a tus cuentas", Pending: "3 sin conciliar",
		Kpis:  []card.BentoKpi{{Value: "$12,300.00", Label: "Hoy"}, {Value: "41", Label: "Pagos"}},
		Trend: []float64{1, 2, 3}, TrendUnit: "miles de pesos", Axis: []string{"hace 7 días", "hoy"},
		LiveTitle: "Por revisar",
		Live: []card.BentoRow{
			{Title: "Pago de Ana", Sub: "SPEI", When: "hace 2 min", Href: "/rekon/payments/1"},
			{Title: "Pago de Luis", Sub: "OXXO", When: "hace 5 min", Href: "#", Attrs: templ.Attributes{"hx-get": "/rekon/payments/2"}},
		},
		Open:  card.BentoLink{Label: "Abrir Pagos", Href: "/rekon/payments"},
		Links: []card.BentoLink{{Label: "Exportar CSV", Href: "/rekon/payments.csv", Attrs: templ.Attributes{"download": ""}}},
	}
}

// La tarjeta completa: color en variables, cifras, tendencia, lista y pie.
func TestBento_TarjetaCompleta(t *testing.T) {
	body := render(t, context.Background(), card.Bento(bento()))

	assert.Contains(t, body, "data-kiban-bento")
	assert.Contains(t, body, `style="--c:#8000FF;--s:#F3E8FF;"`)
	assert.Contains(t, body, `<svg id="pagos"></svg>`)
	assert.Contains(t, body, ">Pagos</h2>")
	assert.Contains(t, body, ">Lo que entra a tus cuentas</p>")
	assert.Contains(t, body, ">3 sin conciliar</span>")
	assert.Contains(t, body, ">$12,300.00</b>")
	assert.Contains(t, body, ">Hoy</small>")
	assert.Contains(t, body, "miles de pesos")
	assert.Contains(t, body, "data-kiban-spark", "la tendencia es un chart.Spark")
	assert.Contains(t, body, `viewBox="0 0 220 56"`, "la chica traza en 220×56")
	assert.Contains(t, body, "<span>hace 7 días</span><span>hoy</span>")
	assert.Contains(t, body, "Por revisar")
	assert.Contains(t, body, `href="/rekon/payments/1"`)
	assert.Contains(t, body, `hx-get="/rekon/payments/2"`)
	assert.Contains(t, body, ">hace 2 min</time>")
	assert.Contains(t, body, `href="/rekon/payments"`)
	assert.Contains(t, body, "Abrir Pagos")
	assert.Contains(t, body, `href="/rekon/payments.csv" download=""`)
	assert.Contains(t, body, "xl:col-span-2", "Cols cero ocupa dos columnas")
	assert.NotContains(t, body, "md:col-span-2", "sólo la grande ocupa dos columnas desde md")
}

// La grande: dos columnas desde md, título más grande y la tendencia más ancha.
func TestBento_Grande(t *testing.T) {
	o := bento()
	o.Big, o.Cols = true, 4
	body := render(t, context.Background(), card.Bento(o))

	assert.Contains(t, body, "xl:col-span-4")
	assert.Contains(t, body, "md:col-span-2")
	assert.Contains(t, body, "text-[22px] leading-7", "título grande")
	assert.Contains(t, body, `viewBox="0 0 420 96"`)
	assert.Contains(t, body, "md:text-[32px]", "cifras grandes")
}

func TestBento_ColumnasEnPantallaAncha(t *testing.T) {
	for cols, want := range map[int]string{0: "xl:col-span-2", 2: "xl:col-span-2", 3: "xl:col-span-3", 4: "xl:col-span-4", 5: "xl:col-span-2", 6: "xl:col-span-6"} {
		o := bento()
		o.Cols = cols
		assert.Contains(t, render(t, context.Background(), card.Bento(o)), want, "Cols=%d", cols)
	}
}

// Sin pendiente, sin renglones, sin botón ni enlaces: nada de eso se pinta.
func TestBento_SinLoOpcional(t *testing.T) {
	body := render(t, context.Background(), card.Bento(card.BentoOptions{Color: "#0047FF", Soft: "#EAF0FF", Title: "Base", Text: "Tus clientes"}))

	assert.NotContains(t, body, "bg-[#CCF8EF]", "sin Pending no hay píldora")
	assert.NotContains(t, body, "animate-ping", "sin renglones no hay lista")
	assert.NotContains(t, body, "data-kiban-spark", "sin tendencia no hay trazo")
	assert.Equal(t, 0, strings.Count(body, "<a "), "sin Open ni Links no hay enlaces")
}

func TestBentoGrid_EnvuelveLasTarjetas(t *testing.T) {
	body := render(t, context.Background(), card.BentoGrid())
	assert.Contains(t, body, "data-kiban-bento-grid")
	assert.Contains(t, body, "xl:grid-cols-6")
	assert.Contains(t, body, "<p>cuerpo</p>")
}
