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

// Como el menú de filtros del proto: cada opción con cuántos quedarían, un
// buscador dentro del menú cuando hay muchas (etiquetas), «Seleccionar
// todos» y «Limpiar» al pie, y los conteos que se actualizan con cada
// respuesta (oob) porque la barra no se vuelve a pintar.
func TestBar_MenuComoElProto(t *testing.T) {
	c := cfg()
	c.MenuSelectAll, c.MenuClear = "Seleccionar todos", "Limpiar"
	c.Filters[0].Options[0].Count, c.Filters[0].Options[0].HasCount = 12, true
	c.Filters[0].Options[1].Count, c.Filters[0].Options[1].HasCount = 0, true
	c.Filters = append(c.Filters, filterbar.Filter{
		Name: "tags", Label: "Etiquetas", Multi: true,
		Search:  filterbar.MenuSearch{Placeholder: "Buscar etiqueta", Empty: "Nada coincide con «{q}»."},
		Options: []filterbar.Option{{Value: "premium", Label: "premium", Selected: true}},
	})
	body := render(t, filterbar.Bar(c))

	assert.Contains(t, body, `<span id="f-n-type-5046" class="text-xs tabular-nums text-[#9096A3]" data-kiban-filter-count>12</span>`)
	assert.Contains(t, body, `<span id="f-n-type-504d" class="text-xs tabular-nums text-[#9096A3]" data-kiban-filter-count>0</span>`, "un cero real también se ve")
	assert.NotContains(t, body, `id="f-n-created`, "sin HasCount no hay conteo")
	assert.Contains(t, body, `placeholder="Buscar etiqueta"`)
	assert.Contains(t, body, `data-kiban-filter-q`)
	assert.Contains(t, body, `[&[hidden]]:hidden" data-kiban-filter-opt`, "el buscador oculta opciones con hidden y la clase flex no debe ganarle")
	assert.Contains(t, body, `data-empty="Nada coincide con «{q}»."`)
	assert.Contains(t, body, `data-kiban-filter-all>Seleccionar todos</button>`)
	assert.Contains(t, body, `data-kiban-filter-none hidden>Limpiar</button>`, "sin nada marcado, limpiar se oculta")
	assert.Contains(t, body, `data-kiban-filter-none>Limpiar</button>`, "con algo marcado se ve")
	assert.Contains(t, body, `change[!event.target.closest('[data-kiban-filter-q],[data-kiban-filter-manual]')]`, "escribir en el buscador del menú no pide la lista")
}

func TestCounts_ActualizaLosConteosOOB(t *testing.T) {
	c := cfg()
	c.Filters[0].Options[0].Count, c.Filters[0].Options[0].HasCount = 3, true
	body := render(t, filterbar.Counts(c))

	assert.Contains(t, body, `<span id="f-n-type-5046" class="text-xs tabular-nums text-[#9096A3]" data-kiban-filter-count hx-swap-oob="true">3</span>`)
	assert.NotContains(t, body, `f-n-created`)
}

func TestBar_SinTextosDeMenuNoHayPie(t *testing.T) {
	body := render(t, filterbar.Bar(cfg()))
	assert.NotContains(t, body, "data-kiban-filter-all")
	assert.NotContains(t, body, "data-kiban-filter-none")
	assert.NotContains(t, body, "data-kiban-filter-count")
}

func TestBar_OpcionConIcono(t *testing.T) {
	c := cfg()
	c.Filters[0].Options[0].Icon = templ.Raw(`<img src="/logo-buro.png" alt="">`)
	body := render(t, filterbar.Bar(c))

	assert.Contains(t, body, `<img src="/logo-buro.png" alt="">`, "el ícono va en la opción")
	assert.Equal(t, 1, bytes.Count([]byte(body), []byte(`aria-hidden="true"><img`)), "sólo la opción con ícono lo lleva")
}

func numCfg(value, value2, op string) filterbar.Config {
	c := cfg()
	c.MenuClear = "Quitar filtro"
	c.Filters = []filterbar.Filter{{Name: "req", Label: "Requests", Numeric: &filterbar.Numeric{
		OpName: "req_op", Name: "req", Name2: "req2", Op: op, Value: value, Value2: value2,
		Unit: "ms", Step: "100",
		Ops: []filterbar.NumericOp{{Value: "eq", Label: "Igual a", Symbol: "="}, {Value: "gt", Label: "Mayor que", Symbol: ">"},
			{Value: "lt", Label: "Menor que", Symbol: "<"}, {Value: "bt", Label: "Entre", Symbol: "entre"}},
		Presets:     []filterbar.NumericPreset{{Op: "gt", Value: "1000", Label: "Más de 1,000"}},
		Placeholder: "Valor", From: "Desde", To: "Hasta", And: "y", PresetsLabel: "Rápidos", Apply: "Aplicar",
		ErrNumber: "Escribe un número.", ErrRange: "Escribe los dos valores del rango.", ErrMax: "El máximo es {max}.",
	}}}
	return c
}

