package chart

import (
	"math"
	"strconv"
	"strings"
)

// SparkOptions es una tendencia chica (Spark): la serie, su color y el
// lienzo en el que se traza.
type SparkOptions struct {
	// Values es la serie de izquierda a derecha (p. ej. los últimos 7 días);
	// el último valor es hoy y lleva un punto. Vacía no se dibuja.
	Values []float64
	// Color es el de la línea, el área tenue y el punto («#8000FF»).
	Color string
	// Width y Height son el lienzo del trazo (viewBox). El SVG se estira al
	// ancho de su contenedor; Class le da la altura («h-14»).
	Width, Height int
	Class         string
}

// sparkGeom es la tendencia ya trazada en un lienzo de w×h: la línea, el área
// debajo y el último punto (hoy), que lleva un círculo.
type sparkGeom struct {
	ViewBox string
	Line    string
	Area    string
	LastX   string
	LastY   string
}

// sparkPoints traza una serie de izquierda a derecha, con el mínimo abajo y
// el máximo arriba y 4 px de margen a los lados. Una serie plana queda abajo;
// una vacía (o un lienzo demasiado chico) no se dibuja.
func sparkPoints(t []float64, w, h int) sparkGeom {
	if len(t) == 0 || w <= 8 || h <= 16 {
		return sparkGeom{}
	}
	if len(t) == 1 {
		t = []float64{t[0], t[0]}
	}
	mx, mn := math.Inf(-1), math.Inf(1)
	for _, v := range t {
		mx, mn = math.Max(mx, v), math.Min(mn, v)
	}
	r := mx - mn
	if r == 0 {
		r = 1
	}
	itoa := strconv.Itoa
	var line strings.Builder
	var lx, ly int
	for i, v := range t {
		lx = int(math.Round(float64(i)/float64(len(t)-1)*float64(w-8))) + 4
		ly = int(math.Round(float64(h-6) - (v-mn)/r*float64(h-16)))
		if i == 0 {
			line.WriteString("M")
		} else {
			line.WriteString(" L")
		}
		line.WriteString(itoa(lx) + " " + itoa(ly))
	}
	l := line.String()
	return sparkGeom{
		ViewBox: "0 0 " + itoa(w) + " " + itoa(h),
		Line:    l,
		Area:    l + " L" + itoa(lx) + " " + itoa(h) + " L4 " + itoa(h) + " Z",
		LastX:   itoa(lx),
		LastY:   itoa(ly),
	}
}
