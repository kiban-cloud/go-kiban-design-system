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

// panelHead es lo que pinta panelHeader, común a FormPanel y ViewPanel.
type panelHead struct {
	ID, Eyebrow, Title, TitlePlaceholder, CloseLabel string
	TitleFrom                                        []string
	Tile, Subtitle, HeaderActions                    templ.Component
}

func (h panelHead) titleFrom() string { return strings.Join(h.TitleFrom, " ") }

func (h panelHead) titleText() string {
	if h.Title != "" {
		return h.Title
	}
	return h.TitlePlaceholder
}

func (c FormPanelConfig) header() panelHead {
	return panelHead{ID: c.ID, Eyebrow: c.Eyebrow, Title: c.Title, TitlePlaceholder: c.TitlePlaceholder, CloseLabel: c.CloseLabel,
		TitleFrom: c.TitleFrom, Tile: c.Tile, Subtitle: c.Subtitle, HeaderActions: c.HeaderActions}
}

// ViewPanelConfig describe el panel de consulta (ViewPanel). Sin Half es el
// de 95 %, el de la ficha del cliente.
type ViewPanelConfig struct {
	ID       string
	Half     bool
	Tile     templ.Component
	Eyebrow  string
	Title    string
	Subtitle templ.Component
	// HeaderActions va antes del botón de cerrar («Editar»).
	HeaderActions templ.Component
	CloseLabel    string
	CloseHref     string
	Standalone    bool
	Static        bool
}

func (c ViewPanelConfig) header() panelHead {
	return panelHead{ID: c.ID, Eyebrow: c.Eyebrow, Title: c.Title, CloseLabel: c.CloseLabel,
		Tile: c.Tile, Subtitle: c.Subtitle, HeaderActions: c.HeaderActions}
}
