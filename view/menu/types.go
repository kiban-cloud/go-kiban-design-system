// Package menu holds the kiban kebab/dropdown action-menu primitive — a
// trigger button that, when clicked, reveals a small floating list of
// per-row or per-element actions (rename, delete, copy value, etc.).
//
// Visual contract: a 3-dot ("more") trigger that drops a right-aligned
// popover of [MenuItem] entries. Each entry is a `<button>` whose raw
// inline JS comes from [MenuItem.OnClick] (same pattern as
// [button.Options.OnClick]). After the user's OnClick runs, the
// popover auto-closes.
//
// Implemented on top of the native `<details>` / `<summary>` pair so
// the open/close toggle works without bespoke JS state. A single
// document-level click listener (registered idempotently the first
// time any [Menu] renders) handles the "click outside to close"
// behavior shared across instances.
//
// Typical use is one [Menu] per row in an action column:
//
//	@menu.Menu(menu.Config{
//	    ID:        "apikey-menu-" + row.ID,
//	    AriaLabel: "Acciones para " + row.Name,
//	    Items:     rowMenuItems(row, canEdit, canDelete),
//	})
package menu

import (
	"strings"

	"github.com/a-h/templ"
)

// MenuItem is one entry in the dropdown. `Label` is the visible text;
// `OnClick` is raw inline JS run when the user picks the entry (same
// pattern as button.Options.OnClick — typically a one-liner like
// `kibanOpenApiKeyEdit('123', "Foo")`). `Variant` swaps the colour
// scheme: empty / unknown values render as the neutral kiban-ink tone;
// `"danger"` switches to red for destructive actions.
type MenuItem struct {
	Label   string
	OnClick string
	Variant string
}

// Config configures one [Menu] invocation. `ID` becomes the trigger
// `<summary id="…">` so callers can reference the trigger from JS or
// tests when needed; pass a unique value per row. `AriaLabel`
// describes the trigger to screen readers ("Acciones para <row name>"
// is the conventional kiban wording). `Items` is rendered top-to-
// bottom in the popover; an empty slice renders the trigger but no
// popover content (use [Items] gating in the caller to skip the menu
// entirely when there are no actions).
type Config struct {
	ID        string
	AriaLabel string
	Items     []MenuItem
}

// menuItemClasses returns the Tailwind class string for a menu item
// row. Kept tiny so the templ stays readable.
func menuItemClasses(variant string) string {
	base := "block w-full text-left px-3 py-2 text-sm transition-colors"
	switch variant {
	case "danger":
		return base + " text-red-600 hover:bg-red-50"
	default:
		return base + " text-kiban-ink hover:bg-kiban-surface"
	}
}

// itemAttrs builds the `onclick` attribute for one [MenuItem]: the
// user's raw OnClick (when set) followed by a small auto-close snippet
// that pops the surrounding `<details>` shut so the menu hides as soon
// as the user picks something.
func itemAttrs(it MenuItem) templ.Attributes {
	close := "var d=this.closest('details'); if(d){d.removeAttribute('open');}"
	user := strings.TrimSpace(it.OnClick)
	if user == "" {
		return templ.Attributes{"onclick": close}
	}
	return templ.Attributes{"onclick": user + "; " + close}
}
