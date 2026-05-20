// Package tabs holds the kiban "in-page tabs" primitive — a horizontal
// strip of links that switches between views inside a single page. It's
// intentionally distinct from `view/layout/SubNav` (the level-2
// navigation that lives in the layout shell): SubNav owns the pill-style
// app navigation chrome, tabs own the underline-style "switch the view
// inside this page" pattern. Keeping the two visually distinct prevents
// users from confusing in-page state with global navigation.
//
// Each TabItem renders as an `<a href={Href}>` so plain browser
// navigation works out of the box. For HTMX-driven tabs that swap a
// content panel in place without a full nav, populate `Attrs` with the
// usual `hx-get` / `hx-target` / `hx-swap` / `hx-push-url` and HTMX
// will own the click; the `Href` attribute then acts as the no-JS
// fallback (or as the URL the browser pushes when `hx-push-url="true"`).
//
// Active state is caller-driven: pass the active item's `Key` as
// `Strip(items, activeKey)`. If no item matches, no tab is highlighted
// — useful while hydrating a fresh page state where no specific tab is
// active yet.
package tabs

import "github.com/a-h/templ"

// TabItem is one entry in a tab strip. `Key` is a stable identifier
// matched against the strip's `activeKey` argument — typically a short
// slug like "summary", "history", "settings". `Label` is the visible
// text. `Href` is the canonical URL the tab points at (browser
// navigates here when no HTMX wiring is present, or HTMX pushes it via
// `hx-push-url="true"` when configured).
//
// Optional affordances:
//
//   - `Icon` is an optional templ component (typically from
//     `view/icons`) rendered to the left of the label. Same convention
//     as `button.Options.IconComponent`.
//
//   - `Count` paired with `HasCount=true` renders a small pill badge
//     to the right of the label. The two-field pattern (rather than
//     just `Count int > 0`) lets a real `0` show up in the badge — for
//     things like "Inbox (0)" — without a sentinel-versus-zero
//     ambiguity.
//
//   - `Disabled=true` greys out the tab and intercepts pointer events
//     so the click target is a no-op. The tab is rendered as an `<a>`
//     without `href` (anchors don't have a real `disabled` attribute,
//     so `pointer-events-none` + `aria-disabled="true"` is the web
//     pattern). Pair with `Title` to surface the reason in a native
//     tooltip.
//
//   - `Title` sets the `title=""` attribute. Useful when `Disabled`
//     and the user needs a hint about why.
//
// `Attrs` is the optional HTMX escape hatch. When set, attributes are
// spread directly onto the `<a>` element so consumers can switch a tab
// from "full-page nav" (default) to "swap a partial in place" without
// the design system needing to know about HTMX. Same pattern used by
// `button.Options`.
type TabItem struct {
	Key      string
	Label    string
	Href     string
	Icon     templ.Component
	Count    int
	HasCount bool
	Disabled bool
	Title    string
	Attrs    templ.Attributes
}

// PanelConfig is the input for Panel — the in-page tabs primitive
// that switches between pre-rendered content blocks WITHOUT a
// network round-trip (unlike Strip, which navigates via anchor
// hrefs). Use Panel when all the tab bodies are cheap to render
// upfront and the user expects instant switching; use Strip when
// tabs map to distinct URLs / lazy-loaded views.
//
//   - ID is unique on the page; the JS in view/layout/base.templ
//     keys off it to scope click handling to this panel only, so
//     two Panels on the same page don't interfere.
//   - ActiveKey is the tab that's visible on first paint. If no
//     tab matches, the panel renders with nothing visible — the
//     first click on the strip will pick whichever tab the user
//     selects.
//   - Tabs lists the strip entries in display order. Each header
//     is `{Key, Label}` — the Body call site for that key supplies
//     the actual content.
type PanelConfig struct {
	ID        string
	ActiveKey string
	Tabs      []TabHeader
}

// TabHeader is one entry in the Panel's strip. Intentionally
// minimal (no Href / Icon / Count / Disabled) — Panel is for
// CSS-only in-page switching, so the navigation affordances that
// matter for Strip don't apply here.
type TabHeader struct {
	Key   string
	Label string
}
