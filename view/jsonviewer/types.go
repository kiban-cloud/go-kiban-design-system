// Package jsonviewer renders an arbitrary JSON-shaped Go value as a
// nested, expandable accordion tree. Mirrors the React `JsonViewer`
// from KDS — keys with object/array values become collapsible
// accordions; primitive values render as `key: value` rows.
//
// The visible state of each accordion is stored in the native
// `<details>` element, so individual node toggles need zero JS. The
// only inline script is the "Expandir todo / Ocultar todo" master
// toggle, scoped per-viewer via Options.ID so multiple viewers on
// the same page don't clobber each other.
//
// Acceptable input shapes:
//
//   - map[string]any / map[string]interface{} (the typical
//     encoding/json decode output for objects).
//   - []any / []interface{} for arrays.
//   - Primitives: string, bool, float64 (json.Number), int*, nil.
//   - Anything else falls through to fmt.Sprint(v) — degrades
//     gracefully instead of panicking on exotic types.
package jsonviewer

// Options drives the View templ. `Data` is the only required field;
// `ID` is needed when multiple viewers coexist on the same page so
// their "Expandir todo" buttons stay scoped to their own subtree.
type Options struct {
	// Data is the value to render. Typically the result of
	// `json.Unmarshal(...)` into `interface{}` — pass it through
	// as-is and the viewer will introspect runtime types.
	Data any

	// ID is the DOM id of the viewer wrapper. Must be unique per
	// page when several viewers are rendered side-by-side. Defaults
	// to "kiban-jsonviewer" when empty, which is fine for the
	// single-viewer common case.
	ID string

	// EmptyMessage is the Spanish placeholder shown when `Data` is
	// nil or an empty container. Defaults to "Sin información." so
	// callers don't have to repeat the same fallback copy on every
	// callsite. Set explicitly to swap wording.
	EmptyMessage string
}

// resolveID picks the wrapper id, falling back to a stable default
// when the caller didn't supply one. Kept as a tiny helper so the
// templ can stay declarative.
func (o Options) resolveID() string {
	if o.ID == "" {
		return "kiban-jsonviewer"
	}
	return o.ID
}

// resolveEmpty picks the empty-state message, falling back to a
// neutral default.
func (o Options) resolveEmpty() string {
	if o.EmptyMessage == "" {
		return "Sin información."
	}
	return o.EmptyMessage
}
