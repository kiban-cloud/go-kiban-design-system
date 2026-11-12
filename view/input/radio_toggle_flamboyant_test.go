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

func TestRadioCardWith_FlamboyantEtiquetasTonoGrisYGrande(t *testing.T) {
	body := renderIn(t, flamCtx, input.RadioCardWith(input.RadioCardOptions{
		Name: "template", Value: "pf", Title: "Onboarding PF", Subtitle: "Incluye:", Tags: []string{"Identidad", "Firma"}, Big: true,
	}))
	assert.Contains(t, body, ">Identidad</span>")
	assert.Contains(t, body, ">Firma</span>")
	assert.Contains(t, body, "grid-cols-[40px_minmax(0,1fr)] items-start p-3")
	assert.Contains(t, body, "bg-[#EAF0FF]", "sin Tone, azul")

	gris := renderIn(t, flamCtx, input.RadioCardWith(input.RadioCardOptions{Name: "template", Value: "empty", Title: "En blanco", Tone: "neutral"}))
	assert.Contains(t, gris, "bg-[#F2F2F2] text-[#474747]")
	assert.NotContains(t, gris, "bg-[#EAF0FF]")
	assert.Contains(t, gris, "grid-cols-[36px_minmax(0,1fr)] items-center", "sin Big, la de siempre")
	assert.NotContains(t, gris, "rounded-full bg-[#F2F4F8]", "sin etiquetas")
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

func TestCheckboxCardWith(t *testing.T) {
	o := input.CheckboxCardOptions{Name: "groups", Value: "g1", Title: "Analistas", Help: "workfloo 3", Checked: true}
	flam := renderIn(t, flamCtx, input.CheckboxCardWith(o))
	assert.Contains(t, flam, `name="groups"`)
	assert.Contains(t, flam, `value="g1"`)
	assert.Contains(t, flam, `checked name="groups"`)
	assert.Contains(t, flam, "has-[:checked]:border-kiban-primary")
	assert.Contains(t, flam, ">Analistas</b>")
	assert.Contains(t, flam, "workfloo 3")

	dis := renderIn(t, flamCtx, input.CheckboxCardWith(input.CheckboxCardOptions{Name: "rights", Value: "X", Title: "Ver", Disabled: true}))
	assert.Contains(t, dis, `disabled name="rights"`)
	assert.NotContains(t, dis, "checked disabled")

	clasico := renderIn(t, context.Background(), input.CheckboxCardWith(o))
	assert.Contains(t, clasico, `value="g1"`, "el diseño de siempre también manda su valor")
	assert.Contains(t, clasico, "border-kiban-border rounded-md")
}
