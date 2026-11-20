package canvas_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/canvas"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
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

func TestCanvas_CompactoSoloEnFlamboyant(t *testing.T) {
	c := canvas.Canvas(canvas.CanvasOptions{ID: "c"})
	assert.NotContains(t, pintar(t, context.Background(), c), "data-kiban-canvas-compact")
	assert.Contains(t, pintar(t, encendido, c), "data-kiban-canvas-compact")
}

func TestNode_TarjetaFlamboyant(t *testing.T) {
	n := canvas.NodeOptions{ID: "n1", Title: "Espera", Subtitle: "5 min", Status: canvas.StatusNotConfigured}
	apagado := pintar(t, context.Background(), canvas.Node(n))
	assert.NotContains(t, apagado, "Falta configurar", "apagado, la tarjeta de siempre")
	assert.NotContains(t, apagado, "border-[#B98CFF]")

	body := pintar(t, encendido, canvas.Node(n))
	assert.Contains(t, body, "rounded-lg px-3 py-2.5")
	assert.Contains(t, body, "border-dashed border-[#B98CFF]")
	assert.Contains(t, body, "Falta configurar")

	n.Outcome, n.OutcomeLabel, n.Status = canvas.OutcomeSuccess, "Resultado", ""
	body = pintar(t, encendido, canvas.Node(n))
	assert.Contains(t, body, "rounded-full px-3.5", "píldora de resultado")
	assert.Contains(t, body, "bg-[#12B76A]")
	assert.Contains(t, body, "Resultado")
}
