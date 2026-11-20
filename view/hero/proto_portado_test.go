package hero_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
	"github.com/kiban-cloud/go-kiban-design-system/view/hero"
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

func TestCycle_ChipSinPlaca_MargenSoloEnFlamboyant(t *testing.T) {
	o := hero.CycleOptions{Space: "Esp", Title: "Hola"}
	assert.Contains(t, pintar(t, context.Background(), hero.Cycle(o)), "py-[3px] pl-1 pr-2.5", "apagado, como siempre")
	assert.Contains(t, pintar(t, encendido, hero.Cycle(o)), "py-[3px] pl-2.5 pr-2.5")
	o.SpaceTile = templ.Raw(`<i></i>`)
	assert.Contains(t, pintar(t, encendido, hero.Cycle(o)), "py-[3px] pl-1 pr-2.5", "con placa, pegado a ella")
}

func TestCycle_AccionTendenciaEIconoDePendientes(t *testing.T) {
	o := hero.CycleOptions{Title: "Hola", PendingCount: 2, PendingLabel: "por revisar", PendingHref: "/p"}
	sin := pintar(t, encendido, hero.Cycle(o))
	assert.NotContains(t, sin, "gap-2.5\"><a", "sin Action, el botón de pendientes va suelto")
	assert.Contains(t, sin, `<h1 class="mb-1 mt-2.5 text-[34px] font-semibold leading-10 tracking-[-.02em] max-[900px]:text-[28px] max-[900px]:leading-[34px]">Hola</h1>`)

	o.Action = &hero.CycleAction{Label: "Nueva ejecución", Href: "/workfloo/nueva"}
	o.TitleTrend, o.TitleTrendDown = "▼ 3 %", true
	o.PendingIcon = "clock"
	body := pintar(t, encendido, hero.Cycle(o))
	assert.Contains(t, body, `href="/workfloo/nueva"`)
	assert.Contains(t, body, "Nueva ejecución")
	assert.Contains(t, body, "▼ 3 %")
	assert.Contains(t, body, "text-[#FFD6D7]", "baja en rojo")
	assert.Contains(t, body, `<polyline points="12 6 12 12 16 14"/>`, "el ícono pedido en vez de la campana")
}

func TestDetail_TituloRenombrableYMeta(t *testing.T) {
	o := hero.DetailOptions{BackHref: "/ab", BackLabel: "A/B tests", Eyebrow: "A/B test", Title: "Prueba 1",
		Rename: &hero.DetailRename{Action: "/ab/1/nombre", Name: "name", Label: "Nombre"},
		Meta:   []templ.Component{hero.DetailText("Creado el 3 sep"), hero.DetailChip(hero.Chip{Label: "Flujo", Href: "/f"})}}
	body := pintar(t, encendido, hero.Detail(o))
	assert.Contains(t, body, `hx-post="/ab/1/nombre"`)
	assert.Contains(t, body, `hx-swap="none"`)
	assert.Contains(t, body, `value="Prueba 1"`)
	assert.Contains(t, body, "Creado el 3 sep")
	assert.Contains(t, body, `href="/f"`)
	assert.Equal(t, 1, strings.Count(body, "·"), "un separador entre dos datos")

	o.Rename = nil
	assert.Contains(t, pintar(t, encendido, hero.Detail(o)), ">Prueba 1</h1>")
}
