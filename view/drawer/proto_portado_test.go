package drawer_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/drawer"
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

func TestFormPanel_EncabezadoSeEnvuelveEnMovilSoloEnFlamboyant(t *testing.T) {
	cfg := drawer.FormPanelConfig{ID: "p", Title: "T"}
	assert.NotContains(t, pintar(t, context.Background(), drawer.FormPanel(cfg)), "max-md:flex-wrap")
	assert.Contains(t, pintar(t, encendido, drawer.FormPanel(cfg)), "max-md:flex-wrap")
}

func TestFormPanel_Pestanas(t *testing.T) {
	cfg := drawer.FormPanelConfig{ID: "p", Title: "T"}
	assert.NotContains(t, pintar(t, encendido, drawer.FormPanel(cfg)), `role="tablist"`, "sin Tabs, sin pestañas")

	cfg.Tabs = []drawer.PanelTab{{Key: "espera", Label: "Espera"}, {Key: "reglas", Label: "Reglas", Count: "4", Badge: "!"}}
	body := pintar(t, encendido, drawer.FormPanel(cfg))
	assert.Contains(t, body, `role="tablist"`)
	assert.Regexp(t, `aria-selected="true" data-kiban-panel-tab="espera"`, body, "la primera va elegida")
	assert.Regexp(t, `aria-selected="false" data-kiban-panel-tab="reglas"`, body)
	assert.Contains(t, body, `data-kiban-panel-tab-count>4</span>`)
	assert.Contains(t, body, "flex-1 justify-end", "con pestañas, las acciones van a la derecha")
}
