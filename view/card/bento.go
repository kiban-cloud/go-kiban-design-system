package card

import "github.com/a-h/templ"

// BentoOptions es la tarjeta de una pieza en la página de inicio de un módulo
// flamboyant (Bento): qué es, sus dos cifras, la tendencia, lo que hay que
// revisar y cómo entrar.
type BentoOptions struct {
	// Big es la tarjeta grande (la pieza principal): título más grande, las
	// cifras a la izquierda de la tendencia y dos columnas de la cuadrícula
	// desde md.
	Big bool
	// Cols son las columnas que ocupa en la cuadrícula de 6 de BentoGrid en
	// pantalla ancha (xl): 2, 3, 4 o 6. Cero: 2.
	Cols int
	// Color es el de la pieza: la franja de arriba, el ícono, la tendencia,
	// la lista y el botón. Soft es su tinte claro: el fondo del ícono y de
	// la lista («#8000FF» y «#F3E8FF»).
	Color string
	Soft  string
	// Icon va en el cuadro de la izquierda del título (22 px); Text es qué
	// es la pieza, debajo del título.
	Icon  templ.Component
	Title string
	Text  string
	// Pending es la píldora aqua a la derecha del título («3 sin
	// conciliar»). Vacía no se pinta.
	Pending string
	// Kpis son las cifras de la tarjeta (dos), ya formateadas.
	Kpis []BentoKpi
	// Trend es la tendencia (chart.Spark) con su unidad arriba a la
	// izquierda («miles de pesos») y las etiquetas del eje debajo, repartidas
	// de punta a punta (los días, o «hace 7 días» y «hoy»).
	Trend     []float64
	TrendUnit string
	Axis      []string
	// LiveTitle y Live son la lista de lo que hay que revisar, con un punto
	// que late junto al título. Sin renglones no se pinta.
	LiveTitle string
	Live      []BentoRow
	// Open es el botón del pie, en el color de la pieza («Abrir Pagos»), y
	// Links los enlaces que le siguen («Exportar CSV»).
	Open  BentoLink
	Links []BentoLink
}

// BentoKpi es una cifra de Bento: el valor y su etiqueta debajo.
type BentoKpi struct {
	Value string
	Label string
}

// BentoRow es un renglón de la lista de Bento: título, línea debajo, cuándo
// (a la derecha) y adónde lleva. Attrs abre un panel con htmx en vez de
// navegar.
type BentoRow struct {
	Title string
	Sub   string
	When  string
	Href  string
	Attrs templ.Attributes
}

// BentoLink es un enlace del pie de Bento. Attrs para descargar
// («download») o abrir un panel con htmx.
type BentoLink struct {
	Label string
	Href  string
	Attrs templ.Attributes
}

// bentoStyle pasa el color y el tinte a las clases (var(--c), var(--s)).
func bentoStyle(o BentoOptions) templ.SafeCSS {
	return templ.SafeCSS("--c:" + o.Color + ";--s:" + o.Soft)
}

// bentoCols es la clase de columnas en pantalla ancha.
func bentoCols(n int) string {
	switch n {
	case 3:
		return "xl:col-span-3"
	case 4:
		return "xl:col-span-4"
	case 6:
		return "xl:col-span-6"
	}
	return "xl:col-span-2"
}

// bentoSpark es el lienzo de la tendencia: más ancho y alto en la grande.
func bentoSpark(o BentoOptions) (w, h int, class string) {
	if o.Big {
		return 420, 96, "h-[110px]"
	}
	return 220, 56, "h-14"
}

// chevronRight es el «>» de los renglones y del botón (Lucide, trazo 1.5,
// como los íconos del shell).
const chevronRight = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"/></svg>`
