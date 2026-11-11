package table

import (
	"embed"

	"github.com/a-h/templ"
	"strconv"

	"golang.org/x/text/language"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"

	"github.com/kiban-cloud/go-kiban-design-system/view/avatar"
)

// La tabla flamboyant (lista de clientes de kiban-proto): tarjeta con borde
// (Frame), encabezado azul claro con columnas ordenables, filas con hover azul
// y checkbox cuadrado, franja azul de selección y pie con paginación y filas
// por página. Cada componente elige su variante con flamboyant.On(ctx); con la
// bandera apagada se ven como siempre.

//go:embed locales/*.toml
var localesFS embed.FS

var catalog = fsi18n.NewCatalog(localesFS, "locales", language.Spanish, language.English)

// HeaderSort hace ordenable una columna. Va en TableConfig.Sort con el mismo
// índice que Headers; un HeaderSort vacío deja la columna sin ordenar.
type HeaderSort struct {
	// Href pide la lista con el orden que sigue al clic: el inverso si ya
	// es la columna del orden, el natural de la columna si no.
	Href string
	// Dir es el orden actual de esta columna: "asc", "desc" o "" si la lista
	// no está ordenada por ella.
	Dir string
}

func (c TableConfig) sortAt(i int) (HeaderSort, bool) {
	if i < len(c.Sort) && c.Sort[i].Href != "" {
		return c.Sort[i], true
	}
	return HeaderSort{}, false
}

func ariaSort(s HeaderSort) string {
	switch s.Dir {
	case "asc":
		return "ascending"
	case "desc":
		return "descending"
	}
	return "none"
}

// MainCellOptions es la celda principal de una fila: avatar, nombre y una
// línea secundaria (referencia, correo, id).
type MainCellOptions struct {
	Avatar   avatar.Options
	Title    string
	Subtitle string
	// SubtitleMono pinta la línea secundaria en monoespaciado (ids).
	SubtitleMono bool
}

func itoa(n int) string { return strconv.Itoa(n) }

// RowOptions configura RowWith.
type RowOptions struct {
	Href      string
	BulkValue string
	// BulkLabel es el nombre de la fila en la confirmación de borrado.
	BulkLabel string
	// BulkDisabled, con BulkValue vacío, pinta la casilla deshabilitada con
	// este texto de ayuda (por qué esa fila no se puede elegir: «Los casos
	// de kiban no se pueden eliminar»), así la columna no se corre.
	BulkDisabled string
	Attrs     templ.Attributes
}
