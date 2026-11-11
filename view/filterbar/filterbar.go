// Package filterbar es la barra de búsqueda y filtros en línea de kiban cloud
// flamboyant (lista de clientes del proto): un buscador, filtros en píldora
// con su menú y, a la derecha, las acciones de la pantalla (children).
//
// Es un <form> GET: cualquier cambio (una casilla, una opción, lo que se
// escribe en el buscador con una pausa) pide la lista con htmx y reemplaza
// Target. La barra no se vuelve a pintar, así que el texto de cada píldora
// («Tipo: Persona moral», «Tipo · 2») y el enlace de limpiar los actualiza
// base.templ en el navegador ([data-kiban-filter]).
package filterbar

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Config describe la barra. Los textos llegan traducidos por quien la usa.
type Config struct {
	ID        string // id del <form>
	Action    string // URL de la lista (hx-get)
	Target    string // hx-target: el contenedor de la tabla
	Indicator string // hx-indicator

	Search  Search
	Filters []Filter

	// Hidden son campos que viajan con cada búsqueda sin verse (el orden,
	// las filas por página). La página no viaja: un filtro nuevo vuelve a 1.
	Hidden []Field

	// ClearHref es la lista sin búsqueda ni filtros. El enlace sólo se ve
	// cuando hay alguno activo y recarga la página entera, para que la barra
	// también quede limpia.
	ClearHref  string
	ClearLabel string

	// MenuSelectAll y MenuClear son el pie del menú de cada filtro, como en el
	// proto: «Seleccionar todos» (sólo en los de varios valores) y «Limpiar»
	// (se ve con algo marcado). Vacíos, el menú no lleva pie.
	MenuSelectAll string
	MenuClear     string
}

type Search struct {
	Name        string
	Value       string
	Placeholder string
	AriaLabel   string
}

// Filter es una píldora con su menú. Multi son casillas (varios valores con
// el mismo Name); si no, una sola opción (radio): conviene incluir una opción
// de valor vacío («Cualquier fecha») para poder quitarla.
type Filter struct {
	Name    string
	Label   string
	Multi   bool
	Options []Option
	// Search pone un buscador arriba del menú, para filtros con muchas
	// opciones (etiquetas). Filtra en el navegador por el texto de la opción.
	Search MenuSearch
	// Note es una aclaración en gris arriba de las opciones («Sólo aplica a
	// Buró de Crédito»). Empty se ve en lugar de las opciones cuando no hay
	// ninguna (un filtro que depende de otro: «Elige primero un servicio»).
	Note  string
	Empty string
	// Numeric cambia las opciones por un número tecleado con su condición
	// (Requests, Éxito, Latencia). Con Numeric, Options y Multi no aplican.
	Numeric *Numeric
}

// Numeric es un filtro de número: igual a, mayor que, menor que o entre dos
// valores, más atajos («Más de 1,000»). Lo tecleado no pide la lista hasta
// «Aplicar» (o Enter); un atajo aplica al momento. Viajan tres campos
// ocultos (OpName, Name, Name2) y sólo cuando hay valor, así la URL no
// lleva filtros vacíos.
type Numeric struct {
	OpName string // la condición: eq, gt, lt o bt
	Name   string // el valor
	Name2  string // el segundo valor de «entre»

	Op, Value, Value2 string // lo aplicado; Value vacío, sin filtro

	Unit string // «%», «ms» o vacío
	Step string
	Max  string // vacío, sin máximo

	Ops     []NumericOp
	Presets []NumericPreset

	// Textos de la ventana. And separa los dos valores de «entre» («y»);
	// ErrMax lleva {max}.
	Placeholder, From, To, And string
	PresetsLabel, Apply        string
	ErrNumber, ErrRange        string
	ErrMax                     string
}

// NumericOp es una condición del filtro numérico: Value (eq, gt, lt, bt),
// Label en el botón («Mayor que») y Symbol en la píldora («>», «entre»).
type NumericOp struct {
	Value, Label, Symbol string
}

// NumericPreset es un atajo que aplica una condición de un clic.
type NumericPreset struct {
	Op, Value, Label string
}

// MenuSearch es el buscador del menú de un filtro. Sin Placeholder no hay.
// Empty es lo que se ve cuando nada coincide; {q} se reemplaza por lo escrito.
type MenuSearch struct {
	Placeholder string
	Empty       string
}

