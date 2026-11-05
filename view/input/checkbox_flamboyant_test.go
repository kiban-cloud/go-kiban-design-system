package input_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func renderCheckboxIn(t *testing.T, ctx context.Context, checked bool) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, input.Checkbox("sel", "Seleccionar", checked, true, nil, "").Render(ctx, &buf))
	return buf.String()
}

func TestCheckbox_SinBanderaEsElDeSiempre(t *testing.T) {
	body := renderCheckboxIn(t, context.Background(), false)

	assert.Contains(t, body, `class="cursor-pointer disabled:cursor-not-allowed"`)
	assert.NotContains(t, body, "appearance-none")
}

func TestCheckbox_FlamboyantEsCuadradoConPalomitaYParcial(t *testing.T) {
	body := renderCheckboxIn(t, feature.Enable(context.Background(), flamboyant.Flag), true)

	assert.Contains(t, body, "appearance-none")
	assert.Contains(t, body, "rounded-[5px]")
	assert.Contains(t, body, "checked:bg-kiban-primary")
	assert.Contains(t, body, "indeterminate:bg-kiban-primary")
	assert.Contains(t, body, "peer-checked:opacity-100")
	assert.Contains(t, body, "peer-indeterminate:opacity-100")
	assert.Contains(t, body, ` checked`)
	assert.Contains(t, body, `value="true"`)
	assert.Contains(t, body, "Seleccionar")
}
