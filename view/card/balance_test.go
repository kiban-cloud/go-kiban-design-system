package card_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/card"
)

// Un cuadro con pie: texto a la izquierda y la acción a la derecha.
func TestBalance_ConPieYAccion(t *testing.T) {
	body := render(t, context.Background(), card.Balance(card.BalanceOptions{
		Key: "bancoDonde", Tile: templ.Raw(`<i>bd</i>`), Title: "Banco Donde", Sub: "Cuenta 0035",
		Amount: "$1,284,903.18", Currency: "MXN", Foot: "Saldo para dispersar",
		Action: card.BalanceLink("Ir a Dispersión", "/rekon/dispersion", nil),
	}))

	assert.Contains(t, body, `data-kiban-balance="bancoDonde"`)
	assert.Contains(t, body, "<i>bd</i>")
	assert.Contains(t, body, ">Banco Donde</div>")
	assert.Contains(t, body, ">Cuenta 0035</div>")
	assert.Contains(t, body, "$1,284,903.18")
	assert.Contains(t, body, ">MXN</small>")
	assert.Contains(t, body, ">Saldo para dispersar</span>")
	assert.Contains(t, body, `href="/rekon/dispersion"`)
	assert.Contains(t, body, ">Ir a Dispersión</a>")
	assert.Contains(t, body, "border-[#E2E5EC] bg-white", "sin Live, el cuadro blanco")
	assert.NotContains(t, body, "data-kiban-balance-line")
	assert.NotContains(t, body, "text-[#9096A3]\">$", "con saldo, la cifra no va gris")
}

// En vivo: cuadro azul con la píldora; Zero pinta el saldo gris.
func TestBalance_EnVivoYSaldoCero(t *testing.T) {
	body := render(t, context.Background(), card.Balance(card.BalanceOptions{
		Key: "stp", Title: "STP", Amount: "$0.00", Currency: "MXN", Live: true, LiveLabel: "En vivo", Zero: true,
	}))

	assert.Contains(t, body, "border-[#DCE5FF] bg-[#F4F7FF]")
	assert.Contains(t, body, "En vivo")
	assert.Contains(t, body, "text-[22px] font-semibold leading-7 tabular-nums text-[#9096A3]")
	assert.NotContains(t, body, "truncate text-xs leading-4", "sin Sub no hay línea debajo del título")
}

// Con renglones no hay pie: cada saldo trae su acción y, sin moneda propia,
// toma la del cuadro.
func TestBalance_RenglonesHeredanLaMoneda(t *testing.T) {
	body := render(t, context.Background(), card.Balance(card.BalanceOptions{
		Key: "stp", Title: "STP", Amount: "$10.00", Currency: "MXN", Foot: "no se ve",
		Lines: []card.BalanceLineOptions{
			{Name: "Centro A", Amount: "$4.00", Action: templ.Raw(`<button>Retirar A</button>`)},
			{Name: "Centro B", Amount: "$6.00", Currency: "USD"},
		},
	}))

	assert.Equal(t, 2, strings.Count(body, "data-kiban-balance-line"))
	assert.Contains(t, body, ">Centro A</b>")
	assert.Contains(t, body, "$4.00 MXN</small>")
	assert.Contains(t, body, "$6.00 USD</small>")
	assert.Contains(t, body, "<button>Retirar A</button>")
	assert.NotContains(t, body, "no se ve", "con Lines no se pinta el pie")
}

func TestBalanceGrid_AnchoYAngosto(t *testing.T) {
	ancho := render(t, context.Background(), card.BalanceGrid(false))
	assert.Contains(t, ancho, "data-kiban-balance-grid")
	assert.Contains(t, ancho, "grid-cols-[repeat(auto-fit,minmax(280px,1fr))]")
	assert.Contains(t, ancho, "<p>cuerpo</p>")

	angosto := render(t, context.Background(), card.BalanceGrid(true))
	assert.Contains(t, angosto, "grid-cols-[minmax(0,380px)]")
	assert.NotContains(t, angosto, "auto-fit")
}

func TestBalanceLink_AtributosHtmx(t *testing.T) {
	body := render(t, context.Background(), card.BalanceLink("Nueva dispersión", "#", templ.Attributes{"hx-get": "/rekon/dispersion/new"}))
	assert.Contains(t, body, `href="#"`)
	assert.Contains(t, body, `hx-get="/rekon/dispersion/new"`)
	assert.Contains(t, body, ">Nueva dispersión</a>")
}
