package drawer_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/button"
	"github.com/kiban-cloud/go-kiban-design-system/view/drawer"
)

func renderPanel(t *testing.T, cfg drawer.FormPanelConfig) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, drawer.FormPanel(cfg).Render(templ.WithChildren(context.Background(), templ.Raw(`<p>campos</p>`)), &buf))
	return buf.String()
}

func panelCfg() drawer.FormPanelConfig {
	primary := button.Options{Label: "Crear cliente", IsSubmit: true}
	return drawer.FormPanelConfig{
		ID: "p", Half: true, Eyebrow: "Nuevo cliente", TitlePlaceholder: "Sin nombre todavía",
		TitleFrom: []string{"firstName", "lastName"}, CloseLabel: "Cerrar",
		FormAttrs:      templ.Attributes{"method": "POST", "action": "/crm/customers/new"},
		SummaryPending: "Completa los campos obligatorios", SummaryReady: "Listo para crear",
		Actions: button.Group{
			SecondaryActions: []button.Options{{Label: "Cancelar", Attrs: templ.Attributes{"data-kiban-panel-close": ""}}},
			PrimaryAction:    &primary,
		},
		Guard:     &drawer.FormPanelGuard{Title: "Tienes cambios sin guardar", Text: "Se pierde", Discard: "Descartar", Stay: "Seguir", SaveClose: "Crear y cerrar"},
		CloseHref: "/crm/customers",
	}
}

func TestFormPanel_EsElFormulario(t *testing.T) {
	body := renderPanel(t, panelCfg())

	assert.Contains(t, body, `id="p"`)
	assert.Contains(t, body, "data-kiban-panel")
	assert.Contains(t, body, `action="/crm/customers/new" method="POST"`)
	assert.Contains(t, body, "<p>campos</p>")
	assert.Contains(t, body, "w-1/2", "Half es el panel del 50 %")
	assert.Contains(t, body, "animate-[kiban-panel-in")
	assert.Contains(t, body, `data-close-href="/crm/customers"`)
	assert.NotContains(t, body, "data-standalone")
}

func TestFormPanel_TituloEnVivoYPie(t *testing.T) {
	body := renderPanel(t, panelCfg())

	assert.Contains(t, body, `data-from="firstName lastName"`)
	assert.Contains(t, body, ">Sin nombre todavía</div>")
	assert.Contains(t, body, "data-empty")
	assert.Contains(t, body, `data-ready="Listo para crear"`)
	assert.Contains(t, body, ">Completa los campos obligatorios</div>")
	assert.Contains(t, body, `type="submit"`)
	assert.Contains(t, body, ">Cancelar")
}

func TestFormPanel_AvisoOcultoConSusBotones(t *testing.T) {
	body := renderPanel(t, panelCfg())

	assert.Contains(t, body, "data-kiban-panel-guard hidden")
	assert.Contains(t, body, "[&[hidden]]:hidden", "la clase flex no le gana al atributo hidden")
	for _, b := range []string{"data-kiban-guard-discard", "data-kiban-guard-stay", "data-kiban-guard-save"} {
		assert.Contains(t, body, b)
	}
}

func TestFormPanel_SinAvisoStandaloneYEstatico(t *testing.T) {
	cfg := panelCfg()
	cfg.Guard, cfg.Standalone, cfg.Static, cfg.Half, cfg.Title = nil, true, true, false, "Ana Pérez"
	body := renderPanel(t, cfg)

	assert.NotContains(t, body, "data-kiban-panel-guard")
	assert.Contains(t, body, "data-standalone")
	assert.NotContains(t, body, "animate-[kiban-panel-in")
	assert.Contains(t, body, "w-[95vw]")
	assert.Contains(t, body, ">Ana Pérez</div>")
	assert.NotContains(t, body, "data-empty")
}

func TestViewPanel_MismoEncabezadoSinFormNiPie(t *testing.T) {
	var buf bytes.Buffer
	cfg := drawer.ViewPanelConfig{ID: "ficha", Eyebrow: "Expediente", Title: "Ana López", CloseLabel: "Cerrar", CloseHref: "/crm/customers", HeaderActions: templ.Raw(`<a>Editar</a>`)}
	require.NoError(t, drawer.ViewPanel(cfg).Render(templ.WithChildren(context.Background(), templ.Raw(`<p>datos</p>`)), &buf))
	body := buf.String()

	assert.Contains(t, body, `id="ficha"`)
	assert.Contains(t, body, "data-kiban-panel")
	assert.Contains(t, body, ">Expediente</div>")
	assert.Contains(t, body, ">Ana López</div>")
	assert.Contains(t, body, "<a>Editar</a>")
	assert.Contains(t, body, "<p>datos</p>")
	assert.Contains(t, body, "w-[95vw]", "sin Half es el panel del 95 %")
	assert.Contains(t, body, `data-close-href="/crm/customers"`)
	assert.NotContains(t, body, "<form")
	assert.NotContains(t, body, "data-kiban-panel-sum")
}

func TestViewPanel_StandaloneYStatic(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, drawer.ViewPanel(drawer.ViewPanelConfig{ID: "f", Half: true, Standalone: true, Static: true}).Render(context.Background(), &buf))
	body := buf.String()

	assert.Contains(t, body, "data-standalone")
	assert.Contains(t, body, "w-1/2")
	assert.NotContains(t, body, "animate-[kiban-panel-in")
}
