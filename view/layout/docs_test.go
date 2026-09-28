package layout_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"

	"github.com/kiban-cloud/go-kiban-design-system/view/icons"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

func TestDocsMenu_OmitsDevelopersWithoutHref(t *testing.T) {
	menu := layout.DocsMenu(context.Background(), layout.DocsOptions{})

	keys := make([]string, len(menu))
	for i, tool := range menu {
		keys[i] = tool.Key
	}
	assert.Equal(t, []string{"docs-kiban", "docs-api", "consultancy"}, keys)
}

func TestDocsMenu_PlacesDevelopersBeforeConsultancy(t *testing.T) {
	menu := layout.DocsMenu(context.Background(), layout.DocsOptions{
		DevelopersHref: "/kiban-cloud/developers",
	})

	keys := make([]string, len(menu))
	for i, tool := range menu {
		keys[i] = tool.Key
	}
	assert.Equal(t, []string{"docs-kiban", "docs-api", "developers", "consultancy"}, keys)
}

// The whole point of the group is that it reads quieter than the tools, and
// that only the developers entry keeps the user in the tab. Both are easy to
// lose when someone adds a fifth entry by copying a neighbour.
func TestDocsMenu_EveryEntryIsSecondaryAndOnlyDevelopersStaysInTab(t *testing.T) {
	menu := layout.DocsMenu(context.Background(), layout.DocsOptions{
		DevelopersHref: "/kiban-cloud/developers",
	})
	require.Len(t, menu, 4)

	for _, tool := range menu {
		assert.True(t, tool.Secondary, "%s must render at the secondary weight", tool.Key)
		assert.Equal(t, tool.Key != layout.DevelopersToolKey, tool.NewTab,
			"%s has the wrong new-tab behaviour", tool.Key)
	}
}

// Two backends used to ship "Knowledge Center" as a bare Go literal, so the
// label was Spanish-only no matter what the user picked. Now that the list is
// in the DS, the regression would be silent in six shells at once.
func TestDocsMenu_LabelsFollowTheRequestLanguage(t *testing.T) {
	es := layout.DocsMenu(fsi18n.WithLanguage(context.Background(), "es"), layout.DocsOptions{
		DevelopersHref: "/kiban-cloud/developers",
	})
	en := layout.DocsMenu(fsi18n.WithLanguage(context.Background(), "en"), layout.DocsOptions{
		DevelopersHref: "/kiban-cloud/developers",
	})

	assert.Equal(t, []string{"Knowledge Center", "Documentación API", "Desarrolladores", "Consultoría"},
		labels(es))
	assert.Equal(t, []string{"Knowledge Center", "API Documentation", "Developers", "Consulting"},
		labels(en))
}

func labels(menu []layout.Tool) []string {
	out := make([]string, len(menu))
	for i, tool := range menu {
		out[i] = tool.Label
	}
	return out
}

// IconRail renders both groups through toolIcon, so the styling difference
// has to come from the Tool itself. Before Secondary/NewTab existed the rail
// hard-coded target=_blank for the whole bottom group, which would have sent
// the developers section to a new tab.
func TestIconRail_SecondaryEntriesAreQuieterAndOnlyNewTabOnesOpenOut(t *testing.T) {
	html := renderIconRail(t, layout.Config{
		ToolKey: "home",
		Tools: []layout.Tool{
			{Key: "home", Icon: icons.Home, Href: "/kiban-cloud/", Label: "Home"},
		},
		Docs: layout.DocsMenu(context.Background(), layout.DocsOptions{
			DevelopersHref: "/kiban-cloud/developers",
		}),
	})

	assert.Contains(t, html, `href="https://docs.kiban.com"`)
	assert.Contains(t, html, `href="/kiban-cloud/developers"`)

	// The tools keep the 14px base; the bottom group drops to 13px.
	assert.Contains(t, html, "py-2 text-sm")
	assert.Contains(t, html, "py-1.5 text-[13px]")
	assert.Contains(t, html, "text-kiban-ink3")

	// Three reference links open out; the developers entry does not, so the
	// rail emits exactly three targets.
	assert.Equal(t, 3, strings.Count(html, `target="_blank"`))
	assert.Equal(t, 3, strings.Count(html, `rel="noopener noreferrer"`))

	developers := html[strings.Index(html, `href="/kiban-cloud/developers"`):]
	developers = developers[:strings.Index(developers, "</a>")]
	assert.NotContains(t, developers, `target="_blank"`)
}

func renderIconRail(t *testing.T, cfg layout.Config) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, layout.IconRail(cfg).Render(context.Background(), &buf))
	return buf.String()
}
