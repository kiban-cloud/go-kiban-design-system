package layout_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

func renderComp(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

// Los íconos nuevos del resumen de rekon existen; uno desconocido no pinta.
func TestFlamboyantIcon_NuevosDelResumen(t *testing.T) {
	for name, trazo := range map[string]string{
		"clock":         `<polyline points="12 6 12 12 16 14"/>`,
		"download":      `<polyline points="7 10 12 15 17 10"/>`,
		"send":          `<polygon points="22 2 15 22 11 13 2 9 22 2"/>`,
		"chevron-right": `<polyline points="9 18 15 12 9 6"/>`,
	} {
		body := renderComp(t, layout.FlamboyantIcon(name))
		assert.Contains(t, body, "<svg", name)
		assert.Contains(t, body, trazo, name)
		assert.Contains(t, body, "</svg>", name)
	}
	assert.Empty(t, renderComp(t, layout.FlamboyantIcon("no-existe")))
}

// IconTile: el ícono en el cuadro azul claro, del tamaño pedido.
func TestIconTile_CuadroConIcono(t *testing.T) {
	body := renderComp(t, layout.IconTile("send", 40))

	assert.Contains(t, body, `style="width:40px;height:40px;"`)
	assert.Contains(t, body, "bg-[#EAF0FF] text-[#0036C4]")
	assert.Contains(t, body, `aria-hidden="true"`)
	assert.Contains(t, body, `<polygon points="22 2 15 22 11 13 2 9 22 2"/>`)
}

// InboxToggleAttrs abre el panel «Pendientes»: sin base usa la de siempre.
func TestInboxToggleAttrs_BasePorDefectoYPropia(t *testing.T) {
	assert.Equal(t, templ.Attributes{
		"data-fs-inbox-toggle": "",
		"aria-controls":        "fs-inbox",
		"hx-get":               layout.DefaultNotificationsBaseURL + "/dropdown",
		"hx-target":            "#fs-inbox-list",
		"hx-swap":              "innerHTML",
		"hx-trigger":           "click",
	}, layout.InboxToggleAttrs(""))
	assert.Equal(t, "/kiban-cloud/notifications/dropdown", layout.InboxToggleAttrs("")["hx-get"])

	assert.Equal(t, "/crm/notifications/dropdown", layout.InboxToggleAttrs("/crm/notifications")["hx-get"])
}

// Los íconos de workfloo flamboyant (kiban-proto) existen.
func TestFlamboyantIcon_NuevosDeWorkfloo(t *testing.T) {
	for _, name := range []string{"play", "circle-check", "circle-x", "workflow", "git-branch", "plug", "history", "plus", "copy", "eye", "eye-off", "trash-2", "zap", "split", "trophy", "external-link", "git-fork", "list-checks", "grip-vertical"} {
		assert.Contains(t, renderComp(t, layout.FlamboyantIcon(name)), "<svg", name)
	}
}
