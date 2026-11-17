package card

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// StepOptions describe una tarjeta de paso de formulario (Step).
type StepOptions struct {
	// Key identifica el paso (data-kiban-step); no se muestra.
	Key   string
	Title string
	// HintText o Hint (si la ayuda lleva marcado, como un asterisco rojo)
	// van debajo del título. Hint gana si vienen los dos.
	HintText string
	Hint     templ.Component
	// Right va a la derecha del título (un botón «Agregar contacto»).
	Right templ.Component
	// Required son los name de los campos que completan el paso: en
	// flamboyant el número se pinta azul cuando todos tienen un valor
	// válido. Vacío: el paso se completa con cualquier campo lleno.
	Required []string
	// NoNumber quita el número en flamboyant: la tarjeta única de un panel
	// de edición, donde no hay pasos que contar.
	NoNumber bool
	// NoProgress deja el número siempre gris, aunque el paso esté completo:
	// los paneles de kiban cloud en el proto no marcan avance.
	NoProgress bool
	// Done pinta el paso como completo desde el servidor: círculo azul con
	// palomita en lugar del número, aunque venga NoNumber. Es para las
	// páginas de configuración (no paneles), donde nadie numera ni revisa
	// los campos en el navegador: la sección ya guardada se ve hecha. Un
	// paso Done no cuenta en la numeración de un panel.
	Done bool
	// Number es el número del círculo cuando lo pone el servidor (páginas
	// sin panel). 0: lo pone el navegador en un panel, o «·».
	Number int
}

func (o StepOptions) hasHint() bool { return o.Hint != nil || o.HintText != "" }

func (o StepOptions) required() string { return strings.Join(o.Required, " ") }

// number es lo que va en el círculo antes de que el navegador numere.
func (o StepOptions) number() string {
	if o.Number > 0 {
		return strconv.Itoa(o.Number)
	}
	return "·"
}

// EmptyOptions describe un estado vacío (Empty).
type EmptyOptions struct {
	Icon  templ.Component
	Title string
	Text  string
}
