package drawer

import (
	"strings"

	"github.com/a-h/templ"

	"github.com/kiban-cloud/go-kiban-design-system/view/button"
)

// FormPanelConfig describe el panel de formulario de kiban cloud flamboyant:
// un panel lateral que ES el <form>, con encabezado (ícono, antetítulo,
// título en vivo), cuerpo en tarjetas (card.Step) y pie con el estado y los
// botones. Se pinta abierto: lo trae htmx al hacer clic («Nuevo cliente») o
// la página entera cuando se entra por su URL (Standalone).
//
// El comportamiento vive en base.templ ([data-kiban-panel]): numerar y
// completar los pasos, el título en vivo, el texto del pie, el aviso de
// cambios sin guardar y cerrar (botón, fondo o Escape).
type FormPanelConfig struct {
	ID string
	// Half es el panel del 50 % (alta, carga masiva); sin Half, 95 %.
	Half bool

	Tile    templ.Component // ícono del encabezado, ya en su cuadro de color
	Eyebrow string          // «Nuevo cliente»
	// Title es el título inicial; vacío muestra TitlePlaceholder. TitleFrom
	// son los name de los campos que lo arman en vivo, unidos con espacio
	// (los que no estén en el form se ignoran).
	Title            string
	TitlePlaceholder string
	TitleFrom        []string
	Subtitle         templ.Component
	HeaderActions    templ.Component
	CloseLabel       string

	// FormAttrs van en el <form>: method, action, hx-post, hx-target…
	FormAttrs templ.Attributes

	FooterLeft templ.Component
	// SummaryPending y SummaryReady son el texto del pie según el form pase
	// o no la validación del navegador (required, type=email…).
	SummaryPending string
	SummaryReady   string
	Actions        button.Group

	// Guard pide confirmación al cerrar con cambios sin guardar. nil: cierra
	// sin preguntar.
	Guard *FormPanelGuard

	// CloseHref es la pantalla de atrás. Standalone: el panel llegó con la
	// página entera (se entró por su URL) y cerrar navega a CloseHref; si
	// no, cerrar sólo quita el panel.
	CloseHref  string
	Standalone bool
	// Static quita la animación de entrada: para volver a pintar el panel
	// con errores sin que parezca que se abre de nuevo.
	Static bool
}

type FormPanelGuard struct {
	Title, Text              string
	Discard, Stay, SaveClose string
}

func (c FormPanelConfig) titleFrom() string { return strings.Join(c.TitleFrom, " ") }

func (c FormPanelConfig) titleText() string {
	if c.Title != "" {
		return c.Title
	}
	return c.TitlePlaceholder
}
