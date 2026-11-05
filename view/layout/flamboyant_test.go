package layout_test

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/icons"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

func flamboyantCfg() layout.Config {
	return layout.Config{
		Title:        "Clientes",
		ProjectName:  "crm",
		SectionLabel: "clientes",
		User:         layout.User{Name: "Ana Pérez", Email: "ana@kiban.com"},
		SpaceID:      "space-123456",
		LogoutAction: "/crm/logout",
		ToolKey:      "customers",
		Tools:        []layout.Tool{{Key: "customers", Label: "CRM", Href: "/crm/customers", Icon: icons.BookOpen}},
	}
}

func renderLayout(t *testing.T, ctx context.Context, cfg layout.Config) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, layout.Layout(cfg).Render(ctx, &buf))
	return buf.String()
}

// shell recorta el HTML del shell flamboyant (desde fs-root hasta su primer
// <script>): Base trae CSS y JS compartidos que mencionan clases y textos de
// los dos shells.
func shell(t *testing.T, body string) string {
	t.Helper()
	i := strings.Index(body, `class="fs-root"`)
	require.NotEqual(t, -1, i, "no se pintó el shell flamboyant")
	j := strings.Index(body[i:], "<script")
	require.NotEqual(t, -1, j)
	return body[i : i+j]
}

func flags(names ...string) context.Context {
	return feature.Enable(context.Background(), names...)
}

func TestLayout_SinBanderasPintaElShellActual(t *testing.T) {
	body := renderLayout(t, context.Background(), flamboyantCfg())

	assert.Contains(t, body, `id="topbar-user-menu-trigger"`)
	assert.NotContains(t, body, `class="fs-root"`)
	assert.NotContains(t, body, `data-kiban-ui=`)
}

func TestLayout_ConLaBanderaPintaElShellFlamboyant(t *testing.T) {
	body := renderLayout(t, flags(layout.FlagShellFlamboyant), flamboyantCfg())

	assert.Contains(t, body, `class="fs-root"`)
	assert.NotContains(t, body, `id="topbar-user-menu-trigger"`)
}

func TestLayout_ElSwitchSoloApareceConSuBandera(t *testing.T) {
	sinSwitch := renderLayout(t, flags(layout.FlagShellFlamboyant), flamboyantCfg())
	assert.NotContains(t, sinSwitch, `data-kiban-ui=`)
}

func TestLayout_ShellActualOfreceProbarElNuevo(t *testing.T) {
	body := renderLayout(t, flags(layout.FlagShellFlamboyantSwitch), flamboyantCfg())

	assert.Contains(t, body, `data-kiban-ui="flamboyant"`)
	assert.Contains(t, body, "Probar el nuevo diseño")
	assert.NotContains(t, body, `class="fs-root"`)
}

func TestLayout_ShellFlamboyantOfreceVolverAlActual(t *testing.T) {
	body := renderLayout(t, flags(layout.FlagShellFlamboyant, layout.FlagShellFlamboyantSwitch), flamboyantCfg())

	assert.Contains(t, body, `data-kiban-ui="clasico"`)
	assert.Contains(t, body, "Volver al diseño actual")
}

func TestLayout_ElSwitchEscribeLaCookieDelMiddleware(t *testing.T) {
	body := renderLayout(t, context.Background(), flamboyantCfg())

	assert.Contains(t, body, "'"+layout.UICookie+"='")
}

func TestFlamboyant_TextosSiguenElIdiomaDelRequest(t *testing.T) {
	ctx := fsi18n.WithLanguage(flags(layout.FlagShellFlamboyant, layout.FlagShellFlamboyantSwitch), "en")
	body := shell(t, renderLayout(t, ctx, flamboyantCfg()))

	for _, en := range []string{"Pending", "Sign out", "Search spaces", "Back to the current design", `data-fs-i18n-collapse="Collapse menu"`, `data-fs-months="Jan,Feb,`} {
		assert.Contains(t, body, en)
	}
	for _, es := range []string{"Pendientes", "Cerrar sesi", "Buscar espacio", "Volver al dise"} {
		assert.NotContains(t, body, es)
	}
}

func TestFlamboyant_DesarrolladoresSoloConHrefEnDocs(t *testing.T) {
	ctx := flags(layout.FlagShellFlamboyant)
	cfg := flamboyantCfg()

	assert.NotContains(t, shell(t, renderLayout(t, ctx, cfg)), "Desarrolladores")

	cfg.Docs = layout.DocsMenu(ctx, layout.DocsOptions{DevelopersHref: "/kiban-cloud/developers"})
	body := shell(t, renderLayout(t, ctx, cfg))
	assert.Contains(t, body, "Desarrolladores")
	assert.Contains(t, body, `href="/kiban-cloud/developers"`)
}

func TestFlamboyantModulo_NombreSegunIdioma(t *testing.T) {
	es, icon, ok := layout.FlamboyantModulo(context.Background(), "crm")
	require.True(t, ok)
	assert.Equal(t, "CRM", es)
	assert.Equal(t, "contact", icon)

	en, _, _ := layout.FlamboyantModulo(fsi18n.WithLanguage(context.Background(), "en"), "rekon")
	assert.Equal(t, "Payments", en)

	_, _, ok = layout.FlamboyantModulo(context.Background(), "otra")
	assert.False(t, ok)
}

func TestFlamboyantFecha_MesDelIdioma(t *testing.T) {
	d := time.Date(2026, time.August, 28, 9, 5, 0, 0, time.UTC)

	assert.Equal(t, "28 ago 2026 09:05", layout.FlamboyantFecha(context.Background(), d, true))
	assert.Equal(t, "28 Aug 2026", layout.FlamboyantFecha(fsi18n.WithLanguage(context.Background(), "en"), d, false))
}

func TestLocales_MismasClavesEnLosDosIdiomas(t *testing.T) {
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
