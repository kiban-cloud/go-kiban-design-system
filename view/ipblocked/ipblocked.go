// Package ipblocked is the full-page screen a kiban tool shows when the
// space's IP access list (Console) doesn't allow the browser's IP for that
// tool.
//
// Like accountblocked, it is served by the shared authcookie middleware right
// after the session is validated, so every tool gets it with no wiring, and it
// is standalone (Base, no Topbar/IconRail) because the middleware doesn't
// have the tool's nav data. The way out is the kiban-cloud console, which is
// never restricted by IP: from there the user can switch space, and an admin
// can fix the list.
package ipblocked

import (
	"context"
	"embed"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"
)

// Host-relative kiban-cloud destinations, served on the same origin as every
// tool.
const (
	SpacesHref   = "/kiban-cloud/"
	IpAccessHref = "/kiban-cloud/developers/ipaccess?lista=consola"
)

//go:embed locales/*.toml
var localesFS embed.FS

var catalog = fsi18n.NewCatalog(localesFS, "locales", language.Spanish, language.English)

// WithRequestLanguage resolves the UI language the same way accountblocked
// does: kiban_lang cookie, then Accept-Language, then Spanish. The middleware
// runs before some tools' locale middleware.
func WithRequestLanguage(c *gin.Context) context.Context {
	lang := catalog.DefaultLanguage()
	if cookie, err := c.Cookie(fsi18n.CookieLang); err == nil && catalog.IsSupported(cookie) {
		lang = cookie
	} else if header := c.GetHeader("Accept-Language"); header != "" {
		lang = catalog.Match(header)
	}
	return fsi18n.WithLanguage(c.Request.Context(), lang, catalog.DefaultLanguage())
}

func t(ctx context.Context, id string) string {
	return catalog.T(ctx, id)
}

func tf(ctx context.Context, id string, data map[string]any) string {
	return catalog.Tf(ctx, id, data)
}

// ToolLabel is how the tool reads on screen; crm shows as "clientes", like in
// the rail.
func ToolLabel(ctx context.Context, tool string) string {
	key := "ip_blocked.tool." + strings.ToLower(tool)
	if label := t(ctx, key); label != key {
		return label
	}
	return strings.ToLower(tool)
}
