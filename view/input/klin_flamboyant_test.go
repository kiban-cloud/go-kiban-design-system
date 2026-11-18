package input_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func TestRadioCardWith_TonosPlacaYMorado(t *testing.T) {
	placa := renderIn(t, flamCtx, input.RadioCardWith(input.RadioCardOptions{Name: "sic", Value: "BC", Title: "Buró de Crédito", Tone: "plate", Icon: templ.Raw(`<img src="/logo.png">`)}))
	assert.Contains(t, placa, "border-[#E2E5EC] bg-white")
	assert.Contains(t, placa, `<img src="/logo.png">`)

	morado := renderIn(t, flamCtx, input.RadioCardWith(input.RadioCardOptions{Name: "t", Value: "PM", Title: "Persona moral", Tone: "purple"}))
	assert.Contains(t, morado, "bg-[#F3EAFF]")
}

func TestDropzone_Multiple(t *testing.T) {
	assert.Contains(t, renderIn(t, flamCtx, input.Dropzone(input.DropzoneOptions{Name: "files", Multiple: true, Required: true})), " multiple")
	assert.NotContains(t, renderIn(t, flamCtx, input.Dropzone(input.DropzoneOptions{Name: "file"})), " multiple")
}

func pillsOpts() input.PillsOptions {
	return input.PillsOptions{Name: "cut", Label: "Base", Selected: "MONTHLY", Required: true, Options: []input.PillOption{
		{Value: "MONTHLY", Label: "Mensual"},
		{Value: "PARTIAL", Label: "Parcial", Attrs: templ.Attributes{"data-x": "1"}},
	}}
}

func TestPills_Flamboyant(t *testing.T) {
	html := renderIn(t, flamCtx, input.Pills(pillsOpts()))

	assert.Contains(t, html, `role="radiogroup"`)
	assert.Contains(t, html, `aria-label="Base"`)
	assert.Contains(t, html, `value="MONTHLY" checked`)
	assert.NotContains(t, html, `value="PARTIAL" checked`)
	assert.Contains(t, html, `data-x="1"`)
	assert.Contains(t, html, "rounded-full")
	assert.Contains(t, html, " required")
	assert.NotContains(t, html, "border-[#F3A6A6]")

	o := pillsOpts()
	o.Error = true
	assert.Contains(t, renderIn(t, flamCtx, input.Pills(o)), "border-[#F3A6A6]")
}

func TestPills_ClasicoYSinSeleccion(t *testing.T) {
	o := pillsOpts()
	o.Selected = ""
	html := renderIn(t, context.Background(), input.Pills(o))

	assert.NotContains(t, html, "rounded-full")
	assert.NotContains(t, html, "checked")
	assert.Contains(t, html, "Mensual")
}

func TestPills_Disabled(t *testing.T) {
	o := pillsOpts()
	o.Disabled = true
	html := renderIn(t, flamCtx, input.Pills(o))
	assert.Contains(t, html, "disabled")
	assert.Contains(t, html, "cursor-not-allowed")
	assert.Contains(t, renderIn(t, context.Background(), input.Pills(o)), "disabled")
}
