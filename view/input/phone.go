package input

// initialDial joins a stored country code + national number into the
// "+<cc><national>" shape that intl-tel-input's setNumber() expects on
// pre-seed. Returns "" when either side is empty so the widget falls back
// to its configured initialCountry.
func initialDial(cc, national string) string {
	if cc != "" && national != "" {
		return "+" + cc + national
	}
	return ""
}
