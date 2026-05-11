package chip_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/chip"
	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, opts chip.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := chip.Chip(opts).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestChip_RendersLabel(t *testing.T) {
	body := render(t, chip.Options{Label: "report.csv"})
	assert.True(t, strings.Contains(body, "report.csv"))
	// Default variant: neutral surface + ink2.
	assert.True(t, strings.Contains(body, "border-kiban-border"))
	assert.True(t, strings.Contains(body, "bg-kiban-surface"))
	// No remove button by default.
	assert.False(t, strings.Contains(body, "data-chip-remove"))
}

func TestChip_RendersSubtext(t *testing.T) {
	body := render(t, chip.Options{Label: "report.csv", Subtext: "12 KB"})
	assert.True(t, strings.Contains(body, "12 KB"))
}

func TestChip_RemovableEmitsButtonWithDataAttr(t *testing.T) {
	body := render(t, chip.Options{Label: "x", Removable: true})
	assert.True(t, strings.Contains(body, "data-chip-remove"))
	assert.True(t, strings.Contains(body, `aria-label="Quitar"`))
	assert.True(t, strings.Contains(body, `type="button"`))
}

func TestChip_RemoveAttrsSpreadOnButton(t *testing.T) {
	body := render(t, chip.Options{
		Label:       "x",
		Removable:   true,
		RemoveAttrs: templ.Attributes{"data-remove": "42"},
	})
	assert.True(t, strings.Contains(body, `data-remove="42"`))
}

func TestChip_RemoveAriaLabelOverride(t *testing.T) {
	body := render(t, chip.Options{
		Label:           "x",
		Removable:       true,
		RemoveAriaLabel: "Drop file",
	})
	assert.True(t, strings.Contains(body, `aria-label="Drop file"`))
}

func TestChip_DangerVariantUsesRedTokens(t *testing.T) {
	body := render(t, chip.Options{Label: "bad.csv", Variant: "danger"})
	assert.True(t, strings.Contains(body, "border-red-400"))
	assert.True(t, strings.Contains(body, "bg-red-50"))
	assert.True(t, strings.Contains(body, "text-red-700"))
}

func TestChip_InfoVariantUsesPrimary(t *testing.T) {
	body := render(t, chip.Options{Label: "x", Variant: "info"})
	assert.True(t, strings.Contains(body, "kiban-primary"))
}

func TestChip_TitleEmitsTooltipAttr(t *testing.T) {
	body := render(t, chip.Options{Label: "x", Title: "explainer"})
	assert.True(t, strings.Contains(body, `title="explainer"`))
}

func TestChip_AttrsSpreadOnOuterSpan(t *testing.T) {
	body := render(t, chip.Options{
		Label: "x",
		Attrs: templ.Attributes{"data-test": "thing"},
	})
	assert.True(t, strings.Contains(body, `data-test="thing"`))
}

func TestChip_UnknownVariantFallsBackToDefault(t *testing.T) {
	// Defensive: an unknown variant string from data shouldn't crash
	// the render — we just fall back to the neutral palette so the chip
	// still appears.
	body := render(t, chip.Options{Label: "x", Variant: "made-up"})
	assert.True(t, strings.Contains(body, "border-kiban-border"))
	assert.True(t, strings.Contains(body, "bg-kiban-surface"))
}
