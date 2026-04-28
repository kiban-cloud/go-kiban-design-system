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
type Options struct {
	Label        string
	Icon         string
	IconPosition string // left | right (default: left)
	// IconComponent optional: any templ component (e.g. icons.Database()). When
	// non-nil it wins over Icon for the active slot (left/right or icon variant).
	IconComponent templ.Component
	Variant       string // primary | secondary | danger | icon

	IsSubmit bool // when true, renders type="submit"
	IsReset  bool // when true, renders type="reset" (wins over IsSubmit)

	Disabled    bool
	ExtraClass  string
	Form        string
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
