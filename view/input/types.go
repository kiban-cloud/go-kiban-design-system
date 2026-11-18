// Package input holds the shared form-control primitives every kiban htmx
// project uses: text, select, textarea, checkbox, toggle, phone (intl-tel-input
// wrapper), hidden, and date. Each component renders a label + control + error
// (or hint) block following the kiban design tokens (kiban-border /
// kiban-primary / red-400 / red-500), so error state is consistent across
// projects.
//
// Conventions:
//   - Components emit `id="f-<name>"` paired with `<label for="f-<name>">` so
//     clicking the label focuses the control.
//   - When `errMsg != ""` the border swaps to red and the message renders
//     under the control. When `errMsg == ""` and `hint != ""` the hint renders
//     instead. They never both appear.
//   - Required is a visual marker only (a red asterisk next to the label).
//     HTML5 `required` is intentionally NOT emitted: validation lives
//     server-side and the label asterisk is the user-facing hint.
package input

import (
	"strings"

	"github.com/a-h/templ"
)

// SelectOption is one row of a `<select>`. Iterated by Select to render
// `<option value=Value selected=...>Label</option>`. Value is what the form
// submits; Label is what the user sees.
type SelectOption struct {
	Value string
	Label string
}

// RadioCardOptions configura RadioCardWith.
type RadioCardOptions struct {
	Name, Value, Title, Subtitle string
	Checked                      bool
	Attrs                        templ.Attributes
	// Icon y Tone sólo se ven en flamboyant: el ícono va en un cuadro azul
	// (Tone vacío), verde (Tone "success"), gris (Tone "neutral"), morado
	// (Tone "purple") o, con Tone "plate", en una placa blanca con borde para
	// un logo (el de un SIC).
	Icon templ.Component
	Tone string
	// Tags son etiquetas bajo el subtítulo (lo que incluye una plantilla) y
	// Big la tarjeta más alta, con el ícono arriba (las plantillas de un
	// flujo). Sólo flamboyant.
	Tags []string
	Big  bool
}

// ToggleOptions configura ToggleWith.
type ToggleOptions struct {
	Name, Label string
	// Hint es la línea de ayuda: qué hace la opción o, deshabilitada, por qué.
	Hint           string
	Value, Enabled bool
	Attrs          templ.Attributes
}

// DropzoneOptions configura Dropzone. Title es la invitación («Arrastra tu
// archivo CSV aquí o elígelo») y Hint, las condiciones.
type DropzoneOptions struct {
	Name, Accept, Title, Hint string
	Required                  bool
	// Pick va después de Title en azul, con Or delante: «Arrastra la
	// captura o elígela» (Title «Arrastra la captura», Or «o», Pick
	// «elígela»), como el proto. Con Pick, Title va en negritas.
	Or, Pick string
	// Image cambia el ícono por el de imagen y, al elegir una imagen, la
	// muestra en miniatura en su lugar.
	Image bool
	// Tall le da a la zona un alto mínimo de 150 px (las capturas de la
	// certificación).
	Tall bool
	// Change, si viene, es lo que dice la zona debajo del nombre del archivo
	// elegido («Cambiar archivo»): con él, al elegir, la zona se pinta en
	// verde y el nombre reemplaza a Title. Sin él, el nombre y el tamaño van
	// debajo, en azul, como siempre.
	Change string
	// Multiple deja elegir varios archivos; la zona lista cuántos y cuáles.
	Multiple bool
}

// CheckboxCardOptions configura CheckboxCardWith: una opción de una lista de
// casillas con el mismo Name (grupos, permisos), cada una con su Value.
type CheckboxCardOptions struct {
	Name, Value, Title, Help string
	Checked, Disabled        bool
	Attrs                    templ.Attributes
}

