// Package drawer holds the kiban "overlay" primitives: SidePanel (right-edge
// slide-in for filters and side detail), Modal (centered confirmation /
// short-form dialog), and Confirm (Modal preset wrapping a destructive or
// high-stakes confirmation prompt).
//
// All three components share:
//
//   - An `id` parameter used by the global JS helpers `kibanOpenOverlay(id)`
//     and `kibanCloseOverlay(id)` (defined in view/layout/base.templ) to
//     toggle the `hidden` class on the outer container. The caller wires
//     their own trigger with `onclick="kibanOpenOverlay('some-id')"`; the
//     component handles close via the close button + backdrop click + the
//     global Escape-key listener (closes only the topmost visible overlay).
//
//   - A `Size` parameter mapped to a Tailwind `max-w-*` class:
//     "sm" → max-w-sm, "" or "md" → max-w-md, "lg" → max-w-lg,
//     "xl" → max-w-xl. Empty string falls back to "md".
//
//   - A FooterActions block (SidePanel and Modal) with a single optional
//     PrimaryAction on the right and zero-or-more SecondaryActions to its
//     left. Footer renders only when at least one slot is non-empty.
//
// Open/close JS lives in base.templ so consumers don't need to ship their
// own toggling logic. Each overlay carries `data-kiban-overlay` so the
// Escape-key listener can find the topmost visible one.
package drawer

import "github.com/a-h/templ"

// Action is one button or link in the footer of a drawer / modal / confirm.
//
// Render rules:
//   - When Href != "", renders as `<a href=Href>`. Otherwise renders as
//     `<button>` with Type defaulting to "button".
//   - Variant ("primary" | "secondary" | "danger") picks the colour
//     scheme; the empty string falls through to a slot-specific default
//     (PrimaryAction → "primary", SecondaryActions → "secondary").
//   - Type ("button" | "submit") only applies to button mode. Use
//     "submit" + Form="external-form-id" to submit a form that lives
//     elsewhere on the page (the typical filter-drawer pattern).
//   - OnClick is raw inline JS; combine with HTMX Attrs as needed
//     (e.g. `OnClick="kibanCloseOverlay('id')"` to dismiss the overlay
//     after submit).
//   - Attrs is the HTMX escape hatch — spread directly onto the element.
//   - Disabled emits the HTML `disabled` attribute and applies the
//     muted-style classes via the variant's `disabled:` modifiers.
type Action struct {
	Label    string
	Variant  string
	Href     string
	Type     string
	Form     string
	OnClick  string
	Attrs    templ.Attributes
	Disabled bool
}

// FooterActions configures the footer row of a SidePanel or Modal.
// PrimaryAction is rendered on the right (rightmost when there are
// secondaries); SecondaryActions render to its left in the order given.
// Footer renders only when at least one slot is set.
type FooterActions struct {
	PrimaryAction    *Action
	SecondaryActions []Action
}

// hasFooter reports whether FooterActions has anything to render.
func (f FooterActions) hasFooter() bool {
	return f.PrimaryAction != nil || len(f.SecondaryActions) > 0
}

// SidePanelConfig drives SidePanel. Body is rendered as templ children;
// when FooterActions is non-empty, a footer row is appended below the
// scrollable body with `border-t border-kiban-border`.
type SidePanelConfig struct {
	ID            string
	Title         string
	Size          string
	FooterActions FooterActions
}

// ModalConfig drives Modal. Body is rendered as templ children. Icon is
// optional and renders to the left of the title in the header — caller
// supplies the fully-styled icon block (typical pattern: a rounded
// w-8/h-8 wrapper with `bg-kiban-primary-soft text-kiban-primary` plus
// the icon SVG inside). Pass nil when no icon.
//
// Pattern for HTMX-form-bound modals (e.g. "Confirm and submit"): wrap
// the entire `@drawer.Modal(...)` call in a `<form hx-post=…>` so the
// PrimaryAction's `Type:"submit"` button submits naturally without
// needing an external `Form` reference. See WithdrawModal in
// rekon-backend for an example.
type ModalConfig struct {
	ID            string
	Title         string
	Size          string
	Icon          templ.Component
	FooterActions FooterActions
}

