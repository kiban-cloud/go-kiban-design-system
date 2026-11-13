package drawer

import (
	"encoding/json"
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
	// TitleInput pone el título como un campo del form que se edita en el
	// encabezado (el nombre del secret): reemplaza a Title y TitleFrom, y el
	// Eyebrow es su label.
	TitleInput    *TitleInput
	Subtitle      templ.Component
	HeaderActions templ.Component
	CloseLabel    string

	// FormAttrs van en el <form>: method, action, hx-post, hx-target…
	FormAttrs templ.Attributes

	FooterLeft templ.Component
	// SummaryPending y SummaryReady son el texto del pie según el form pase
	// o no la validación del navegador (required, type=email…).
	// SummaryClean, si viene, reemplaza a SummaryReady mientras no se haya
	// tocado nada («Sin cambios» en una edición).
	SummaryPending string
	// SummaryMissing, si viene, cambia SummaryPending por lo que falta de
	// verdad, contando los campos obligatorios vacíos de cada parte, y se
	// actualiza mientras se llena el form: «Faltan la URL del proceso y 4
	// capturas», «Falta 1 captura» (el pie de la certificación del proto).
	// Sin nada vacío que contar, queda SummaryPending.
	SummaryMissing *SummaryMissing
	SummaryReady   string
	SummaryClean   string
	Actions        button.Group
	// PrimarySplit reemplaza a Actions.PrimaryAction por un botón dividido
	// con variantes de la misma acción («Crear», «Crear y cerrar», «Crear y
	// nuevo»). Nil: el botón de Actions.
	PrimarySplit *button.SplitOptions

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

// TitleInput es el campo del título editable de FormPanel.
type TitleInput struct {
	Name, Value, Placeholder string
	// Help es para qué sirve el campo, sólo para lectores de pantalla (el
	// encabezado no tiene lugar para una ayuda visible).
	Help string
	// Required lo exige el navegador: el pie dice SummaryPending sin él.
	Required bool
}

// FormPanelGuard son los textos del aviso de cambios sin guardar. SaveClose
// vacío quita el botón «Guardar y cerrar» (en ViewPanel no hay qué guardar).
type FormPanelGuard struct {
	Title, Text              string
	Discard, Stay, SaveClose string
	// Dirty es el aviso del encabezado, junto a cerrar, en cuanto hay un
	// cambio («Cambios sin guardar», el .fsp-dirty del proto). Vacío: sin
	// aviso.
	Dirty string
}

// panelHead es lo que pinta panelHeader, común a FormPanel y ViewPanel.
type panelHead struct {
	ID, Eyebrow, Title, TitlePlaceholder, CloseLabel, Dirty string
	TitleFrom                                               []string
	TitleInput                                              *TitleInput
	Tile, Subtitle, HeaderActions                           templ.Component
}

func (h panelHead) titleFrom() string { return strings.Join(h.TitleFrom, " ") }

func (h panelHead) titleText() string {
	if h.Title != "" {
		return h.Title
	}
	return h.TitlePlaceholder
}

func (c FormPanelConfig) summaryInitial() string {
	if c.SummaryClean != "" {
		return c.SummaryClean
	}
	return c.SummaryPending
}

func (c FormPanelConfig) header() panelHead {
	return panelHead{ID: c.ID, Eyebrow: c.Eyebrow, Title: c.Title, TitlePlaceholder: c.TitlePlaceholder, CloseLabel: c.CloseLabel,
		TitleFrom: c.TitleFrom, TitleInput: c.TitleInput, Tile: c.Tile, Subtitle: c.Subtitle, HeaderActions: c.HeaderActions, Dirty: c.dirtyLabel()}
}

func (c FormPanelConfig) dirtyLabel() string {
	if c.Guard == nil {
		return ""
	}
	return c.Guard.Dirty
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
	// Footer es el pie del panel (estado y botones); nil, sin pie.
	Footer templ.Component
	// Guard pide confirmar al cerrar si se escribió en un form del cuerpo
	// (data-dirty, como en FormPanel). Sin SaveClose no hay «Guardar y
	// cerrar». nil: cierra sin preguntar.
	Guard *FormPanelGuard
}

func (c ViewPanelConfig) header() panelHead {
	return panelHead{ID: c.ID, Eyebrow: c.Eyebrow, Title: c.Title, CloseLabel: c.CloseLabel,
		Tile: c.Tile, Subtitle: c.Subtitle, HeaderActions: c.HeaderActions, Dirty: c.dirtyLabel()}
}

func (c ViewPanelConfig) dirtyLabel() string {
	if c.Guard == nil {
		return ""
	}
	return c.Guard.Dirty
}

// ClosePanelEvent es el evento con el que el servidor cierra el panel desde
// el que se hizo la acción, con la animación de salida: c.Header("HX-Trigger",
// drawer.ClosePanelEvent). Lo atiende base.templ. Sirve cuando la respuesta
// cambia otra parte de la página (la lista tras borrar) y el panel sobra.
const ClosePanelEvent = "kiban-panel-close"

// SummaryMissing arma el pie pendiente con lo que falta: One y Many son el
// verbo según falte uno o más («Falta», «Faltan»), y And une las partes («y»).
type SummaryMissing struct {
	One, Many, And string
	Parts          []MissingPart
}

// MissingPart es una parte del pie pendiente: los name de sus campos y el
// texto según falte uno o más («1 captura», «{n} capturas»; {n} es cuántos).
// Un grupo de radios falta mientras no haya uno elegido. From, en vez de
// Names, es el selector de un elemento del form cuyo
// data-kiban-missing-names (name separados por coma) dice qué contar: para
// campos que dependen de una elección y se repintan (las credenciales
// obligatorias del tipo de secret). Si no está en el form, la parte no cuenta.
type MissingPart struct {
	Names     []string
	One, Many string
	From      string
}

func (c FormPanelConfig) missingJSON() string {
	if c.SummaryMissing == nil {
		return ""
	}
	m := c.SummaryMissing
	type part struct {
		Names []string `json:"names"`
		One   string   `json:"one"`
		Many  string   `json:"many"`
		From  string   `json:"from,omitempty"`
	}
	out := struct {
		One   string `json:"one"`
		Many  string `json:"many"`
		And   string `json:"and"`
		Parts []part `json:"parts"`
	}{One: m.One, Many: m.Many, And: m.And}
	for _, p := range m.Parts {
		out.Parts = append(out.Parts, part{Names: p.Names, One: p.One, Many: p.Many, From: p.From})
	}
	b, _ := json.Marshal(out)
	return string(b)
}
