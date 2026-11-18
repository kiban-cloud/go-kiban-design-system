package layout

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
)

func renderAdmin(t *testing.T, ctx context.Context) string {
	t.Helper()
	var buf bytes.Buffer
	cfg := AdminConfig{Title: "Entregas", ProjectName: "klin internal", LogoutAction: "/salir", User: User{Name: "Ana", Email: "ana@kiban.com"}}
	require.NoError(t, AdminLayout(cfg).Render(ctx, &buf))
	return buf.String()
}

func TestAdminLayout_SwitchDeDiseno(t *testing.T) {
	t.Run("sin la bandera del switch no aparece", func(t *testing.T) {
		body := renderAdmin(t, context.Background())
		assert.NotContains(t, body, `data-kiban-ui="`)
	})
	t.Run("con el switch y el diseño de siempre lleva al nuevo", func(t *testing.T) {
		body := renderAdmin(t, feature.Enable(context.Background(), FlagShellFlamboyantSwitch))
		assert.Contains(t, body, `data-kiban-ui="flamboyant"`)
		assert.Contains(t, body, "Probar el nuevo diseño")
	})
	t.Run("con el diseño nuevo regresa al de siempre", func(t *testing.T) {
		body := renderAdmin(t, feature.Enable(context.Background(), FlagShellFlamboyantSwitch, FlagShellFlamboyant))
		assert.Contains(t, body, `data-kiban-ui="clasico"`)
		assert.Contains(t, body, "Volver al diseño actual")
	})
}