// ConfirmConfig drives Confirm — a fixed-shape Modal preset for "are you
// sure?" prompts that need richer markup than HTMX's native `hx-confirm`.
// PrimaryAction is the destructive / proceed button; SecondaryActionLabel
// is the cancel button (always closes the overlay, no HTMX wiring).
//
// For destructive confirms, set `PrimaryAction.Variant = "danger"`.
//
// SecondaryActionLabel defaults to "Cancelar" when empty. Size defaults
// to "sm" when empty.
type ConfirmConfig struct {
	ID                   string
	Title                string
	Message              string
	Size                 string
	PrimaryAction        Action
	SecondaryActionLabel string
}

// sizeClass maps the public size keyword to its Tailwind max-width class.
// Falls back to max-w-md (the most common case) on empty / unknown.
func sizeClass(size string) string {
	switch size {
	case "sm":
		return "max-w-sm"
	case "lg":
		return "max-w-lg"
	case "xl":
		return "max-w-xl"
	default:
		return "max-w-md"
	}
}

// confirmSizeClass mirrors sizeClass but defaults to max-w-sm — confirmation
// dialogs are tighter than general modals.
func confirmSizeClass(size string) string {
	if size == "" {
		return "max-w-sm"
	}
	return sizeClass(size)
}

// actionClass returns the Tailwind class string for an Action button/link.
// `defaultVariant` is the per-slot fallback ("primary" for PrimaryAction,
// "secondary" for SecondaryActions, etc.); applied only when the action's
// Variant field is empty.
func actionClass(variant, defaultVariant string) string {
	v := variant
	if v == "" {
		v = defaultVariant
	}
	const base = "inline-flex items-center justify-center rounded-md px-4 py-2 text-sm font-medium"
	switch v {
	case "primary":
		return base + " bg-kiban-primary text-white hover:opacity-90 disabled:opacity-50 disabled:cursor-wait"
	case "danger":
		return base + " bg-red-600 text-white hover:opacity-90 disabled:opacity-50 disabled:cursor-wait"
	default: // "secondary" + unknown values
		return base + " bg-white border border-kiban-border text-kiban-ink hover:border-kiban-ink3 disabled:opacity-50"
	}
}

// buttonType maps a possibly-empty Action.Type to the actual HTML attribute
// value: empty string defaults to "button" so a button never accidentally
// inherits the surrounding form's submit behaviour.
func buttonType(t string) string {
	if t == "" {
		return "button"
	}
	return t
}

// overlayCloseAttrs builds the spread-able attribute set for the chrome's
// close button: `type=button`, the `kibanCloseOverlay(id)` inline JS, and
// the muted hover styling that matches every existing kiban close button.
//
// We bundle these via templ.Attributes (instead of `onclick={…}`
// interpolation) because templ's strict `onclick` typing requires a
// `templ.ComponentScript`, which would force every consumer through the
// `templ.JSFuncCall` wrapper. The Attributes spread bypasses that check
// — safe here because `id` is a developer-supplied CSS identifier, not
// user input.
func overlayCloseAttrs(id string) templ.Attributes {
	return templ.Attributes{
		"type":       "button",
		"onclick":    "kibanCloseOverlay('" + id + "')",
		"class":      "text-kiban-ink3 hover:text-kiban-ink shrink-0",
		"aria-label": "Cerrar",
	}
}

// overlayBackdropAttrs builds the spread-able attribute set for the
// click-to-dismiss backdrop. `bgClass` lets the caller pick the opacity
// — SidePanel uses `bg-black/30` (lighter, since the panel covers only
// part of the screen), Modal/Confirm use `bg-black/40` (darker, focusing
// attention on the centered card).
func overlayBackdropAttrs(id, bgClass string) templ.Attributes {
	return templ.Attributes{
		"class":   "absolute inset-0 " + bgClass,
		"onclick": "kibanCloseOverlay('" + id + "')",
	}
}

// actionAttrs merges an Action's user-supplied HTMX `Attrs` with its
// `OnClick` and `Form` fields into a single attribute set ready for
// `{ … }` spread. Same rationale as overlayCloseAttrs: avoids templ's
// onclick-type check by routing through Attributes instead of the
// `onclick={…}` syntax.
func actionAttrs(a Action) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range a.Attrs {
		out[k] = v
	}
	if a.OnClick != "" {
		out["onclick"] = a.OnClick
	}
	if a.Form != "" {
		out["form"] = a.Form
	}
	return out
}

// confirmCancelLabel returns the text for the cancel button on a Confirm
// dialog, defaulting to "Cancelar" when the caller leaves it blank.
func confirmCancelLabel(label string) string {
	if label == "" {
		return "Cancelar"
	}
	return label
}
