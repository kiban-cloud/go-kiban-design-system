package chart

import (
	"strconv"

	"github.com/a-h/templ"
)

// StackPart es un tramo de una barra apilada (Stack) o de un renglón de
// Rank: su valor, su color y cómo se llama en la leyenda. Display es el valor
// ya formateado para la leyenda («1,510»); vacío, la leyenda sólo muestra el
// nombre.
type StackPart struct {
	Label   string
	Value   float64
	Display string
	Color   string
}

// StackOptions es una barra apilada con su leyenda debajo (Stack): cómo se
// reparte un total entre pocas partes (exitosas, en curso, con error).
type StackOptions struct {
	Parts []StackPart
	// HideLegend quita la leyenda.
	HideLegend bool
}

// RankOptions es un ranking de barras horizontales (Rank): un renglón por
// cosa, del más grande al más chico, cada barra proporcional al mayor. Un
// renglón con Parts se parte en tramos (p. ej. por workfloo); sin Parts es
// una barra de Color.
type RankOptions struct {
	Rows []RankRow
	// Color es el de las barras sin tramos.
	Color string
	// Legend son los tramos que se explican debajo (nombre y color).
	Legend []StackPart
}

// RankRow es un renglón de Rank: el nombre, el total (Value para la barra,
// Display ya formateado) y, si viene, su parte del total («34 %»).
type RankRow struct {
	Name    string
	Value   float64
	Display string
	Share   string
	Parts   []StackPart
}

// pctStyle es el ancho de un tramo como porcentaje de total.
func pctStyle(v, total float64) templ.SafeCSS {
	if total <= 0 || v <= 0 {
		return templ.SafeCSS("width:0")
	}
	return templ.SafeCSS("width:" + strconv.FormatFloat(v/total*100, 'f', 2, 64) + "%")
}

// partStyle es un tramo: su ancho dentro de la barra y su color.
func partStyle(p StackPart, total float64) templ.SafeCSS {
	return pctStyle(p.Value, total) + templ.SafeCSS(";background:"+p.Color)
}

// colorStyle es el fondo de un punto de la leyenda o de una barra lisa.
func colorStyle(c string) templ.SafeCSS { return templ.SafeCSS("background:" + c) }

// partsTotal suma los tramos.
func partsTotal(ps []StackPart) float64 {
	var t float64
	for _, p := range ps {
		t += p.Value
	}
	return t
}

// rankMax es el renglón más grande: la barra que llena el ancho.
func rankMax(rows []RankRow) float64 {
	var m float64
	for _, r := range rows {
		if r.Value > m {
			m = r.Value
		}
	}
	return m
}
