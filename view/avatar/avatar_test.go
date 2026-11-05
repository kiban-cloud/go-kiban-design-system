package avatar_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/avatar"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

func renderAvatar(t *testing.T, ctx context.Context, o avatar.Options) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, avatar.Avatar(o).Render(ctx, &buf))
	return buf.String()
}

func TestAvatar_PorDefectoSigueCircular(t *testing.T) {
	body := renderAvatar(t, context.Background(), avatar.Options{Name: "Ana Pérez"})

	assert.Contains(t, body, "rounded-full")
	assert.Contains(t, body, "h-8 w-8 text-sm")
	assert.Contains(t, body, "bg-kiban-primary-soft text-kiban-primary")
	assert.Contains(t, body, ">AP")
}

func TestAvatar_CuadradoConTono(t *testing.T) {
	body := renderAvatar(t, context.Background(), avatar.Options{Name: "Ana Pérez", Shape: "square", Tone: "success"})

	assert.Contains(t, body, "rounded-lg")
	assert.Contains(t, body, "bg-emerald-50 text-emerald-700")
	assert.NotContains(t, body, "rounded-full")
}

func TestAvatar_FlamboyantCambiaTonos(t *testing.T) {
	ctx := feature.Enable(context.Background(), flamboyant.Flag)

	pf := renderAvatar(t, ctx, avatar.Options{Name: "Ana Pérez", Shape: "square", Tone: "success"})
	assert.Contains(t, pf, "bg-[#E3FBF5] text-[#007A63]")
	assert.Contains(t, pf, "font-semibold")
	assert.Contains(t, pf, "h-8 w-8 text-[11px]")

	pm := renderAvatar(t, ctx, avatar.Options{Name: "Acme SA", Shape: "square", Size: "lg"})
	assert.Contains(t, pm, "bg-[#EAF0FF] text-[#0036C4]")
	assert.Contains(t, pm, "h-9 w-9")
}
