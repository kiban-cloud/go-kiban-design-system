package input

import (
	"context"

	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// Los campos de formulario con flamboyant prendido siguen los .fsp-fld de las
// maquetas de kiban-proto: etiqueta de 12 px con interlineado de 18 px en
// gris #474F66, 5 px entre etiqueta y campo, y el campo de 36 px de alto con
// 10 px de relleno a los lados. Apagada la bandera, el campo es el de siempre.

// fieldWrapClass es el contenedor de etiqueta, campo y ayuda.
func fieldWrapClass(ctx context.Context) string {
	if flamboyant.On(ctx) {
		return "flex flex-col gap-[5px]"
	}
	return "space-y-1"
}

// fieldLabelClass es la etiqueta del campo.
func fieldLabelClass(ctx context.Context) string {
	if flamboyant.On(ctx) {
		return "text-xs font-medium leading-[18px] text-[#474F66]"
	}
	return "text-xs font-medium text-kiban-ink2"
}

// fieldControlClass es la base del <input>, <select> o <textarea>; el color
// del borde según el error lo pone cada campo.
func fieldControlClass(ctx context.Context) string {
	if flamboyant.On(ctx) {
		return "w-full border rounded-md px-2.5 py-2 text-sm leading-[18px] text-kiban-ink focus:outline-none"
	}
	return "w-full border rounded-md px-3 py-2 text-sm focus:outline-none"
}
