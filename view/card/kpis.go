package card

import "github.com/a-h/templ"

// KpiOptions es una cifra de Kpis: la etiqueta en versalitas con su ícono,
// el valor grande con su unidad en chico («$11,109.16 MXN»), una línea de
// detalle debajo («3 órdenes pendientes · 2 vencen en 48 h») y Featured, la
// tarjeta destacada (fondo azul kiban): lo primero que hay que ver.
type KpiOptions struct {
	Label    string
	Icon     templ.Component
	Value    string
	Unit     string
	Detail   string
	Featured bool
}
