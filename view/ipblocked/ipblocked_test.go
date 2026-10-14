package ipblocked_test

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/ipblocked"
)

func render(t *testing.T, tool, ip string) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, ipblocked.Page(tool, ip).Render(context.Background(), &buf))
	return buf.String()
}

func TestPage_DiceLaIpLaHerramientaYLasSalidas(t *testing.T) {
	body := render(t, "WORKFLOO", "203.0.113.9")

	assert.Contains(t, body, "<!doctype html>")
	assert.Contains(t, body, `data-ip-blocked="WORKFLOO"`)
	assert.Contains(t, body, "Acceso no permitido desde esta IP")
	assert.Contains(t, body, "203.0.113.9")
	assert.Contains(t, body, "workfloo")
	assert.Contains(t, body, `href="/kiban-cloud/developers/ipaccess?lista=consola"`)
	assert.Contains(t, body, `href="/kiban-cloud/"`)
}

// crm se llama "clientes" en la interfaz.
func TestToolLabel(t *testing.T) {
	assert.Equal(t, "clientes", ipblocked.ToolLabel(context.Background(), "CRM"))
	assert.Equal(t, "link", ipblocked.ToolLabel(context.Background(), "LINK"))
}

// Los dos bundles tienen las mismas claves.
func TestBundleParity(t *testing.T) {
	claves := func(path string) []string {
		raw, err := os.ReadFile(path)
		require.NoError(t, err)
		var out []string
		for _, m := range regexp.MustCompile(`(?m)^"([^"]+)"`).FindAllStringSubmatch(string(raw), -1) {
			out = append(out, m[1])
		}
		sort.Strings(out)
		return out
	}
	assert.Equal(t, claves("locales/es.toml"), claves("locales/en.toml"))
}
