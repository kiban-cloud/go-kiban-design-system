// Package binding centralizes the translation of validator errors produced
// by gin's ShouldBind into Spanish, per-field messages the templ templates
// can render next to each input. Shared across every kiban htmx project so
// form error wording stays consistent.
//
// Import as `binding "github.com/kiban-cloud/go-kiban-design-system/binding"`
// — the package name is `binding`. Calling `binding.FieldErrors(err)` returns
// `map[formFieldName]spanishMessage`.
//
// # Built-in tags
//
// `MessageFor` knows how to translate the standard validator tags shipped
// by the `validator` library: required, email, min, max, len, url, oneof,
// gt, gte, eqfield. Unknown tags fall through to "Valor inválido".
//
// # Custom tags
//
// Projects that register their own validator tags (e.g. `regexRFC`,
// `regexCURP`, `regexCLABE` for Mexican fiscal IDs) call
// `binding.RegisterMessage` at startup to plug in the Spanish message
// for those tags:
//
//	func init() {
//	    binding.RegisterMessage("regexRFC", func(param string) string {
//	        return "RFC inválido"
//	    })
//	    binding.RegisterMessage("regexCURP", func(param string) string {
//	        return "CURP inválido"
//	    })
//	}
//
// Custom registrations override built-ins — projects that want to tweak
// the wording of a standard tag (e.g. switch "Este campo es obligatorio"
// to "Campo requerido") can re-register `required` with their own
// formatter. The default policy is to leave the built-ins alone so
// wording stays uniform across kiban projects; override only when you
// have a deliberate reason.
package binding

import (
	"reflect"
	"strings"
	"sync"

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

// customMessages holds project-registered tag formatters. Reads happen on
// every validation error, writes happen at startup; the RWMutex makes
// concurrent reads cheap while still being safe if a project registers
// late (e.g. from a lazy init path).
var (
	customMessagesMu sync.RWMutex
	customMessages   = map[string]func(param string) string{}
)

// RegisterMessage plugs in the Spanish message for a validator tag. Used
// for project-specific tags that aren't part of the built-in set
// (`regexRFC`, `regexCURP`, `regexCLABE`, etc.). The formatter receives
// the tag's `param` value (whatever's after the `=` in
// `binding:"tagName=param"`); ignore it when the tag doesn't take a
// parameter.
//
// Custom registrations override built-ins, so projects can re-register
// `required` / `email` / `min` / etc. with their own wording when
// product copy demands it. Default policy: leave the built-ins alone so
// the kiban projects don't drift on standard messages.
//
// Panics if `tag` is empty or `formatter` is nil — both are programming
// bugs, not runtime conditions, and silent misregistration would be far
// worse than a startup panic.
func RegisterMessage(tag string, formatter func(param string) string) {
	if tag == "" {
		panic("binding.RegisterMessage: empty tag")
	}
	if formatter == nil {
		panic("binding.RegisterMessage: nil formatter for tag " + tag)
	}
	customMessagesMu.Lock()
	defer customMessagesMu.Unlock()
	customMessages[tag] = formatter
}

// MessageFor returns the Spanish message for a validator tag and its
// optional `param`. Lookup order:
//
//  1. Project-registered formatters (via `RegisterMessage`).
//  2. Built-in tags shipped with this package (required, email, min,
//     max, len, url, oneof, gt, gte, eqfield).
//  3. "Valor inválido" — generic fallback so unknown tags never crash a
//     render; the worst outcome is a generic error next to the field.
//
// Use this directly when you need to translate a tag outside of
// `FieldErrors` (e.g. building error messages for an API JSON
// response). Inside HTMX handlers, prefer `FieldErrors(err)` — it walks
// a `validator.ValidationErrors` slice and calls `MessageFor` for each
// entry, returning the field-keyed map the templ templates expect.
func MessageFor(tag, param string) string {
	customMessagesMu.RLock()
	formatter, ok := customMessages[tag]
	customMessagesMu.RUnlock()
	if ok {
		return formatter(param)
	}

	switch tag {
	case "required":
		return "Este campo es obligatorio"
	case "email":
		return "Correo electrónico inválido"
	case "min":
		return "Valor demasiado corto (mínimo " + param + ")"
	case "max":
		return "Valor demasiado largo (máximo " + param + ")"
	case "len":
		return "Debe tener exactamente " + param + " caracteres"
	case "url":
		return "URL inválida"
	case "oneof":
		return "Valor no permitido"
	case "gt":
		return "Debe ser mayor que " + param
	case "gte":
		return "Debe ser mayor o igual que " + param
	case "eqfield":
		return "No coincide con " + param
	default:
		return "Valor inválido"
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
		out[fe.Field()] = MessageFor(fe.Tag(), fe.Param())
	}
	return out
}
