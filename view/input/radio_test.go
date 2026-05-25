package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/input"
	"github.com/stretchr/testify/assert"
)

func renderRadioCard(t *testing.T, name, value, title, subtitle string, checked bool, attrs templ.Attributes) string {
	t.Helper()
	var buf bytes.Buffer
	err := input.RadioCard(name, value, title, subtitle, checked, attrs).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestRadioCard_RendersTitleAndSubtitle(t *testing.T) {
	body := renderRadioCard(t, "type", "PF", "Persona Física", "Una persona individual", false, nil)
	assert.True(t, strings.Contains(body, "Persona F"))
	assert.True(t, strings.Contains(body, "Una persona individual"))
}

func TestRadioCard_EmitsDataMarker(t *testing.T) {
	// The data attribute is the contract that lets CSS in base.templ
	// flip the highlight via :has(input:checked). Removing it would
	// silently break the dynamic visual state in HTMX flows.
	body := renderRadioCard(t, "type", "PF", "Persona Física", "", false, nil)
	assert.True(t, strings.Contains(body, "data-kiban-radio-card"))
}

func TestRadioCard_NoServerRenderedHighlightClasses(t *testing.T) {
	// Highlight comes from :has(input:checked) CSS in base.templ, NOT
	// from server-side classes. If anyone re-adds border-kiban-primary
	// / bg-kiban-primary-soft to the label, the highlight would stick
	// to the originally-checked card across HTMX partial swaps.
	bodyChecked := renderRadioCard(t, "type", "PF", "Persona Física", "", true, nil)
	bodyUnchecked := renderRadioCard(t, "type", "PM", "Persona Moral", "", false, nil)
	for _, body := range []string{bodyChecked, bodyUnchecked} {
		assert.False(t, strings.Contains(body, "border-kiban-primary "))
		assert.False(t, strings.Contains(body, "bg-kiban-primary-soft"))
	}
}

func TestRadioCard_CheckedAttributeOnInput(t *testing.T) {
	bodyChecked := renderRadioCard(t, "type", "PF", "Persona Física", "", true, nil)
	bodyUnchecked := renderRadioCard(t, "type", "PM", "Persona Moral", "", false, nil)
	assert.True(t, strings.Contains(bodyChecked, "checked"))
	// Unchecked rendering omits the `checked` attribute entirely (templ's
	// `checked?={…}` syntax). Browser will then pick the first checked
	// sibling, or none, by native radio-group rules.
	assert.False(t, strings.Contains(bodyUnchecked, ` checked`))
}

func TestRadioCard_SubtitleHiddenWhenEmpty(t *testing.T) {
	body := renderRadioCard(t, "type", "PF", "Persona Física", "", false, nil)
	// No `<div class="text-xs text-kiban-ink3">` row when subtitle is empty.
	assert.False(t, strings.Contains(body, "text-kiban-ink3"))
}

func TestRadioCard_AttrsSpreadOnInput(t *testing.T) {
	body := renderRadioCard(t, "type", "PF", "Persona Física", "", false, templ.Attributes{
		"hx-post":   "/preview",
		"hx-target": "#form-section",
	})
	assert.True(t, strings.Contains(body, `hx-post="/preview"`))
	assert.True(t, strings.Contains(body, `hx-target="#form-section"`))
}

func TestRadioCard_EmitsNameAndValue(t *testing.T) {
	body := renderRadioCard(t, "type", "PM", "Persona Moral", "", false, nil)
	assert.True(t, strings.Contains(body, `name="type"`))
	assert.True(t, strings.Contains(body, `value="PM"`))
	assert.True(t, strings.Contains(body, `type="radio"`))
}
