package input_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func picker() input.PickerOptions {
	return input.PickerOptions{
		ID: "clientes", Name: "customer", Selected: "c2", Label: "Cliente",
		Search: "Buscar cliente", Empty: "Ningún cliente coincide", Required: true,
		Attrs: templ.Attributes{"hx-post": "/preview"},
		Options: []input.PickerOption{
			{Value: "c1", Title: "Ana López", Sub: "ana@x.mx", Tile: templ.Raw(`<i>AL</i>`), Search: "RFC LOAA800101"},
			{Value: "c2", Title: "Luis Pérez"},
		},
	}
}

// Radios con el mismo name, la elegida marcada, buscador sin name y el texto
// de búsqueda en minúsculas.
func TestPicker_RadiosBuscadorYElegida(t *testing.T) {
	body := renderIn(t, context.Background(), input.Picker(picker()))

	assert.Contains(t, body, `id="clientes"`)
	assert.Contains(t, body, "data-kiban-picker")
	assert.Contains(t, body, `type="search" placeholder="Buscar cliente" aria-label="Buscar cliente"`)
	assert.Contains(t, body, "data-kiban-picker-q data-nodirty", "el buscador no cuenta como cambio")
	assert.NotContains(t, body, `type="search" name=`, "el buscador no viaja con el form")
	assert.Contains(t, body, `role="radiogroup" aria-label="Cliente"`)
	assert.Equal(t, 2, strings.Count(body, `type="radio" name="customer"`))
	assert.Contains(t, body, `value="c2" checked`)
	assert.NotContains(t, body, `value="c1" checked`)
	assert.Equal(t, 2, strings.Count(body, " required"), "cada radio pide elegir")
	assert.Equal(t, 2, strings.Count(body, `hx-post="/preview"`), "Attrs en cada radio")
	assert.Contains(t, body, `data-q="ana lópez ana@x.mx rfc loaa800101"`)
	assert.Contains(t, body, `data-q="luis pérez  "`)
	assert.Contains(t, body, "<i>AL</i>")
	assert.Contains(t, body, "grid-cols-[40px_minmax(0,1fr)]", "con placa")
	assert.Contains(t, body, "grid-cols-[minmax(0,1fr)]", "sin placa")
	assert.Contains(t, body, ">ana@x.mx</small>")
	assert.Equal(t, 1, strings.Count(body, "<small"), "sin Sub no hay línea")
	assert.Contains(t, body, `data-kiban-picker-empty hidden>Ningún cliente coincide</p>`)
	assert.NotContains(t, body, "grid-cols-2", "una columna por defecto")
}

func TestPicker_SinBuscadorNiVacioConErrorYDosColumnas(t *testing.T) {
	o := picker()
	o.ID, o.Search, o.Empty, o.Label = "", "", "", ""
	o.Columns, o.Error = 2, "Elige un cliente"
	body := renderIn(t, context.Background(), input.Picker(o))

	assert.NotContains(t, body, `type="search"`)
	assert.NotContains(t, body, "data-kiban-picker-empty")
	assert.NotContains(t, body, "id=")
	assert.NotContains(t, body, "aria-label=")
	assert.Contains(t, body, "grid-cols-2 max-sm:grid-cols-1")
	assert.Contains(t, body, ">Elige un cliente</p>")
}
