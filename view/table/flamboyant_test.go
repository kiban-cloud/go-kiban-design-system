package table_test

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

	"github.com/kiban-cloud/go-kiban-design-system/view/avatar"
	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
	"github.com/kiban-cloud/go-kiban-design-system/view/table"
)

var flam = feature.Enable(context.Background(), flamboyant.Flag)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

func withChild(ctx context.Context, html string) context.Context {
	return templ.WithChildren(ctx, templ.Raw(html))
}

func TestTable_SinBanderaEsLaDeSiempre(t *testing.T) {
	body := render(t, context.Background(), table.Table(table.TableConfig{Headers: []string{"Nombre"}, BulkSelect: true}))

	assert.Contains(t, body, "border-t border-kiban-border")
	assert.Contains(t, body, `onchange="var checked=this.checked`)
	assert.NotContains(t, body, "data-kiban-select-all")
}

func TestTable_FlamboyantEncabezadoAzulYSeleccionarTodo(t *testing.T) {
	body := render(t, flam, table.Table(table.TableConfig{Headers: []string{"Cliente", "Creado"}, BulkSelect: true}))

	assert.Contains(t, body, "bg-[#F4F7FF]")
	assert.Contains(t, body, "text-[#0036C4]")
	assert.Contains(t, body, `data-kiban-select-all="ids"`)
	assert.Contains(t, body, "rounded-[5px]", "casilla cuadrada")
	assert.Contains(t, body, ">Cliente<")
}

func TestTable_ColumnaOrdenable(t *testing.T) {
	cfg := table.TableConfig{
		Headers:   []string{"Cliente", "Creado"},
		Sort:      []table.HeaderSort{{}, {Href: "/crm/customers?sort=created_asc", Dir: "desc"}},
		SortAttrs: templ.Attributes{"hx-target": "#lista"},
	}
	body := render(t, flam, table.Table(cfg))

	assert.Contains(t, body, `aria-sort="descending"`)
	assert.Contains(t, body, `href="/crm/customers?sort=created_asc"`)
	assert.Contains(t, body, `hx-target="#lista"`)
	assert.Equal(t, 1, bytes.Count([]byte(body), []byte("aria-sort")), "sólo la columna con Href se ordena")
}

func TestRow_FlamboyantCasillaDentroDeLabel(t *testing.T) {
	body := render(t, withChild(flam, "<td>x</td>"), table.Row("/crm/customers/1", "1"))

	assert.Contains(t, body, `data-href="/crm/customers/1"`)
	assert.Contains(t, body, `name="ids"`)
	assert.Contains(t, body, `value="1"`)
	assert.Contains(t, body, "<label", "el label evita que el clic junto a la casilla abra la fila")
	assert.Contains(t, body, `aria-label="Seleccionar fila"`)
	assert.Contains(t, body, "<td>x</td>")
}

func TestBulkActionBar_FlamboyantConteoEnVivo(t *testing.T) {
	body := render(t, flam, table.BulkActionBar(table.BulkActionBarConfig{CountOne: "1 seleccionado", CountMany: "{n} seleccionados"}))

	assert.Contains(t, body, "data-kiban-bulk-count")
	assert.Contains(t, body, `data-one="1 seleccionado"`)
	assert.Contains(t, body, `data-many="{n} seleccionados"`)
	assert.Contains(t, body, "data-kiban-bulk-clear")
	assert.Contains(t, body, "group-has-[input[name=ids]:checked]/bulk:flex")
}

