// Package flash holds the kiban "banner" primitive — the full-width
// status messages rendered above forms / above list pages after a
// mutation, an authorization failure, or any other actionable
// notification the user needs to see before continuing.
//
// Banners are server-rendered. Client-side they support two ways to
// disappear: a click on the trailing × (all variants), and an
// auto-dismiss timer (success / info only). Warning and error stay
// up until the user acknowledges them. The handlers live in
// `view/layout/base.templ`; consumers don't need any wiring.
//
// Variants follow the same vocabulary used in card and badge:
//
//	"success"   emerald — completed actions, positive confirmations
//	"error"     red     — operations that failed, hard blocks
//	"warning"   amber   — soft warnings, attention-required states
//	"info"      blue / kiban-primary — neutral informative messages
//
// Empty variant falls back to "info".
//
// Two API shapes:
//
//   - Banner(variant, msg): generic, takes a variant string. Use this
//     when the variant is computed at runtime.
//   - Success(msg) / Error(msg) / Warning(msg) / Info(msg): typed
//     wrappers, the more idiomatic call when the variant is known at
//     compile time. They delegate to Banner.
package flash

import (
	"strconv"

	"github.com/a-h/templ"
)

// Auto-dismiss timings per variant (in milliseconds). Success and
// info are positive / neutral notices the user shouldn't have to
// chase off the page — short enough that they don't linger, long
// enough to read a sentence or two. Warning and error require a
// deliberate dismiss (return 0 → no timer attribute rendered).
//
// Tuning rationale: 5s is the toast-library standard for short
// success copy; info gets a slightly longer 6s because its text
// tends to be more substantive ("Recordá que…", "El reporte está
// generándose…"). If a specific call site needs different timing,
// a future iteration can accept a per-banner override — today the
// matrix is small enough that hard-coding by variant is the
// simpler contract.
const (
	autoDismissSuccessMs = 5000
	autoDismissInfoMs    = 6000
)

// autoDismissMs returns the auto-dismiss delay in milliseconds for
// a given variant, or 0 when the variant should stay until manually
// dismissed.
func autoDismissMs(variant string) int {
	switch variant {
	case "success":
		return autoDismissSuccessMs
	case "info", "":
		return autoDismissInfoMs
	}
	return 0
}

// autoDismissMsAttr formats the auto-dismiss delay for emission as
// an HTML attribute value. Always paired with `autoDismissMs > 0`
// at the call site, so we don't need a special-case for the
// zero-timer variants.
func autoDismissMsAttr(variant string) string {
	return strconv.Itoa(autoDismissMs(variant))
}

// Options is the full contract for a banner. Banner and the typed wrappers
// build one of these; call Alert directly when you need an action inside the
// box or a banner that can't be closed.
type Options struct {
	Variant string
	Message string

	// Dismissible turns on the × and the auto-dismiss timer. It defaults to
	// false because the banners that need this struct are the ones explaining
	// why something can't happen: closing one hides the only explanation on
	// the page, and nothing brings it back until the next render.
	Dismissible bool

	// ActionLabel + ActionHref render a link inside the banner. With Attrs but
	// no Href it renders a <button> instead, so an HTMX action doesn't need a
	// link that goes nowhere.
	ActionLabel string
	ActionHref  string
	ActionAttrs templ.Attributes
}

func (o Options) hasAction() bool {
	return o.ActionLabel != "" && (o.ActionHref != "" || len(o.ActionAttrs) > 0)
}

// autoDismissMs follows Dismissible: a banner the user can't close by hand
// shouldn't vanish on its own either.
func (o Options) autoDismissMs() int {
	if !o.Dismissible {
		return 0
	}
	return autoDismissMs(o.Variant)
}

func (o Options) autoDismissMsAttr() string {
	return strconv.Itoa(o.autoDismissMs())
}

// NonEmptyAttrs reports whether attrs should be spread onto the element.
func NonEmptyAttrs(attrs templ.Attributes) bool {
	return len(attrs) > 0
}
