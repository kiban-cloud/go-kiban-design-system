package card_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

// Cifras del resultado de una carga, como el proto: número, etiqueta e ícono
// en su cuadro con el tono del resultado.
func TestStats_CifrasConTono(t *testing.T) {
	body := render(t, context.Background(), card.Stats([]card.StatOptions{
		{Value: "12", Label: "clientes importados", Icon: templ.Raw(`<i>ok</i>`), Tone: "success"},
		{Value: "3", Label: "filas duplicadas", Icon: templ.Raw(`<i>dup</i>`), Tone: "warning"},
		{Value: "1", Label: "filas con error", Icon: templ.Raw(`<i>err</i>`), Tone: "danger"},
	}))

	assert.Equal(t, 3, strings.Count(body, "data-kiban-stat"))
	assert.Contains(t, body, ">12</div>")
	assert.Contains(t, body, "clientes importados")
	assert.Contains(t, body, "bg-[#E3FBF5] text-[#007A63]", "importados en verde")
	assert.Contains(t, body, "bg-[#FFF4E5] text-[#9A4F00]", "duplicados en ámbar")
	assert.Contains(t, body, "bg-[#FDECEC] text-[#B42318]", "errores en rojo")
	assert.Contains(t, body, "<i>err</i>")
}

func TestStats_Vacio(t *testing.T) {
	assert.Empty(t, strings.TrimSpace(render(t, context.Background(), card.Stats(nil))))
}
