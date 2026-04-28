// Package card holds the kiban "card" chrome primitives — the bordered,
// padded white surfaces that group form rows, detail blocks, and action
// rails on every backoffice page.
//
// Two components, two roles:
//
//   - Card: the bare chrome wrapper. Renders the white surface, kiban
//     border, rounded corners, and standard padding. Accepts any content
//     via templ children. Pass a `variant` to swap the border colour for
//     status-coloured cards (success/warning/danger/info) — used for
//     full-card warnings like CSV-upload duplicate/error summaries.
//
//   - Section: a heading + subtitle block paired with body content,
//     designed to live inside a Card. Multiple Sections in one Card
//     produce the divided-card layout (border-t between sub-sections)
//     used on long forms and detail pages. A single-section card just
//     uses `Card` directly — no need for `Section` when there's only one.
//
// Variant strings:
//
//	""                           // default — kiban-border
//	"success"                    // emerald-200
//	"warning"                    // amber-200
//	"danger"                     // red-200
//	"info"                       // kiban-primary-tinted
//
// Variant is intentionally a `string` (not a typed enum) for parity with
// the rest of the input/badge/flash packages and so callers can pass an
// empty string to mean "default" without an extra import.
package card
