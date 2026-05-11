package file_chip_input

import (
	"fmt"
	"strconv"
)

// effectiveHint returns the hint string under the file picker. When
// the caller supplied [Options.Hint], that wins. Otherwise we auto-build
// a one-liner combining the size cap and multi-file affordance — which
// is what 90% of callers want and otherwise duplicate verbatim.
func effectiveHint(opts Options) string {
	if opts.Hint != "" {
		return opts.Hint
	}
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

// humanSize formats a byte count using the largest unit ≥ 1 so the
// hint stays short. Mirrors comment_input's helper one-for-one — they
// share the same unit thresholds so chips and hints read consistently.
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

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
