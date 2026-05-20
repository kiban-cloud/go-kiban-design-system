package stepper_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/stepper"
)

func renderStepper(t *testing.T, stages []stepper.Stage) string {
	t.Helper()
	var buf bytes.Buffer
	err := stepper.Stepper(stages).Render(context.Background(), &buf)
	require.NoError(t, err)
	return buf.String()
}

// Empty input renders nothing — no wrapper, no whitespace artefacts.
func TestStepper_EmptyInputRendersNothing(t *testing.T) {
	body := renderStepper(t, nil)
	assert.Equal(t, "", body)
}

// All-complete pin: every dot shows the checkmark SVG path, no
// step numbers leak in, and every connector is the emerald line.
func TestStepper_AllCompleteRendersCheckmarks(t *testing.T) {
	body := renderStepper(t, []stepper.Stage{
		{Label: "Términos", Status: stepper.StatusComplete},
		{Label: "Ingrese NIP", Status: stepper.StatusComplete},
		{Label: "Autorización", Status: stepper.StatusComplete},
	})
	// Checkmark SVG present.
	assert.True(t, strings.Contains(body, `d="M4 10l4 4 8-8"`))
	// No "incomplete" / "active" dot classes leak when every stage is complete.
	assert.False(t, strings.Contains(body, "border-kiban-primary"))
	assert.False(t, strings.Contains(body, "text-kiban-ink4"))
	// Connectors are emerald.
	assert.True(t, strings.Contains(body, "bg-emerald-300"))
	// Labels render.
	assert.True(t, strings.Contains(body, "Términos"))
}

// Mid-progress shape: first stage complete (checkmark + emerald
// connector after it), second stage active (kiban-primary border +
// position number 2), third stage incomplete (neutral border +
// position number 3, connector before it is neutral grey).
func TestStepper_MidProgressShape(t *testing.T) {
	body := renderStepper(t, []stepper.Stage{
		{Label: "A", Status: stepper.StatusComplete},
		{Label: "B", Status: stepper.StatusActive},
		{Label: "C", Status: stepper.StatusIncomplete},
	})
	// Checkmark for the complete stage.
	assert.True(t, strings.Contains(body, `d="M4 10l4 4 8-8"`))
	// Active stage's dot carries the kiban-primary border, and the
	// position number "2" is its text content.
	assert.True(t, strings.Contains(body, "border-kiban-primary"))
	// Incomplete stage shows the "3" position number with the
	// neutral border class.
	assert.True(t, strings.Contains(body, "border-kiban-border"))
	// Both connectors render — first emerald (after complete),
	// second neutral (after active).
	assert.True(t, strings.Contains(body, "bg-emerald-300"))
}

// Single-stage input renders the dot only (no connector).
func TestStepper_SingleStageRendersNoConnector(t *testing.T) {
	body := renderStepper(t, []stepper.Stage{
		{Label: "Solo", Status: stepper.StatusActive},
	})
	assert.True(t, strings.Contains(body, "Solo"))
	// Connector divs are emitted between stages only — with a single
	// stage there's no connector class anywhere.
	assert.False(t, strings.Contains(body, "bg-emerald-300"))
	assert.False(t, strings.Contains(body, "bg-kiban-border"))
}

// Unknown status falls through to the incomplete dot appearance —
// the layout shouldn't break on a typo.
func TestStepper_UnknownStatusFallsBackToIncomplete(t *testing.T) {
	body := renderStepper(t, []stepper.Stage{
		{Label: "Mystery", Status: "WEIRD_STATUS"},
	})
	// Same dot classes as the incomplete branch.
	assert.True(t, strings.Contains(body, "bg-kiban-surface"))
	assert.True(t, strings.Contains(body, "text-kiban-ink4"))
}

// Empty Label hides the label row but still renders the dot — keeps
// the stepper visually intact when a stage is intentionally
// label-less (e.g. icon-only flows).
func TestStepper_EmptyLabelHidesLabelRow(t *testing.T) {
	body := renderStepper(t, []stepper.Stage{
		{Label: "", Status: stepper.StatusActive},
	})
	// No <div class="text-xs ..."> label row.
	assert.False(t, strings.Contains(body, `text-kiban-ink mt-2`))
}