// radioCardToneClass es el cuadro del ícono de RadioCardWith según su Tone.
func radioCardToneClass(tone string) string {
	switch tone {
	case "success":
		return "bg-[#E3FBF5] text-[#007A63]"
	case "neutral":
		return "bg-[#F2F2F2] text-[#474747]"
	case "purple":
		return "bg-[#F3EAFF] text-[#6A00D6]"
	case "plate":
		return "overflow-hidden border border-[#E2E5EC] bg-white p-1 [&_img]:h-full [&_img]:w-full [&_img]:object-contain"
	}
	return "bg-[#EAF0FF] text-[#0036C4]"
}

// PillsOptions configura Pills. Selected es el Value marcado; Label es el
// aria-label del grupo. Attrs de cada opción van en su <input>.
type PillsOptions struct {
	Name     string
	Label    string
	Options  []PillOption
	Selected string
	Required bool
	Error    bool
	Disabled bool
}

// PillOption es una píldora de Pills.
type PillOption struct {
	Value, Label string
	Attrs        templ.Attributes
}

// MoneyOptions configura Money.
type MoneyOptions struct {
	// Name es el name del campo; el id es "f-" + Name, como los demás.
	Name string
	// Label es la etiqueta visible arriba de la caja; vacía, la caja va sola
	// (el título de la tarjeta ya dice qué es) y AriaLabel la nombra.
	Label     string
	AriaLabel string
	Value     string
	// Prefix es el signo de adelante; vacío, «$». Currency es la moneda de
	// atrás («MXN»); vacía, no se pinta.
	Prefix   string
	Currency string
	// Placeholder vacío es «0.00».
	Placeholder string
	// Hint es la ayuda de abajo y Error, el error (gana sobre Hint). Below
	// reemplaza a los dos: para una ayuda que el servidor repinta (htmx oob)
	// con su propio id.
	Hint, Error string
	Below       templ.Component
	// Required pone el asterisco en la etiqueta y el required del navegador.
	Required bool
	// Max es el monto más alto que acepta («10000», sin comas); MaxError, lo
	// que dice el navegador si se pasa. Vacío: sin tope.
	Max, MaxError string
	// Compact es la caja de 36 px para ir en línea con otros campos.
	Compact bool
	Attrs   templ.Attributes
}

func (o MoneyOptions) id() string { return "f-" + o.Name }

func (o MoneyOptions) prefix() string {
	if o.Prefix == "" {
		return "$"
	}
	return o.Prefix
}

func (o MoneyOptions) placeholder() string {
	if o.Placeholder == "" {
		return "0.00"
	}
	return o.Placeholder
}

func (o MoneyOptions) ariaLabel() string {
	if o.AriaLabel != "" {
		return o.AriaLabel
	}
	return ""
}

// PickerOptions configura Picker.
type PickerOptions struct {
	// ID es el id del bloque entero: para repintarlo (htmx oob).
	ID string
	// Name es el name de los radios; Selected, el Value elegido.
	Name     string
	Selected string
	Options  []PickerOption
	// Label nombra el grupo para lectores de pantalla («Cliente»).
	Label string
	// Search es el placeholder del buscador («Buscar banco»); vacío, sin
	// buscador. Empty es lo que se dice cuando nada coincide («Ningún
	// cliente coincide»); vacío, la lista sólo queda vacía.
	Search, Empty string
	// Columns 2 pone las tarjetas en dos columnas (una en móvil); cualquier
	// otro valor, una.
	Columns int
	// Required pide elegir una (el required del navegador en los radios).
	Required bool
	// Error va debajo de la lista.
	Error string
	// Attrs van en cada radio.
	Attrs templ.Attributes
}

// PickerOption es una opción de Picker. Tile es su placa de 40 px (iniciales,
// logo); nil, sin placa. Search es texto de más para el buscador que no se
// ve (la clave de un banco, un id).
type PickerOption struct {
	Value, Title, Sub string
	Tile              templ.Component
	Search            string
}

func (p PickerOption) searchText() string {
	return strings.ToLower(strings.Join([]string{p.Title, p.Sub, p.Search}, " "))
}
