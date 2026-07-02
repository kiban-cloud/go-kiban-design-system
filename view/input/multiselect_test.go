package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func renderMultiSelect(t *testing.T, opts input.MultiSelectOptions) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, input.MultiSelect(opts).Render(context.Background(), &buf))
	return buf.String()
}

func TestMultiSelect_RendersWrapperWithDataAttributes(t *testing.T) {
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:  "estado",
		Label: "Estado del domicilio",
		Items: []input.SelectOption{
			{Value: "otro", Label: "Otro"},
			{Value: "oax", Label: "Oaxaca"},
		},
		AllowCreate: true,
	})
	require.Contains(t, out, `data-kiban-multiselect`, "wrapper carries the multiselect data marker for JS init")
	require.Contains(t, out, `data-name="estado"`, "submit name flows onto the wrapper for JS-built chips")
	require.Contains(t, out, `data-allow-create="true"`)
	require.Contains(t, out, `data-kiban-multiselect-control`)
	require.Contains(t, out, `data-kiban-multiselect-input`)
	require.Contains(t, out, `data-kiban-multiselect-list`)
	require.Contains(t, out, `data-kiban-multiselect-create`, "create row present when AllowCreate=true")
}

func TestMultiSelect_OptionsCarryValueAndLabel(t *testing.T) {
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name: "x",
		Items: []input.SelectOption{
			{Value: "otro", Label: "Otro"},
			{Value: "oax", Label: "Oaxaca"},
		},
	})
	require.Contains(t, out, `data-value="otro"`)
	require.Contains(t, out, `data-label="Otro"`)
	require.Contains(t, out, `data-value="oax"`)
	require.Contains(t, out, `data-label="Oaxaca"`)
}

func TestMultiSelect_SelectedRenderChipsWithHiddenInputsAndLabels(t *testing.T) {
	// Selected keys render as chips: description shown, key submitted
	// via a repeated hidden input under the same name.
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:     "estado",
		Selected: []string{"otro", "oax"},
		Items: []input.SelectOption{
			{Value: "otro", Label: "Otro"},
			{Value: "oax", Label: "Oaxaca"},
			{Value: "pue", Label: "Puebla"},
		},
	})
	require.Contains(t, out, `<input type="hidden" name="estado" value="otro"`,
		"selected key submits via a hidden input under the field name")
	require.Contains(t, out, `<input type="hidden" name="estado" value="oax"`,
		"a second selection repeats the same name (PostFormArray on the server)")
	// Chip shows the human label, not the raw key.
	require.Contains(t, out, `data-kiban-multiselect-chip-label`)
	require.Contains(t, out, `>Otro</span>`)
	require.Contains(t, out, `>Oaxaca</span>`)
}

func TestMultiSelect_SelectedOptionsAreHiddenInDropdown(t *testing.T) {
	// An already-selected option is hidden in the dropdown so it can't
	// be picked twice.
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:     "x",
		Selected: []string{"otro"},
		Items: []input.SelectOption{
			{Value: "otro", Label: "Otro"},
			{Value: "oax", Label: "Oaxaca"},
		},
	})
	// Find the dropdown <li> for "otro" (not the chip <span>) and
	// assert it carries hidden.
	var otroOption string
	for _, seg := range strings.Split(out, "<li") {
		if strings.Contains(seg, `data-kiban-multiselect-option`) && strings.Contains(seg, `data-value="otro"`) {
			otroOption = seg
			break
		}
	}
	require.NotEmpty(t, otroOption, "dropdown must render an option li for the selected value")
	require.Contains(t, otroOption, "hidden", "already-selected option is hidden in the dropdown")
}

func TestMultiSelect_LabelForUnknownSelectedFallsBackToValue(t *testing.T) {
	// A selected key with no matching option (free-add case) shows the
	// raw key as its chip label.
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:     "x",
		Selected: []string{"CUSTOM_VAL"},
		Items:    []input.SelectOption{{Value: "a", Label: "Alpha"}},
	})
	require.Contains(t, out, `value="CUSTOM_VAL"`)
	require.Contains(t, out, `>CUSTOM_VAL</span>`)
}

func TestMultiSelect_EmptyNameRendersUnnamedHiddenInputs(t *testing.T) {
	// Inactive-pane idiom: empty Name → hidden inputs have no name and
	// the browser skips them on submit.
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:     "",
		Selected: []string{"otro"},
		Items:    []input.SelectOption{{Value: "otro", Label: "Otro"}},
	})
	require.Contains(t, out, `data-name=""`)
	require.NotContains(t, out, `name="otro"`, "value must not leak into the name attribute")
}

func TestMultiSelect_DoesNotEmitCreateRowWhenAllowCreateFalse(t *testing.T) {
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:        "x",
		AllowCreate: false,
	})
	require.NotContains(t, out, `data-kiban-multiselect-create`)
	require.Contains(t, out, `data-allow-create="false"`)
}

func TestMultiSelect_RendersErrorMessageAndStyling(t *testing.T) {
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:   "x",
		ErrMsg: "Selecciona al menos un valor.",
	})
	require.Contains(t, out, "Selecciona al menos un valor.")
	require.Contains(t, out, "border-red-400")
}

func TestMultiSelect_RendersHintWhenNoError(t *testing.T) {
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name: "x",
		Hint: "Puedes elegir una o varias opciones.",
	})
	require.Contains(t, out, "Puedes elegir una o varias opciones.")
}

func TestMultiSelect_RequiredAsterisk(t *testing.T) {
	out := renderMultiSelect(t, input.MultiSelectOptions{
		Name:     "x",
		Label:    "Estado",
		Required: true,
	})
	idx := strings.Index(out, ">Estado")
	require.GreaterOrEqual(t, idx, 0)
	require.Contains(t, out[idx:idx+200], "text-red-500")
}
