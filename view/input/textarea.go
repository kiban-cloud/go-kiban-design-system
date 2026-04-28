package input

import "strconv"

// rowsAttr stringifies the visible-row count for Textarea, falling back to
// "3" when the caller passes 0 (the templ zero-value, common when the field
// isn't relevant for a given usage).
func rowsAttr(rows int) string {
	if rows <= 0 {
		rows = 3
	}
	return strconv.Itoa(rows)
}