func TestBar_FiltroNumericoSinValor(t *testing.T) {
	body := render(t, filterbar.Bar(numCfg("", "", "")))

	assert.Contains(t, body, `data-kiban-filter-numeric`)
	assert.Contains(t, body, `<span data-kiban-filter-text>Requests</span>`, "sin filtro, sólo el nombre")
	assert.Contains(t, body, `name="req_op" value="" disabled`, "sin valor los ocultos no viajan")
	assert.Contains(t, body, `name="req" value="" disabled`)
	assert.Contains(t, body, `data-op="gt"`, "sin filtro, mayor que como el proto")
	assert.Contains(t, body, `data-kiban-num-opbtn="gt" data-symbol="&gt;" aria-pressed="true"`)
	assert.Contains(t, body, `data-kiban-filter-manual`, "lo tecleado no pide la lista")
	assert.Contains(t, body, `data-kiban-num-preset data-op="gt" data-v="1000">Más de 1,000</button>`)
	assert.Contains(t, body, `data-kiban-num-clear hidden>Quitar filtro</button>`)
	assert.NotContains(t, body, `data-kiban-filter-clear>`, "sin filtros, limpiar oculto")
}

func TestBar_FiltroNumericoAplicado(t *testing.T) {
	body := render(t, filterbar.Bar(numCfg("1000", "", "gt")))
	assert.Contains(t, body, `Requests: <b class="font-semibold">&gt; 1,000 ms</b>`)
	assert.Contains(t, body, `name="req_op" value="gt" data-kiban-num-op`)
	assert.Contains(t, body, `name="req" value="1000" data-kiban-num-v`)
	assert.Contains(t, body, `name="req2" value="" disabled`, "el segundo valor sólo con «entre»")
	assert.Contains(t, body, `data-kiban-filter-pill data-on`)
	assert.Contains(t, body, `data-kiban-num-clear>Quitar filtro</button>`)

	body = render(t, filterbar.Bar(numCfg("2500.5", "1000", "bt")))
	assert.Contains(t, body, `entre 1,000 y 2,500.5 ms`, "el menor primero")
	assert.Contains(t, body, `name="req2" value="1000" data-kiban-num-v2`)
	assert.Contains(t, body, `placeholder="Desde"`)
}

func TestBar_NotaVacioYPista(t *testing.T) {
	c := cfg()
	c.Filters = []filterbar.Filter{
		{Name: "renew", Label: "Autorenovación", Multi: true, Note: "Sólo aplica a Buró de Crédito", Options: []filterbar.Option{{Value: "on", Label: "Prendida"}}},
		{Name: "sub", Label: "Funcionalidad", Multi: true, Empty: "Elige primero un servicio"},
		{Name: "sub2", Label: "Funcionalidad", Multi: true, Options: []filterbar.Option{{Value: "sat|cfdi", Label: "CFDI", Hint: "SAT"}}},
	}
	body := render(t, filterbar.Bar(c))
	assert.Contains(t, body, `Sólo aplica a Buró de Crédito</p>`)
	assert.Contains(t, body, `Elige primero un servicio</div>`)
	assert.Contains(t, body, `CFDI<span class="ml-1 text-[11px] font-normal text-[#9096A3] [&[hidden]]:hidden">SAT</span>`)
}

// Los dos menús (opciones y numérico) llevan data-kiban-filter-menu: es lo que
// base.templ mide al abrir para correrlo dentro de la pantalla. Sin la marca,
// el menú de una píldora cerca del borde derecho se sale del viewport.
func TestBar_LosMenusSeMarcanParaAcomodarseEnLaPantalla(t *testing.T) {
	assert.Contains(t, render(t, filterbar.Bar(cfg())), `<div data-kiban-filter-menu class="absolute left-0`)
	assert.Contains(t, render(t, filterbar.Bar(numCfg("", "", ""))), `data-kiban-filter-menu class="absolute left-0`)
}
