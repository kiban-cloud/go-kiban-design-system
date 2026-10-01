package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Un campo marcado como obligatorio tiene que DECIRLO en el HTML, no sólo
// pintar el asterisco.
//
// Existe porque pasó: Text, Textarea y Select recibían `required` y lo único
// que hacían con él era el asterisco. Sin el atributo, ni el navegador ni htmx
// paran el envío, así que un formulario incompleto salía igual, el servidor lo
// rechazaba y el usuario veía el error después de haber esperado.
func TestInputs_ObligatorioEmiteElAtributo(t *testing.T) {
	render := func(t *testing.T, c templ.Component) string {
		t.Helper()
		var buf bytes.Buffer
		require.NoError(t, c.Render(context.Background(), &buf))
		return buf.String()
	}

	t.Run("text", func(t *testing.T) {
		assert.Contains(t, render(t, input.Text("nombre", "Nombre", "", "", "", true)), " required")
		assert.NotContains(t, render(t, input.Text("nombre", "Nombre", "", "", "", false)), " required")
	})

	t.Run("textarea", func(t *testing.T) {
		con := render(t, input.Textarea("notas", "Notas", "", "", "", true, 4, "", false))
		sin := render(t, input.Textarea("notas", "Notas", "", "", "", false, 4, "", false))
		assert.Contains(t, con, " required")
		assert.NotContains(t, sin, " required")
	})

	t.Run("select", func(t *testing.T) {
		ops := []input.SelectOption{{Value: "a", Label: "A"}}
		assert.Contains(t, render(t, input.Select("tipo", "Tipo", "", "", "", ops, true)), " required")
		assert.NotContains(t, render(t, input.Select("tipo", "Tipo", "", "", "", ops, false)), " required")
	})
}

// El asterisco y el atributo van juntos: enseñar uno sin el otro es prometerle
// al usuario una validación que no existe, o al revés, bloquearlo sin decirle
// por qué.
func TestInputs_AsteriscoYAtributoVanJuntos(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, input.Text("x", "Campo", "", "", "", true).Render(context.Background(), &buf))
	html := buf.String()

	assert.True(t, strings.Contains(html, "*"), "falta el asterisco")
	assert.True(t, strings.Contains(html, " required"), "falta el atributo")
}
