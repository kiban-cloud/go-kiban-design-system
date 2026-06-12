package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func renderCheckbox(t *testing.T, name, label string, value, enabled bool, attrs templ.Attributes, hint string) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, input.Checkbox(name, label, value, enabled, attrs, hint).Render(context.Background(), &buf))
	return buf.String()
}

// Default (no `value` in attrs): the load-bearing value="true" must be
// emitted so a single bool checkbox binds correctly under Gin.
func TestCheckbox_DefaultValueTrue(t *testing.T) {
	out := renderCheckbox(t, "acceptTerms", "Acepto", false, true, nil, "")
	require.Contains(t, out, `value="true"`, "default checkbox emits value=\"true\" for bool binding")
}

// Override: passing `value` in attrs replaces the default instead of
// emitting a SECOND value attribute. Regression guard for the duplicate-
// `value` bug that made multi-value checkbox groups (roles, members,
// rights) submit the literal "true" for every checked box.
func TestCheckbox_AttrsValueOverridesDefault(t *testing.T) {
	out := renderCheckbox(t, "roles", "Admin", true, true,
		templ.Attributes{"value": "SPACE_ADMIN"}, "")

	require.Contains(t, out, `value="SPACE_ADMIN"`, "override value reaches the rendered input")
	require.NotContains(t, out, `value="true"`, "default value=\"true\" must not also be emitted")
	require.Equal(t, 1, strings.Count(out, " value="), "exactly one value attribute is rendered")
}
