package badge_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/badge"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

func renderBadge(t *testing.T, ctx context.Context, variant, size string) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, badge.Variant("Pagado", variant, size).Render(ctx, &buf))
	return buf.String()
}

func TestVariant_SinBanderaEsElRectanguloDeSiempre(t *testing.T) {
	body := renderBadge(t, context.Background(), "success", "")

	assert.Contains(t, body, "rounded-md")
	assert.Contains(t, body, "bg-emerald-50 text-emerald-700")
	assert.NotContains(t, body, "rounded-full")
}

func TestVariant_FlamboyantEsPildoraConPunto(t *testing.T) {
	ctx := feature.Enable(context.Background(), flamboyant.Flag)
	cases := map[string]string{
		"success": "text-[#007A63]",
		"warning": "text-[#9A4F00]",
		"danger":  "text-[#B42318]",
		"info":    "text-[#0036C4]",
		"":        "text-[#474F66]",
	}
	for v, color := range cases {
		body := renderBadge(t, ctx, v, "")
		assert.Contains(t, body, "rounded-full", v)
		assert.Contains(t, body, "before:bg-current", v)
		assert.Contains(t, body, color, v)
		assert.Contains(t, body, ">Pagado<", v)
	}
	assert.Contains(t, renderBadge(t, ctx, "info", "md"), "text-sm")
}

func TestStatus_SigueElDisenoDelRequest(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, badge.Status("PAID", "Pagado", "").Render(feature.Enable(context.Background(), flamboyant.Flag), &buf))
	assert.Contains(t, buf.String(), "rounded-full")
}
