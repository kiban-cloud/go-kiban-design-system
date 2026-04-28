// Package button holds shared kiban button primitives used by HTMX/templ
// projects (rekon, crm, and future tools).
package button

import "strings"

const defaultType = "button"

// Opts defines the shared options used by all button variants.
type Opts struct {
	// Content
	Label        string
	Icon         string
	IconPosition string // left | right (default: left)

	// HTML base attributes
	Type     string // button | submit | reset
	Disabled bool
	Class    string
	Form     string

	// Accessibility
	AriaLabel string
	Title     string

	// HTMX
	HxPost        string
	HxGet         string
	HxTarget      string
	HxSwap        string
	HxConfirm     string
	HxIndicator   string
	HxDisabledElt string
}

func (o Opts) normalizedType() string {
	if strings.TrimSpace(o.Type) == "" {
		return defaultType
	}
	return o.Type
}

func (o Opts) effectiveAriaLabel() string {
	if strings.TrimSpace(o.AriaLabel) != "" {
		return o.AriaLabel
	}
	return o.Label
}

func (o Opts) effectiveTitle() string {
	if strings.TrimSpace(o.Title) != "" {
		return o.Title
	}
	return o.effectiveAriaLabel()
}

func (o Opts) hasIcon() bool {
	return strings.TrimSpace(o.Icon) != ""
}

func (o Opts) iconAtRight() bool {
	return strings.EqualFold(strings.TrimSpace(o.IconPosition), "right")
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
