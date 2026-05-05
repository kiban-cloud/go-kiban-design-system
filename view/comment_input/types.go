// Package comment_input renders a self-contained "post a comment" widget:
// textarea + chip-style file uploader (with per-chip remove + invalid
// marking) + submit button. Encapsulates the look-and-feel + the file
// chip JS so consumers (klin's delivery comments today, future
// rekon/crm comment flows tomorrow) get the same UX with a single
// `@comment_input.Field(opts)` callsite.
//
// Architecture:
//   - The form is a plain HTML POST to Options.Action. The DS itself is
//     HTMX-agnostic; consumers that want HTMX wire it via Options.FormAttrs.
//   - The chip uploader logic (DataTransfer-append on each pick, remove
//     individual files, mark invalid by size, disable submit while invalid)
//     is JS that lives in view/layout/base.templ — re-runs on
//     DOMContentLoaded and after each htmx:afterSwap, mirroring the
//     intl-tel-input init pattern. Each instance is scoped by its
//     `[data-kiban-comment-input]` wrapper.
//   - All visual primitives are reused from existing DS subpackages
//     (card, flash, input.Textarea, input.File, button) — this package
//     adds no new colors/classes, only the chip markup + the JS contract.
package comment_input

import "github.com/a-h/templ"

// Options configures a single [Field]. Only Action is required; every
// other field has a sensible default so simple use cases are one line.
type Options struct {
	// Action is the form-submission URL. Required.
	Action string

	// Method overrides the form method. Empty defaults to "POST".
	Method string

	// ID is used to derive unique ids for the wrapper, form, textarea, and
	// file input so multiple Field instances on the same page don't
	// collide. Empty defaults to "comment-input".
	ID string

	// TextName is the form field name for the textarea. Defaults to "text".
	TextName string

	// FilesName is the form field name for the file input. Defaults to
	// "files".
	FilesName string

	// Title is the card heading. Defaults to "Nuevo comentario". Pass ""
	// explicitly with WithoutCard:true if rendering inside a custom chrome.
	Title string

	// Subtitle is the muted line under the title. Defaults to
	// "Comparte una actualización o adjunta archivos.".
	Subtitle string

	// Placeholder for the textarea. Defaults to "Escribe un comentario…".
	Placeholder string

	// SubmitLabel is the text on the send button. Defaults to "Enviar".
	SubmitLabel string

	// MaxChars hints the user-facing character limit. When > 0, renders
	// "Máximo N caracteres" as the textarea hint. Server-side validation
	// is the consumer's responsibility — this is purely informative.
	MaxChars int

	// MaxSizeBytes is the per-file size cap enforced client-side. Files
	// exceeding this are rendered as red/invalid chips with a tooltip and
	// the submit button is disabled until they're removed. 0 disables the
	// client-side check (server still validates).
	MaxSizeBytes int64

	// Accept restricts the file picker via the `accept` attribute (e.g.
	// ".pdf,.csv,image/*"). Empty allows anything. Browsers honor this as
	// a filter; the server should still validate.
	Accept string

	// Multiple emits the `multiple` attribute on the file input so the
	// user can select more than one file at a time. The chip uploader
	// also accumulates across multiple picker opens (DataTransfer trick).
	Multiple bool

	// DisableFiles hides the file uploader entirely (text-only comment
	// flow). Defaults to false (file input visible). Inverse-named so
	// the zero-value gives the most common case (files enabled).
	DisableFiles bool

	// TextValue round-trips the textarea content on validation failure
	// so the user doesn't lose their draft.
	TextValue string

	// TextError is the per-field error for the textarea (red border + msg).
	TextError string

	// GlobalError renders a flash.Error banner above the form (for
	// non-field-specific failures like "no pudimos guardar el comentario").
	GlobalError string

	// Success renders a flash.Success banner above the form (post-submit
	// confirmation, "Comentario publicado.").
	Success string

	// WithoutCard skips the card.Card chrome and emits just the form bits.
	// Use when wrapping the field in a different container (custom card,
	// drawer, etc.).
	WithoutCard bool

	// FormAttrs are spread onto the `<form>` element. The escape hatch
	// for HTMX wiring (hx-post, hx-target, hx-swap, …) or any extra
	// attribute the consumer needs.
	FormAttrs templ.Attributes
}

func (o Options) effectiveID() string {
	if o.ID != "" {
		return o.ID
	}
	return "comment-input"
}

func (o Options) effectiveMethod() string {
	if o.Method != "" {
		return o.Method
	}
	return "POST"
}

func (o Options) effectiveTextName() string {
	if o.TextName != "" {
		return o.TextName
	}
	return "text"
}

func (o Options) effectiveFilesName() string {
	if o.FilesName != "" {
		return o.FilesName
	}
	return "files"
}

func (o Options) effectiveTitle() string {
	if o.Title != "" {
		return o.Title
	}
	return "Nuevo comentario"
}

func (o Options) effectiveSubtitle() string {
	if o.Subtitle != "" {
		return o.Subtitle
	}
	return "Comparte una actualización o adjunta archivos."
}

func (o Options) effectivePlaceholder() string {
	if o.Placeholder != "" {
		return o.Placeholder
	}
	return "Escribe un comentario…"
}

func (o Options) effectiveSubmitLabel() string {
	if o.SubmitLabel != "" {
		return o.SubmitLabel
	}
	return "Enviar"
}

func (o Options) filesEnabled() bool {
	return !o.DisableFiles
}
