package drawer

import "github.com/a-h/templ"

// Tonos de HtmxConfirm.
const (
	// ConfirmToneDanger es el de siempre: bote de basura y botón rojo, para
	// borrar o cancelar algo.
	ConfirmToneDanger = ""
	// ConfirmToneInfo es para una acción que no destruye nada pero conviene
	// confirmar (retirar saldo, generar llaves): ícono de aviso, botón azul
	// y el foco en confirmar.
	ConfirmToneInfo = "info"
)

// HtmxConfirm es la confirmación flamboyant de un botón con hx-post (o
// cualquier petición htmx): en vez del confirm() del navegador, el runtime
// del design system (base.templ, evento htmx:confirm) abre un diálogo chico
// centrado con el título, el texto y los botones, y sólo al confirmar manda
// la petición. Es el mismo diálogo que arman a mano los data-kiban-confirm-*;
// Attrs los escribe por ti.
//
// Para listas («¿Borrar 3 órdenes?») siguen valiendo los atributos
// data-kiban-confirm-list / -item / -title-one / -ok-one / -more, que se
// pueden sumar al mapa que regresa Attrs.
type HtmxConfirm struct {
	// Title es la pregunta («¿Retirar $184,320.55 MXN?»). Obligatorio: sin
	// título no hay diálogo.
	Title string
	// Text es el párrafo debajo, en texto plano.
	Text string
	// Strong es un pedazo de Text que va en negritas (el nombre del medio).
	// Si no aparece en Text, se ignora.
	Strong string
	// OK es la etiqueta del botón que confirma («Retirar $184,320.55»).
	OK string
	// Cancel es la del botón que cierra sin hacer nada; vacío, «Cancel».
	Cancel string
	// Tone es ConfirmToneDanger (por defecto) o ConfirmToneInfo.
	Tone string
}

// Attrs son los data-kiban-confirm-* para ponerle al botón, junto a sus
// hx-post / hx-target / hx-swap.
func (c HtmxConfirm) Attrs() templ.Attributes {
	a := templ.Attributes{"data-kiban-confirm-title": c.Title}
	if c.Text != "" {
		a["data-kiban-confirm-text"] = c.Text
	}
	if c.Strong != "" {
		a["data-kiban-confirm-strong"] = c.Strong
	}
	if c.OK != "" {
		a["data-kiban-confirm-ok"] = c.OK
	}
	if c.Cancel != "" {
		a["data-kiban-confirm-cancel"] = c.Cancel
	}
	if c.Tone != "" {
		a["data-kiban-confirm-tone"] = c.Tone
	}
	return a
}
