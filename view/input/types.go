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

import "github.com/a-h/templ"

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

// ChipsOptions configura Chips. Selected es el Value marcado; Label es el
// aria-label del grupo. Attrs de cada opción van en su <input>.
type ChipsOptions struct {
	Name     string
	Label    string
	Options  []ChipOption
	Selected string
	Required bool
	Error    bool
	Disabled bool
}

type ChipOption struct {
	Value, Label string
	Attrs        templ.Attributes
}
