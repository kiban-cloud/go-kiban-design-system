package detail_row_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/detail_row"
)

func TestColumns_GruposCamposYVacios(t *testing.T) {
	var buf bytes.Buffer
	groups := []detail_row.Group{
		{Title: "Contacto", Fields: []detail_row.Field{
			{Label: "Correo", Text: "ana@ejemplo.com"},
			{Label: "Teléfono"},
			{Label: "Id", Value: templ.Raw(`<code>abc</code>`)},
		}},
		{Title: "Cobro", Empty: "Sin métodos de cobro"},
	}
	require.NoError(t, detail_row.Columns(groups).Render(context.Background(), &buf))
	body := buf.String()

	assert.Contains(t, body, ">Contacto</h3>")
	assert.Contains(t, body, ">Correo</dt>")
	assert.Contains(t, body, "ana@ejemplo.com")
	assert.Contains(t, body, "—", "un campo sin valor muestra la raya")
	assert.Contains(t, body, "<code>abc</code>")
	assert.Contains(t, body, ">Sin métodos de cobro</p>")
}
