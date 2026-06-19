package input

import "strings"

// initialDial joins a stored country code + national number into the
// "+<cc><national>" shape that intl-tel-input's setNumber() expects on
// pre-seed. The stored cc may or may not already carry a leading "+"
// (the widget now emits "+52", but older records / other consumers store
// "52"), so we normalize to exactly one. Returns "" when either side is
// empty so the widget falls back to its configured initialCountry.
func initialDial(cc, national string) string {
	if cc != "" && national != "" {
		return "+" + strings.TrimPrefix(cc, "+") + national
	}
	return ""
}
