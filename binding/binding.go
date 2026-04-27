// Package binding centralizes the translation of validator errors produced
// by gin's ShouldBind into Spanish, per-field messages the templ templates
// can render next to each input. Shared across every kiban htmx project so
// form error wording stays consistent.
//
// Import as `binding "github.com/kiban-cloud/go-kiban-design-system/binding"`
// — the package name is `binding`. Calling `binding.FieldErrors(err)` returns
// `map[formFieldName]spanishMessage`.
package binding

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	// Ensure validator reports the `form:"..."` tag name rather than the Go
	// struct field name, so FieldErrors keys match the HTML input name
	// attribute. Runs once at package import time — handlers don't need to
	// call anything to opt in.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// FieldErrors maps a validator-produced error into map[formFieldName]spanishMessage.
// Non-validator errors (ShouldBind can also fail with json/time/strconv
// errors) return an empty map — the caller shows a global form error in
// that case.
func FieldErrors(err error) map[string]string {
	out := map[string]string{}
	if err == nil {
		return out
	}
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return out
	}
	for _, fe := range ve {
		out[fe.Field()] = messageFor(fe)
	}
	return out
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio"
	case "email":
		return "Correo electrónico inválido"
	case "min":
		return "Valor demasiado corto (mínimo " + fe.Param() + ")"
	case "max":
		return "Valor demasiado largo (máximo " + fe.Param() + ")"
	case "len":
		return "Debe tener exactamente " + fe.Param() + " caracteres"
	case "url":
		return "URL inválida"
	case "oneof":
		return "Valor no permitido"
	case "gt":
		return "Debe ser mayor que " + fe.Param()
	case "gte":
		return "Debe ser mayor o igual que " + fe.Param()
	case "eqfield":
		return "No coincide con " + fe.Param()
	default:
		return "Valor inválido"
	}
}
