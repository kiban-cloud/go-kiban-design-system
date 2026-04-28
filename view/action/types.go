// Package action holds the shared "button or link" primitive used by every
// kiban component that surfaces caller-driven controls in a footer / bar /
// row layout. Today that's `view/drawer` (SidePanel + Modal + Confirm
// footer rows) and `view/table` (BulkActionBar). Future components
// (action sheets, inline action menus, etc.) reuse the same Action +
// Group structs so projects don't deal with three slightly-different
// "action button" shapes.
//
// Why its own package: drawer originally owned `Action` + `FooterActions`,
// but as soon as a second consumer (BulkActionBar) showed up it became
// clear the types weren't drawer-specific. Extracting them up here:
//
//   - keeps drawer / table / future components from depending on each
//     other for type definitions,
//   - lets the type names be neutral ("Group" instead of the
//     drawer-flavoured "FooterActions"),
//   - leaves room for a future `view/button/` package to share the same
//     variant vocabulary.
//
// Variant strings used across the design system:
//
//	"primary"    bg kiban-primary, text white            (typical default for the lead action)
//	"secondary"  outline border-kiban-border             (typical default for cancel-like actions)
//	"danger"     bg red-600, text white                  (destructive — "Eliminar", "Cancelar orden", …)
//
// Empty string falls through to the per-slot default supplied by the
// rendering helper (e.g. "primary" for `Group.PrimaryAction`,
// "secondary" for entries in `Group.SecondaryActions`).
package action

import "github.com/a-h/templ"

// Action is one button or link rendered in a row of controls. The render
// rules:
//
//   - When `Href != ""`, renders as `<a href=Href>`. Otherwise renders as
//     a `<button>` with `Type` defaulting to "button" so the element
//     never accidentally inherits an enclosing form's submit behaviour.
//   - `Variant` ("primary" | "secondary" | "danger") picks the colour
//     scheme; empty string defers to the per-slot default chosen by the
//     rendering helper.
//   - `Type="submit"` + `Form="external-form-id"` is the canonical
//     filter-drawer pattern: the button submits a form rendered
//     elsewhere on the page (e.g. inside the drawer body) instead of
//     wrapping the action row in a form.
//   - `OnClick` is raw inline JS — combine with HTMX `Attrs` for the
//     "submit then close the overlay" pattern by setting
//     `OnClick: "kibanCloseOverlay('id')"`.
//   - `Attrs` is the HTMX escape hatch (hx-post / hx-target / hx-swap /
//     hx-confirm / hx-indicator …); the design system stays
//     HTMX-agnostic and just spreads the attributes onto the element.
//   - `Disabled` emits the HTML `disabled` attribute and applies the
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

// Group is a small, ordered collection of actions sharing a row of
// controls — typically a footer, a bulk-action bar, or an inline action
// menu. `PrimaryAction` is rendered last (rightmost) and defaults to the
// "primary" variant. `SecondaryActions` render to its left in the order
// given and default to the "secondary" variant. Renderers should treat a
// fully-empty Group (no PrimaryAction, no SecondaryActions) as "render
// nothing" so consumers can pass an empty Group when a layout has no
// actions in a particular state.
type Group struct {
	PrimaryAction    *Action
	SecondaryActions []Action
}

// IsEmpty reports whether the Group has nothing to render. Renderers use
// this to skip the surrounding chrome (e.g. drawer's `border-t pt-4`
// divider) when there are no actions to show.
func (g Group) IsEmpty() bool {
	return g.PrimaryAction == nil && len(g.SecondaryActions) == 0
}

// actionClass returns the Tailwind class string for an Action's button or
// anchor element. `defaultVariant` is the per-slot fallback applied when
// the action's own Variant field is empty — typically "primary" for a
// PrimaryAction slot, "secondary" for entries in SecondaryActions.
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

// buttonType maps a possibly-empty Action.Type to the actual HTML
// attribute value. Empty defaults to "button" so a button never
// accidentally inherits an enclosing form's submit behaviour just
// because it sits inside one.
func buttonType(t string) string {
	if t == "" {
		return "button"
	}
	return t
}

// actionAttrs merges an Action's user-supplied HTMX `Attrs` with its
// `OnClick` and `Form` fields into a single attribute set ready for the
// `{ … }` spread in the templ. Bundling them via templ.Attributes lets
// us emit `onclick` as a plain string attribute — the alternative
// (`onclick={…}`) would force every consumer through `templ.JSFuncCall`
// or similar, since templ's strict typing for script attributes
// requires `templ.ComponentScript` values.
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
