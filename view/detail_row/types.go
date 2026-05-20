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

// Item is one row in a List call: a label / value pair that the
// view renders as `<dt>label</dt><dd>value</dd>`. Callers that
// already iterate themselves can use Row directly without ever
// touching this type.
type Item struct {
	Label string
	Value string
}
