// Package graphic_bars renders a titled card of horizontal labelled
// bars where each bar's label sits *inside* the coloured fill and the
// total (plus an optional percentage) lives in a right-hand column.
//
// It replaces a charting library for the "share of executions"
// visuals (A/B Testing statistics, dashboard breakdowns…): the data is
// a handful of categories with a percent each, so a CSS bar reads
// better than a full chart canvas and keeps the label legible inside
// the bar instead of crowding an external axis. Mirrors the look of
// kiban's React `AlphaGraphicCard`.
package graphic_bars

import "strconv"

// formatPercent renders a percent without trailing zeros ("60", not
// "60.00"; "12.5" stays "12.5") so the width string and the visible
// "%" label stay tidy.
func formatPercent(percent float64) string {
	return strconv.FormatFloat(percent, 'f', -1, 64)
}

// percentLabel is the right-column "X%" text.
func percentLabel(percent float64) string {
	return formatPercent(percent) + "%"
}

// Bar is one horizontal bar. Total is a caller-formatted display
// string (a count "5", a money "$12.50", a duration "1h 30m"); the
// component never formats it. Percent (0..100) drives the fill width
// and — unless Options.HidePercent — the right-hand "%" label.
type Bar struct {
	Label   string
	Total   string
	Percent float64
}

// Options configures one card. Variant selects the 6-colour palette
// (fill / border / label colour); HidePercent drops the trailing "%"
// for metrics where a percentage is meaningless (durations, costs).
type Options struct {
	Title       string
	Variant     string
	HidePercent bool
	Bars        []Bar
}

// Variant constants — exported so callers avoid stringly-typed values.
const (
	VariantError    = "error"
	VariantSuccess  = "success"
	VariantWarning  = "warning"
	VariantWorkfloo = "workfloo"
	VariantPrimary  = "primary"
	VariantNeutral  = "neutral"
)

// palette is the resolved colour triple for a variant.
type palette struct {
	fill   string
	border string
	label  string
}

// paletteFor maps a variant to its colour triple. Unknown variants
// collapse to neutral so an unexpected value never renders an
// invisible (uncoloured) bar.
func paletteFor(variant string) palette {
	switch variant {
	case VariantError:
		return palette{fill: "#ffe3e3", border: "#ffc5c5", label: "#b70000"}
	case VariantSuccess:
		return palette{fill: "#ecfdf5", border: "#aefcdd", label: "#047857"}
	case VariantWarning:
		return palette{fill: "#fff4d9", border: "#ffce5f", label: "#a76100"}
	case VariantWorkfloo:
		return palette{fill: "#f6efff", border: "#d0b4ff", label: "#2d02bb"}
	case VariantPrimary:
		return palette{fill: "#f1f7ff", border: "#c7dfff", label: "#0000cc"}
	default:
		return palette{fill: "#f5f6f7", border: "#dee1e5", label: "#2f3946"}
	}
}

// barFillStyle is the inline style map for the coloured fill div: the
// dynamic width + the variant's fill/border colours. A map keeps each
// property individually CSS-sanitized by templ.
func barFillStyle(variant string, percent float64) map[string]string {
	p := paletteFor(variant)
	return map[string]string{
		"width":         clampPercent(percent),
		"background":    p.fill,
		"border":        "1px solid " + p.border,
		"border-radius": "0.25rem",
	}
}

// labelStyle colours the in-bar label per variant.
func labelStyle(variant string) map[string]string {
	return map[string]string{"color": paletteFor(variant).label}
}

// clampPercent renders the width string, pinning to [0,100] so a bad
// upstream value can't overflow the row.
func clampPercent(percent float64) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	return formatPercent(percent) + "%"
}
