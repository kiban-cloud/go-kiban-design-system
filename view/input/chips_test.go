package input_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

// Los valores viajan en un solo hidden, uno por renglón; el campo de texto
// no tiene name.
func TestChips_ValoresEnElHidden(t *testing.T) {
	body := renderIn(t, context.Background(), input.Chips(input.ChipsOptions{
		Name: "emails", Label: "Correos", Required: true, Values: []string{"a@x.mx", "b@x.mx"},
		Placeholder: "cobranza@tuempresa.mx", AddPlaceholder: "Agregar otro correo", RemoveLabel: "Eliminar",
	}))

	assert.Contains(t, body, `<input type="hidden" name="emails" value="a@x.mx
b@x.mx" data-kiban-chips-value>`)
	assert.Contains(t, body, `<label for="f-emails"`)
	assert.Contains(t, body, `id="f-emails" type="text"`)
	assert.Equal(t, 1, strings.Count(body, `name="`), "sólo el hidden lleva name")
	assert.Contains(t, body, `data-kiban-chip="a@x.mx"`)
	assert.Contains(t, body, `data-kiban-chip="b@x.mx"`)
	assert.Contains(t, body, `aria-label="Eliminar a@x.mx"`)
	assert.Contains(t, body, `placeholder="Agregar otro correo"`, "con chips, el placeholder de agregar")
	assert.Contains(t, body, `data-placeholder="cobranza@tuempresa.mx" data-add-placeholder="Agregar otro correo" data-remove-label="Eliminar"`)
	assert.Contains(t, body, `<span class="ml-0.5 text-red-500">*</span>`)
	assert.Contains(t, body, `aria-label="Correos"`, "sin AriaLabel, el Label")
	assert.Contains(t, body, `<template data-kiban-chips-tpl>`, "el molde para el navegador")
	assert.Contains(t, body, `inputmode="text" autocomplete="off"`)
}

// Sin valores: placeholder normal, sin AddPlaceholder se usa Placeholder, y
// «Quitar» por defecto en la ×.
func TestChips_SinValoresYDefaults(t *testing.T) {
	body := renderIn(t, context.Background(), input.Chips(input.ChipsOptions{Name: "tags", ID: "etiquetas", AriaLabel: "Etiquetas", Placeholder: "Agregar"}))

	assert.Contains(t, body, `name="tags" value=""`)
	assert.Contains(t, body, `id="etiquetas"`)
	assert.Contains(t, body, `placeholder="Agregar"`)
	assert.Contains(t, body, `data-add-placeholder="Agregar"`)
	assert.Contains(t, body, `data-remove-label="Quitar"`)
	assert.Contains(t, body, `aria-label="Etiquetas"`)
	assert.NotContains(t, body, "<label for=\"etiquetas\" class=\"text-xs", "sin Label no hay etiqueta visible")
	assert.Equal(t, 1, strings.Count(body, "data-kiban-chip="), "sólo el molde")
}

// Type «email»: los chips sin forma de correo van en rojo con su tooltip.
func TestChips_CorreosInvalidosEnRojo(t *testing.T) {
	body := renderIn(t, context.Background(), input.Chips(input.ChipsOptions{
		Name: "emails", Type: "email", Values: []string{"ok@x.mx", "no-es-correo"}, InvalidTitle: "No es un correo",
	}))

	assert.Contains(t, body, `data-kiban-chips="email"`)
	assert.Contains(t, body, `inputmode="email" autocomplete="email"`)
	assert.Contains(t, body, `data-kiban-chip="no-es-correo" data-bad title="No es un correo"`)
	assert.NotContains(t, body, `data-kiban-chip="ok@x.mx" data-bad`)
	assert.Equal(t, 1, strings.Count(body, "data-bad "), "sólo el inválido (el molde no)")

	libre := renderIn(t, context.Background(), input.Chips(input.ChipsOptions{Name: "tags", Values: []string{"no-es-correo"}}))
	assert.NotContains(t, libre, "data-bad ", "sin Type acepta cualquier texto")
}

func TestChips_ErrorGanaSobreAyuda(t *testing.T) {
	err := renderIn(t, context.Background(), input.Chips(input.ChipsOptions{Name: "e", Error: "Correo inválido", Hint: "Separa con comas"}))
	assert.Contains(t, err, ">Correo inválido</div>")
	assert.NotContains(t, err, "Separa con comas")
	assert.Contains(t, err, "border-red-400")

	ayuda := renderIn(t, context.Background(), input.Chips(input.ChipsOptions{Name: "e", Hint: "Separa con comas"}))
	assert.Contains(t, ayuda, ">Separa con comas</div>")
	assert.Contains(t, ayuda, "border-[#CDD3DF]")
}
