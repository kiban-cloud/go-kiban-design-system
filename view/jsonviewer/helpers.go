package jsonviewer

import (
	"fmt"
	"sort"
	"strconv"
)

// entry pairs a stable display key with the raw value at that key.
// Used by the renderer to split a node's children into primitives
// (rendered inline) and containers (rendered as nested accordions),
// preserving deterministic order.
type entry struct {
	key   string
	value any
}

// classify splits a container value's children into primitives + a
// sub-tree of nested objects/arrays. Maps are ordered alphabetically
// (numeric keys sort numerically), and slices preserve their index
// order. The split mirrors the React JsonViewer's behaviour: at each
// level you see flat key/value rows for primitives, then nested
// accordions for the object children.
func classify(v any) (primitives []entry, children []entry) {
	switch t := v.(type) {
	case map[string]any:
		for _, k := range sortedKeys(t) {
			if isContainer(t[k]) {
				children = append(children, entry{k, t[k]})
			} else {
				primitives = append(primitives, entry{k, t[k]})
			}
		}
	case []any:
		for i, item := range t {
			key := strconv.Itoa(i)
			if isContainer(item) {
				children = append(children, entry{key, item})
			} else {
				primitives = append(primitives, entry{key, item})
			}
		}
	}
	return
}

// sortedKeys returns the map keys ordered numerically when every key
// parses as an int (a stringified array index), alphabetically
// otherwise. Matches the React component's `getObjectTemplate` rule.
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	allNumeric := len(keys) > 0
	for _, k := range keys {
		if _, err := strconv.Atoi(k); err != nil {
			allNumeric = false
			break
		}
	}
	if allNumeric {
		sort.SliceStable(keys, func(i, j int) bool {
			a, _ := strconv.Atoi(keys[i])
			b, _ := strconv.Atoi(keys[j])
			return a < b
		})
	} else {
		sort.Strings(keys)
	}
	return keys
}

// isContainer reports whether v is a JSON object/array (i.e. has
// children worth pushing into a nested accordion). Empty containers
// also count — the viewer still renders an empty accordion for them
// so the structure is visible.
func isContainer(v any) bool {
	switch v.(type) {
	case map[string]any, []any:
		return true
	}
	return false
}

// isEmpty reports whether v is the JSON equivalent of "no data" —
// nil, empty containers. Used by View to decide between rendering
// the tree and falling back to the EmptyMessage placeholder.
func isEmpty(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case map[string]any:
		return len(t) == 0
	case []any:
		return len(t) == 0
	}
	return false
}

// primitiveString converts a primitive JSON value to its display
// form. nil renders as "-" so empty fields are visually obvious;
// bools as Spanish-friendly "true"/"false" (matches the React MF's
// behaviour); numbers fall through to fmt.Sprint's default which
// handles ints and floats predictably.
func primitiveString(v any) string {
	switch t := v.(type) {
	case nil:
		return "-"
	case string:
		if t == "" {
			return "-"
		}
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		// Integer-valued floats render without trailing .0
		// (the typical case when JSON numbers are actually ints).
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}
