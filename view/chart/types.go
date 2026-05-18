// Package chart wraps Chart.js (loaded via Config.LoadChartJS in base.templ)
// in templ-renderable components. The server emits a `<canvas>` whose
// `data-chart` attribute carries the Chart.js config as JSON; the small init
// script in this package's `chart.templ` reads it and instantiates the chart
// on DOMContentLoaded + htmx:afterSwap.
//
// Why this shape: the templ side stays HTMX-agnostic (no chart code in
// handlers) and HTMX swaps that re-emit a chart canvas Just Work — the
// scanner runs on every swap. A caller wanting to update a chart in place
// returns a fresh `Line(...)` fragment from their handler with the same
// outer `id`, and HTMX swaps it; the JS detects the new `data-chart-init=""`
// flag missing and re-runs.
package chart

// Options is the high-level config a caller hands to the chart templs.
// It mirrors the subset of Chart.js v4 options every kiban chart uses;
// projects with niche needs can compose extra config via `Extra` (a raw
// JSON blob merged into the Chart.js config at init time).
type Options struct {
	// ID for the canvas element. Must be unique on the page so HTMX
	// swaps and Chart.js's own internal handle don't collide.
	ID string

	// Labels for the x-axis (line/bar) or each slice (donut/pie).
	Labels []string

	// Datasets is one or more named series. For donut/pie use a single
	// dataset; for line/bar multiple datasets stack as separate series.
	Datasets []Dataset

	// HeightPx fixes the canvas height. Default 240. Chart.js sizes the
	// width to the parent; height needs an explicit number to avoid the
	// canvas growing forever inside a flex container.
	HeightPx int

	// HideLegend suppresses the Chart.js legend block. Default (false) =
	// the legend renders below the chart. Pair with HideYAxis when the
	// dataset is self-evident (e.g. a single series whose bars are
	// labeled by date on the x-axis) and the chrome would just be noise.
	HideLegend bool

	// HideYAxis suppresses the y-axis entirely (ticks, labels, line, AND
	// the horizontal gridlines that extend from it). Useful for compact
	// "sparkline-style" bars on dashboard tiles. Stacked + HideYAxis is
	// allowed but unusual — without the y scale, the absolute magnitude
	// of stacks is hard to read.
	HideYAxis bool

	// HideXAxis is the horizontal mirror of HideYAxis: collapses the
	// x-axis (ticks, labels, line, gridlines). On horizontal bar charts
	// this is the *value* axis (0/0.5/1/…) — paired with HideGridLines
	// it produces a chart that's just category labels + bars.
	HideXAxis bool

	// HideGridLines suppresses the lines that run across the plot area
	// on both axes — both the grid lines themselves and the axis border
	// lines they extend from. The axis tick labels remain visible
	// unless HideXAxis / HideYAxis is also set. Use this for clean,
	// dashboard-tile-style charts where the chrome would compete with
	// the data.
	HideGridLines bool

	// Stacked sets `stacked: true` on both axes for bar charts. Ignored
	// by other chart types.
	Stacked bool

	// ExtraClass is applied to the chart's outer wrapper div. Use this
	// for layout integration — most commonly `"flex-1 min-h-0"` when the
	// chart sits inside a flex column and should grow to fill the
	// remaining vertical space (so the canvas height tracks the card
	// height rather than the fixed HeightPx). HeightPx still applies as
	// the initial pre-init height; Chart.js's `responsive: true` resizes
	// the canvas to its container after init.
	ExtraClass string
}

// Dataset is one series. Color is a CSS color string ("#0047FF",
// "rgba(0,71,255,0.4)"); callers pick the palette themselves so charts
// can match the rest of the page.
type Dataset struct {
	Label string
	Data  []float64
	// Color is used for both line stroke and bar fill. For donut/pie
	// charts pass a slice of colors via Colors instead (one per slice).
	Color  string
	Colors []string
}

// LegendEnabled returns the resolved legend flag — default true unless
// the caller opted into HideLegend.
func (o Options) LegendEnabled() bool {
	return !o.HideLegend
}

// HeightPxOrDefault returns Options.HeightPx or 240 when zero.
func (o Options) HeightPxOrDefault() int {
	if o.HeightPx <= 0 {
		return 240
	}
	return o.HeightPx
}
