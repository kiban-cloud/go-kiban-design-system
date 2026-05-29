package layout_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
	"github.com/stretchr/testify/assert"
)

// renderTopbar runs layout.Topbar to a string so we can assert on the
// emitted markup. The Topbar templ is not exported as a standalone
// component in tests usually, but it's a normal templ.Component so we
// can call .Render directly.
func renderTopbar(t *testing.T, cfg layout.Config) string {
	t.Helper()
	var buf bytes.Buffer
	if err := layout.Topbar(cfg).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Topbar render: %v", err)
	}
	return buf.String()
}

func TestSpaceSwitcher_ReadonlyWithNoSpaces(t *testing.T) {
	body := renderTopbar(t, layout.Config{
		User:   layout.User{Email: "u@example.com"},
		Spaces: nil, // user has zero spaces
	})

	// Readonly fallback shows the "—" sentinel inside a non-interactive chip.
	assert.Contains(t, body, `cursor-default`)
	assert.Contains(t, body, "—")
	// No clickable trigger / dropdown panel.
	assert.NotContains(t, body, `data-kiban-menu="topbar-space-switcher"`)
	assert.NotContains(t, body, `kibanToggleMenu('topbar-space-switcher')`)
}

func TestSpaceSwitcher_ReadonlyWithSingleSpace(t *testing.T) {
	body := renderTopbar(t, layout.Config{
		User:              layout.User{Email: "u@example.com"},
		SpaceID:           "sp1",
		CurrentSpaceName:  "Mi Espacio",
		Spaces:            []layout.SpaceOption{{Id: "sp1", Name: "Mi Espacio"}},
		SwitchSpaceAction: "/cloud-htmx/spaces/switch",
	})

	// With a single space the chip is still readonly (no reason to open
	// a dropdown when there's nothing to choose).
	assert.Contains(t, body, `cursor-default`)
	// Current name shown, not the bare ObjectId.
	assert.Contains(t, body, "Mi Espacio")
	assert.NotContains(t, body, `data-kiban-menu="topbar-space-switcher"`)
}

func TestSpaceSwitcher_DropdownWithMultipleSpaces(t *testing.T) {
	body := renderTopbar(t, layout.Config{
		User:             layout.User{Email: "u@example.com"},
		SpaceID:          "sp1",
		CurrentSpaceName: "Espacio Uno",
		Spaces: []layout.SpaceOption{
			{Id: "sp1", Name: "Espacio Uno"},
			{Id: "sp2", Name: "Espacio Dos"},
		},
		SwitchSpaceAction: "/cloud-htmx/spaces/switch",
	})

	// Clickable trigger with kibanToggleMenu wiring.
	assert.Contains(t, body, `data-kiban-menu="topbar-space-switcher"`)
	assert.Contains(t, body, `kibanToggleMenu('topbar-space-switcher')`)
	// ARIA contract for the dropdown.
	assert.Contains(t, body, `aria-haspopup="true"`)
	assert.Contains(t, body, `role="menu"`)
	// Both spaces rendered as POST forms with the right action.
	assert.Equal(t, 2, strings.Count(body, `action="/cloud-htmx/spaces/switch"`))
	// Each carries the right spaceId.
	assert.Contains(t, body, `value="sp1"`)
	assert.Contains(t, body, `value="sp2"`)
	// Current space marked with "Actual".
	assert.Contains(t, body, "Actual")
}

func TestSpaceSwitcher_ReadonlyWhenActionMissing(t *testing.T) {
	// Defensive: even with multiple spaces, if SwitchSpaceAction is
	// empty (misconfigured consumer) the chip falls back to readonly so
	// no form ever posts to "".
	body := renderTopbar(t, layout.Config{
		User: layout.User{Email: "u@example.com"},
		Spaces: []layout.SpaceOption{
			{Id: "sp1", Name: "A"},
			{Id: "sp2", Name: "B"},
		},
		SwitchSpaceAction: "",
	})

	assert.Contains(t, body, `cursor-default`)
	assert.NotContains(t, body, `data-kiban-menu="topbar-space-switcher"`)
}

func TestSpaceSwitcher_FallbackToSpaceIDWhenNameEmpty(t *testing.T) {
	// CurrentSpaceName empty + SpaceID set → chip shows the ObjectId so
	// the user has at least some context, not "—".
	body := renderTopbar(t, layout.Config{
		User:    layout.User{Email: "u@example.com"},
		SpaceID: "6851bc072a3ed11fcb6e0429",
	})

	assert.Contains(t, body, "6851bc072a3ed11fcb6e0429")
}
