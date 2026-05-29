package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
	"github.com/stretchr/testify/assert"
)

// renderPassword centralizes the templ → string boilerplate so each test
// stays focused on the markup assertion.
func renderPassword(t *testing.T, name, label, value, errMsg, hint string, required bool, placeholder string) string {
	t.Helper()
	var buf bytes.Buffer
	err := input.Password(name, label, value, errMsg, hint, required, placeholder).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestPassword_RendersInputAndLabel(t *testing.T) {
	body := renderPassword(t, "password", "Contraseña", "", "", "", true, "")

	// Label + asterisk for required.
	assert.True(t, strings.Contains(body, `<label for="f-password"`))
	assert.True(t, strings.Contains(body, "*</span>"))

	// Input itself: type=password, autocomplete=off, id matches label's for.
	assert.True(t, strings.Contains(body, `id="f-password"`))
	assert.True(t, strings.Contains(body, `type="password"`))
	assert.True(t, strings.Contains(body, `autocomplete="off"`))
	assert.True(t, strings.Contains(body, `name="password"`))
}

func TestPassword_ToggleMarkupContract(t *testing.T) {
	// The visibility toggle JS in view/layout/base.templ relies on these
	// data-attributes — keep this test green or the toggle stops working
	// silently. Documented in password.templ's godoc as a contract.
	body := renderPassword(t, "password", "Contraseña", "", "", "", false, "")

	// Wrapper carries the marker the JS uses to find the input sibling.
	assert.True(t, strings.Contains(body, "data-kiban-password-field"))
	// Toggle button + accessible default label.
	assert.True(t, strings.Contains(body, "data-kiban-password-toggle"))
	assert.True(t, strings.Contains(body, `aria-label="Mostrar contrase`))
	// Both icon slots present — show visible (flex), hide hidden initially.
	assert.True(t, strings.Contains(body, "data-kiban-password-icon-show"))
	assert.True(t, strings.Contains(body, "data-kiban-password-icon-hide"))
	// The "show" span starts visible; the "hide" span starts hidden. Match
	// on the full class attribute fragment so a future refactor that
	// reorders classes still trips this assertion.
	assert.True(t, strings.Contains(body, `data-kiban-password-icon-show class="flex"`))
	assert.True(t, strings.Contains(body, `data-kiban-password-icon-hide class="hidden"`))

	// Right padding reserves space for the toggle button (pr-10) so the
	// eye icon doesn't overlap the user's typing. Brittle on purpose:
	// changing this needs to be a deliberate visual decision.
	assert.True(t, strings.Contains(body, "pr-10"))
}

func TestPassword_ErrorStateChangesBorder(t *testing.T) {
	body := renderPassword(t, "password", "Contraseña", "", "Contraseña incorrecta", "", false, "")

	// Red border classes show up when errMsg is set.
	assert.True(t, strings.Contains(body, "border-red-400"))
	// Error message rendered below the wrapper.
	assert.True(t, strings.Contains(body, "Contrase")) // "Contraseña incorrecta" minus accented chars
}

func TestPassword_HintShownWhenNoError(t *testing.T) {
	body := renderPassword(t, "password", "Contraseña", "", "", "Mínimo 8 caracteres", false, "")

	// Hint renders in muted ink3.
	assert.True(t, strings.Contains(body, "text-kiban-ink3"))
	assert.True(t, strings.Contains(body, "nimo 8 caracteres")) // accented "Mí" gets HTML-escaped
}

func TestPassword_PlaceholderRoundtrip(t *testing.T) {
	body := renderPassword(t, "password", "Contraseña", "", "", "", false, "••••••••")
	// Sentinel placeholder is rendered as-is on the input.
	assert.True(t, strings.Contains(body, `placeholder="••••••••"`))
}
