package layout_test

import (
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
	"github.com/stretchr/testify/assert"
)

// renderTopbar is defined in space_switcher_test.go (same package).

// The "Ver todas" affordance must live in the notifications panel HEADER
// (next to the "Notificaciones" title), always visible without scrolling
// past the list — not at the bottom of the (potentially long) fragment.
// It points to the notifications base page and opens in a new tab.
func TestTopbar_NotificationsSeeAllInHeader(t *testing.T) {
	body := renderTopbar(t, layout.Config{
		Title:                "T",
		ProjectName:          "kiban",
		NotificationsBaseURL: "/kiban-cloud/notifications",
	})

	assert.Contains(t, body, ">Ver todas<", "header must expose a 'Ver todas' link")
	assert.Contains(t, body, `href="/kiban-cloud/notifications"`, "link points to the notifications page")
	assert.Contains(t, body, `target="_blank"`, "opens the full page in a new tab")

	// It must sit in the header, before the scroll container the fragment
	// swaps into — that's what guarantees "no scroll to reach it".
	seeAll := strings.Index(body, ">Ver todas<")
	list := strings.Index(body, `id="topbar-notifications-list"`)
	assert.NotEqual(t, -1, seeAll)
	assert.NotEqual(t, -1, list)
	assert.Less(t, seeAll, list, "'Ver todas' must precede the scrollable list")
}
