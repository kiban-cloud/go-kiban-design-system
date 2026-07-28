package layout_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
	"github.com/stretchr/testify/assert"
)

// renderBase runs layout.Base (the HTML shell shared by every kiban module)
// to a string so we can assert on the global head/body wiring. Base takes
// children; we pass an empty component since the assertions here are about
// the shell itself.
func renderBase(t *testing.T, cfg layout.Config) string {
	t.Helper()
	empty := templ.ComponentFunc(func(_ context.Context, _ io.Writer) error { return nil })
	var buf bytes.Buffer
	if err := layout.Base(cfg).Render(templ.WithChildren(context.Background(), empty), &buf); err != nil {
		t.Fatalf("Base render: %v", err)
	}
	return buf.String()
}

// The global HTMX busy-region rule must be wired into the shared shell so
// every module (link / rekon / crm / klin) gets consistent "blocked +
// something is happening" feedback on a partial swap with no per-button
// wiring. This guards the three moving parts: the CSS class, the spinner
// overlay, and the JS that flags the swap target on htmx:beforeRequest /
// clears it on htmx:afterRequest.
func TestBase_DataHrefRowShowsNavLoader(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	// The data-href row-click delegation must trigger the shared
	// #nav-loader overlay before navigating, so list→detail row clicks
	// aren't a frozen wait.
	assert.Contains(t, body, "getElementById('nav-loader')")
	assert.Contains(t, body, "navLoader.classList.add('is-visible')")
	assert.Contains(t, body, "window.location.href = el.getAttribute('data-href')")
}

// The shared HTMX error-response pattern (see the htmxerror package) needs
// three things wired into the shell so every module gets non-200 errors that
// Cloud Monitoring can see: the response-targets extension loaded and enabled
// on <body>, the #global-error-banner landing zone for 5xx fragments, and the
// beforeRequest cleanup that clears stale feedback so a new error never shows
// next to an old one.
func TestBase_HtmxErrorResponsePattern(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	// Extension script loaded and enabled on the body so hx-target-4xx /
	// hx-target-5xx resolve.
	assert.Contains(t, body, "htmx-ext-response-targets@2.0.4")
	assert.Contains(t, body, `hx-ext="response-targets"`)

	// Landing zone for 5xx fragments routed via hx-target-5xx.
	assert.Contains(t, body, `id="global-error-banner"`)

	// Cleanup: clear the banner and any data-error-type / data-success-type
	// feedback before each request.
	assert.Contains(t, body, "htmx:beforeRequest")
	assert.Contains(t, body, "[data-error-type],[data-success-type]")
}

func TestBase_GlobalHtmxBusyRule(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:       "T",
		ProjectName: "kiban",
	})

	// CSS: the busy region is non-interactive and gets a spinner overlay.
	assert.Contains(t, body, ".kiban-busy")
	assert.Contains(t, body, "pointer-events: none")
	assert.Contains(t, body, ".kiban-busy::after")

	// Trigger gets disabled during its own request (anti double-click),
	// inherited from <body> by every htmx control.
	assert.Contains(t, body, `hx-disabled-elt="this"`)

	// JS: flag the swap target on request start, clear it when the request
	// finishes (success or error).
	assert.Contains(t, body, "htmx:beforeRequest")
	assert.Contains(t, body, "htmx:afterRequest")
	assert.Contains(t, body, "classList.add('kiban-busy')")
	assert.Contains(t, body, "classList.remove('kiban-busy')")
}

// Non-HTMX navigations (a form submit that posts then redirects, or a
// [data-block-ui] button) must show a full-screen blocking overlay so the
// user gets feedback + can't double-submit during the server round-trip.
// The overlay lives in Base (not just Layout) so login-style bare pages get
// it too. This guards the element, the CSS, and the submit/click handlers —
// including the skip for HTMX forms (those use the per-region block).
func TestBase_GlobalActionOverlay(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:       "T",
		ProjectName: "kiban",
	})

	// The overlay element + its CSS.
	assert.Contains(t, body, `id="kiban-action-overlay"`)
	assert.Contains(t, body, ".kiban-action-overlay")
	assert.Contains(t, body, "position: fixed")

	// JS: show on a real form submit, opt-in via [data-block-ui].
	assert.Contains(t, body, "addEventListener('submit'")
	assert.Contains(t, body, "data-block-ui")

	// HTMX forms / boosted forms are skipped (region block handles them);
	// forms can also opt out with data-no-block-ui.
	assert.Contains(t, body, "hx-boost")
	assert.Contains(t, body, "data-no-block-ui")
}

func TestBase_CustomStylesheetHref(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:                "T",
		ProjectName:          "kiban",
		CustomStylesheetHref: "/portal/theme.css",
	})

	// The per-client stylesheet is emitted as a <link> and must appear after
	// the design-system <style> block so it wins the cascade.
	assert.Contains(t, body, `<link rel="stylesheet" href="/portal/theme.css">`)
	styleEnd := strings.LastIndex(body, "</style>")
	linkAt := strings.Index(body, `href="/portal/theme.css"`)
	assert.Greater(t, linkAt, styleEnd, "custom stylesheet must come after the base <style> block")
}

func TestBase_NoCustomStylesheetByDefault(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.NotContains(t, body, `rel="stylesheet" href="/portal`)
}

// The code editor must hide its source <textarea> from CSS, never from the
// inline style CodeMirror.fromTextArea() writes.
//
// htmx's settle phase calls cloneAttributes() on every element it matches by
// id across a swap, dropping any attribute the server did not send. `style`
// is never sent, so the display:none CodeMirror sets right after
// htmx:afterSwap is stripped ~20ms later and the raw textarea reappears next
// to its own widget — two editors stacked on the same field. The wrapper
// carries no id, so a rule keyed on its ready flag survives the settle.
func TestBase_CodeEditorHidesTextareaViaCSS(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:          "T",
		ProjectName:    "kiban",
		LoadCodeMirror: true,
	})

	assert.Contains(t, body, `[data-kiban-code-editor-ready="1"] > textarea`)
	// The flag the rule hangs off must actually be set once a widget mounts.
	assert.Contains(t, body, `setAttribute('data-kiban-code-editor-ready', '1')`)
}
