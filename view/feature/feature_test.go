package feature_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
)

func texto(s string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	})
}

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	assert.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

// bloque envuelve el hijo como lo haría `@feature.Block("x") { … }` en un templ.
func bloque(name string, hijo templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return feature.Block(name).Render(templ.WithChildren(ctx, hijo), w)
	})
}

// Sin banderas en el contexto no se ve nada nuevo: una vista renderizada por
// un módulo que todavía no puso el middleware muestra sólo lo que ya existe.
func TestBlock_SinContextoNoEmiteNada(t *testing.T) {
	body := render(t, context.Background(), bloque("cuotas", texto("NUEVO")))
	assert.Equal(t, "", body)
}

// Prendida, el bloque sale tal cual: sin marco, sin etiqueta, sin atributo. En
// producción el flag no deja rastro en el HTML.
func TestBlock_PrendidaEmiteSoloLosHijos(t *testing.T) {
	ctx := feature.With(context.Background(), feature.Flags{Enabled: map[string]bool{"cuotas": true}})
	body := render(t, ctx, bloque("cuotas", texto("NUEVO")))
	assert.Equal(t, "NUEVO", body)
}

// En Comparar se ve todo, apagado incluido, con el marco y la etiqueta de la
// bandera (o su nombre si no tiene etiqueta).
func TestBlock_CompararResaltaConEtiqueta(t *testing.T) {
	ctx := feature.With(context.Background(), feature.Flags{
		Compare: true,
		Labels:  map[string]string{"cuotas": "Fase 2"},
	})
	body := render(t, ctx, bloque("cuotas", texto("NUEVO")))
	assert.Contains(t, body, `data-feature="cuotas"`)
	assert.Contains(t, body, "outline-dashed")
	assert.Contains(t, body, "Fase 2")
	assert.Contains(t, body, "NUEVO")

	sinEtiqueta := render(t, feature.With(context.Background(), feature.Flags{Compare: true}),
		bloque("dispersion-masiva", texto("X")))
	assert.Contains(t, sinEtiqueta, "dispersion-masiva")
}

func TestReplace_EligePorBandera(t *testing.T) {
	apagada := render(t, context.Background(), feature.Replace("cuotas", texto("ANTES"), texto("DESPUES")))
	assert.Equal(t, "ANTES", apagada)

	ctx := feature.With(context.Background(), feature.Flags{Enabled: map[string]bool{"cuotas": true}})
	prendida := render(t, ctx, feature.Replace("cuotas", texto("ANTES"), texto("DESPUES")))
	assert.Equal(t, "DESPUES", prendida)
}

// Comparar muestra los dos, el actual primero y atenuado como «antes».
func TestReplace_CompararMuestraAmbos(t *testing.T) {
	ctx := feature.With(context.Background(), feature.Flags{Compare: true, Labels: map[string]string{"cuotas": "Fase 2"}})
	body := render(t, ctx, feature.Replace("cuotas", texto("ANTES"), texto("DESPUES")))
	iAntes, iDespues := strings.Index(body, "ANTES"), strings.Index(body, "DESPUES")
	assert.True(t, iAntes >= 0 && iDespues > iAntes, "el actual va antes del nuevo: %s", body)
	assert.Contains(t, body, "antes · Fase 2")
	assert.Contains(t, body, "opacity-60")
	assert.Equal(t, 2, strings.Count(body, `data-feature="cuotas"`))
}

func TestLabel_CaeAlNombre(t *testing.T) {
	assert.Equal(t, "cuotas", feature.Label(context.Background(), "cuotas"))
	ctx := feature.With(context.Background(), feature.Flags{Labels: map[string]string{"cuotas": "Fase 2"}})
	assert.Equal(t, "Fase 2", feature.Label(ctx, "cuotas"))
}
