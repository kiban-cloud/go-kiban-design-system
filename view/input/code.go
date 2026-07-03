package input

import (
	"strconv"
	"strings"
)

// codeLength normalises the requested one-time-code length, falling back to 6
// (the templ zero-value, the usual NIP/OTP size) when the caller passes 0.
func codeLength(length int) int {
	if length <= 0 {
		return 6
	}
	return length
}

// codeMaxLength stringifies the code length for the input's maxlength attr.
func codeMaxLength(length int) string {
	return strconv.Itoa(codeLength(length))
}

// codePattern builds the HTML5 `pattern` that enforces exactly `length` digits.
func codePattern(length int) string {
	return "\\d{" + strconv.Itoa(codeLength(length)) + "}"
}

// codePlaceholder renders one bullet per expected digit ("••••••" for 6).
func codePlaceholder(length int) string {
	return strings.Repeat("•", codeLength(length))
}
