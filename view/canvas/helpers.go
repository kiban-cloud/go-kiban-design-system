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

// flNodeClass es la tarjeta de flamboyant según estado y forma.
func flNodeClass(o NodeOptions) string {
	base := "group relative w-[280px] border bg-white transition-[border-color,box-shadow] hover:border-[#0047FF] hover:shadow-[0_4px_14px_-6px_rgba(0,71,255,0.35)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#0047FF]"
	if o.Outcome != "" {
		base += " flex h-10 items-center rounded-full px-3.5"
	} else {
		base += " rounded-lg px-3 py-2.5"
	}
	switch o.Status {
	case StatusNotConfigured:
		return base + " border-dashed border-[#B98CFF]"
	case StatusError:
		return base + " border-[#F3A6A6]"
	}
	return base + " border-[#CDD3DF]"
}

func flOutcomeDot(o string) string {
	switch o {
	case OutcomeSuccess:
		return "bg-[#12B76A]"
	case OutcomeDanger:
		return "bg-[#E5484D]"
	}
	return "bg-[#9096A3]"
}
