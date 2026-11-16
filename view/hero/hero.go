// Package hero es el encabezado grande de la página de inicio de un módulo en
// kiban cloud flamboyant: el bloque azul con degradado que saluda, dice qué
// hay pendiente y presenta el ciclo del módulo como pasos que llevan a cada
// pantalla. Es el mismo lenguaje en Base, rekon y Monitoreo (kiban-proto).
package hero

import (
	"strconv"

	"github.com/a-h/templ"
)

// CycleOptions es el hero de un módulo (Cycle).
type CycleOptions struct {
	// Space es el nombre del espacio en el chip de arriba y SpaceTile su
	// placa de 20 px a la izquierda. Sin Space no hay chip.
	Space     string
	SpaceTile templ.Component
	// Title es el saludo grande («Buenas tardes») y Text la línea debajo.
	Title string
	Text  string
	// PendingCount > 0 pinta a la derecha el botón blanco con campana:
	// «<b>3</b> por revisar», que lleva a PendingHref.
	PendingCount int
	PendingLabel string
	PendingHref  string
	// PendingAttrs, si viene, convierte el botón en un <button> con esos
	// atributos en lugar de un enlace a PendingHref: con
	// layout.InboxToggleAttrs abre el panel «Pendientes» del shell, como en
	// la maqueta.
	PendingAttrs templ.Attributes
	// StepsLabel es el aria-label de la lista de pasos («Ciclo del dinero»).
	StepsLabel string
	// Steps son los pasos del ciclo, de izquierda a derecha, unidos por una
	// línea. Pensado para 3 a 6; en móvil se desplazan de lado.
	Steps []CycleStep
	// Animated pone tres destellos aqua que recorren la línea del ciclo,
	// como el dinero que pasa de un paso al siguiente (rekon). Apagado, la
	// línea queda quieta (Base). Con prefers-reduced-motion no se mueven
	// nunca. Los estilos viven en layout/base.templ (.kiban-cy-flow).
	Animated bool
}

// CycleStep es un paso del ciclo: el ícono en su cuadro blanco (con un
// contador si Badge > 0), el verbo en versalitas («Cobra»), el nombre de la
// pantalla («Órdenes de pago») y la cifra del día debajo. Lleva a Href.
type CycleStep struct {
	Icon  templ.Component
	Verb  string
	Name  string
	Sub   string
	Badge int
	Href  string
}

// cycleVars son las medidas que dependen de cuántos pasos hay: en móvil cada
// paso mide 150 px (más 4 de separación) y la línea va del centro del
// primero al del último; desde 900 px los pasos se reparten el ancho y la
// línea deja medio paso a cada lado.
func cycleVars(n int) templ.SafeCSS {
	if n < 1 {
		n = 1
	}
	ns := strconv.Itoa(n)
	inset := strconv.FormatFloat(50/float64(n), 'f', -1, 64) + "%"
	return templ.SafeCSS("--cy-cols:repeat(" + ns + ",150px);--cy-cols-lg:repeat(" + ns + ",minmax(120px,1fr));--cy-line:" + strconv.Itoa((n-1)*154) + "px;--cy-inset:" + inset)
}

// heroStyle es el fondo: el degradado azul de Base, con el brillo aqua abajo
// y, arriba a la derecha, las tres formas del isotipo de kiban en degradado
// aqua → azul → morado (en las maquetas era una imagen; aquí, tres span con
// esquinas redondeadas y degradados).
const heroStyle = templ.SafeCSS("background:radial-gradient(120% 140% at 0% 0%,#2F6BFF 0%,#0047FF 38%,#0A2A8F 100%)")

const glowStyle = templ.SafeCSS("background:radial-gradient(closest-side,rgba(0,221,178,.35),transparent 70%)")

// Las tres formas de la malla, en un cuadro de 300 px: arriba y abajo a la
// izquierda, dos cuartos con la esquina de adentro muy redonda; a la
// derecha, medio óvalo.
const (
	meshTopStyle    = templ.SafeCSS("left:0;top:0;width:61%;height:48%;border-bottom-right-radius:77% 100%;background:linear-gradient(45deg,#14CDB4,#3A8BF2)")
	meshBottomStyle = templ.SafeCSS("left:0;bottom:0;width:59%;height:48%;border-top-right-radius:77% 100%;background:linear-gradient(135deg,#14CDB4 0%,#5B8DEF 55%,#8B5CF6 100%)")
	meshRightStyle  = templ.SafeCSS("right:0;top:0;width:49%;height:100%;border-radius:100% 0 0 100% / 50% 0 0 50%;background:linear-gradient(200deg,#3A8BF2 0%,#7C6CF6 50%,#9B55FF 100%)")
)

// stepStyle escalona la entrada de cada paso (60 ms por posición).
func stepStyle(i int) templ.SafeCSS {
	return templ.SafeCSS("--i:" + strconv.Itoa(i))
}
