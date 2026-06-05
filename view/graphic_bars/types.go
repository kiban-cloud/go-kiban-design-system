// Package graphic_bars renders a small "label + value + horizontal bar"
// list inside a card — used by workfloo's A/B testing dashboard for
// per-workfloo metric breakdowns (label execution counts, time, cost,
// etc.). The intent is a compact comparison view: each row shows a
// label, a formatted total, an optional percent readout, and a colored
// bar whose width tracks the percent.
//
// Variant tints the bars across a row group so the user can pick out
// which kind of metric they're looking at at a glance. Values follow
// the broader kiban variant vocabulary (primary / success / warning /
// error / neutral) with two domain-specific tints — `workfloo` for
// per-workfloo identifier rows and a shared neutral default.
package graphic_bars

// Variant identifies the bar's tint. Use the package-level constants
// rather than raw strings so the renderer's switch stays exhaustive.
const (
	VariantPrimary  = "primary"
	VariantSuccess  = "success"
	VariantWarning  = "warning"
	VariantError    = "error"
	VariantNeutral  = "neutral"
	VariantWorkfloo = "workfloo"
)

// Bar is one row in the chart. Total is the already-formatted string
// the consumer wants displayed in the value column (e.g. "1,234",
// "$12.50", "1h 30m") so the renderer stays formatting-agnostic.
// Percent (0..100) drives the bar width.
type Bar struct {
	Label   string
	Total   string
	Percent float64
}

// Options configures one chart instance.
//
// Title is the card-internal heading (set when the chart is rendered
// standalone inside a card.Section that already shows the section
// title, this can be empty).
//
// Variant picks the bar tint applied to every row in this chart.
//
// HidePercent omits the trailing "%" column for charts whose Total is
// the meaningful value (duration / cost) and percent is just shown via
// bar width.
//
// Bars is the row list. Empty Bars renders an "(sin datos)" placeholder.
type Options struct {
	Title       string
	Variant     string
	HidePercent bool
	Bars        []Bar
}

// barTrackClass returns the structural classes for the bar's track —
// height, rounding, overflow clipping. The actual background color
// comes from `barTrackStyle` per variant so the visual palette is
// explicit (and easy to assert on from a test).
func barTrackClass() string {
	return "h-2 w-full rounded-full overflow-hidden"
}

// barTrackStyle returns the inline style for the bar's track
// (background) — a light tint of the variant's color. Inline rather
// than via Tailwind utility so the exact hex is stable across Tailwind
// purge configurations and pinned in tests.
func barTrackStyle(variant string) string {
	switch variant {
	case VariantPrimary:
		return "background-color: #eff6ff"
	case VariantSuccess:
		return "background-color: #ecfdf5"
	case VariantWarning:
		return "background-color: #fffbeb"
	case VariantError:
		return "background-color: #fef2f2"
	case VariantWorkfloo:
		return "background-color: #eef2ff"
	default:
		return "background-color: #f1f5f9"
	}
}

// barFillStyle returns the inline style for the bar's filled portion:
// a saturated tint of the variant's color paired with the percent
// width. Inline so the exact hex is stable + visible from the rendered
// HTML.
func barFillStyle(variant string, percent float64) string {
	color := "#94a3b8"
	switch variant {
	case VariantPrimary:
		color = "#2563eb"
	case VariantSuccess:
		color = "#10b981"
	case VariantWarning:
		color = "#f59e0b"
	case VariantError:
		color = "#ef4444"
	case VariantWorkfloo:
		color = "#6366f1"
	}
	return "background-color: " + color + "; width: " + formatPercent(clampPercent(percent)) + "%"
}

// barFillClass returns the structural classes for the bar fill.
// Color comes from barFillStyle so a per-variant inline style
// covers both fill color and width in one attribute.
func barFillClass() string {
	return "h-full rounded-full"
}

// clampPercent guards against malformed Percent values so a row can't
// overflow the bar track or invert it. 0..100 in, 0..100 out.
func clampPercent(p float64) float64 {
	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}

// formatPercent renders the percent value as a string with no decimal
// when it's a whole number, or one decimal otherwise.
func formatPercent(p float64) string {
	whole := int(p)
	if float64(whole) == p {
		return itoa(whole)
	}
	tenths := int(p*10) - whole*10
	return itoa(whole) + "." + itoa(tenths)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}
