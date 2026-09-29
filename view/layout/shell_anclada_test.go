package layout_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// renderLayoutAnclada renderiza la shell completa para inspeccionar el CSS que
// emite Base, que es donde vive el anclaje.
func renderLayoutAnclada(t *testing.T) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, layout.Layout(layout.Config{
		Title:       "t",
		ProjectName: "p",
	}).Render(context.Background(), &buf))
	return buf.String()
}

// El grupo rail + sub-nav se anclaba mal de dos formas distintas y las dos son
// silenciosas: sin align-self el flex lo estira a toda el área scrolleable y el
// sticky no tiene recorrido, así que no pega nada; sin la altura explícita no
// hay dónde recortar y cada rail crece sin límite. Se ve igual de bien en una
// pantalla alta y falla sólo en las bajas, que es donde nadie prueba.
func TestLayout_GrupoDeNavegacionSeAnclaBajoElTopbar(t *testing.T) {
	html := renderLayoutAnclada(t)

	i := strings.Index(html, "@media (min-width: 768px)")
	require.NotEqual(t, -1, i, "falta la media query de escritorio del anclaje")
	bloque := html[i:min(i+700, len(html))]

	assert.Contains(t, bloque, "position: sticky")
	assert.Contains(t, bloque, "top: 3.5rem", "debe quedar justo debajo del topbar h-14")
	assert.Contains(t, bloque, "align-self: flex-start", "sin esto el sticky no tiene recorrido")
	assert.Contains(t, bloque, "overflow-y: auto", "cada rail scrollea por su cuenta")
}

// El topbar y los rails van acoplados: el grupo se pega a top:3.5rem, así que un
// topbar que se va con el scroll deja una franja vacía de su alto encima.
func TestTopbar_SeAnclaAlTopeParaNoDejarFranja(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, layout.Topbar(layout.Config{ProjectName: "p"}).Render(context.Background(), &buf))
	// Sobre el atributo class y no sobre "sticky top-0" a secas: templ emite los
	// comentarios HTML al output, y el comentario que explica este anclaje contiene
	// esa misma cadena — asertarla sola daba un test que jamás se ponía rojo.
	assert.Contains(t, buf.String(), "shrink-0 sticky top-0 z-30")
}

// El drawer móvil pone position:fixed sobre el MISMO selector. Si el anclaje de
// escritorio no estuviera en su propia media query competirían por orden de
// cascada y el drawer podría dejar de salir.
func TestLayout_ElAnclajeNoInvadeElDrawerMovil(t *testing.T) {
	html := renderLayoutAnclada(t)

	movil := strings.Index(html, "@media (max-width: 767px)")
	require.NotEqual(t, -1, movil)
	assert.Contains(t, html[movil:min(movil+900, len(html))], "position: fixed",
		"el drawer móvil debe seguir siendo fixed")

	escritorio := strings.Index(html, "@media (min-width: 768px)")
	assert.Less(t, escritorio, movil, "el bloque de escritorio va antes y separado del móvil")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
