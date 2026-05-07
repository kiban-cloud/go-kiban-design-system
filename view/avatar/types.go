// Package avatar renders a circular profile picture with an initials
// fallback for when the image URL is empty, broken, or slow. Used by the
// admin layout's user menu and any other "user identity" surface.
package avatar

import "strings"

type Options struct {
	Src     string // image URL; empty falls back to initials
	Name    string // used for initials and default aria-label
	Size    string // "sm" (24px) | "md" (32px, default) | "lg" (40px)
	AltText string // optional override of Name for aria-label
}

func (o Options) sizeClass() string {
	switch strings.ToLower(strings.TrimSpace(o.Size)) {
	case "sm":
		return "h-6 w-6 text-xs"
	case "lg":
		return "h-10 w-10 text-base"
	default:
		return "h-8 w-8 text-sm"
	}
}

func (o Options) altText() string {
	if strings.TrimSpace(o.AltText) != "" {
		return o.AltText
	}
	return o.Name
}

func (o Options) hasSrc() bool {
	return strings.TrimSpace(o.Src) != ""
}

// Initials returns up to 2 uppercase letters derived from Name.
// "Antonio Blancas" -> "AB", "Antonio" -> "A", "" -> "?".
// Public so consumers (e.g. server-rendered notification lists) can
// reuse the same derivation outside the templ.
func Initials(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "?"
	}
	firstRunes := []rune(parts[0])
	first := strings.ToUpper(string(firstRunes[0]))
	if len(parts) == 1 {
		return first
	}
	lastRunes := []rune(parts[len(parts)-1])
	last := strings.ToUpper(string(lastRunes[0]))
	return first + last
}
