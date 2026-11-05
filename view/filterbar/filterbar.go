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
}

type Option struct {
	Value    string
	Label    string
	Selected bool
}

type Field struct {
	Name  string
	Value string
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

func (c Config) active() bool {
	if c.Search.Value != "" {
		return true
	}
	for _, f := range c.Filters {
		if len(f.selected()) > 0 {
			return true
		}
	}
	return false
}

func (f Filter) inputType() string {
	if f.Multi {
		return "checkbox"
	}
	return "radio"
}
