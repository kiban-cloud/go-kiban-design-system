package input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

// Piezas portadas de kiban-proto: apagado flamboyant, lo de siempre;
// prendido, lo del proto.

var encendido = feature.Enable(context.Background(), flamboyant.Flag)

func pintar(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

var _ = strings.Contains

func TestPicker_DenseYTono(t *testing.T) {
	o := input.PickerOptions{Name: "r", Options: []input.PickerOption{{Value: "ok", Title: "Exitoso", Tone: "success"}, {Value: "x", Title: "Sin tono"}}}
	normal := pintar(t, context.Background(), input.Picker(o))
	assert.Contains(t, normal, "relative grid max-h-[360px] gap-2 overflow-auto p-0.5")
	assert.Contains(t, normal, `data-tone="success"`)
	assert.Contains(t, normal, "bg-[#00A884]", "el punto verde")
	assert.Equal(t, 1, strings.Count(normal, "data-tone="), "sin Tone, sin data-tone")

	o.Dense = true
	dense := pintar(t, context.Background(), input.Picker(o))
	assert.Contains(t, dense, "grid-cols-[repeat(auto-fill,minmax(200px,1fr))]")
	assert.NotContains(t, dense, "max-h-[360px]", "Dense no tiene tope de alto")
	assert.Contains(t, dense, "text-[13px] font-normal")
}
