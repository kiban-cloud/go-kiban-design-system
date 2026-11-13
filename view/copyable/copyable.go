// Package copyable pinta un valor que se copia con un clic: el id de un
// cliente, una referencia, una CLABE. Es de kiban cloud flamboyant (encabezados
// de panel y de la ficha) y no tiene versión de siempre.
//
// El clic lo atiende base.templ ([data-kiban-copy]): copia al portapapeles y
// marca el botón con data-copied un momento, que cambia el ícono por una
// palomita. Con IDWith{Toast} también muestra el aviso flotante
// (window.kibanToast, en base.templ).
package copyable

import (
	"embed"

	"golang.org/x/text/language"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"
)

//go:embed locales/*.toml
var localesFS embed.FS

var catalog = fsi18n.NewCatalog(localesFS, "locales", language.Spanish, language.English)

// Options configura IDWith. Value es lo que se ve y se copia; Toast, si
// viene, es el aviso flotante que confirma la copia.
type Options struct {
	Value string
	Toast string
}
