// Package copyable pinta un valor que se copia con un clic: el id de un
// cliente, una referencia, una CLABE. Es de kiban cloud flamboyant (encabezados
// de panel y de la ficha) y no tiene versión de siempre.
//
// El clic lo atiende base.templ ([data-kiban-copy]): copia al portapapeles y
// marca el botón con data-copied un momento, que cambia el ícono por una
// palomita.
package copyable

import (
	"embed"

	"golang.org/x/text/language"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"
)

//go:embed locales/*.toml
var localesFS embed.FS

var catalog = fsi18n.NewCatalog(localesFS, "locales", language.Spanish, language.English)
