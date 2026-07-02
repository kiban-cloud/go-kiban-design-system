package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func renderPhone(t *testing.T, label, ccName, ccValue, phoneName, phoneValue, errMsg string, required bool) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, input.Phone(label, ccName, ccValue, phoneName, phoneValue, errMsg, required).Render(context.Background(), &buf))
	return buf.String()
}

// The visible tel input carries the client-side numeric constraints:
// inputmode="numeric" (mobile numeric keyboard) + maxlength="12" (national
// number cap, mirroring the backend [0-9]{9,12} validation). The digit-only
// stripping itself lives in base.templ's sync() JS.
func TestPhone_VisibleInputNumericConstraints(t *testing.T) {
	out := renderPhone(t, "Teléfono", "country_code", "", "phone_number", "", "", true)

	require.Contains(t, out, `inputmode="numeric"`, "visible tel input requests the numeric keyboard")
	require.Contains(t, out, `maxlength="12"`, "visible tel input caps the national number at 12 digits")
	require.Contains(t, out, `data-tel-visible`, "the visible input keeps its widget hook")
}

// The two hidden inputs (dial code + national number) still bind under the
// caller-supplied names — the constraints must not disturb the form payload.
func TestPhone_HiddenInputsBindByName(t *testing.T) {
	out := renderPhone(t, "Teléfono", "country_code", "+52", "phone_number", "5512345678", "", false)

	require.Contains(t, out, `name="country_code"`, "dial-code hidden input keeps its name")
	require.Contains(t, out, `name="phone_number"`, "national-number hidden input keeps its name")
	require.Contains(t, out, `data-tel-cc`, "dial-code hidden input keeps its widget hook")
	require.Contains(t, out, `data-tel-national`, "national-number hidden input keeps its widget hook")
}

// Error state renders the red message under the field.
func TestPhone_ErrorMessage(t *testing.T) {
	out := renderPhone(t, "Teléfono", "country_code", "", "phone_number", "", "Campo inválido.", true)

	require.Contains(t, out, "Campo inv", "error message is rendered")
	require.True(t, strings.Contains(out, "text-red-600"), "error message uses the red style")
}
