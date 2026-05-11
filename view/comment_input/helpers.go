package comment_input

import (
	"fmt"

	"github.com/a-h/templ"
)

// formAction is a thin wrapper that templ requires for `action={…}` —
// the type-checker on `<form action>` insists on templ.SafeURL.
func formAction(s string) templ.SafeURL {
	return templ.URL(s)
}

// charsHint returns the hint string under the textarea. Empty when
// MaxChars <= 0 so the textarea hint renders nothing.
func charsHint(maxChars int) string {
	if maxChars <= 0 {
		return ""
	}
	return fmt.Sprintf("Máximo %d caracteres", maxChars)
}

// submitAttrs returns the extra HTML attributes spread onto the submit
// button. The `data-kiban-file-chip-submit` attribute connects the
// button to the file_chip_input so the latter's JS auto-disables it
// while any selected file violates the size cap. The value matches the
// file_chip_input wrapper's `data-kiban-file-chip-input` (its ID).
func submitAttrs(filesID string) templ.Attributes {
	return templ.Attributes{
		"data-kiban-file-chip-submit": filesID,
	}
}
