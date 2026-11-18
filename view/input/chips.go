package input

import (
	"regexp"
	"strings"
)

// ChipsOptions describe un campo de lista de valores como chips (Chips):
// los correos que reciben un aviso, etiquetas, dominios. Cada valor es un
// chip con su ×; al final de la misma caja, un campo para agregar otro.
type ChipsOptions struct {
	// Name es el campo del form: un hidden con los valores, uno por renglón
	// («a@x.mx\nb@x.mx»). El handler lo parte por renglón, coma o espacio.
	Name string
	// ID es el id del campo de texto (para el label). Vacío: «f-<Name>».
	ID string
	// Label va arriba de la caja; vacío no pinta label (la tarjeta ya
	// tiene título) y AriaLabel nombra el campo para lectores de pantalla.
	Label     string
	AriaLabel string
	Required  bool
	// Values son los valores con los que llega (los guardados).
	Values []string
	// Type "email" marca en rojo, con InvalidTitle de tooltip, los chips
	// sin forma de correo, y pide el teclado de correo en el celular. Vacío
	// acepta cualquier texto.
	Type string
	// Placeholder es el del campo cuando no hay chips
	// («cobranza@tuempresa.mx»); AddPlaceholder cuando ya hay
	// («Agregar otro correo»). Sin AddPlaceholder se usa Placeholder.
	Placeholder    string
	AddPlaceholder string
	// RemoveLabel es el aria-label de la × antes del valor («Quitar»).
	RemoveLabel  string
	InvalidTitle string
	// Error pinta la caja en rojo con el mensaje debajo; Hint es la ayuda
	// cuando no hay error.
	Error string
	Hint  string
}

var chipsEmail = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func (o ChipsOptions) id() string {
	if o.ID != "" {
		return o.ID
	}
	return "f-" + o.Name
}

func (o ChipsOptions) joined() string { return strings.Join(o.Values, "\n") }

func (o ChipsOptions) addPlaceholder() string {
	if o.AddPlaceholder != "" {
		return o.AddPlaceholder
	}
	return o.Placeholder
}

func (o ChipsOptions) placeholder() string {
	if len(o.Values) > 0 {
		return o.addPlaceholder()
	}
	return o.Placeholder
}

func (o ChipsOptions) removeLabel() string {
	if o.RemoveLabel != "" {
		return o.RemoveLabel
	}
	return "Quitar"
}

// bad dice si un valor se pinta en rojo: sólo con Type "email".
func (o ChipsOptions) bad(v string) bool {
	return o.Type == "email" && !chipsEmail.MatchString(v)
}

func (o ChipsOptions) inputMode() string {
	if o.Type == "email" {
		return "email"
	}
	return "text"
}

func (o ChipsOptions) autocomplete() string {
	if o.Type == "email" {
		return "email"
	}
	return "off"
}

func (o ChipsOptions) ariaLabel() string {
	if o.AriaLabel != "" {
		return o.AriaLabel
	}
	return o.Label
}
