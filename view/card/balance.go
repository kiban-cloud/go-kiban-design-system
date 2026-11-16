package card

import "github.com/a-h/templ"

// BalanceOptions es un cuadro de saldo (Balance): de quién es el saldo
// (placa y título), cuánto hay y qué se puede hacer con él.
type BalanceOptions struct {
	// Key identifica el cuadro (data-kiban-balance); no se muestra. Sirve
	// para encontrarlo en pruebas («bancoDonde», «stp»).
	Key string
	// Tile es la placa cuadrada a la izquierda del título (32 px); Title el
	// nombre de la cuenta o del medio y Sub la línea debajo («Cuenta 0035»,
	// «3 centros de costos»).
	Tile  templ.Component
	Title string
	Sub   string
	// Amount es el saldo ya formateado («$1,284,903.18») y Currency la
	// moneda en chico al lado («MXN»).
	Amount   string
	Currency string
	// Zero pinta el saldo en gris: no hay nada que mover.
	Zero bool
	// Live es un saldo que se consulta en vivo: el cuadro va en azul y
	// lleva la píldora LiveLabel («En vivo»).
	Live      bool
	LiveLabel string
	// Foot y Action son el pie: el texto a la izquierda («Saldo para
	// dispersar») y el botón o enlace a la derecha (BalanceLink, un botón).
	// Con Lines no hay pie: cada renglón trae su acción.
	Foot   string
	Action templ.Component
	// Lines son los saldos que agrupa el cuadro (los centros de costos de
	// una cuenta), uno por renglón con su acción.
	Lines []BalanceLineOptions
}

// BalanceLineOptions es un renglón de un cuadro de saldo (BalanceLine): el
// nombre, su saldo ya formateado y la acción a la derecha. Currency vacío
// toma la del cuadro.
type BalanceLineOptions struct {
	Name     string
	Amount   string
	Currency string
	Action   templ.Component
}

// line es el renglón con la moneda del cuadro si no trae la suya.
func (o BalanceOptions) line(l BalanceLineOptions) BalanceLineOptions {
	if l.Currency == "" {
		l.Currency = o.Currency
	}
	return l
}
