package hero_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/hero"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

func cycle() hero.CycleOptions {
	return hero.CycleOptions{
		Space: "Tienda Ana", SpaceTile: templ.Raw(`<i>ta</i>`),
		Title: "Buenas tardes", Text: "Esto es lo que pasa hoy",
		PendingCount: 3, PendingLabel: "por revisar", PendingHref: "/rekon/pendientes",
		StepsLabel: "Ciclo del dinero",
		Steps: []hero.CycleStep{
			{Icon: templ.Raw(`<svg id="cobra"></svg>`), Verb: "Cobra", Name: "Órdenes de pago", Sub: "$1,200", Badge: 2, Href: "/rekon/orders"},
			{Verb: "Recibe", Name: "Pagos", Href: "/rekon/payments"},
			{Verb: "Retira", Name: "Retiros", Href: "/rekon/withdrawals"},
			{Verb: "Dispersa", Name: "Dispersión", Href: "/rekon/dispersion"},
		},
	}
}

func TestCycle_SaludoChipYPasos(t *testing.T) {
	body := render(t, hero.Cycle(cycle()))

	assert.Contains(t, body, "data-kiban-hero")
	assert.Contains(t, body, "<i>ta</i>")
	assert.Contains(t, body, "Tienda Ana")
	assert.Contains(t, body, ">Buenas tardes</h1>")
	assert.Contains(t, body, ">Esto es lo que pasa hoy</p>")
	assert.Contains(t, body, `role="list" aria-label="Ciclo del dinero"`)
	assert.Equal(t, 4, strings.Count(body, `role="listitem"`))
	assert.Contains(t, body, `href="/rekon/orders"`)
	assert.Contains(t, body, `<svg id="cobra"></svg>`)
	assert.Contains(t, body, ">Cobra</span>")
	assert.Contains(t, body, ">Órdenes de pago</b>")
	assert.Contains(t, body, ">$1,200</small>")
	assert.Contains(t, body, ">2</i>", "el contador del paso")
	assert.Equal(t, 1, strings.Count(body, "bg-[#00DDB2]"), "sólo el paso con Badge > 0 lleva contador")
	assert.Contains(t, body, "--cy-cols:repeat(4,150px)")
	assert.Contains(t, body, "--cy-line:462px", "la línea mide (n-1)×154 px")
	assert.Contains(t, body, "--cy-inset:12.5%")
	assert.Contains(t, body, "--i:3", "cada paso entra escalonado")
	assert.NotContains(t, body, "kiban-cy-flow", "sin Animated la línea queda quieta")
}

// El botón de pendientes es un enlace; con PendingAttrs, un <button> que
// abre el panel «Pendientes» del shell.
func TestCycle_PendientesEnlaceOBoton(t *testing.T) {
	body := render(t, hero.Cycle(cycle()))
	assert.Contains(t, body, `<a href="/rekon/pendientes"`)
	assert.Contains(t, body, `<b class="font-bold">3</b> por revisar`)
	assert.NotContains(t, body, "data-fs-inbox-toggle")

	o := cycle()
	o.PendingAttrs = layout.InboxToggleAttrs("")
	body = render(t, hero.Cycle(o))
	assert.Contains(t, body, `<button type="button"`)
	assert.Contains(t, body, "data-fs-inbox-toggle")
	assert.Contains(t, body, `hx-get="/kiban-cloud/notifications/dropdown"`)
	assert.NotContains(t, body, `href="/rekon/pendientes"`)
}

func TestCycle_SinPendientesNiChipNiPasos(t *testing.T) {
	body := render(t, hero.Cycle(hero.CycleOptions{Title: "Hola"}))

	assert.Contains(t, body, ">Hola</h1>")
	assert.NotContains(t, body, "por revisar")
	assert.NotContains(t, body, "rounded-full bg-white pl-3", "sin PendingCount no hay botón")
	assert.NotContains(t, body, "border-white/20", "sin Space no hay chip")
	assert.NotContains(t, body, `role="list"`, "sin pasos no hay ciclo")
	assert.NotContains(t, body, "<p ", "sin Text no hay línea")
}

func TestCycle_Animado(t *testing.T) {
	o := cycle()
	o.Animated = true
	assert.Equal(t, 3, strings.Count(render(t, hero.Cycle(o)), `<i class="kiban-cy-flow">`))
}
