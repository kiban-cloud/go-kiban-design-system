package tabs_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/tabs"
)

// bodyEntry is one slot of test content: a tab key + the literal
// HTML/text the test wants to surface inside that tab's Body.
type bodyEntry struct {
	Key     string
	Content string
}

// renderPanel composes a Panel + Body children block via
// templ.WithChildren. Panel takes children, each Body takes its
// own children — `bodies` flattens both levels: for each entry,
// we render `Body(key)` with literal content as its child.
func renderPanel(t *testing.T, cfg tabs.PanelConfig, bodies []bodyEntry) string {
	t.Helper()
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, b := range bodies {
			content := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, b.Content)
				return err
			})
			bodyCtx := templ.WithChildren(ctx, content)
			if err := tabs.Body(b.Key).Render(bodyCtx, w); err != nil {
				return err
			}
		}
		return nil
	})
	ctx := templ.WithChildren(context.Background(), children)
	var buf bytes.Buffer
	require.NoError(t, tabs.Panel(cfg).Render(ctx, &buf))
	return buf.String()
}

// Strip buttons + bodies both render. Switching is CSS-driven (no
// template-level "only render the active body"), so every body
// shows up in the markup with its key attribute for the CSS rule
// to dispatch on.
func TestPanel_RendersStripAndBodies(t *testing.T) {
	body := renderPanel(t, tabs.PanelConfig{
		ID:        "test-panel",
		ActiveKey: "consulta",
		Tabs: []tabs.TabHeader{
			{Key: "consulta", Label: "Consulta"},
			{Key: "detalles", Label: "Detalles"},
		},
	}, []bodyEntry{
		{Key: "consulta", Content: "RFC content here"},
		{Key: "detalles", Content: "Detalles content here"},
	})

	// Strip buttons carry the right keys + labels.
	assert.True(t, strings.Contains(body, `data-kiban-tabs-key="consulta"`))
	assert.True(t, strings.Contains(body, `data-kiban-tabs-key="detalles"`))
	assert.True(t, strings.Contains(body, ">Consulta<"))
	assert.True(t, strings.Contains(body, ">Detalles<"))
	// Panel root carries the active-key attribute.
	assert.True(t, strings.Contains(body, `data-kiban-tabs-active-key="consulta"`))
	// Both bodies render — visibility is CSS-driven, not template-driven.
	assert.True(t, strings.Contains(body, "RFC content here"))
	assert.True(t, strings.Contains(body, "Detalles content here"))
}

// HasError stamps a red dot on the offending tab's pill so a
// validation error living in an inactive (hidden) body is still
// signalled. Tabs without HasError get no dot.
func TestPanel_HasErrorRendersDot(t *testing.T) {
	body := renderPanel(t, tabs.PanelConfig{
		ID:        "err-panel",
		ActiveKey: "prod",
		Tabs: []tabs.TabHeader{
			{Key: "prod", Label: "Producción"},
			{Key: "sandbox", Label: "Sandbox", HasError: true},
		},
	}, []bodyEntry{
		{Key: "prod", Content: "prod body"},
		{Key: "sandbox", Content: "sandbox body"},
	})
	// Exactly one dot, on the Sandbox pill.
	assert.Equal(t, 1, strings.Count(body, "bg-red-500"))
}

// Empty Tabs renders just the body container (no strip). The
// caller's children still render so a degenerate one-tab layout
// can still flow through Panel without visual chrome.
func TestPanel_EmptyTabsRendersBodyContainerOnly(t *testing.T) {
	body := renderPanel(t, tabs.PanelConfig{
		ID:        "empty-panel",
		ActiveKey: "",
		Tabs:      nil,
	}, []bodyEntry{
		{Key: "anything", Content: "lone content"},
	})
	assert.False(t, strings.Contains(body, "data-kiban-tabs-tab"))
	assert.True(t, strings.Contains(body, "lone content"))
	assert.True(t, strings.Contains(body, "data-kiban-tabs-panel"))
}

// Each Body wrapper carries [data-kiban-tabs-body data-kiban-tabs-key="…"]
// — the exact selectors the CSS in base.templ matches against the
// panel's active-key to show only one body at a time.
func TestPanel_BodiesCarryKeyForCSSDispatch(t *testing.T) {
	body := renderPanel(t, tabs.PanelConfig{
		ID:        "p",
		ActiveKey: "a",
		Tabs:      []tabs.TabHeader{{Key: "a", Label: "A"}, {Key: "b", Label: "B"}},
	}, []bodyEntry{
		{Key: "a", Content: "A body"},
		{Key: "b", Content: "B body"},
	})
	assert.True(t, strings.Contains(body, `data-kiban-tabs-body data-kiban-tabs-key="a"`))
	assert.True(t, strings.Contains(body, `data-kiban-tabs-body data-kiban-tabs-key="b"`))
}

// Panel root uses the caller-supplied ID so JS scoping (when
// multiple panels coexist on a page) targets exactly one panel
// at a time.
func TestPanel_RootCarriesCallerID(t *testing.T) {
	body := renderPanel(t, tabs.PanelConfig{
		ID:        "step-tabs-n-1",
		ActiveKey: "data",
		Tabs:      []tabs.TabHeader{{Key: "data", Label: "Datos"}},
	}, nil)
	assert.True(t, strings.Contains(body, `id="step-tabs-n-1"`))
}
