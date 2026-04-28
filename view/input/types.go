// Package input holds the shared form-control primitives every kiban htmx
// project uses: text, select, textarea, checkbox, toggle, phone (intl-tel-input
// wrapper), hidden, and date. Each component renders a label + control + error
// (or hint) block following the kiban design tokens (kiban-border /
// kiban-primary / red-400 / red-500), so error state is consistent across
// projects.
//
// Conventions:
//   - Components emit `id="f-<name>"` paired with `<label for="f-<name>">` so
//     clicking the label focuses the control.
//   - When `errMsg != ""` the border swaps to red and the message renders
//     under the control. When `errMsg == ""` and `hint != ""` the hint renders
//     instead. They never both appear.
//   - Required is a visual marker only (a red asterisk next to the label).
//     HTML5 `required` is intentionally NOT emitted: validation lives
//     server-side and the label asterisk is the user-facing hint.
package input

// SelectOption is one row of a `<select>`. Iterated by Select to render
// `<option value=Value selected=...>Label</option>`. Value is what the form
// submits; Label is what the user sees.
type SelectOption struct {
	Value string
	Label string
}
