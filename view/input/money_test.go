package input_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

// Por defecto: signo «$», placeholder «0.00», teclado decimal sin
// type=number y la caja grande.
func TestMoney_PorDefecto(t *testing.T) {
	body := renderIn(t, context.Background(), input.Money(input.MoneyOptions{Name: "amount", AriaLabel: "Monto"}))

	assert.Contains(t, body, "data-kiban-money")
	assert.Contains(t, body, `id="f-amount" name="amount"`)
	assert.Contains(t, body, `inputmode="decimal"`)
	assert.NotContains(t, body, `type="number"`)
	assert.Contains(t, body, `placeholder="0.00"`)
	assert.Contains(t, body, `aria-hidden="true">$</span>`)
	assert.Contains(t, body, `aria-label="Monto"`)
	assert.Contains(t, body, "h-14", "la caja grande")
	assert.NotContains(t, body, "<em", "sin moneda no se pinta")
	assert.NotContains(t, body, "data-max", "sin tope")
	assert.NotContains(t, body, " required")
	assert.NotContains(t, body, "aria-invalid")
	assert.NotContains(t, body, `<label for="f-amount" class="text-xs`, "sin Label, sólo la caja")
}

func TestMoney_EtiquetaMonedaRequeridoYCompacto(t *testing.T) {
	body := renderIn(t, context.Background(), input.Money(input.MoneyOptions{
		Name: "amount", Label: "Monto", Value: "1,500.50", Prefix: "US$", Currency: "MXN", Placeholder: "100",
		Required: true, Compact: true, Attrs: templ.Attributes{"hx-post": "/preview"},
	}))

	assert.Contains(t, body, ">Monto")
	assert.Contains(t, body, `<span class="ml-0.5 text-red-500">*</span>`)
	assert.Contains(t, body, `value="1,500.50"`)
	assert.Contains(t, body, ">US$</span>")
	assert.Contains(t, body, ">MXN</em>")
	assert.Contains(t, body, `placeholder="100"`)
	assert.Contains(t, body, " required")
	assert.Contains(t, body, "h-9", "la caja compacta")
	assert.NotContains(t, body, "h-14")
	assert.Contains(t, body, `hx-post="/preview"`)
	assert.NotContains(t, body, "aria-label=", "con Label no hace falta aria-label")
}

// Max pone el tope que revisa el runtime (data-max y su mensaje).
func TestMoney_Tope(t *testing.T) {
	body := renderIn(t, context.Background(), input.Money(input.MoneyOptions{Name: "amount", Max: "10000", MaxError: "Hasta $10,000 con OXXO Pay"}))

	assert.Contains(t, body, `data-max="10000"`)
	assert.Contains(t, body, `data-max-error="Hasta $10,000 con OXXO Pay"`)
}

// Error gana sobre Hint; Below reemplaza a los dos.
func TestMoney_ErrorAyudaYBelow(t *testing.T) {
	ayuda := renderIn(t, context.Background(), input.Money(input.MoneyOptions{Name: "a", Hint: "Sin comisión"}))
	assert.Contains(t, ayuda, ">Sin comisión</p>")
	assert.Contains(t, ayuda, "border-[#CDD3DF]")

	err := renderIn(t, context.Background(), input.Money(input.MoneyOptions{Name: "a", Hint: "Sin comisión", Error: "Monto inválido"}))
	assert.Contains(t, err, ">Monto inválido</p>")
	assert.NotContains(t, err, "Sin comisión")
	assert.Contains(t, err, `aria-invalid="true"`)
	assert.Contains(t, err, "border-[#F3A6A6]")

	below := renderIn(t, context.Background(), input.Money(input.MoneyOptions{Name: "a", Hint: "h", Error: "e", Below: templ.Raw(`<p id="amount-help">oob</p>`)}))
	assert.Contains(t, below, `<p id="amount-help">oob</p>`)
	assert.NotContains(t, below, ">e</p>")
	assert.NotContains(t, below, ">h</p>")
}
