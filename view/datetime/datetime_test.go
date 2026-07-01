package datetime_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/datetime"
)

func renderLocal(t *testing.T, ts time.Time) string {
	t.Helper()
	var buf bytes.Buffer
	err := datetime.Local(ts).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestLocal_RendersUTCInstantAndLocalizerHook(t *testing.T) {
	// 09:30 in a +02:00 zone → 07:30 UTC. The datetime attribute must carry
	// the UTC instant so the browser localizer starts from an unambiguous
	// value, regardless of the server's zone.
	ts := time.Date(2026, 7, 1, 9, 30, 0, 0, time.FixedZone("CEST", 2*3600))
	out := renderLocal(t, ts)

	assert.Contains(t, out, `datetime="2026-07-01T07:30:00Z"`)
	assert.Contains(t, out, "data-ds-localtime")
	// Server fallback is the UTC wall-clock, suffixed so it is never
	// silently read as local time before the script runs.
	assert.Contains(t, out, "01/07/2026 07:30 UTC")
}

func TestLocal_ZeroTimeRendersNothing(t *testing.T) {
	out := renderLocal(t, time.Time{})
	assert.Equal(t, "", strings.TrimSpace(out))
}

func TestTimezoneBanner_RendersFillablePlaceholder(t *testing.T) {
	var buf bytes.Buffer
	err := datetime.TimezoneBanner().Render(context.Background(), &buf)
	assert.NoError(t, err)
	// Empty on the server; the base-layout script fills the text from the
	// browser's resolved zone. The hook attribute must be present.
	assert.Contains(t, buf.String(), "data-ds-tzlabel")
}
