package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
	"github.com/stretchr/testify/assert"
)

func renderCode(t *testing.T, name, label, value, errMsg string, length int, required bool) string {
	t.Helper()
	var buf bytes.Buffer
	err := input.Code(name, label, value, errMsg, length, required).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestCode_EnforcesFixedLength(t *testing.T) {
	body := renderCode(t, "nip", "Código", "", "", 6, true)

	assert.True(t, strings.Contains(body, `name="nip"`))
	assert.True(t, strings.Contains(body, `maxlength="6"`))
	assert.True(t, strings.Contains(body, `pattern="\d{6}"`))
	assert.True(t, strings.Contains(body, `inputmode="numeric"`))
	assert.True(t, strings.Contains(body, `autocomplete="one-time-code"`))
	assert.True(t, strings.Contains(body, "required"))
	// Placeholder shows one bullet per expected digit.
	assert.True(t, strings.Contains(body, `placeholder="••••••"`))
}

func TestCode_LengthZeroDefaultsToSix(t *testing.T) {
	body := renderCode(t, "nip", "", "", "", 0, true)

	assert.True(t, strings.Contains(body, `maxlength="6"`))
	assert.True(t, strings.Contains(body, `pattern="\d{6}"`))
}

func TestCode_CustomLength(t *testing.T) {
	body := renderCode(t, "nip", "", "", "", 4, false)

	assert.True(t, strings.Contains(body, `maxlength="4"`))
	assert.True(t, strings.Contains(body, `pattern="\d{4}"`))
	assert.True(t, strings.Contains(body, `placeholder="••••"`))
	// Not required → no required attr, no asterisk.
	assert.False(t, strings.Contains(body, "*</span>"))
}

func TestCode_LabelOptional(t *testing.T) {
	withLabel := renderCode(t, "nip", "Código", "", "", 6, true)
	assert.True(t, strings.Contains(withLabel, `<label for="f-nip"`))

	noLabel := renderCode(t, "nip", "", "", "", 6, true)
	assert.False(t, strings.Contains(noLabel, `<label for="f-nip"`))
}

func TestCode_ErrorStateChangesBorder(t *testing.T) {
	body := renderCode(t, "nip", "Código", "", "Código inválido", 6, true)

	assert.True(t, strings.Contains(body, "border-red-400"))
	assert.True(t, strings.Contains(body, "digo inv")) // "Código inválido" minus accented chars
}

func TestCode_ValueRoundtrip(t *testing.T) {
	body := renderCode(t, "nip", "Código", "123456", "", 6, true)
	assert.True(t, strings.Contains(body, `value="123456"`))
}
