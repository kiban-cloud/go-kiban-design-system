// Package table holds the shared list-page primitives every kiban htmx
// project uses: pagination buttons, empty-state card, the table chrome
// itself + a Row helper for the standard row markup, and the bulk-action
// bar that appears above the table when rows are selected.
//
// Per-row content (cells: badges, dates, formatted amounts, action menus)
// stays per-project — each consumer fills the row body via templ
// children. The package is intentionally chrome-only.
package table

import "github.com/kiban-cloud/go-kiban-design-system/view/action"

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

// TableConfig drives the Table component. Headers are rendered as plain
// `<th>` text; for richer headers (icons, sort indicators) extend this
// struct with a richer Header type later. BulkSelect=true prepends a
// `select-all` checkbox header column; the matching per-row checkbox
// cell is the caller's responsibility (typically via the Row helper's
// `bulkValue` parameter).
//
// The body is rendered via templ children — caller writes their own
// `<tr>` elements (or, more commonly, `@table.Row(href, bulkValue) {…}`
// for the standard kiban hover/cursor markup).
type TableConfig struct {
	Headers    []string
	BulkSelect bool
}

// BulkActionBarConfig drives the BulkActionBar component — the row of
// actions that appears above the table when at least one row checkbox
// is selected. Visibility is purely CSS-driven via Tailwind's
// `group-has-[…]/bulk:flex` variant: the consumer's enclosing form must
// carry `class="group/bulk"`, the bar carries the variant that watches
// for any `input[name=ids]:checked` inside the same named group.
//
// Layout: optional `Message` on the left, `Actions` on the right
// (rendered through the shared `action.Group`, so primary/secondary
// variants and HTMX wiring work the same as in drawer footers).
type BulkActionBarConfig struct {
	// Message is the muted helper text on the left of the bar
	// (e.g. "Selección múltiple — la acción afecta sólo a los marcados").
	// Optional; pass empty to omit.
	Message string

	// Actions is the right-side action row. Typically a single primary
	// action (e.g. "Eliminar seleccionados" with `Variant:"danger"` +
	// HTMX wiring), but any combination of primary + secondaries is
	// supported.
	Actions action.Group
}
