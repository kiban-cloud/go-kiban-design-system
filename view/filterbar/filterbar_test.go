package filterbar_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/filterbar"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

func cfg() filterbar.Config {
	return filterbar.Config{
		ID: "f", Action: "/crm/customers", Target: "#lista", Indicator: "#cargando",
		Search: filterbar.Search{Name: "q", Placeholder: "Buscar", AriaLabel: "Buscar clientes"},
		Filters: []filterbar.Filter{
			{Name: "type", Label: "Tipo", Multi: true, Options: []filterbar.Option{{Value: "PF", Label: "Física"}, {Value: "PM", Label: "Moral"}}},
			{Name: "created", Label: "Creado", Options: []filterbar.Option{{Value: "", Label: "Cualquier fecha", Selected: true}, {Value: "7", Label: "Últimos 7 días"}}},
		},
		Hidden:     []filterbar.Field{{Name: "itemsPerPage", Value: "25"}},
		ClearHref:  "/crm/customers",
		ClearLabel: "Limpiar filtros",
	}
}

func TestBar_FormularioQuePideLaLista(t *testing.T) {
	body := render(t, filterbar.Bar(cfg()))

	assert.Contains(t, body, `hx-get="/crm/customers"`)
	assert.Contains(t, body, `hx-target="#lista"`)
	assert.Contains(t, body, `type="search" name="q"`)
	assert.Contains(t, body, `id="f-state"`)
	assert.Contains(t, body, `name="itemsPerPage" value="25"`)
	assert.Contains(t, body, `type="checkbox" name="type" value="PF"`)
	assert.Contains(t, body, `type="radio" name="created" value="7"`)
}

func TestBar_SinFiltrosLaPildoraSoloDiceSuNombreYLimpiarSeOculta(t *testing.T) {
	body := render(t, filterbar.Bar(cfg()))

	assert.Contains(t, body, `data-kiban-filter-text>Tipo</span>`)
	assert.Contains(t, body, `data-kiban-filter-text>Creado</span>`, "la opción vacía no cuenta como filtro")
	assert.NotContains(t, body, "data-on")
	assert.Contains(t, body, `data-kiban-filter-clear hidden`)
}

func TestBar_UnaOpcionMarcada(t *testing.T) {
	c := cfg()
	c.Filters[0].Options[1].Selected = true
	body := render(t, filterbar.Bar(c))

	assert.Contains(t, body, `Tipo: <b class="font-semibold">Moral</b>`)
	assert.Contains(t, body, "data-on")
	assert.NotContains(t, body, `data-kiban-filter-clear hidden`)
}

func TestBar_VariasOpciones(t *testing.T) {
	c := cfg()
	c.Filters[0].Options[0].Selected = true
	c.Filters[0].Options[1].Selected = true
	assert.Contains(t, render(t, filterbar.Bar(c)), `Tipo <b class="font-semibold">· 2</b>`)
}

func TestBar_BusquedaActivaMuestraLimpiar(t *testing.T) {
	c := cfg()
	c.Search.Value = "acme"
	body := render(t, filterbar.Bar(c))

	assert.Contains(t, body, `value="acme"`)
	assert.NotContains(t, body, `data-kiban-filter-clear hidden`)
}

func TestState_FueraDeBanda(t *testing.T) {
	assert.Contains(t, render(t, filterbar.State(cfg(), true)), `hx-swap-oob="true"`)
	assert.NotContains(t, render(t, filterbar.State(cfg(), false)), "hx-swap-oob")
}
