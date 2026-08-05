package button_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/button"
	"github.com/stretchr/testify/assert"
)

func renderSplit(t *testing.T, cfg button.SplitOptions) string {
	t.Helper()
	var buf bytes.Buffer
	err := button.Split(cfg).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func saveSplit() button.SplitOptions {
	return button.SplitOptions{
		ID:      "save-split",
		Primary: button.Options{Label: "Guardar", IsSubmit: true, Form: "detail-form"},
		Items: []button.SplitItem{
			{Label: "Guardar y cerrar", Name: "_action", Value: "save_and_close", Form: "detail-form"},
			{Label: "Guardar y nuevo", Name: "_action", Value: "save_and_new", Form: "detail-form"},
		},
	}
}

// Los items mandan el form nativamente: el server distingue la variante
// por el name/value posteado, sin JS de por medio.
func TestSplit_ItemsSubmitTheFormWithTheirAction(t *testing.T) {
	body := renderSplit(t, saveSplit())

	assert.Contains(t, body, `type="submit"`)
	assert.Contains(t, body, `name="_action"`)
	assert.Contains(t, body, `value="save_and_close"`)
	assert.Contains(t, body, `value="save_and_new"`)
	assert.Contains(t, body, `form="detail-form"`)
}

// Reusa el contrato del menu kebab (mismo JS compartido en base.templ):
// wrapper data-kiban-menu + trigger/panel keyed por ID.
func TestSplit_ReusesSharedMenuContract(t *testing.T) {
	body := renderSplit(t, saveSplit())

	assert.Contains(t, body, `data-kiban-menu="save-split"`)
	assert.Contains(t, body, `id="save-split-trigger"`)
	assert.Contains(t, body, `id="save-split-panel"`)
	// templ escapa las comillas simples dentro del atributo.
	assert.Contains(t, body, "window.kibanToggleMenu(&#39;save-split&#39;)")
	assert.Contains(t, body, "window.kibanCloseMenu(&#39;save-split&#39;)")
}

// El OnClick del consumidor corre ANTES del cierre: devolver false ahí
// es lo que permite abortar el submit con un confirm().
func TestSplit_UserOnClickRunsBeforeClose(t *testing.T) {
	cfg := saveSplit()
	cfg.Items[0].OnClick = "return confirm('¿Seguro?')"
	body := renderSplit(t, cfg)

	assert.Contains(t, body, "return confirm(&#39;¿Seguro?&#39;);window.kibanCloseMenu(&#39;save-split&#39;)")
}

// Sin items degrada a un botón normal, así el call site puede gatear las
// variantes por permisos sin ramificar.
func TestSplit_FallsBackToPlainButtonWithoutItems(t *testing.T) {
	body := renderSplit(t, button.SplitOptions{
		ID:      "save-split",
		Primary: button.Options{Label: "Guardar", IsSubmit: true},
	})

	assert.Contains(t, body, "Guardar")
	assert.NotContains(t, body, "data-kiban-menu")
	assert.NotContains(t, body, `role="menu"`)
}

// El caret hereda el fill del primario y pierde el redondeo del lado que
// se une, para que los dos segmentos lean como un solo control.
func TestSplit_SegmentsJoinVisually(t *testing.T) {
	body := renderSplit(t, saveSplit())

	assert.Contains(t, body, "rounded-r-none")
	assert.Contains(t, body, "rounded-l-none")
}

// Sin TriggerAriaLabel explícito, el caret se etiqueta con el primario.
func TestSplit_TriggerFallsBackToPrimaryLabel(t *testing.T) {
	body := renderSplit(t, saveSplit())

	assert.Contains(t, body, `aria-label="Guardar"`)
}
