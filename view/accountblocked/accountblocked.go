// Package accountblocked is the full-page screen a kiban tool shows when the
// space's billing account is suspended or inactive.
//
// Before this screen each tool reacted to a cut account on its own, and each
// one got it differently wrong: crm said "No pudimos cargar los clientes",
// link and klin said the user had no permission, rekon spun forever and
// workfloo showed an empty console. None of them said what was going on.
//
// The page is served by the shared authcookie middleware (it runs right after
// the session is validated), so every tool gets the same screen with no
// wiring. It is deliberately standalone — Base, no Topbar/IconRail — because
// the middleware doesn't have the tool's nav data, and because a cut account
// has nothing to navigate to inside the tool anyway. The way out is the
// kiban-cloud console (billing to pay, the spaces list to switch space), which
// keeps its own session middleware and is never blocked.
package accountblocked

import (
	"context"
	"embed"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"
	domain_core_authorization_model "github.com/kiban-cloud/go-kiban/domain/authorization/model"
)

// Kind is which cut the account is under. It picks the copy and the primary
// action: a suspended account can pay its way back, an inactive one (client
// dado de baja) can only talk to support.
type Kind string

const (
	Suspended Kind = "suspended"
	Inactive  Kind = "inactive"
)

// Global kiban destinations, host-relative like layout.DefaultDevelopersHref:
// kiban-cloud serves them on the same origin as every tool.
const (
	BillingHref  = "/kiban-cloud/billing"
	SpacesHref   = "/kiban-cloud/"
	SupportEmail = "help@kiban.com"
)

// KindFor maps the billing account status go-kiban puts on the Authorization
// to a Kind. ok is false for any status that doesn't cut the service — the
// predicate is go-kiban's, so the screen and the API gate
// (AbortIfBillingAccountBlocked) can never disagree on who is blocked.
func KindFor(billingAccountStatus string) (kind Kind, ok bool) {
	if !domain_core_authorization_model.BillingAccountBlocked(billingAccountStatus) {
		return "", false
	}
	if billingAccountStatus == domain_core_authorization_model.BILLING_ACCOUNT_STATUS_INACTIVE {
		return Inactive, true
	}
	return Suspended, true
}

//go:embed locales/*.toml
var localesFS embed.FS

var catalog = fsi18n.NewCatalog(localesFS, "locales", language.Spanish, language.English)

// WithRequestLanguage puts the request's UI language on ctx for this screen.
//
// The middleware that renders the page runs before some tools' own locale
// middleware (link mounts auth first), so the language can't be assumed to be
// on the context yet. It is resolved the same way the tools resolve it for a
// logged-out request: kiban_lang cookie, then Accept-Language, then Spanish.
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

func key(kind Kind, field string) string {
	return "account_blocked." + string(kind) + "." + field
}
