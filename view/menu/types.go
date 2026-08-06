// Package menu renders a kebab-style action menu: an icon trigger that
// toggles a dropdown panel of items. The runtime JS (toggle, close,
// outside-click, Escape) lives in view/layout/base.templ so a single set
// of handlers serves any number of menus on a page; this package only
// emits the markup wired by ID.
package menu

import (
	"strings"

	"github.com/a-h/templ"
)

// MenuItem is one row inside a Menu panel.
//
// OnClick is raw inline JS for the item's `onclick`, mirroring
// button.Options.OnClick. The menu always appends a close call after
// the user's expression so consumers don't have to remember to close
// the menu themselves; pass an empty OnClick when the click should
// only close the menu.
type MenuItem struct {
	Label   string
	OnClick string
	// Variant: "" or "default" → kiban-ink text; "danger" → red text.
	// Used for destructive items (delete, revoke, etc.).
	Variant string
	// Icon is an optional left-side icon. Use any templ.Component;
	// matches the IconComponent slot pattern from button.Options.
	Icon templ.Component
	// Attrs is the escape hatch for HTMX, data-*, custom event
	// handlers, etc. Same pattern as button.Options.Attrs.
	Attrs templ.Attributes
	// KeepOpen suppresses the automatic close-on-click. Use it for an
	// item whose feedback lives INSIDE the panel — a "copy" row that
	// flips to a check for a moment, say: with the default close the
	// panel is gone before the user can see it.
	//
	// The item then owns the dismissal. Close it yourself
	// (window.kibanCloseMenu('<id>')) once the feedback has been seen,
	// or the panel stays until the user clicks elsewhere / hits Escape
	// — both still work, they're global handlers.
	KeepOpen bool
}

// Config drives the Menu component. ID is required and must be unique
// per page — it keys the JS toggle/close handlers and the trigger ↔
// panel relationship.
type Config struct {
	ID        string
	AriaLabel string // aria-label on the trigger button (e.g. "Acciones para X")
	Items     []MenuItem
}

// resolvedAttrs builds the per-item <button> attribute set: user Attrs +
// merged onclick that runs the user expression and then closes the menu.
func (it MenuItem) resolvedAttrs(menuID string) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range it.Attrs {
		out[k] = v
	}
	user := strings.TrimSpace(it.OnClick)
	if it.KeepOpen {
		// El item se encarga de cerrar cuando corresponda.
		if user != "" {
			out["onclick"] = user
		}
		return out
	}
	closeCall := "window.kibanCloseMenu('" + menuID + "')"
	if user != "" {
		out["onclick"] = user + ";" + closeCall
	} else {
		out["onclick"] = closeCall
	}
	return out
}

func (it MenuItem) class() string {
	base := "flex w-full items-center gap-2 px-3 py-2 text-sm text-left transition-colors hover:bg-kiban-surface"
	if strings.EqualFold(strings.TrimSpace(it.Variant), "danger") {
		return base + " text-red-600"
	}
	return base + " text-kiban-ink"
}

// triggerAttrs returns the attribute set for the kebab trigger button.
// Bundling these in a single map matches the pattern used by button —
// templ's strict typing for `onclick` makes a spread-friendly map the
// least-friction path.
func triggerAttrs(menuID string) templ.Attributes {
	return templ.Attributes{
		"id":            menuID + "-trigger",
		"aria-haspopup": "true",
		"aria-expanded": "false",
		"aria-controls": menuID + "-panel",
		"onclick":       "window.kibanToggleMenu('" + menuID + "')",
	}
}

func hasItem(items []MenuItem) bool {
	return len(items) > 0
}
