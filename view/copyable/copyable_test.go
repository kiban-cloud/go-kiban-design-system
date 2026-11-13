package copyable_test

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"sort"
	"testing"

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
