package button_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/button"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

func renderIn(t *testing.T, ctx context.Context, opts button.Options) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, button.Button(opts).Render(ctx, &buf))
	return buf.String()
}

var flam = feature.Enable(context.Background(), flamboyant.Flag)

func TestButton_SinBanderaUsaLasClasesDeSiempre(t *testing.T) {
	for _, v := range []string{"primary", "secondary", "danger", "icon"} {
		body := renderIn(t, context.Background(), button.Options{Label: "x", Variant: v})
		assert.Contains(t, body, `class="`+button.BuildClass(v, "")+`"`, v)
	}
}

func TestButton_FlamboyantCambiaElEstilo(t *testing.T) {
	cases := map[string]string{
		"primary":        "bg-kiban-primary",
		"secondary":      "border-[#CDD3DF]",
		"danger":         "bg-[#D92D20]",
		"danger-outline": "border-[#F3A6A6]",
		"icon":           "h-[34px]",
	}
	for v, marca := range cases {
		body := renderIn(t, flam, button.Options{Label: "x", Variant: v})
		assert.Contains(t, body, marca, v)
		if v != "icon" {
			assert.Contains(t, body, "text-[13px]", v)
		}
		assert.NotContains(t, body, button.BuildClass(v, ""), v)
	}
}

func TestButton_FlamboyantConservaExtraClassYEnlace(t *testing.T) {
	body := renderIn(t, flam, button.Options{Label: "Editar", Variant: "secondary", Href: "/crm/customers/1", ExtraClass: "w-full"})

	assert.Contains(t, body, `<a href="/crm/customers/1"`)
	assert.Contains(t, body, "w-full")
}

func TestButton_DangerOutlineExisteEnLosDosDisenos(t *testing.T) {
	assert.Equal(t, "danger-outline", button.Options{Variant: "danger-outline"}.EffectiveVariant())
	body := renderIn(t, context.Background(), button.Options{Label: "Descartar", Variant: "danger-outline"})
	assert.Contains(t, body, "border-red-300")
}
