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
//   - A `FooterActions action.Group` — primary + secondaries — rendered as
//     a row at the bottom. `Action` and `Group` live in `view/action/`
//     because they are reused outside the drawer (e.g. by
//     `table.BulkActionBar`); this package re-exposes neither so callers
//     should import `view/action` directly.
//
// Open/close JS lives in base.templ so consumers don't need to ship their
// own toggling logic. Each overlay carries `data-kiban-overlay` so the
// Escape-key listener can find the topmost visible one.
package drawer

import (
	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/action"
)

// SidePanelConfig drives SidePanel. Body is rendered as templ children;
// when FooterActions is non-empty, a footer row is appended below the
// scrollable body with `border-t border-kiban-border`.
type SidePanelConfig struct {
	ID            string
	Title         string
	Size          string
	FooterActions action.Group
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
	FooterActions action.Group
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
	PrimaryAction        action.Action
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

// confirmGroup builds the action.Group for a Confirm dialog: a single
// cancel SecondaryAction (auto-wired to kibanCloseOverlay) and the
// caller-supplied PrimaryAction. Kept as a Go helper instead of
// inlining inside the templ so the address-of on the local PrimaryAction
// copy is unambiguous to the templ-generated code.
func confirmGroup(cfg ConfirmConfig) action.Group {
	primary := cfg.PrimaryAction
	return action.Group{
		SecondaryActions: []action.Action{
			{
				Label:   confirmCancelLabel(cfg.SecondaryActionLabel),
				OnClick: "kibanCloseOverlay('" + cfg.ID + "')",
			},
		},
		PrimaryAction: &primary,
	}
}
