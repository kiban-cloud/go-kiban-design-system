package timeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/timeline"
)

func renderTimeline(t *testing.T, events []timeline.Event) string {
	t.Helper()
	var buf bytes.Buffer
	err := timeline.Timeline(events).Render(context.Background(), &buf)
	require.NoError(t, err)
	return buf.String()
}

// Empty input renders nothing — no <ul>, no whitespace artefact.
func TestTimeline_EmptyInputRendersNothing(t *testing.T) {
	body := renderTimeline(t, nil)
	assert.Equal(t, "", body)
}

// Populated input: each event surfaces label + date, kinds map
// to expected dot colours, and the order is preserved.
func TestTimeline_RendersEventsInOrder(t *testing.T) {
	body := renderTimeline(t, []timeline.Event{
		{Label: "Preparando", Kind: timeline.KindWarning, Date: "15/05/2026 15:18"},
		{Label: "Enviado", Kind: timeline.KindDefault, Date: "15/05/2026 15:19"},
		{Label: "Recibido", Kind: timeline.KindSuccess, Date: "15/05/2026 15:20"},
	})

	assert.True(t, strings.Contains(body, "<ul"))
	// Labels in order.
	idxA := strings.Index(body, "Preparando")
	idxB := strings.Index(body, "Enviado")
	idxC := strings.Index(body, "Recibido")
	assert.True(t, idxA >= 0 && idxB > idxA && idxC > idxB)
	// Date column.
	assert.True(t, strings.Contains(body, "15/05/2026 15:18"))
	// Dot colours.
	assert.True(t, strings.Contains(body, "bg-amber-500"))   // warning
	assert.True(t, strings.Contains(body, "bg-emerald-500")) // success
	assert.True(t, strings.Contains(body, "bg-kiban-ink4"))  // default
}

// Empty Date hides the date column without breaking the row layout.
func TestTimeline_EmptyDateHidesDateColumn(t *testing.T) {
	body := renderTimeline(t, []timeline.Event{
		{Label: "Sin fecha", Kind: timeline.KindInfo, Date: ""},
	})
	assert.True(t, strings.Contains(body, "Sin fecha"))
	// No date-styling div emitted (the text-xs text-kiban-ink3 class
	// is unique to the date column).
	assert.False(t, strings.Contains(body, "text-xs text-kiban-ink3"))
}

// Every Kind constant maps to a dot class — pins the kind →
// colour table so a refactor doesn't accidentally drop one.
func TestTimeline_KindToDotClass(t *testing.T) {
	cases := []struct {
		kind     string
		wantClass string
	}{
		{timeline.KindSuccess, "bg-emerald-500"},
		{timeline.KindWarning, "bg-amber-500"},
		{timeline.KindInfo, "bg-kiban-primary"},
		{timeline.KindDanger, "bg-red-500"},
		{timeline.KindDefault, "bg-kiban-ink4"},
		// Unknown kinds collapse to the default appearance.
		{"WHATEVER", "bg-kiban-ink4"},
		{"", "bg-kiban-ink4"},
	}
	for _, tc := range cases {
		t.Run("kind="+tc.kind, func(t *testing.T) {
			body := renderTimeline(t, []timeline.Event{{Label: "x", Kind: tc.kind}})
			assert.True(t, strings.Contains(body, tc.wantClass),
				"expected %q dot class for kind %q", tc.wantClass, tc.kind)
		})
	}
}
