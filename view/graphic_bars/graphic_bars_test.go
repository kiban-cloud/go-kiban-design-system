package graphic_bars_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/graphic_bars"
)

func render(t *testing.T, opts graphic_bars.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := graphic_bars.GraphicBars(opts).Render(context.Background(), &buf)
	require.NoError(t, err)
	return buf.String()
}

// Empty input renders nothing — no title, no frame, no whitespace.
func TestGraphicBars_EmptyInputRendersNothing(t *testing.T) {
	assert.Equal(t, "", render(t, graphic_bars.Options{Title: "Estatus"}))
}

// A populated card renders the title, each bar's label inside the
// fill, the caller-formatted total, and the percent width on the fill.
func TestGraphicBars_RendersTitleLabelsTotalsAndWidth(t *testing.T) {
	body := render(t, graphic_bars.Options{
		Title:   "Satisfactorios",
		Variant: graphic_bars.VariantSuccess,
		Bars: []graphic_bars.Bar{
			{Label: "Alta", Total: "5", Percent: 60},
			{Label: "Baja", Total: "3", Percent: 40},
		},
	})
	assert.Contains(t, body, "Satisfactorios")
	// Labels in order.
	idxA := strings.Index(body, "Alta")
	idxB := strings.Index(body, "Baja")
	assert.True(t, idxA >= 0 && idxB > idxA)
	// Caller-formatted totals + percent labels.
	assert.Contains(t, body, ">5<")
	assert.Contains(t, body, "60%")
	// The success palette colours the fill + the in-bar label.
	assert.Contains(t, body, "#ecfdf5") // fill
	assert.Contains(t, body, "#047857") // label colour
	// Fill width tracks the percent.
	assert.Contains(t, body, "width:60%")
}

// HidePercent drops the right-column "%" (durations/costs) while still
// driving the fill width from Percent.
func TestGraphicBars_HidePercentDropsPercentLabel(t *testing.T) {
	body := render(t, graphic_bars.Options{
		Variant:     graphic_bars.VariantPrimary,
		HidePercent: true,
		Bars:        []graphic_bars.Bar{{Label: "Alta", Total: "1h 30m", Percent: 70}},
	})
	assert.Contains(t, body, "1h 30m")
	// The percent-label column is unique to the !HidePercent branch.
	assert.NotContains(t, body, "text-xs text-kiban-ink3")
	assert.Contains(t, body, "width:70%") // width still tracks Percent
}

// Each variant maps to its fill colour; an unknown variant collapses
// to the neutral palette so a bad value never renders an invisible bar.
func TestGraphicBars_VariantToFill(t *testing.T) {
	cases := []struct{ variant, fill string }{
		{graphic_bars.VariantError, "#ffe3e3"},
		{graphic_bars.VariantSuccess, "#ecfdf5"},
		{graphic_bars.VariantWarning, "#fff4d9"},
		{graphic_bars.VariantWorkfloo, "#f6efff"},
		{graphic_bars.VariantPrimary, "#f1f7ff"},
		{graphic_bars.VariantNeutral, "#f5f6f7"},
		{"WHATEVER", "#f5f6f7"}, // unknown → neutral
	}
	for _, tc := range cases {
		t.Run("variant="+tc.variant, func(t *testing.T) {
			body := render(t, graphic_bars.Options{
				Variant: tc.variant,
				Bars:    []graphic_bars.Bar{{Label: "x", Total: "1", Percent: 10}},
			})
			assert.Contains(t, body, tc.fill)
		})
	}
}

// Percent width is clamped to [0,100] so a bad upstream value can't
// overflow the row.
func TestGraphicBars_ClampsWidth(t *testing.T) {
	over := render(t, graphic_bars.Options{Bars: []graphic_bars.Bar{{Label: "x", Total: "1", Percent: 150}}})
	assert.Contains(t, over, "width:100%")
	under := render(t, graphic_bars.Options{Bars: []graphic_bars.Bar{{Label: "x", Total: "1", Percent: -20}}})
	assert.Contains(t, under, "width:0%")
}
