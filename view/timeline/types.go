// Package timeline renders a vertical list of dated events with a
// status-coloured dot per row. Generic across kiban backends:
// status logs (NIP send/validation), audit trails (workfloo
// execution events), payment lifecycle ticks (rekon), delivery
// tracking (klin).
//
// Each row is `[dot] [label]                              [date]`,
// flexed so the date sticks to the right and the label takes the
// middle. Date is optional — when empty, only the dot + label
// render.
package timeline

// Event is one row in the timeline. Kind drives the dot colour;
// Label is the user-visible status copy (already-localised by the
// caller); Date is the pre-formatted timestamp string (caller
// owns timezone + locale formatting).
type Event struct {
	Label string
	Kind  string // "success" / "warning" / "info" / "danger" / "default"
	Date  string // optional; empty hides the date column
}

// Kind constants — exported so callers can avoid stringly-typed
// values at the call site.
const (
	KindSuccess = "success"
	KindWarning = "warning"
	KindInfo    = "info"
	KindDanger  = "danger"
	KindDefault = "default"
)
