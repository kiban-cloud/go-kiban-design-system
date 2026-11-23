// Package avatar renders a circular profile picture with an initials
// fallback for when the image URL is empty, broken, or slow. Used by the
// admin layout's user menu and any other "user identity" surface.
package avatar

import (
	"context"
	"hash/fnv"
	"strings"

	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

type Options struct {
	Src     string // image URL; empty falls back to initials
	Name    string // used for initials and default aria-label
	Size    string // "sm" (24px) | "md" (32px, default) | "lg" (40px; 36px square)
	AltText string // optional override of Name for aria-label
	// Shape "square" es el avatar de una entidad (un cliente, una empresa):
	// cuadrado de esquinas redondeadas. Vacío o "circle", el de una persona
	// usuaria.
	Shape string
	// Tone colorea el fondo y las iniciales: "" (azul kiban) o "success"
	// (verde). Sirve para distinguir tipos, como persona física y moral.
	// "auto" (sólo flamboyant) saca el color del nombre con AutoTint: la
	// misma persona sale del mismo color en todas las pantallas.
	Tone string
	// Ring dibuja el borde blanco de 2px que separa a un avatar del anterior
	// cuando van encimados (-ml-1.5). Va en el avatar y no en un envoltorio:
	// el avatar es relative y se pinta encima del aro de un envoltorio
	// estático, así que el aro del vecino quedaría tapado.
	Ring bool
}

// autoTints son los pares fondo/texto de los avatares de personas en
// flamboyant, los del proto: azul, verde, morado, naranja, gris azulado y rojo.
var autoTints = [][2]string{
	{"#E8EEFF", "#0036C4"}, {"#E6FBF6", "#007A63"}, {"#F2E8FF", "#6A00D6"},
	{"#FFF1E6", "#9A4A00"}, {"#EEF3FF", "#33406B"}, {"#FDECEC", "#B42318"},
}

// AutoTint devuelve el fondo y el color de las iniciales para key (el nombre
// de la persona): FNV-1a de 32 bits sobre key, módulo la paleta. Pública para
// las pilas de avatares que no usan Avatar (p. ej. el Resumen de kiban cloud).
func AutoTint(key string) (bg, fg string) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(strings.TrimSpace(key)))
	t := autoTints[h.Sum32()%uint32(len(autoTints))]
	return t[0], t[1]
}

// autoTintClass son las clases de Tailwind de AutoTint.
func autoTintClass(key string) string {
	bg, fg := AutoTint(key)
	return "bg-[" + bg + "] text-[" + fg + "]"
}

func (o Options) sizeClass() string {
	switch strings.ToLower(strings.TrimSpace(o.Size)) {
	case "sm":
		return "h-6 w-6 text-xs"
	case "lg":
		return "h-10 w-10 text-base"
	default:
		return "h-8 w-8 text-sm"
	}
}

// flamboyantSizeClass es el tamaño del círculo en flamboyant: la letra mide
// el 38% del diámetro, como en el proto (24px → 9px, 32px → 12px, 40px → 15px).
func (o Options) flamboyantSizeClass() string {
	switch strings.ToLower(strings.TrimSpace(o.Size)) {
	case "sm":
		return "h-6 w-6 text-[9px] tracking-[.02em]"
	case "lg":
		return "h-10 w-10 text-[15px] tracking-[.02em]"
	default:
		return "h-8 w-8 text-[12px] tracking-[.02em]"
	}
}

func (o Options) square() bool {
	return strings.EqualFold(strings.TrimSpace(o.Shape), "square")
}

// class arma las clases del contenedor: forma, tamaño y tono, en su versión
// flamboyant cuando el request la pide (view/flamboyant).
func (o Options) class(ctx context.Context) string {
	fl := flamboyant.On(ctx)
	c := "relative inline-flex items-center justify-center overflow-hidden shrink-0 "
	if fl {
		c += "font-semibold "
	} else {
		c += "font-medium "
	}
	if o.square() {
		c += "rounded-lg " + o.squareSizeClass()
	} else if fl {
		c += "rounded-full " + o.flamboyantSizeClass()
	} else {
		c += "rounded-full " + o.sizeClass()
	}
	if o.Ring {
		c += " shadow-[0_0_0_2px_#fff]"
	}
	tone := strings.ToLower(strings.TrimSpace(o.Tone))
	success := tone == "success"
	switch {
	case fl && tone == "auto":
		return c + " " + autoTintClass(o.Name)
	case fl && success:
		return c + " bg-[#E3FBF5] text-[#007A63]"
	case fl:
		return c + " bg-[#EAF0FF] text-[#0036C4]"
	case success:
		return c + " bg-emerald-50 text-emerald-700"
	default:
		return c + " bg-kiban-primary-soft text-kiban-primary"
	}
}

func (o Options) squareSizeClass() string {
	switch strings.ToLower(strings.TrimSpace(o.Size)) {
	case "sm":
		return "h-6 w-6 text-[10px]"
	case "lg":
		return "h-9 w-9 text-xs"
	default:
		return "h-8 w-8 text-[11px]"
	}
}

func (o Options) altText() string {
	if strings.TrimSpace(o.AltText) != "" {
		return o.AltText
	}
	return o.Name
}

func (o Options) hasSrc() bool {
	return strings.TrimSpace(o.Src) != ""
}

// Initials returns up to 2 uppercase letters derived from Name.
// "Antonio Blancas" -> "AB", "Antonio" -> "A", "" -> "?".
// Public so consumers (e.g. server-rendered notification lists) can
// reuse the same derivation outside the templ.
func Initials(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "?"
	}
	firstRunes := []rune(parts[0])
	first := strings.ToUpper(string(firstRunes[0]))
	if len(parts) == 1 {
		return first
	}
	lastRunes := []rune(parts[len(parts)-1])
	last := strings.ToUpper(string(lastRunes[0]))
	return first + last
}
