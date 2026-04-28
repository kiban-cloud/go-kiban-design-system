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
// `Attrs` is the optional HTMX escape hatch. When set, attributes are
// spread directly onto the `<a>` element so consumers can switch a tab
// from "full-page nav" (default) to "swap a partial in place" without
// the design system needing to know about HTMX. Same pattern used by
// `action.Action`.
type TabItem struct {
	Key   string
	Label string
	Href  string
	Attrs templ.Attributes
}
