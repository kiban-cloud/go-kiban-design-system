package hero

import "github.com/a-h/templ"

// DetailOptions es el encabezado de la ficha de una cosa (Detail): la liga
// de regreso, el antetítulo, el título (que se puede renombrar en línea),
// una línea de datos y las acciones a la derecha. Lo que va debajo (un
// reparto, unas cifras) llega como children. Lo estrenaron el A/B test y el
// backtesting de workfloo flamboyant (kiban-proto).
type DetailOptions struct {
	BackHref, BackLabel string
	Eyebrow             string
	Title               string
	// Rename, si viene, hace del título un campo: al cambiarlo se manda
	// solo (htmx, change) a Rename.Action con el valor en Rename.Name.
	Rename *DetailRename
	// Meta son los datos debajo del título (el ID copiable, chips, la
	// fecha), separados por un punto medio.
	Meta []templ.Component
	// Actions van a la derecha del título (botones).
	Actions templ.Component
}

// DetailRename es el renombre en línea del título de Detail.
type DetailRename struct {
	Action string
	Name   string
	Label  string // para lectores de pantalla («Nombre del A/B test»)
	// Attrs van en el <form> (hx-target, hx-swap…); por defecto el
	// resultado no se pinta (hx-swap none).
	Attrs templ.Attributes
}

func (r DetailRename) formAttrs() templ.Attributes {
	a := templ.Attributes{"hx-post": r.Action, "hx-trigger": "change", "hx-swap": "none"}
	for k, v := range r.Attrs {
		a[k] = v
	}
	return a
}

// Chip es un dato con ícono para Meta de Detail (el flujo, el árbol).
type Chip struct {
	Icon  templ.Component
	Label string
	Href  string
}
