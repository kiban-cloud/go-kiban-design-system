// Package canvas renders a workfloo-style tree of nodes connected by
// orthogonal edges. The visual is built from plain HTML boxes (full Tailwind
// chrome) laid out by flex/grid, with a single SVG overlay drawn client-side
// for the edges between them. There is no canvas library, no panning, no
// zooming, and no node dragging — the layout is locked by the DOM. Branching
// (e.g. decision-tree children) renders as parallel columns; sequential
// chains render as a single vertical column. Nodes never overlap because the
// underlying layout is flex/grid; the JS only paints the connecting lines.
//
// The package is intentionally small: it exposes the data shapes a caller
// needs (NodeOptions / EdgeOptions / etc.) plus four templ components
// (Canvas / Column / Node / EdgeButton) that compose into the workfloo
// graph. The JS that measures the rendered nodes and draws the SVG paths
// lives in view/layout/base.templ, scoped to [data-kiban-canvas] so multiple
// canvases on the same page stay independent.
package canvas

import "github.com/a-h/templ"

// StatusOK / StatusError / StatusNotConfigured are the three visual states
// of a node. The string values are exported so callers can avoid stringly
// typed code at the call site.
const (
	StatusOK            = "ok"
	StatusError         = "error"
	StatusNotConfigured = "notConfigured"
)

// NodeOptions configures one Node card. The card has a fixed-width shape so
// every node on the canvas lines up vertically regardless of label length.
//
// ID is required — it must be unique within the enclosing Canvas because the
// edge-drawing JS uses it to look up DOM positions when computing SVG paths.
//
// Status defaults to StatusOK when empty. Unknown values fall back to the
// OK appearance so an unexpected value never explodes the layout.
//
// Href makes the entire card a link (renders as <a>). When empty the card is
// a plain <div>; if the caller wants HTMX-driven click behavior they spread
// it via Attrs (mirroring button.Options.Attrs / input.*.attrs).
//
// ActionMenu is rendered in the card's top-right corner. Use it for kebab
// menus (view/menu) or any small affordance the caller wants on each node.
// Pass nil to omit.
type NodeOptions struct {
	ID         string
	Title      string
	Subtitle   string
	Icon       templ.Component
	Status     string
	Href       string
	Attrs      templ.Attributes
	ActionMenu templ.Component
}

// EdgeOptions is one straight-line connection between two nodes. From/To are
// the NodeOptions.ID values of the source and destination cards. Label is
// optional small text rendered at the midpoint of the edge (used for
// decision-tree branch labels like "Sí" / "No"). Variant tints the line:
// "" or "default" → neutral grey, "error" → red (used to mark the
// NextErrorNodeId path in a workfloo), "success" → green (used to mark
// the Verdadero / true branch of a decision-tree step so it pairs visually
// with the "error" Falso branch).
type EdgeOptions struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Label   string `json:"label,omitempty"`
	Variant string `json:"variant,omitempty"`
}

// CanvasOptions is the outer wrapper around a workfloo graph.
//
// ID defaults to "kiban-canvas"; override when multiple canvases share a
// page so each instance's edge SVG stays scoped.
//
// Edges is the full list of lines to draw. The JS reads it from a
// data-edges JSON attribute and paints SVG paths in a single overlay.
//
// Empty drives the empty-state slot. When the caller passes no Node
// children, the canvas renders EmptyMessage centered instead of an empty
// frame. The caller can also supply EmptyAction (a templ.Component, usually
// a button) rendered below the message.
type CanvasOptions struct {
	ID           string
	Edges        []EdgeOptions
	EmptyMessage string
	EmptyAction  templ.Component
}

// EdgeButtonOptions renders a small "+" affordance between two nodes (or
// before the first node, when a canvas is empty). It's the "add a node
// here" trigger used by the workfloo editor.
//
// AriaLabel is required for accessibility (e.g. "Agregar nodo después de X").
// Attrs is spread onto the underlying <button>; callers wire HTMX or onclick
// handlers there.
type EdgeButtonOptions struct {
	AriaLabel string
	Attrs     templ.Attributes
}
