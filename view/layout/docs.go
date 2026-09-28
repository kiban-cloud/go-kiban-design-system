package layout

import (
	"context"
	"embed"

	"golang.org/x/text/language"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"

	"github.com/kiban-cloud/go-kiban-design-system/view/icons"
)

// The bottom-of-rail links used to be re-declared in every backend's
// view/layout/types.go: the same hrefs six times, under three different i18n
// key schemes, and with "Knowledge Center" written as a bare literal in two
// of them (so it never translated, in either language). DocsMenu replaces all
// six with one list, which is also what makes the group additive — a new
// entry lands in every shell on a DS bump instead of six near-identical PRs.
//
// The labels need their own bundle because the machinery in
// go-kiban-fullstack/pkg/i18n keeps the request language on the context, not
// in the catalog: the language each backend's locale middleware puts there is
// the one this catalog reads back. So the DS owns only this handful of nav
// strings, and every project keeps owning its own copy.

//go:embed locales/*.toml
var localesFS embed.FS

var docsCatalog = fsi18n.NewCatalog(localesFS, "locales", language.Spanish, language.English)

// Knowledge center, API reference and the consultancy booking calendar are
// global kiban destinations, identical for every tool and every environment,
// so they are constants here rather than Config fields nobody would ever set
// differently.
const (
	KnowledgeCenterHref = "https://docs.kiban.com"
	APIDocsHref         = "https://docs.kiban.cloud"
	ConsultancyHref     = "https://calendar.app.google/hrPyz7YiNkKsbjnLA"
)

// DevelopersToolKey is what a page under the developers section sets as
// Config.ToolKey so the rail entry lights up. Exported so the hosting backend
// doesn't have to repeat the string.
const DevelopersToolKey = "developers"

// DefaultDevelopersHref is where the developers section is served. Like
// DefaultLogoHref and DefaultNotificationsBaseURL it is a global,
// host-relative kiban path: one backend (kiban-cloud) serves it and every
// shell on the origin links to it, so no project has to know the route.
//
// It is NOT the default of DocsOptions.DevelopersHref — an empty href has to
// keep meaning "leave the entry out", which is what a deploy needs while the
// section is not serving yet. A shell that wants the entry passes this
// constant.
const DefaultDevelopersHref = "/kiban-cloud/developers"

// DocsOptions carries the parts of the bottom menu that are not the same
// everywhere.
type DocsOptions struct {
	// DevelopersHref points at the developers section. Empty omits the entry
	// entirely, which is the correct state for any deploy where the section
	// isn't serving yet — a rail link to a 404 is worse than no link. It is
	// a path rather than a constant because the section is served by one
	// backend and linked from all of them.
	DevelopersHref string
}

// DocsMenu returns the bottom-of-rail entries, in the order they are shown.
// Every entry is Secondary: this group is reference material sitting under
// the tools the user works in all day, and it should read that way.
//
// Only the developers entry stays in the current tab. The other three leave
// the product, and opening them in place would strand the user outside the
// shell with the app's own back-stack behind them.
func DocsMenu(ctx context.Context, opts DocsOptions) []Tool {
	menu := []Tool{
		{
			Key:       "docs-kiban",
			Icon:      icons.BookOpen,
			Href:      KnowledgeCenterHref,
			External:  true,
			NewTab:    true,
			Secondary: true,
			Label:     docsCatalog.T(ctx, "nav.docs.knowledge"),
		},
		{
			Key:       "docs-api",
			Icon:      icons.FileText,
			Href:      APIDocsHref,
			External:  true,
			NewTab:    true,
			Secondary: true,
			Label:     docsCatalog.T(ctx, "nav.docs.api"),
		},
	}

	if opts.DevelopersHref != "" {
		menu = append(menu, Tool{
			Key:       DevelopersToolKey,
			Icon:      icons.Terminal,
			Href:      opts.DevelopersHref,
			External:  true,
			NewTab:    false,
			Secondary: true,
			Label:     docsCatalog.T(ctx, "nav.docs.developers"),
		})
	}

	return append(menu, Tool{
		Key:       "consultancy",
		Icon:      icons.MessageSquare,
		Href:      ConsultancyHref,
		External:  true,
		NewTab:    true,
		Secondary: true,
		Label:     docsCatalog.T(ctx, "nav.docs.consultancy"),
	})
}
