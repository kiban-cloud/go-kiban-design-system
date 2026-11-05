// Package avatar renders a circular profile picture with an initials
// fallback for when the image URL is empty, broken, or slow. Used by the
// admin layout's user menu and any other "user identity" surface.
package avatar

import (
	"context"
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
	Tone string
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
	} else {
		c += "rounded-full " + o.sizeClass()
	}
	success := strings.EqualFold(strings.TrimSpace(o.Tone), "success")
	switch {
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
