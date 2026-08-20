package tabs_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/tabs"
	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, items []tabs.TabItem, activeKey string) string {
	t.Helper()
	var buf bytes.Buffer
	err := tabs.Strip(items, activeKey).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

// TestStrip_ContainerBaselineMatchesActiveBorderWidth is the visual
// contract test for the polished SaaS-style strip: the container's
// bottom rule and the active tab's underline must share the same border
// width (1px each — `border-b` Tailwind class) so they visually align.
// If anyone bumps one to `border-b-2` without bumping the other, this
// guard catches it before merge.
func TestStrip_ContainerBaselineMatchesActiveBorderWidth(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
		{Key: "b", Label: "Orange", Href: "/b"},
	}, "a")

	// Container baseline: 1-pixel `border-b` (same width as the active
	// underline below). The container may carry extra utilities (e.g.
	// `overflow-x-auto` so a long tab row scrolls on narrow viewports),
	// so match the border classes without pinning the exact class string.
	assert.True(t, strings.Contains(body, `class="border-b border-kiban-border`))
	// Active inner span: 1-pixel `border-b` (NOT border-b-2 — same width
	// as the container above) coloured `border-kiban-primary`. Asserted
	// as separate class-presence checks because the inner span also
	// carries `pb-0.5` (small gap between text and the underline) which
	// otherwise breaks a contiguous substring match.
	assert.True(t, strings.Contains(body, "border-b"))
	assert.True(t, strings.Contains(body, "border-kiban-primary"))
	assert.False(t, strings.Contains(body, "border-b-2"))
	// Inactive label inner span: same 1-pixel reserved height with a
	// transparent colour so the row doesn't shift when active changes.
	assert.True(t, strings.Contains(body, "border-transparent"))
}

func TestStrip_TightGap(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
	}, "a")
	// Tight 12px (gap-3) between tabs — replacing the prior gap-2.
	assert.True(t, strings.Contains(body, "flex gap-3"))
}

// TestStrip_LabelHasGapAboveBorder guards the breathing room between
// the label text and the underline / baseline. The class lives on the
// inner span — the same element that carries the border — so the gap
// applies to both the active blue underline and the inactive
// (transparent) reserved space.
func TestStrip_LabelHasGapAboveBorder(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
	}, "a")
	// Span class list looks like "inline-block border-b pb-2 border-…".
	assert.True(t, strings.Contains(body, "border-b pb-2"), "inner span must carry pb-2 between text and border")
}

// TestStrip_ActiveUnderlineOverlapsContainerBaseline guards the
// 1-pixel overlap mechanism. The nav must carry `-mb-px` so the
// active tab's `border-b border-kiban-primary` sits exactly on top of
// the container's `border-b border-kiban-border` instead of floating
// above it. If the negative margin disappears, the blue underline
// drifts above the gray rule and the tabs look detached from their
// rail. Tab anchors must also use `pt-2` (padding above the label),
// not `pb-2` — padding below would push the inner span up off the
// baseline and break the overlap even with `-mb-px` in place.
func TestStrip_ActiveUnderlineOverlapsContainerBaseline(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
	}, "a")
	assert.True(t, strings.Contains(body, "-mb-px"), "nav must pull up by border width to overlap container baseline")
	// Anchor class list looks like "inline-flex items-center gap-2 pt-2 text-sm …".
	// We anchor the assertions to `gap-2` so they target the anchor and
	// not the inner span (which legitimately carries pb-2 between the
	// label and the border).
	assert.True(t, strings.Contains(body, "gap-2 pt-2"), "anchor must use pt-2 so the inner span sits at the bottom edge")
	assert.False(t, strings.Contains(body, "gap-2 pb-2"), "pb-2 on the anchor would push the underline up off the container baseline")
}

func TestStrip_ActiveLabelEmphasized(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
		{Key: "b", Label: "Orange", Href: "/b"},
	}, "a")
	// Active anchor styling: kiban-primary text + medium weight.
	assert.True(t, strings.Contains(body, "text-kiban-primary font-medium"))
	// Inactive anchor styling: muted ink3 + hover restores ink.
	assert.True(t, strings.Contains(body, "text-kiban-ink3 hover:text-kiban-ink"))
}

func TestStrip_AttrsSpreadOnAnchor(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a", Attrs: templ.Attributes{
			"hx-get":      "/tab/a",
			"hx-target":   "#content",
			"hx-push-url": "/a",
		}},
	}, "a")
	assert.True(t, strings.Contains(body, `hx-get="/tab/a"`))
	assert.True(t, strings.Contains(body, `hx-target="#content"`))
	assert.True(t, strings.Contains(body, `hx-push-url="/a"`))
}

