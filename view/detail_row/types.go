// Package detail_row renders read-only label/value pairs as a
// two-column list. Used by detail / "Detalles" tabs across kiban
// htmx surfaces (workfloo histórico, rekon payment detail,
// reportalos invoice detail …) where the same shape — a vertical
// stack of "Label: value" rows — keeps appearing.
//
// Two entry points:
//
//   - Row(label, value) renders one row. Use it when you want to
//     control the wrapping container yourself (e.g. interleaving
//     conditional sections, mixing in flash banners, etc.).
//   - List(items) wraps multiple rows in a styled <dl>. Use it
//     when you have a flat list of static rows.
//
// Both render the same markup; List is just a convenience.
package detail_row

import "github.com/a-h/templ"

// Item is one row in a List call: a label / value pair that the
// view renders as `<dt>label</dt><dd>value</dd>`. Callers that
// already iterate themselves can use Row directly without ever
// touching this type.
type Item struct {
	Label string
	Value string
}

// Group es una columna de Columns: un antetítulo y sus datos.
type Group struct {
	Title  string
	Fields []Field
	// Empty se muestra en lugar de los datos cuando no hay ninguno.
	Empty string
}

// Field es un dato de Columns. Value gana sobre Text cuando el valor lleva
// marcado (una CLABE copiable, una píldora); sin los dos se muestra un guion.
type Field struct {
	Label string
	Text  string
	Value templ.Component
}

// GridItem es un dato de Grid. Mono lo pinta en monoespaciado (CLABE, claves
// de rastreo, ids). Un Value vacío no se muestra.
type GridItem struct {
	Label string
	Value string
	Mono  bool
}

// GridOptions configura GridWith. Sin nada más que Items es Grid.
type GridOptions struct {
	Items []GridItem
	// Cols son las columnas desde sm: 2, o 3 con cualquier otro valor (0
	// incluido), como Grid. En móvil siempre una.
	Cols int
	// FieldLabels pinta la etiqueta como la de un campo de formulario (12 px,
	// gris oscuro, sin mayúsculas) y el valor a 12 px: datos de sólo lectura
	// junto a formularios, como el .fsp-fld del proto. Sin él, la etiqueta
	// chica en mayúsculas grises de Grid.
	FieldLabels bool
	// Empty, si viene, se muestra en lugar de un valor vacío (una raya), así
	// la rejilla no se descuadra. Vacío: los datos sin valor no se muestran.
	Empty string
}
