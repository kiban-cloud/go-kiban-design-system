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
//     "xl" → max-w-xl, "2xl" → max-w-2xl, "3xl" → max-w-3xl,
//     "4xl" → max-w-4xl, "5xl" → max-w-5xl, "6xl" → max-w-6xl,
//     "7xl" → max-w-7xl, "8xl" → max-w-[min(88rem,92vw)] (wide but
//     still modal — keeps a backdrop strip), "full" → max-w-none
//     (full-viewport). Empty string falls back to "md".
//
//     4xl–7xl are intended for multi-column editor experiences
//     (e.g. the RULESET wizard's 4-column layout). Use sparingly
//     — short-form modals look best at lg or below.
//
//   - A `FooterActions button.Group` — primary + secondaries — rendered as
//     a row at the bottom. `Group` lives in `view/button/` and is reused
//     outside the drawer (e.g. by `table.BulkActionBar`); this package
//     re-exposes nothing, so callers should import `view/button` directly.
//
// Open/close JS lives in base.templ so consumers don't need to ship their
// own toggling logic. Each overlay carries `data-kiban-overlay` so the
// Escape-key listener can find the topmost visible one.
package drawer

import (
	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/button"
)

// SidePanelConfig drives SidePanel. Body is rendered as templ children;
// when FooterActions is non-empty, a footer row is appended below the
// scrollable body with `border-t border-kiban-border`.
type SidePanelConfig struct {
	ID            string
	Title         string
	Size          string
	FooterActions button.Group
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
	FooterActions button.Group
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
	PrimaryAction        button.Options
	SecondaryActionLabel string
}

// sizeClass maps the public size keyword to its Tailwind max-width class.
// Falls back to max-w-md (the most common case) on empty / unknown.
//
// The 2xl/3xl sizes were added for grid-style modals (e.g. workfloo's
// templates picker) whose content doesn't fit at xl without forcing
// each card to a too-narrow column. Use sparingly — most modals are
// "short form" and look best at lg or below.
func sizeClass(size string) string {
	switch size {
	case "sm":
		return "max-w-sm"
	case "lg":
		return "max-w-lg"
	case "xl":
		return "max-w-xl"
	case "2xl":
		return "max-w-2xl"
	case "3xl":
		return "max-w-3xl"
	case "4xl":
		return "max-w-4xl"
	case "5xl":
		return "max-w-5xl"
	case "6xl":
		return "max-w-6xl"
	case "7xl":
		return "max-w-7xl"
	case "8xl":
		// One step past 7xl for wide two-pane editors (e.g. workfloo's
		// FORM node preview + field picker) that need more room than
		// 7xl but should still read as a modal — the `92vw` cap always
		// leaves a strip of backdrop on the left instead of spanning
		// the whole viewport like `full`. Rendered via Tailwind's
		// arbitrary-value + min() (supported by the Play CDN runtime).
		return "max-w-[min(88rem,92vw)]"
	case "full":
		// Full-viewport modal — drops the max-width cap so the
		// dialog spans the available space (minus the modal's
		// own `mx-4` horizontal gutter). Used by long-form
		// multi-column editors (workfloo RULESET, future
		// wide-canvas dialogs) where 7xl still leaves the
		// right-side column too narrow on big monitors.
		return "max-w-none"
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

// confirmCancelLabel returns the text for the cancel button on a Confirm
// dialog, defaulting to "Cancelar" when the caller leaves it blank.
func confirmCancelLabel(label string) string {
	if label == "" {
		return "Cancelar"
	}
	return label
}

// confirmGroup builds the button.Group for a Confirm dialog: a single
// cancel SecondaryAction (auto-wired to kibanCloseOverlay) and the
// caller-supplied PrimaryAction. Kept as a Go helper instead of
// inlining inside the templ so the address-of on the local PrimaryAction
// copy is unambiguous to the templ-generated code.
func confirmGroup(cfg ConfirmConfig) button.Group {
	primary := cfg.PrimaryAction
	return button.Group{
		SecondaryActions: []button.Options{
			{
				Label:   confirmCancelLabel(cfg.SecondaryActionLabel),
				OnClick: "kibanCloseOverlay('" + cfg.ID + "')",
			},
		},
		PrimaryAction: &primary,
	}
}
