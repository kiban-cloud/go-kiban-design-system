// Package breadcrumbs renders a nav path: a list of steps where each
// step is either a link (Href set) or plain text (Href empty, typically
// the current page).
package breadcrumbs

// Item is one step in the breadcrumb trail.
//
// Convention: the last Item has Href = "" so it renders as plain text
// (the current page). Earlier items carry Href to act as links back up
// the trail. If a non-last item has Href = "" it just renders as plain
// text — useful for grouping segments that aren't navigable on their
// own.
type Item struct {
	Label string
	Href  string
}
