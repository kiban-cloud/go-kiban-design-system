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
	"context"
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

// FieldOptions configura Field. Label es el nombre accesible del botón de
// copiar («Copiar link»); vacío, «Copiar» y el valor. Mono pinta el valor en
// monoespaciado (una CLABE, un token).
type FieldOptions struct {
	Value string
	Label string
	Muted bool
	Mono  bool
}

func (o FieldOptions) label(ctx context.Context) string {
	if o.Label != "" {
		return o.Label
	}
	return catalog.Tf(ctx, "copyable.copy", map[string]any{"Value": o.Value})
}

// InlineOptions configura Inline. Label va antes del valor, en el mismo
// texto («ID 991ab76c…»); vacío, sólo el valor. Sep pone adelante el punto
// gris que lo separa de lo anterior en el subtítulo de un panel. CopyLabel es
// el nombre accesible del botón («Copiar ID»); vacío, «Copiar» y el valor.
type InlineOptions struct {
	Label     string
	Value     string
	CopyLabel string
	Sep       bool
}

func (o InlineOptions) text() string {
	if o.Label != "" {
		return o.Label + " " + o.Value
	}
	return o.Value
}

func (o InlineOptions) copyLabel(ctx context.Context) string {
	if o.CopyLabel != "" {
		return o.CopyLabel
	}
	return catalog.Tf(ctx, "copyable.copy", map[string]any{"Value": o.Value})
}
