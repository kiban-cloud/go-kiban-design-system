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
}

func (o StepOptions) hasHint() bool { return o.Hint != nil || o.HintText != "" }

func (o StepOptions) required() string { return strings.Join(o.Required, " ") }