type Option struct {
	Value    string
	Label    string
	Selected bool
	// Count es cuántos registros quedarían con esta opción (aplicando los
	// demás filtros). HasCount lo muestra, para que un 0 real se vea.
	Count    int
	HasCount bool
	// Icon va antes del texto de la opción: el logo de un SIC, el punto de
	// color de un estado. Nil, nada.
	Icon templ.Component
	// Hint va en gris junto al texto: de qué servicio es una funcionalidad
	// cuando el menú mezcla varios.
	Hint string
}

type Field struct {
	Name  string
	Value string
}

// countID es el id del conteo de una opción; Counts lo usa para
// reemplazarlo oob sin volver a pintar la barra. Lleva el valor en hex y no
// la posición: si la lista de opciones cambió desde que se pintó la barra
// (una etiqueta nueva), un conteo nunca cae en otra opción, y el id queda
// válido aunque el valor traiga espacios o acentos.
func countID(formID, filter, value string) string {
	return formID + "-n-" + filter + "-" + hex.EncodeToString([]byte(value))
}

// selected son las opciones marcadas con valor (la opción vacía de un radio
// no cuenta como filtro).
func (f Filter) selected() []Option {
	var out []Option
	for _, o := range f.Options {
		if o.Selected && o.Value != "" {
			out = append(out, o)
		}
	}
	return out
}

// on dice si el filtro está aplicado.
func (f Filter) on() bool {
	if f.Numeric != nil {
		return f.Numeric.Value != ""
	}
	return len(f.selected()) > 0
}

func (c Config) active() bool {
	if c.Search.Value != "" {
		return true
	}
	for _, f := range c.Filters {
		if f.on() {
			return true
		}
	}
	return false
}

// op es la condición marcada en la ventana: la aplicada o, sin filtro, la
// primera que no sea «igual a» (mayor que, como el proto).
func (n Numeric) op() string {
	if n.Value != "" && n.Op != "" {
		return n.Op
	}
	for _, o := range n.Ops {
		if o.Value != "eq" {
			return o.Value
		}
	}
	if len(n.Ops) > 0 {
		return n.Ops[0].Value
	}
	return ""
}

func (f Filter) inputType() string {
	if f.Multi {
		return "checkbox"
	}
	return "radio"
}

// placeholder1 es el del primer valor: «Desde» con «entre», si no el normal.
func (n Numeric) placeholder1() string {
	if n.op() == "bt" && n.From != "" {
		return n.From
	}
	return n.Placeholder
}

func boolAttr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// text es lo que dice la píldora con el filtro aplicado: «> 1,000»,
// «= 100%», «entre 1,000 y 2,000 ms». base.templ lo arma igual en el
// navegador (kibanNumText).
func (n Numeric) text() string {
	sym := n.Op
	for _, o := range n.Ops {
		if o.Value == n.Op {
			sym = o.Symbol
		}
	}
	unit := n.Unit
	if unit != "" && unit != "%" {
		unit = " " + unit
	}
	if n.Op == "bt" {
		a, b := n.Value, n.Value2
		if fa, fb := parseNum(a), parseNum(b); fb < fa {
			a, b = b, a
		}
		return sym + " " + groupNum(a) + " " + n.And + " " + groupNum(b) + unit
	}
	return sym + " " + groupNum(n.Value) + unit
}

func parseNum(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// groupNum pone separador de miles («1,000»; «99.5» queda igual), como
// toLocaleString("es-MX") en el navegador.
func groupNum(s string) string {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	neg := f < 0
	if neg {
		f = -f
	}
	str := strconv.FormatFloat(f, 'f', -1, 64)
	if i := strings.IndexByte(str, '.'); i >= 0 {
		dec := str[i+1:]
		if len(dec) > 2 {
			str = strconv.FormatFloat(f, 'f', 2, 64)
			str = strings.TrimRight(strings.TrimRight(str, "0"), ".")
		}
	}
	intPart, frac := str, ""
	if i := strings.IndexByte(str, '.'); i >= 0 {
		intPart, frac = str[:i], str[i:]
	}
	var b strings.Builder
	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	out := b.String() + frac
	if neg {
		out = "-" + out
	}
	return out
}
