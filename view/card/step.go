package card

import (
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
}

func (o StepOptions) hasHint() bool { return o.Hint != nil || o.HintText != "" }

func (o StepOptions) required() string { return strings.Join(o.Required, " ") }

// EmptyOptions describe un estado vacío (Empty).
type EmptyOptions struct {
	Icon  templ.Component
	Title string
	Text  string
}
