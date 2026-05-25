package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func renderAutocomplete(t *testing.T, opts input.AutocompleteOptions) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, input.Autocomplete(opts).Render(context.Background(), &buf))
	return buf.String()
}

func TestAutocomplete_RendersWrapperWithDataAttributes(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:  "category",
		Label: "Categoría",
		Value: "alta",
		Items: []input.SelectOption{
			{Value: "alta", Label: "Alta"},
			{Value: "media", Label: "Media"},
		},
		AllowCreate: true,
	})
	require.Contains(t, out, `data-kiban-autocomplete`, "wrapper carries the autocomplete data marker for JS init")
	require.Contains(t, out, `data-allow-create="true"`, "AllowCreate=true ships as the data attribute")
	require.Contains(t, out, `data-kiban-autocomplete-input`, "visible text input is marked")
	require.Contains(t, out, `data-kiban-autocomplete-value`, "hidden submit input is marked")
	require.Contains(t, out, `data-kiban-autocomplete-list`, "dropdown list marker")
	require.Contains(t, out, `data-kiban-autocomplete-create`, "create-new row marker present when AllowCreate=true")
}

func TestAutocomplete_VisibleInputShowsLabelForCurrentValue(t *testing.T) {
	// Value matches one of the items → visible input shows the
	// human-readable label rather than the raw value.
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:  "category",
		Value: "alta",
		Items: []input.SelectOption{
			{Value: "alta", Label: "Alta prioridad"},
			{Value: "media", Label: "Media"},
		},
	})
	require.Contains(t, out, `value="Alta prioridad"`, "visible input shows the matching item's label")
	require.Contains(t, out, `value="alta"`, "hidden input keeps the raw value")
}

func TestAutocomplete_NameRoundTripsAsHiddenInput(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:  "rule_r1_outputLabel",
		Value: "premium",
	})
	require.Contains(t, out,
		`<input type="hidden" name="rule_r1_outputLabel" value="premium"`,
		"the form-submitted value rides on the named hidden input",
	)
	// The visible input is unnamed (it's UX-only).
	visibleIdx := strings.Index(out, `data-kiban-autocomplete-input`)
	if visibleIdx < 0 {
		t.Fatalf("missing visible input: %s", out)
	}
	// Grab the opening tag of the visible input.
	openIdx := strings.LastIndex(out[:visibleIdx], "<input")
	tagEnd := strings.Index(out[openIdx:], ">")
	openTag := out[openIdx : openIdx+tagEnd+1]
	require.NotContains(t, openTag, `name="rule_r1_outputLabel"`,
		"the visible input must NOT carry the form-submit name (only the hidden does)")
}

func TestAutocomplete_DoesNotEmitCreateRowWhenAllowCreateFalse(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:        "x",
		AllowCreate: false,
	})
	require.NotContains(t, out, `data-kiban-autocomplete-create`,
		"create-new row should not render when AllowCreate=false")
	require.Contains(t, out, `data-allow-create="false"`)
}

func TestAutocomplete_RendersItemsWithValueAndLabelAttributes(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name: "x",
		Items: []input.SelectOption{
			{Value: "a", Label: "Alpha"},
			{Value: "b", Label: "Beta"},
		},
	})
	require.Contains(t, out, `data-value="a"`)
	require.Contains(t, out, `data-label="Alpha"`)
	require.Contains(t, out, `data-value="b"`)
	require.Contains(t, out, `data-label="Beta"`)
}

func TestAutocomplete_RendersErrorMessage(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:   "x",
		ErrMsg: "Etiqueta obligatoria.",
	})
	require.Contains(t, out, "Etiqueta obligatoria.")
	// error styling on the visible input.
	require.Contains(t, out, "border-red-400")
}

func TestAutocomplete_RendersHintWhenNoError(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name: "x",
		Hint: "Elige o crea una etiqueta.",
	})
	require.Contains(t, out, "Elige o crea una etiqueta.")
}

func TestAutocomplete_CreateLabelOverride(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:        "x",
		AllowCreate: true,
		CreateLabel: "Añadir nuevo",
	})
	require.Contains(t, out, `data-create-label="Añadir nuevo"`,
		"CreateLabel override flows into the data attribute")
}

func TestAutocomplete_RequiredAsterisk(t *testing.T) {
	out := renderAutocomplete(t, input.AutocompleteOptions{
		Name:     "x",
		Label:    "Mensaje",
		Required: true,
	})
	// Required renders a red asterisk inside the <label>.
	idx := strings.Index(out, ">Mensaje")
	if idx < 0 {
		t.Fatalf("missing label: %s", out)
	}
	chunk := out[idx : idx+200]
	require.Contains(t, chunk, "text-red-500")
}
