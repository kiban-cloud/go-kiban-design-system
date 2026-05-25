package canvas

import (
	"encoding/json"
	"strings"
)

// canvasID returns the wrapper id with a sensible default so callers can
// leave CanvasOptions.ID empty when they only have one canvas on the page.
func canvasID(id string) string {
	if strings.TrimSpace(id) == "" {
		return "kiban-canvas"
	}
	return id
}

// nodeShellClass picks the border + background tint of a node card based on
// status. Unknown values fall back to the OK appearance so an unexpected
// value never breaks the layout.
func nodeShellClass(status string) string {
	base := "group relative w-[280px] rounded-md border bg-white shadow-sm transition-colors"
	switch status {
	case StatusError:
		return base + " border-red-300 bg-red-50/60"
	case StatusNotConfigured:
		return base + " border-amber-300 bg-amber-50/60"
	default:
		return base + " border-kiban-border hover:border-kiban-primary"
	}
}

// statusPillClass styles the small badge rendered at the bottom of a node
// card when status != ok. Returns "" for ok (caller skips rendering).
func statusPillClass(status string) string {
	switch status {
	case StatusError:
		return "text-xs text-red-700 bg-red-100 px-2 py-0.5 rounded"
	case StatusNotConfigured:
		return "text-xs text-amber-700 bg-amber-100 px-2 py-0.5 rounded"
	default:
		return ""
	}
}

// statusPillLabel returns the Spanish label rendered inside the bottom-of-
// card pill. Mirrors the React WorkflooEdit's "Error" / "Pendiente" badges.
func statusPillLabel(status string) string {
	switch status {
	case StatusError:
		return "Error"
	case StatusNotConfigured:
		return "Pendiente"
	default:
		return ""
	}
}

// encodeEdges JSON-encodes the edges list so the templ can drop it into a
// data-edges attribute. The runtime JS parses it back and paints SVG paths.
// Errors collapse to "[]" — an empty edges list — so a serialization bug
// degrades to "no edges drawn" rather than blank-page failure.
func encodeEdges(edges []EdgeOptions) string {
	if len(edges) == 0 {
		return "[]"
	}
	b, err := json.Marshal(edges)
	if err != nil {
		return "[]"
	}
	return string(b)
}