func TestStrip_RendersIconBeforeLabel(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a", Icon: templ.Raw(`<svg data-test-icon></svg>`)},
	}, "a")
	// Icon is inside a wrapper span before the label span.
	assert.True(t, strings.Contains(body, "data-test-icon"))
	iconIdx := strings.Index(body, "data-test-icon")
	labelIdx := strings.Index(body, "Apple")
	assert.True(t, iconIdx >= 0 && labelIdx >= 0 && iconIdx < labelIdx, "icon must render before label")
}

func TestStrip_CountBadgeOnlyWhenHasCount(t *testing.T) {
	// HasCount=false: badge is omitted even if Count is non-zero.
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Inbox", Href: "/a", Count: 5},
	}, "a")
	assert.False(t, strings.Contains(body, "rounded-full"))
	// HasCount=true with Count=0: badge IS rendered (this is the
	// reason the explicit toggle exists — let "Inbox (0)" be a real
	// state without a sentinel-vs-zero ambiguity).
	body = render(t, []tabs.TabItem{
		{Key: "a", Label: "Inbox", Href: "/a", Count: 0, HasCount: true},
	}, "a")
	assert.True(t, strings.Contains(body, "rounded-full"))
	assert.True(t, strings.Contains(body, ">0<"))
	// Active state colours the badge with the primary tint.
	body = render(t, []tabs.TabItem{
		{Key: "a", Label: "Inbox", Href: "/a", Count: 7, HasCount: true},
	}, "a")
	assert.True(t, strings.Contains(body, "bg-kiban-primary-soft text-kiban-primary"))
	// Inactive tab gets the muted variant.
	body = render(t, []tabs.TabItem{
		{Key: "a", Label: "Inbox", Href: "/a", Count: 7, HasCount: true},
		{Key: "b", Label: "Sent", Href: "/b"},
	}, "b")
	assert.True(t, strings.Contains(body, "bg-kiban-surface text-kiban-ink3"))
}

func TestStrip_DisabledIsNotClickable(t *testing.T) {
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
		{Key: "b", Label: "Próximamente", Href: "/b", Disabled: true, Title: "Aún no disponible"},
	}, "a")
	// Disabled tab: no href, aria-disabled, pointer-events-none, opacity.
	assert.True(t, strings.Contains(body, `aria-disabled="true"`))
	assert.True(t, strings.Contains(body, "pointer-events-none"))
	assert.True(t, strings.Contains(body, "opacity-50"))
	// Title surfaces the reason as a native tooltip.
	assert.True(t, strings.Contains(body, `title="Aún no disponible"`) ||
		strings.Contains(body, `title="A&uacute;n no disponible"`))
	// The first tab still rendered with its href intact (regression
	// guard: disabled handling must not bleed into siblings).
	assert.True(t, strings.Contains(body, `href="/a"`))
}

func TestStrip_NoActiveKeyHighlightsNothing(t *testing.T) {
	// If activeKey doesn't match any item, the strip still renders
	// without crashing and no item is highlighted.
	body := render(t, []tabs.TabItem{
		{Key: "a", Label: "Apple", Href: "/a"},
		{Key: "b", Label: "Orange", Href: "/b"},
	}, "")
	assert.False(t, strings.Contains(body, "text-kiban-primary font-medium"))
	assert.False(t, strings.Contains(body, "border-b border-kiban-primary"))
	assert.True(t, strings.Contains(body, "Apple"))
	assert.True(t, strings.Contains(body, "Orange"))
}

// El riel gris y el scroll NO pueden vivir en la misma caja: un elemento con
// overflow recorta su contenido al padding box, que excluye su propio borde, y
// ahí es exactamente donde cae el subrayado del tab activo — con lo que
// desaparece. El scroller se solapa con el riel por su padding para que el
// subrayado de 1px se pinte encima del riel de 1px.
func TestStrip_RailAndScrollerAreSeparateBoxes(t *testing.T) {
	body := render(t, []tabs.TabItem{{Key: "a", Label: "Uno", Href: "/a"}}, "a")

	assert.Contains(t, body, `<div class="border-b border-kiban-border">`,
		"el riel va en una caja que no recorta")
	assert.Contains(t, body, `<div class="overflow-x-auto pb-px -mb-px">`,
		"el scroller aloja el subrayado en su padding y se solapa con el riel")
	assert.NotContains(t, body, `border-b border-kiban-border overflow-x-auto`,
		"juntarlos recorta el subrayado del tab activo")
}
