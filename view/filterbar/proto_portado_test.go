package filterbar_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/filterbar"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// Piezas portadas de kiban-proto: apagado flamboyant, lo de siempre;
// prendido, lo del proto.

var encendido = feature.Enable(context.Background(), flamboyant.Flag)

func pintar(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

var _ = strings.Contains

func TestBar_SinNombreDeBusqueda_ApagadoLaPintaPrendidoNo(t *testing.T) {
	cfg := filterbar.Config{ID: "f", Action: "/x", Target: "#t"}
	assert.Contains(t, pintar(t, context.Background(), filterbar.Bar(cfg)), `type="search"`, "apagado, la barra de siempre")
	assert.NotContains(t, pintar(t, encendido, filterbar.Bar(cfg)), `type="search"`)
	cfg.Search.Name = "q"
	assert.Contains(t, pintar(t, encendido, filterbar.Bar(cfg)), `type="search"`)
}

func TestBar_HxSyncSoloEnFlamboyant(t *testing.T) {
	cfg := filterbar.Config{ID: "f", Action: "/x", Search: filterbar.Search{Name: "q"}}
	assert.NotContains(t, pintar(t, context.Background(), filterbar.Bar(cfg)), "hx-sync")
	assert.Contains(t, pintar(t, encendido, filterbar.Bar(cfg)), `hx-sync="this:queue last"`)
}

func TestBar_FiltroDeTextoLibre(t *testing.T) {
	f := filterbar.Filter{Name: "id", Label: "ID", Text: &filterbar.Text{Name: "id", Placeholder: "Pega el ID", Apply: "Aplicar", Mono: true}}
	cfg := filterbar.Config{ID: "f", Action: "/x", MenuClear: "Quitar filtro", Filters: []filterbar.Filter{f}}
	body := pintar(t, encendido, filterbar.Bar(cfg))
	assert.Contains(t, body, "data-kiban-filter-free")
	assert.Contains(t, body, "data-kiban-txt-apply")
	assert.Regexp(t, `<input type="hidden" name="id" value="" disabled`, body, "sin valor, el oculto no viaja")
	assert.Contains(t, body, "font-mono")

	cfg.Filters[0].Text.Value = "abc123"
	body = pintar(t, encendido, filterbar.Bar(cfg))
	assert.Contains(t, body, `ID: <b class="font-semibold">abc123</b>`)
	assert.NotRegexp(t, `name="id" value="abc123" disabled`, body)
}
