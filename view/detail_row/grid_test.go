package detail_row_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/detail_row"
)

func renderGrid(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

var gridItems = []detail_row.GridItem{
	{Label: "Banco", Value: "STP"},
	{Label: "CLABE", Value: "646180000000000001", Mono: true},
	{Label: "Referencia"},
}

// Grid: tres columnas, etiqueta en versalitas y los datos sin valor no se
// muestran.
func TestGrid_TresColumnasSinVacios(t *testing.T) {
	body := renderGrid(t, detail_row.Grid(gridItems))

	assert.Contains(t, body, "data-kiban-detail-grid")
	assert.Contains(t, body, "grid-cols-3")
	assert.NotContains(t, body, "grid-cols-2")
	assert.Contains(t, body, ">Banco</dt>")
	assert.Contains(t, body, "STP</dd>")
	assert.Contains(t, body, "uppercase", "la etiqueta chica en mayúsculas")
	assert.Equal(t, 2, strings.Count(body, "<dt"))
	assert.NotContains(t, body, "Referencia", "sin valor no se muestra")
	assert.Equal(t, 1, strings.Count(body, "font-mono"), "sólo la CLABE en monoespaciado")
}

// GridWith: dos columnas, etiqueta de campo y una raya en lugar del vacío.
func TestGridWith_DosColumnasEtiquetaDeCampoYVacio(t *testing.T) {
	body := renderGrid(t, detail_row.GridWith(detail_row.GridOptions{Items: gridItems, Cols: 2, FieldLabels: true, Empty: "—"}))

	assert.Contains(t, body, "grid-cols-2")
	assert.NotContains(t, body, "grid-cols-3")
	assert.Contains(t, body, `<dt class="text-xs font-medium leading-[18px] text-[#474F66]">Banco</dt>`)
	assert.NotContains(t, body, "uppercase")
	assert.Equal(t, 3, strings.Count(body, "<dt"), "con Empty se muestran todos")
	assert.Contains(t, body, ">Referencia</dt>")
	assert.Contains(t, body, "—</dd>")
}

// Cualquier Cols que no sea 2 son tres columnas.
func TestGridWith_ColsDistintoDeDosSonTres(t *testing.T) {
	for _, n := range []int{0, 1, 3, 4} {
		assert.Contains(t, renderGrid(t, detail_row.GridWith(detail_row.GridOptions{Items: gridItems, Cols: n})), "grid-cols-3", "Cols=%d", n)
	}
}

func TestGrid_SinDatosEsLaListaVacia(t *testing.T) {
	body := renderGrid(t, detail_row.Grid(nil))
	assert.Contains(t, body, "<dl")
	assert.NotContains(t, body, "<dt")
}
