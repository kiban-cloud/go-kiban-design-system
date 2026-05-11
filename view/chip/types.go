// Package chip is the atomic chip primitive — a small pill-shaped tag
// with an optional close button. Use it anywhere a list of removable
// labels makes sense: file uploaders, filter pills, multi-select
// readouts, tag inputs, etc.
//
// What it is NOT:
//   - A status badge (use [badge.Variant] / [badge.Status] for those —
//     a chip is interactive/removable, a badge is a status display).
//   - A button (use [button.Button]).
//   - A standalone widget — chips are always rendered as items inside a
//     parent that owns the list state. The chip itself emits markup
//     only; the consumer wires `RemoveAttrs` to whatever JS / HTMX
//     contract removes the chip.
//
// The chip's remove button (when [Options.Removable] is true) emits a
// `<button type="button" data-chip-remove>` so a parent's delegated
// click listener can find it without per-chip wiring. Consumers can
// still attach extra attributes via [Options.RemoveAttrs] (e.g., a
// `data-remove="<key>"` to identify which chip).
package chip

import "github.com/a-h/templ"

// Options configures one [Chip]. Only [Options.Label] is required.
type Options struct {
	// Label is the visible chip text. Required.
	Label string

	// Subtext is an optional muted secondary line rendered to the right
	// of the label (e.g. file size next to a filename). Empty omits.
	Subtext string

	// Title sets the native `title=""` attribute (browser tooltip).
	// Useful for explaining why a chip is in the `danger` variant
	// ("Archivo inválido: supera 5 MB"), or showing the full label
	// when truncated.
	Title string

	// Variant picks the color scheme:
	//   - "" / "default" — neutral border + surface background
	//   - "danger" — red border / soft red background (invalid items)
	//   - "info" — kiban-primary tint (selected / featured)
	//   - "success" — emerald tint
	//   - "warning" — amber tint
	// Unknown variants fall back to "default" rather than crashing the
	// render.
	Variant string

	// Removable, when true, renders an "×" button on the right edge of
	// the chip. The button has `data-chip-remove` so a parent's
	// delegated click handler can find it. Wire any extra identifying
	// attribute (e.g. `data-remove="<idx>"`) via [Options.RemoveAttrs].
	Removable bool

	// RemoveAttrs are spread onto the remove `<button>` when
	// [Options.Removable] is true. Use this to pass a chip-specific
	// identifier the parent's click handler reads (`data-remove="<id>"`,
	// HTMX `hx-delete`, etc.). Ignored when Removable is false.
	RemoveAttrs templ.Attributes

	// RemoveAriaLabel overrides the remove button's `aria-label`. Empty
	// defaults to "Quitar" — Spanish to match the rest of the kit.
	RemoveAriaLabel string

	// Attrs are spread onto the chip's outer `<span>` element. The
	// escape hatch for any extra attribute the consumer needs (custom
	// `data-*`, `id`, etc.).
	Attrs templ.Attributes
}

// effectiveRemoveAriaLabel returns the aria-label to render on the
// remove button. Defaults to Spanish "Quitar" so chips read naturally
// to a screen reader without requiring the caller to remember.
func (o Options) effectiveRemoveAriaLabel() string {
	if o.RemoveAriaLabel != "" {
		return o.RemoveAriaLabel
	}
	return "Quitar"
}
