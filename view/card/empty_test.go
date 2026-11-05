package card_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

func TestEmpty_TituloTextoIconoYAcciones(t *testing.T) {
	body := render(t, context.Background(), card.Empty(card.EmptyOptions{Icon: templ.Raw(`<i>ico</i>`), Title: "Sin pagos", Text: "La integración llega pronto"}))

	assert.Contains(t, body, "<i>ico</i>")
	assert.Contains(t, body, ">Sin pagos</div>")
	assert.Contains(t, body, "La integración llega pronto")
	assert.Contains(t, body, "<p>cuerpo</p>", "children son las acciones")
}

func TestEmpty_SinTextoNiIcono(t *testing.T) {
	body := render(t, context.Background(), card.Empty(card.EmptyOptions{Title: "Nada"}))

	assert.Contains(t, body, ">Nada</div>")
	assert.NotContains(t, body, "max-w-md")
}