func TestPagination_FlamboyantConFilasPorPagina(t *testing.T) {
	cfg := table.PaginationConfig{
		Page: 2, HasPrev: true, HasNext: false,
		PageURL:     func(p int) string { return "/l?page=" + string(rune('0'+p)) },
		Target:      "#lista",
		PageSizes:   []int{10, 25, 50},
		PageSize:    25,
		PageSizeURL: "/l?page=1",
	}
	body := render(t, flam, table.Pagination(cfg))

	assert.Contains(t, body, "Página 2")
	assert.Contains(t, body, `name="itemsPerPage"`)
	assert.Contains(t, body, `hx-get="/l?page=1"`)
	assert.Contains(t, body, `<option value="25" selected>`)
	assert.Contains(t, body, `hx-get="/l?page=1"`)
	assert.Contains(t, body, `aria-label="Página siguiente" class="grid h-[30px] w-[30px] cursor-not-allowed`)

	en := render(t, fsi18n.WithLanguage(flam, "en"), table.Pagination(cfg))
	assert.Contains(t, en, "Page 2")
	assert.Contains(t, en, "Rows")
}

func TestEmptyState_AceptaAcciones(t *testing.T) {
	ctx := withChild(context.Background(), `<a href="/nuevo">Nuevo</a>`)
	assert.Contains(t, render(t, ctx, table.EmptyState("Sin clientes", "")), `<a href="/nuevo">Nuevo</a>`)
	assert.Contains(t, render(t, withChild(flam, `<a href="/nuevo">Nuevo</a>`), table.EmptyState("Sin clientes", "")), `<a href="/nuevo">Nuevo</a>`)
}

func TestFrame_SoloPintaTarjetaEnFlamboyant(t *testing.T) {
	assert.Equal(t, "<p>t</p>", render(t, withChild(context.Background(), "<p>t</p>"), table.Frame()))
	assert.Contains(t, render(t, withChild(flam, "<p>t</p>"), table.Frame()), "rounded-[10px]")
}

func TestMainCell(t *testing.T) {
	body := render(t, context.Background(), table.MainCell(table.MainCellOptions{
		Avatar:   avatar.Options{Name: "Acme Servicios", Shape: "square"},
		Title:    "Acme Servicios",
		Subtitle: "cli_01", SubtitleMono: true,
	}))

	assert.Contains(t, body, ">AS")
	assert.Contains(t, body, "Acme Servicios")
	assert.Contains(t, body, "font-mono")
	assert.Contains(t, body, ">cli_01<")
}

func TestLocales_MismasClaves(t *testing.T) {
	assert.Equal(t, localeKeys(t, "locales/es.toml"), localeKeys(t, "locales/en.toml"))
}

func localeKeys(t *testing.T, path string) []string {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	var out []string
	for _, m := range regexp.MustCompile(`(?m)^"([^"]+)"\s*=`).FindAllStringSubmatch(string(b), -1) {
		out = append(out, m[1])
	}
	sort.Strings(out)
	return out
}

func TestRowWith_AtributosEtiquetaYCursor(t *testing.T) {
	body := render(t, withChild(flam, `<td>Ana</td>`), table.RowWith(table.RowOptions{BulkValue: "c1", BulkLabel: "Ana López", Attrs: templ.Attributes{"hx-get": "/crm/customers/c1"}}))

	assert.Contains(t, body, `hx-get="/crm/customers/c1"`)
	assert.Contains(t, body, "cursor-pointer")
	assert.Contains(t, body, `value="c1"`)
	assert.Contains(t, body, `data-label="Ana López"`)
	assert.NotContains(t, body, "data-href")
}

// La casilla de «seleccionar todo» se alinea con las de las filas: celda a la
// izquierda (un <th> centra por defecto) y etiqueta de bloque, que no cuelga
// de la línea base del texto y queda centrada en la altura de la fila.
func TestTable_CasillaDelEncabezadoAlineadaConLasFilas(t *testing.T) {
	head := render(t, flam, table.Table(table.TableConfig{Headers: []string{"Cliente"}, BulkSelect: true}))
	row := render(t, withChild(flam, `<td>Ana</td>`), table.RowWith(table.RowOptions{BulkValue: "c1"}))

	assert.Contains(t, head, `pl-4 pr-0 text-left align-middle">`)
	label := `<label class="-m-1.5 flex w-fit cursor-pointer items-center justify-center p-1.5">`
	assert.Contains(t, head, label)
	assert.Contains(t, row, label)
	assert.NotContains(t, head+row, "inline-grid cursor-pointer")
}
