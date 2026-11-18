package input_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

// Con la bandera flamboyant, los campos siguen los .fsp-fld del proto
// (etiqueta gris de 18 px, 5 px de separación, relleno de 10 px); sin ella,
// el campo de siempre.
func TestCampos_FlamboyantCambiaEtiquetaContenedorYCaja(t *testing.T) {
	clasico := renderIn(t, context.Background(), input.Text("rfc", "RFC", "", "", "", false))
	assert.Contains(t, clasico, `<div class="space-y-1">`)
	assert.Contains(t, clasico, `class="text-xs font-medium text-kiban-ink2"`)
	assert.Contains(t, clasico, "w-full border rounded-md px-3 py-2 text-sm focus:outline-none")

	flam := renderIn(t, flamCtx, input.Text("rfc", "RFC", "", "", "", false))
	assert.Contains(t, flam, `<div class="flex flex-col gap-[5px]">`)
	assert.Contains(t, flam, `class="text-xs font-medium leading-[18px] text-[#474F66]"`)
	assert.Contains(t, flam, "w-full border rounded-md px-2.5 py-2 text-sm leading-[18px] text-kiban-ink focus:outline-none")
	assert.NotContains(t, flam, "space-y-1")
}

// El mismo trío de clases en los demás campos que lo adoptaron.
func TestCampos_FlamboyantEnSelectYTextarea(t *testing.T) {
	sel := renderIn(t, flamCtx, input.Select("regimen", "Régimen", "", "", "", []input.SelectOption{{Value: "601", Label: "General"}}, false))
	assert.Contains(t, sel, "flex flex-col gap-[5px]")
	assert.Contains(t, sel, "text-[#474F66]")
	assert.Contains(t, sel, "px-2.5 py-2 text-sm leading-[18px]")
	assert.Contains(t, sel, "bg-white")

	ta := renderIn(t, flamCtx, input.Textarea("notas", "Notas", "", "", "", false, 3, "", false))
	assert.Contains(t, ta, "flex flex-col gap-[5px]")
	assert.Contains(t, ta, "text-[#474F66]")
}

// La etiqueta de Money y de Chips también sigue la bandera.
func TestCampos_FlamboyantEnMoneyYChips(t *testing.T) {
	assert.Contains(t, renderIn(t, flamCtx, input.Money(input.MoneyOptions{Name: "a", Label: "Monto"})), "text-[#474F66]")
	assert.Contains(t, renderIn(t, context.Background(), input.Money(input.MoneyOptions{Name: "a", Label: "Monto"})), "text-kiban-ink2")

	assert.Contains(t, renderIn(t, flamCtx, input.Chips(input.ChipsOptions{Name: "e", Label: "Correos"})), "flex flex-col gap-[5px]")
	assert.Contains(t, renderIn(t, context.Background(), input.Chips(input.ChipsOptions{Name: "e", Label: "Correos"})), "space-y-1")
}
