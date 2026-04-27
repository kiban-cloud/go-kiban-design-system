// Package table holds the shared list-page primitives every kiban htmx
// project uses: pagination buttons, empty-state card, and (eventually)
// table chrome. Per-row rendering stays per-project since columns differ.
package table

// PaginationConfig drives the shared Pagination component. PageURL is a
// callback so the caller controls how filter state and other context is
// preserved across page navigations — the shared component just renders the
// buttons with the right HTMX attrs.
type PaginationConfig struct {
	Page    int
	HasPrev bool
	HasNext bool

	// PageURL returns the href / hx-get target for the given page number.
	// Caller-provided so each list page can preserve its own filter state.
	PageURL func(page int) string

	// Target is the CSS selector for `hx-target` (e.g. "#customers-content").
	Target string
	// Indicator is the CSS selector for `hx-indicator` (e.g. "#customers-loading").
	Indicator string
}
