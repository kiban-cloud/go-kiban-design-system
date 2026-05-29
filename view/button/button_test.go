package button_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/button"
	"github.com/stretchr/testify/assert"
)

func renderButton(t *testing.T, opts button.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := button.Button(opts).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestButton_RendersButtonWhenNoHref(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Guardar"})
	assert.True(t, strings.Contains(body, "<button"))
	assert.True(t, strings.Contains(body, "Guardar"))
	assert.False(t, strings.Contains(body, "<a "))
}

func TestButton_RendersAnchorWhenHref(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Volver", Href: "/customers", Variant: "secondary"})
	assert.True(t, strings.Contains(body, `<a href="/customers"`))
	assert.True(t, strings.Contains(body, "Volver"))
}

func TestButton_HttpsHrefPassesThrough(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Docs", Href: "https://example.com/docs"})
	assert.True(t, strings.Contains(body, `href="https://example.com/docs"`))
}

// templ.URL replaces non-allowlisted schemes with about:invalid#... to
// block XSS via untrusted URLs. By default (AllowDataURL=false) the
// Button should refuse to emit a literal data: URL.
func TestButton_DataURLBlockedByDefault(t *testing.T) {
	body := renderButton(t, button.Options{
		Label: "Descargar",
		Href:  "data:text/csv;base64,YWJjZA==",
	})
	assert.False(t, strings.Contains(body, "data:text/csv"))
	assert.True(t, strings.Contains(body, "about:invalid"))
}

// Opt-in: when the caller explicitly trusts the URL (server-built CSV
// download, generated report), AllowDataURL bypasses sanitization via
// templ.SafeURL.
func TestButton_DataURLAllowedWhenOptIn(t *testing.T) {
	body := renderButton(t, button.Options{
		Label:        "Descargar",
		Href:         "data:text/csv;base64,YWJjZA==",
		AllowDataURL: true,
		Attrs:        templ.Attributes{"download": "report.csv"},
	})
	assert.True(t, strings.Contains(body, "data:text/csv;base64,YWJjZA=="))
	assert.False(t, strings.Contains(body, "about:invalid"))
	assert.True(t, strings.Contains(body, `download="report.csv"`))
}

// Sanity check: AllowDataURL doesn't affect non-data hrefs — the same
// http(s) URL renders identically with or without the flag.
func TestButton_AllowDataURLPreservesHttpsHref(t *testing.T) {
	body := renderButton(t, button.Options{
		Label:        "Docs",
		Href:         "https://example.com/docs",
		AllowDataURL: true,
	})
	assert.True(t, strings.Contains(body, `href="https://example.com/docs"`))
}

func TestButton_SubmitTypeRendered(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Guardar", IsSubmit: true})
	assert.True(t, strings.Contains(body, `type="submit"`))
}

func TestButton_DisabledWhenButton(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Guardar", Disabled: true})
	assert.True(t, strings.Contains(body, "disabled"))
}

func TestButton_DisabledIgnoredOnAnchor(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Volver", Href: "/x", Disabled: true})
	// `<a disabled>` isn't a thing — Disabled is silently dropped for anchors.
	// (We can't simply look for " disabled" because Tailwind variants like
	//  `disabled:opacity-50` appear in the class attribute regardless.)
	assert.False(t, strings.Contains(body, `disabled=""`))
	assert.False(t, strings.Contains(body, "<a disabled"))
}

func TestButton_OnClickEmittedAsAttribute(t *testing.T) {
	body := renderButton(t, button.Options{Label: "Cerrar", OnClick: "kibanCloseOverlay('x')"})
	assert.True(t, strings.Contains(body, `onclick="kibanCloseOverlay(&#39;x&#39;)"`))
}
