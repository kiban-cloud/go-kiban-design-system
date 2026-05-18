// Package table holds the shared list-page primitives every kiban htmx
// project uses: pagination buttons, empty-state card, the table chrome
// itself + a Row helper for the standard row markup, and the bulk-action
// bar that appears above the table when rows are selected.
//
// Per-row content (cells: badges, dates, formatted amounts, action menus)
// stays per-project — each consumer fills the row body via templ
// children. The package is intentionally chrome-only.
package table

import "github.com/kiban-cloud/go-kiban-design-system/view/button"

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
	// NextVariant overrides the styling of the "Siguiente" button.
	// Default ("" / "secondary") renders both buttons in the same
	// outlined secondary style — the project-wide standard, since
	// list paging is a navigation affordance, not a destructive or
	// primary action. Pass "primary" to opt into the legacy
	// kiban-primary solid-blue look (kept for callers that already
	// rely on it). Anterior is always secondary regardless.
	NextVariant string
}

// paginationNextEnabledClass picks the Tailwind class string for the
// active "Siguiente" button. Default ("" / "secondary") = outlined
// secondary, matching the "Anterior" button. "primary" = legacy
// solid kiban-primary with white text, kept for opt-in callers.
func paginationNextEnabledClass(variant string) string {
	if variant == "primary" {
		return "bg-kiban-primary text-white rounded-md px-4 py-2 hover:opacity-90 disabled:opacity-50 disabled:cursor-wait"
	}
	return "bg-white border border-kiban-border rounded-md px-4 py-2 hover:border-kiban-ink3 disabled:opacity-50 disabled:cursor-wait"
}

// paginationNextDisabledClass mirrors paginationNextEnabledClass for
// the disabled (HasNext == false) edge.
func paginationNextDisabledClass(variant string) string {
	if variant == "primary" {
		return "bg-kiban-primary text-white rounded-md px-4 py-2 opacity-40 cursor-not-allowed"
	}
	return "bg-white border border-kiban-border rounded-md px-4 py-2 text-kiban-ink4 cursor-not-allowed"
}

// TableConfig drives the Table component. Headers are rendered as plain
// `<th>` text; for richer headers (icons, sort indicators) extend this
// struct with a richer Header type later. BulkSelect=true prepends a
// `select-all` checkbox header column; the matching per-row checkbox
// cell is the caller's responsibility (typically via the Row helper's
// `bulkValue` parameter).
//
// HeaderAlignRight is an optional parallel slice keyed by header index;
// when `HeaderAlignRight[i]` is true, the i-th `<th>` renders with
// `text-right` instead of the default `text-left`. Used for action
// columns whose body cells are right-aligned (kebab menus, button
// rows). A nil or shorter-than-Headers slice falls back to left-align
// — pre-existing callers don't need to change.
//
// The body is rendered via templ children — caller writes their own
// `<tr>` elements (or, more commonly, `@table.Row(href, bulkValue) {…}`
// for the standard kiban hover/cursor markup).
type TableConfig struct {
	Headers          []string
	HeaderAlignRight []bool
	BulkSelect       bool
	// HeaderNoWrap forces every `<th>` to render on a single line via
	// `whitespace-nowrap`. Default (`false`) keeps the historical
	// behaviour where multi-word column titles ("Fecha de creación")
	// can wrap when the column is narrow. Useful in dense admin
	// tables where wrapped headers misalign visually with no-wrap
	// body cells.
	HeaderNoWrap bool
}

// HeaderAlignClass returns the alignment class for a given header
// column index, picking `text-right` when the parallel slice marks it
// or `text-left` otherwise. Lives next to the struct so the templ can
// stay declarative.
func (c TableConfig) HeaderAlignClass(i int) string {
	if i < len(c.HeaderAlignRight) && c.HeaderAlignRight[i] {
		return "text-right"
	}
	return "text-left"
}

// HeaderWrapClass returns `whitespace-nowrap` when [TableConfig.HeaderNoWrap]
// is true and an empty string otherwise — concatenated into each
// `<th>`'s class list by the Table templ.
func (c TableConfig) HeaderWrapClass() string {
	if c.HeaderNoWrap {
		return "whitespace-nowrap"
	}
	return ""
}

// BulkActionBarConfig drives the BulkActionBar component — the row of
// actions that appears above the table when at least one row checkbox
// is selected. Visibility is purely CSS-driven via Tailwind's
// `group-has-[…]/bulk:flex` variant: the consumer's enclosing form must
// carry `class="group/bulk"`, the bar carries the variant that watches
// for any `input[name=ids]:checked` inside the same named group.
//
// Layout: optional `Message` on the left, `Actions` on the right
// (rendered through the shared `button.Group`, so primary/secondary
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
	Actions button.Group
}
