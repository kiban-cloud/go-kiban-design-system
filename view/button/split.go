package button

import (
	"strings"

	"github.com/a-h/templ"
)

// SplitItem is one row inside a [Split] dropdown: a variant of the
// primary action rather than an unrelated command.
//
// Items submit their form natively — no JS. Give the item a Name/Value
// pair and the server reads which variant was picked from the posted
// body (`_action=save_and_close`). Form names the target form when the
// split button lives outside it.
type SplitItem struct {
	Label string
	Name  string
	Value string
	Form  string
	// OnClick is raw inline JS, mirroring Options.OnClick. The dropdown
	// appends its own close call, so consumers don't have to. Returning
	// false from the expression aborts the submit — that's how a
	// confirm() guard is wired.
	OnClick string
	// Variant: "" or "default" → kiban-ink text; "danger" → red text.
	Variant  string
	Icon     templ.Component
	Disabled bool
	Attrs    templ.Attributes
}

// SplitOptions drives [Split]. ID is required and must be unique per
// page: it keys the shared dropdown handlers in view/layout/base.templ,
// exactly like menu.Config.ID.
//
// Primary is the left segment — the default action, run when the user
// clicks the button body instead of opening the dropdown. It takes a
// full Options so callers keep submit/form/disabled behaviour.
type SplitOptions struct {
	ID      string
	Primary Options
	Items   []SplitItem
	// TriggerAriaLabel labels the caret segment for screen readers.
	// Falls back to the primary label when empty.
	TriggerAriaLabel string
}

func (o SplitOptions) hasItems() bool { return len(o.Items) > 0 }

// primaryOptions strips the right-side rounding so the two segments read
// as one control, and forces the shared variant when unset.
func (o SplitOptions) primaryOptions() Options {
	p := o.Primary
	p.ExtraClass = joinClasses("rounded-r-none", p.ExtraClass)
	return p
}

func (o SplitOptions) triggerAriaLabel() string {
	if strings.TrimSpace(o.TriggerAriaLabel) != "" {
		return o.TriggerAriaLabel
	}
	return o.Primary.Label
}

// triggerClass mirrors the primary variant so both halves share a fill,
// minus the left rounding and with a divider between the segments.
func (o SplitOptions) triggerClass() string {
	return BuildClass(o.Primary.EffectiveVariant(),
		"rounded-l-none border-l border-white/25 px-2")
}

// splitTriggerAttrs wires the caret to the shared menu handlers. Same
// id/aria contract as menu.Config so base.templ needs no new JS.
func splitTriggerAttrs(id string) templ.Attributes {
	return templ.Attributes{
		"id":            id + "-trigger",
		"aria-haspopup": "true",
		"aria-expanded": "false",
		"aria-controls": id + "-panel",
		"onclick":       "window.kibanToggleMenu('" + id + "')",
	}
}

// resolvedAttrs merges the caller's Attrs with the close-on-click
// behaviour, mirroring menu.MenuItem.resolvedAttrs.
func (it SplitItem) resolvedAttrs(splitID string) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range it.Attrs {
		out[k] = v
	}
	closeCall := "window.kibanCloseMenu('" + splitID + "')"
	if user := strings.TrimSpace(it.OnClick); user != "" {
		// The user expression runs first: returning false aborts the
		// submit, and the close still happens either way.
		out["onclick"] = user + ";" + closeCall
	} else {
		out["onclick"] = closeCall
	}
	if it.Name != "" {
		out["name"] = it.Name
		out["value"] = it.Value
	}
	if it.Form != "" {
		out["form"] = it.Form
	}
	if it.Disabled {
		out["disabled"] = "disabled"
	}
	return out
}

func (it SplitItem) class() string {
	base := "flex w-full items-center gap-2 px-3 py-2 text-sm text-left transition-colors hover:bg-kiban-surface disabled:opacity-50 disabled:cursor-not-allowed"
	if strings.EqualFold(strings.TrimSpace(it.Variant), "danger") {
		return base + " text-red-600"
	}
	return base + " text-kiban-ink"
}
