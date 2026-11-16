package drawer_test

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/drawer"
)

// Attrs escribe sólo los data-kiban-confirm-* que vienen; el título siempre.
func TestHtmxConfirm_SoloTitulo(t *testing.T) {
	assert.Equal(t, templ.Attributes{"data-kiban-confirm-title": "¿Borrar la orden?"},
		drawer.HtmxConfirm{Title: "¿Borrar la orden?"}.Attrs(), "el tono por defecto (danger) no escribe atributo")
}

func TestHtmxConfirm_TodoConTonoInfo(t *testing.T) {
	a := drawer.HtmxConfirm{
		Title: "¿Retirar $184,320.55 MXN?", Text: "Se transfiere todo el saldo de SPEI by kiban", Strong: "SPEI by kiban",
		OK: "Retirar $184,320.55", Cancel: "Cancelar", Tone: drawer.ConfirmToneInfo,
	}.Attrs()

	assert.Equal(t, templ.Attributes{
		"data-kiban-confirm-title":  "¿Retirar $184,320.55 MXN?",
		"data-kiban-confirm-text":   "Se transfiere todo el saldo de SPEI by kiban",
		"data-kiban-confirm-strong": "SPEI by kiban",
		"data-kiban-confirm-ok":     "Retirar $184,320.55",
		"data-kiban-confirm-cancel": "Cancelar",
		"data-kiban-confirm-tone":   "info",
	}, a)
}

func TestHtmxConfirm_TonoDangerEsElDefault(t *testing.T) {
	assert.Equal(t, "", drawer.ConfirmToneDanger)
	_, ok := drawer.HtmxConfirm{Title: "x", Tone: drawer.ConfirmToneDanger}.Attrs()["data-kiban-confirm-tone"]
	assert.False(t, ok)
}
