package input_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func renderIn(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

var flamCtx = feature.Enable(context.Background(), flamboyant.Flag)

func TestRadioCardWith_ClasicoIgualQueRadioCard(t *testing.T) {
	attrs := templ.Attributes{"hx-post": "/x"}
	antes := renderIn(t, context.Background(), input.RadioCard("type", "PF", "Física", "Persona", true, attrs))
	con := renderIn(t, context.Background(), input.RadioCardWith(input.RadioCardOptions{
		Name: "type", Value: "PF", Title: "Física", Subtitle: "Persona", Checked: true, Attrs: attrs, Icon: templ.Raw("<svg></svg>"),
	}))
	assert.Equal(t, antes, con, "el ícono sólo existe en flamboyant")
}

func TestRadioCardWith_FlamboyantIconoYTono(t *testing.T) {
	body := renderIn(t, flamCtx, input.RadioCardWith(input.RadioCardOptions{
		Name: "type", Value: "PF", Title: "Física", Checked: true, Icon: templ.Raw(`<svg id="i"></svg>`), Tone: "success",
	}))
	assert.Contains(t, body, `<svg id="i"></svg>`)
	assert.Contains(t, body, "bg-[#E3FBF5]")
	assert.Contains(t, body, "has-[:checked]:border-kiban-primary")
	assert.Contains(t, body, ` checked`)
}

func TestToggleWith_ClasicoIgualQueToggle(t *testing.T) {
	assert.Equal(t,
		renderIn(t, context.Background(), input.Toggle("isSTP", "STP", true, true, nil)),
		renderIn(t, context.Background(), input.ToggleWith(input.ToggleOptions{Name: "isSTP", Label: "STP", Hint: "ayuda", Value: true, Enabled: true})))
}

func TestToggleWith_FlamboyantConAyudaYDeshabilitado(t *testing.T) {
	body := renderIn(t, flamCtx, input.ToggleWith(input.ToggleOptions{Name: "isSTP", Label: "STP", Hint: "no configurado en rekon", Enabled: false}))

	assert.Contains(t, body, "no configurado en rekon")
	assert.Contains(t, body, " disabled")
	assert.Contains(t, body, "cursor-not-allowed")
	assert.Contains(t, body, `value="true"`)
}
