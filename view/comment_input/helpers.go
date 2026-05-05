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

// filesHint returns the hint under the file input. Combines the size
// limit (when set) and a hint about multi-file selection (when enabled)
// into one line — keeping the visual quiet.
func filesHint(opts Options) string {
	switch {
	case opts.MaxSizeBytes > 0 && opts.Multiple:
		return fmt.Sprintf("Hasta %s por archivo. Puedes seleccionar varios.", humanSize(opts.MaxSizeBytes))
	case opts.MaxSizeBytes > 0:
		return fmt.Sprintf("Hasta %s por archivo.", humanSize(opts.MaxSizeBytes))
	case opts.Multiple:
		return "Puedes seleccionar varios archivos."
	}
	return ""
}

// submitAttrs returns the extra HTML attributes spread onto the submit
// button. We tag it with `data-kiban-comment-submit` so the chip JS
// can find and disable it while any chip is invalid.
func submitAttrs(id string) templ.Attributes {
	return templ.Attributes{
		"data-kiban-comment-submit": id,
	}
}

// humanSize formats a byte count as "5 MB" / "512 KB". Picks the
// largest unit where the value is >= 1 to keep the hint short.
func humanSize(bytes int64) string {
	const (
		kb = 1024
		mb = 1024 * 1024
		gb = 1024 * 1024 * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%d GB", bytes/gb)
	case bytes >= mb:
		return fmt.Sprintf("%d MB", bytes/mb)
	case bytes >= kb:
		return fmt.Sprintf("%d KB", bytes/kb)
	}
	return fmt.Sprintf("%d B", bytes)
}
