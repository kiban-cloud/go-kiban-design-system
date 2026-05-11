// Package file_chip_input renders a `<input type="file">` paired with a
// chip-style readout of the selected files (one [chip.Chip] per file).
// Selecting more files appends instead of replacing (DataTransfer trick),
// each chip has a remove "×" that drops just that file, and oversize
// files render as `danger` chips with a tooltip + the consumer's submit
// button auto-disabled until they're removed.
//
// Architecture:
//   - The component emits markup only: an `<input type="file">` (using
//     [input.File]) plus a sibling `<ul data-chip-list>`. The chip list
//     starts empty server-side; JS populates it from the user's pick.
//   - The DataTransfer / dedup / per-chip-remove / disable-on-invalid
//     JS lives in this package's templ ([fileChipInputScript]) so each
//     instance pulls its behavior in via a once-per-page script tag.
//     The script is re-entrant: it scans `[data-kiban-file-chip-input]`
//     on DOMContentLoaded *and* `htmx:afterSwap`, mirroring the rest of
//     the design system's init pattern.
//   - Each instance is identified by a unique `[data-kiban-file-chip-input]`
//     wrapper. Multiple instances on the same page work without
//     collisions.
//
// Caller integration:
//   - Wrap the field in a `<form>` (or HTMX equivalent). The component
//     does not render the form; it's a sub-widget.
//   - For "disable submit while invalid file present", the consumer's
//     submit button must carry `data-kiban-file-chip-submit="<id>"`
//     (where `<id>` matches [Options.ID]) — same contract
//     [comment_input] uses for its own submit. Without that attr the
//     JS still renders chips correctly; only the disable-while-invalid
//     auto-wire is skipped.
package file_chip_input

import "github.com/a-h/templ"

// Options configures a single [Field]. Only [Options.Name] is required;
// every other field has a sensible default.
type Options struct {
	// Name is the form field name on the underlying `<input type="file">`.
	// Required.
	Name string

	// Label is the text rendered above the file picker (via
	// [input.File]'s label). Empty omits the label entirely. Default
	// is empty — most consumers render their own label outside the
	// component.
	Label string

	// ID identifies this instance. Used as the wrapper id and as the
	// value the submit button references via
	// `data-kiban-file-chip-submit="<id>"`. Empty defaults to
	// "kiban-file-chip-input"; collisions across multiple instances on
	// the same page are the consumer's problem to avoid.
	ID string

	// Hint is the line under the file picker. When empty, [filesHint]
	// auto-builds one from [Options.MaxSizeBytes] and [Options.Multiple]
	// so the simple cases ("Hasta 5 MB por archivo. Puedes seleccionar
	// varios.") need no caller code.
	Hint string

	// Accept restricts the picker via the native `accept` attribute
	// (e.g. `.csv,text/csv`). Empty means "any file". Browsers honor
	// this as a filter; the server is still the authority on what's
	// allowed.
	Accept string

	// Multiple emits the `multiple` attribute on the `<input>`. When
	// false, the picker only takes one file at a time and re-picking
	// replaces the prior selection (the chip JS still handles that
	// correctly — the staging list resets to the new file).
	Multiple bool

	// Required emits HTML5 `required` on the input — the browser blocks
	// submit when nothing's selected. Server-side validation should
	// still own the source of truth.
	Required bool

	// MaxSizeBytes is the per-file size cap enforced client-side. Files
	// over the cap render as `danger` chips with a tooltip; while any
	// such chip is present, the consumer's submit (when wired with
	// `data-kiban-file-chip-submit`) is auto-disabled. 0 disables the
	// client-side check.
	MaxSizeBytes int64

	// FileVariant overrides the chip variant for valid (under-cap)
	// files. Default "" / "default" renders neutral chips matching the
	// rest of the kit. Use "info" to call out the selection visually
	// when the field is the page's focus (rare).
	FileVariant string
}

// effectiveID returns the wrapper id, falling back to a stable default
// so a caller passing zero options still works on a one-instance page.
func (o Options) effectiveID() string {
	if o.ID != "" {
		return o.ID
	}
	return "kiban-file-chip-input"
}

// dataAttrs collects the wrapper's data-* attributes. Done as a helper
// so the templ file stays declarative — the templ for-loops `Attrs`
// onto the element and we don't have an `if` ladder per data-attr.
func (o Options) dataAttrs() templ.Attributes {
	a := templ.Attributes{
		"data-kiban-file-chip-input": o.effectiveID(),
	}
	if o.MaxSizeBytes > 0 {
		a["data-max-size"] = formatInt64(o.MaxSizeBytes)
	}
	if o.FileVariant != "" {
		a["data-file-variant"] = o.FileVariant
	}
	return a
}
