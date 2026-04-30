// Package button holds shared kiban button primitives used by HTMX/templ
// projects (rekon, crm, and future tools).
package button

import (
	"strings"

	"github.com/a-h/templ"
)

const (
	defaultVariant = "primary"
)

// Options configures a single [Button]. HTMX and other ad-hoc attributes go
// through Attrs (same pattern as view/input) so the component stays
// transport-agnostic.
//
// Href makes [Button] render as an `<a>` instead of a `<button>`. This is the
// "button or link" pattern used by drawer footers, table BulkActionBar, etc.
// — a control slot that can either submit/click locally or navigate to a
// route. When Href is set, IsSubmit/IsReset/Form/Disabled are silently
// ignored because they don't apply to anchors; consumers that need a
// disabled-link visual should drive it via ExtraClass + aria-disabled in
// Attrs.
type Options struct {
	Label        string
	Icon         string
	IconPosition string // left | right (default: left)
	// IconComponent optional: any templ component (e.g. icons.Database()). When
	// non-nil it wins over Icon for the active slot (left/right or icon variant).
	IconComponent templ.Component
	Variant       string // primary | secondary | danger | icon

	// Href, when non-empty, makes Button render as an `<a>` and skips the
	// button-only attributes (type/disabled/form). Empty -> renders <button>.
	Href string

	IsSubmit bool // when true, renders type="submit" (button only)
	IsReset  bool // when true, renders type="reset" (wins over IsSubmit; button only)

	// OnClick is raw inline JS for the `onclick` attribute. First-class field
	// rather than a key in Attrs because the "submit then close overlay"
	// pattern (`OnClick: "kibanCloseOverlay('id')"`) is common enough across
	// drawer / modal callsites that the ergonomic shortcut earns its keep.
	OnClick string

	Disabled    bool // ignored when Href is set (anchors have no disabled attr)
	ExtraClass  string
	Form        string // ignored when Href is set
	AriaLabel   string
	Title       string
	Attrs       templ.Attributes
}

func (o Options) normalizedVariant() string {
	v := strings.ToLower(strings.TrimSpace(o.Variant))
	switch v {
	case "destructive":
		return "danger"
	case "primary", "secondary", "danger", "icon":
		return v
	default:
		return defaultVariant
	}
}

// EffectiveVariant is the resolved variant string after defaults.
func (o Options) EffectiveVariant() string {
	return o.normalizedVariant()
}

func (o Options) effectiveAriaLabel() string {
	if strings.TrimSpace(o.AriaLabel) != "" {
		return o.AriaLabel
	}
	return o.Label
}

func (o Options) effectiveTitle() string {
	if strings.TrimSpace(o.Title) != "" {
		return o.Title
	}
	return o.effectiveAriaLabel()
}

func (o Options) hasIcon() bool {
	return strings.TrimSpace(o.Icon) != ""
}

func hasIconComponent(p Options) bool {
	return p.IconComponent != nil
}

func (o Options) iconAtRight() bool {
	return strings.EqualFold(strings.TrimSpace(o.IconPosition), "right")
}

// htmlButtonType returns the HTML type attribute.
func htmlButtonType(p Options) string {
	if p.IsReset {
		return "reset"
	}
	if p.IsSubmit {
		return "submit"
	}
	return "button"
}

// BuildClass merges the variant base classes with ExtraClass.
func BuildClass(variant, extraClass string) string {
	return joinClasses(baseClasses(variant), extraClass)
}

func baseClasses(variant string) string {
	switch strings.ToLower(strings.TrimSpace(variant)) {
	case "secondary":
		return "inline-flex items-center justify-center gap-2 rounded-md border border-kiban-border bg-white px-4 py-2 text-sm font-medium text-kiban-ink transition-colors hover:border-kiban-ink3 hover:bg-kiban-surface disabled:cursor-not-allowed disabled:opacity-50"
	case "danger":
		return "inline-flex items-center justify-center gap-2 rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
	case "icon":
		return "inline-flex h-9 w-9 items-center justify-center rounded-md border border-kiban-border bg-white text-kiban-ink transition-colors hover:border-kiban-ink3 hover:bg-kiban-surface disabled:cursor-not-allowed disabled:opacity-50"
	default:
		return "inline-flex items-center justify-center gap-2 rounded-md bg-kiban-primary px-4 py-2 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
	}
}

