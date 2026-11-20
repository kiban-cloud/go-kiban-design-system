package chart

import "strconv"

// CompareOptions compara dos series por categoría (Compare): para cada una,
// la barra de la primera serie y la de la segunda, una debajo de la otra, y
// la diferencia. Las barras son proporcionales al valor más grande de las
// dos series.
type CompareOptions struct {
	// Series son las dos series: su nombre (leyenda) y su color.
	Series [2]StackPart
	Rows   []CompareRow
	// Delta muestra la diferencia (segunda menos primera) en una píldora
	// verde (sube), roja (baja) o gris (igual).
	Delta bool
}

// CompareRow es una categoría: su nombre y el valor de cada serie (A, la
// primera; B, la segunda), con el texto ya formateado (ADisplay, BDisplay).
type CompareRow struct {
	Label              string
	A, B               float64
	ADisplay, BDisplay string
}

func compareMax(rows []CompareRow) float64 {
	var m float64
	for _, r := range rows {
		m = max(m, r.A, r.B)
	}
	return m
}

// delta es «+66», «−42» o «0».
func (r CompareRow) delta() string {
	d := r.B - r.A
	s := strconv.FormatFloat(d, 'f', -1, 64)
	switch {
	case d > 0:
		return "+" + s
	case d < 0:
		return "−" + s[1:]
	}
	return "0"
}

func (r CompareRow) deltaClass() string {
	switch {
	case r.B > r.A:
		return "bg-[#DDFBF3] text-[#04332A]"
	case r.B < r.A:
		return "bg-[#FDECEC] text-[#B42318]"
	}
	return "bg-[#F1F3F7] text-[#474F66]"
}
