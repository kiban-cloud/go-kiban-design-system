package copyable_test

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"sort"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"

	"github.com/kiban-cloud/go-kiban-design-system/view/copyable"
)

func render(t *testing.T, ctx context.Context, value string) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, copyable.ID(value).Render(ctx, &buf))
	return buf.String()
}

func TestID_CopiaElValorQueMuestra(t *testing.T) {
	body := render(t, context.Background(), "cli_01HX")

	assert.Contains(t, body, `data-kiban-copy="cli_01HX"`)
	assert.Contains(t, body, `>cli_01HX`)
	assert.Contains(t, body, `type="button"`)
}

func TestID_TextosSegunIdioma(t *testing.T) {
	es := render(t, context.Background(), "cli_1")
	assert.Contains(t, es, `aria-label="Copiar cli_1"`)
	assert.Contains(t, es, `data-kiban-copied-label="Copiado"`)

	en := render(t, fsi18n.WithLanguage(context.Background(), "en"), "cli_1")
	assert.Contains(t, en, `aria-label="Copy cli_1"`)
	assert.Contains(t, en, `data-kiban-copied-label="Copied"`)
}

func TestLocales_MismasClaves(t *testing.T) {
	keys := func(path string) []string {
		b, err := os.ReadFile(path)
		require.NoError(t, err)
		var out []string
		for _, m := range regexp.MustCompile(`(?m)^"([^"]+)"\s*=`).FindAllStringSubmatch(string(b), -1) {
			out = append(out, m[1])
		}
		sort.Strings(out)
		return out
	}
	assert.Equal(t, keys("locales/es.toml"), keys("locales/en.toml"))
}

func TestIDWith_ToastOpcional(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, copyable.IDWith(copyable.Options{Value: "abc", Toast: "ID copiado"}).Render(context.Background(), &buf))
	body := buf.String()
	assert.Contains(t, body, `data-kiban-copy="abc"`)
	assert.Contains(t, body, `data-kiban-copy-toast="ID copiado"`)

	buf.Reset()
	require.NoError(t, copyable.ID("abc").Render(context.Background(), &buf))
	assert.NotContains(t, buf.String(), "data-kiban-copy-toast", "ID solo no lleva aviso")
}

func renderIn(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

func TestField_CopiaElValorConSuEtiqueta(t *testing.T) {
	body := renderIn(t, context.Background(), copyable.Field(copyable.FieldOptions{Value: "https://pay.kiban.com/o/1", Label: "Copiar link"}))

	assert.Contains(t, body, "data-kiban-copy-field")
	assert.Contains(t, body, `data-kiban-copy="https://pay.kiban.com/o/1"`)
	assert.Contains(t, body, ">https://pay.kiban.com/o/1</span>")
	assert.Contains(t, body, `aria-label="Copiar link"`)
	assert.Contains(t, body, `title="Copiar link"`)
	assert.Contains(t, body, `data-kiban-copied-label="Copiado"`)
	assert.NotContains(t, body, "line-through")
	assert.NotContains(t, body, "font-mono")
}

// Sin Label el nombre accesible es «Copiar» y el valor, según el idioma;
// Muted lo tacha (igual se copia) y Mono lo pinta en monoespaciado.
func TestField_EtiquetaPorDefectoMutedYMono(t *testing.T) {
	body := renderIn(t, context.Background(), copyable.Field(copyable.FieldOptions{Value: "abc", Muted: true, Mono: true}))
	assert.Contains(t, body, `aria-label="Copiar abc"`)
	assert.Contains(t, body, "line-through")
	assert.Contains(t, body, "font-mono")
	assert.Contains(t, body, `data-kiban-copy="abc"`)

	en := renderIn(t, fsi18n.WithLanguage(context.Background(), "en"), copyable.Field(copyable.FieldOptions{Value: "abc"}))
	assert.Contains(t, en, `aria-label="Copy abc"`)
}

// Inline: la etiqueta va antes del valor en el mismo texto, pero se copia
// sólo el valor.
func TestInline_EtiquetaSeparadorYCopia(t *testing.T) {
	body := renderIn(t, context.Background(), copyable.Inline(copyable.InlineOptions{Label: "ID", Value: "991ab76c", CopyLabel: "Copiar ID", Sep: true}))

	assert.Contains(t, body, "data-kiban-copy-inline")
	assert.Contains(t, body, "ID 991ab76c")
	assert.Contains(t, body, `data-kiban-copy="991ab76c"`)
	assert.Contains(t, body, `aria-label="Copiar ID"`)
	assert.Contains(t, body, "before:mr-2", "Sep pone el punto gris")
}

func TestInline_SoloValor(t *testing.T) {
	body := renderIn(t, context.Background(), copyable.Inline(copyable.InlineOptions{Value: "991ab76c"}))

	assert.Contains(t, body, `data-kiban-copy-inline>991ab76c`)
	assert.Contains(t, body, `aria-label="Copiar 991ab76c"`)
	assert.NotContains(t, body, "before:mr-2")
}