func joinClasses(classes ...string) string {
	out := make([]string, 0, len(classes))
	for _, className := range classes {
		trimmed := strings.TrimSpace(className)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, " ")
}

// NonEmptyAttrs reports whether attrs should be spread onto the element.
func NonEmptyAttrs(attrs templ.Attributes) bool {
	return len(attrs) > 0
}

// resolvedAttrs merges Options.Attrs with the dedicated OnClick field into a
// single attribute set ready to be spread on the element. Bundling onclick
// here lets the templ emit it as a plain string attribute — the alternative
// (`onclick={…}`) would force every consumer through templ.JSFuncCall since
// templ's strict typing for script attributes requires templ.ComponentScript
// values.
func resolvedAttrs(p Options) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range p.Attrs {
		out[k] = v
	}
	if strings.TrimSpace(p.OnClick) != "" {
		out["onclick"] = p.OnClick
	}
	return out
}

func hasHref(p Options) bool {
	return strings.TrimSpace(p.Href) != ""
}

// Group is a small ordered collection of buttons sharing a row of controls
// — typically a drawer/modal footer, a bulk-action bar, or an inline action
// menu. PrimaryAction renders rightmost (default variant "primary");
// SecondaryActions render to its left in order (default variant
// "secondary"). Renderers should treat an empty Group as "render nothing"
// so consumers can pass an empty Group when a layout has no actions in a
// particular state.
//
// The Group lives in this package (rather than its own) because it's a thin
// composition over Button — same vocabulary, same variants, same Attrs.
// Drawer footers and table BulkActionBar consume Group directly.
type Group struct {
	PrimaryAction    *Options
	SecondaryActions []Options
}

// IsEmpty reports whether the Group has nothing to render. Renderers use
// this to skip the surrounding chrome (e.g. drawer's `border-t pt-4`
// divider) when there are no actions to show.
func (g Group) IsEmpty() bool {
	return g.PrimaryAction == nil && len(g.SecondaryActions) == 0
}

// WithDefaultVariant returns a copy of o with Variant filled in when the
// caller left it empty. Used by RenderGroup so PrimaryAction defaults to
// "primary" and entries in SecondaryActions default to "secondary" without
// requiring callers to repeat the variant on every action.
func WithDefaultVariant(o Options, defaultVariant string) Options {
	if strings.TrimSpace(o.Variant) == "" {
		o.Variant = defaultVariant
	}
	return o
}

func shouldRenderLeftIcon(p Options) bool {
	if p.normalizedVariant() == "icon" {
		return true
	}
	return (p.hasIcon() || hasIconComponent(p)) && !p.iconAtRight()
}

func shouldRenderRightIcon(p Options) bool {
	if p.normalizedVariant() == "icon" {
		return false
	}
	return (p.hasIcon() || hasIconComponent(p)) && p.iconAtRight()
}

// renderCustomIconLeft is true when the left (or sole icon-variant) slot should
// render IconComponent instead of the Icon string registry.
func renderCustomIconLeft(p Options) bool {
	if !hasIconComponent(p) {
		return false
	}
	if p.normalizedVariant() == "icon" {
		return true
	}
	return !p.iconAtRight()
}

func renderCustomIconRight(p Options) bool {
	return hasIconComponent(p) && p.normalizedVariant() != "icon" && p.iconAtRight()
}

func effectiveAriaLabelForVariant(p Options) string {
	if p.normalizedVariant() == "icon" {
		return p.effectiveAriaLabel()
	}
	return ""
}

func effectiveTitleForVariant(p Options) string {
	if p.normalizedVariant() == "icon" {
		return p.effectiveTitle()
	}
	return ""
}
