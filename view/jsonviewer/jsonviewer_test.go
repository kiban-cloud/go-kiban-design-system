package jsonviewer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/jsonviewer"
)

func render(t *testing.T, opts jsonviewer.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := jsonviewer.View(opts).Render(context.Background(), &buf)
	require.NoError(t, err)
	return buf.String()
}

// decode is the realistic input path: callers typically hand in the
// result of `json.Unmarshal` into `interface{}`, which produces
// `map[string]any` / `[]any` / primitives. Tests build the inputs
// the same way so we exercise the same runtime types the viewer
// sees in production.
func decode(t *testing.T, payload string) any {
	t.Helper()
	var v any
	require.NoError(t, json.Unmarshal([]byte(payload), &v))
	return v
}

func TestView_EmptyMessageWhenNil(t *testing.T) {
	body := render(t, jsonviewer.Options{Data: nil})
	assert.Contains(t, body, "Sin información.")
	// No tree chrome when there's nothing to show.
	assert.NotContains(t, body, "data-kiban-jsonviewer")
	assert.NotContains(t, body, "Expandir todo")
}

func TestView_EmptyMessageOverride(t *testing.T) {
	body := render(t, jsonviewer.Options{
		Data:         nil,
		EmptyMessage: "Sin solicitud.",
	})
	assert.Contains(t, body, "Sin solicitud.")
	assert.NotContains(t, body, "Sin información.")
}

func TestView_EmptyContainerCountsAsEmpty(t *testing.T) {
	body := render(t, jsonviewer.Options{Data: map[string]any{}})
	assert.Contains(t, body, "Sin información.")
}

func TestView_PrimitivesRenderAsKeyValueRows(t *testing.T) {
	body := render(t, jsonviewer.Options{
		Data: decode(t, `{"rfc": "MAMJ900118Y71", "edad": 35, "activo": true}`),
	})
	// Top-level wrapper + master toggle.
	assert.Contains(t, body, "data-kiban-jsonviewer")
	assert.Contains(t, body, "Expandir todo")
	// Each primitive renders as a key + value pair. The exact DOM
	// is flexible; assert on the rendered strings.
	assert.Contains(t, body, "rfc")
	assert.Contains(t, body, "MAMJ900118Y71")
	assert.Contains(t, body, "edad")
	assert.Contains(t, body, "35")
	assert.Contains(t, body, "activo")
	assert.Contains(t, body, "true")
	// Primitive-only top level → no nested <details>.
	assert.NotContains(t, body, "<details")
}

func TestView_NestedObjectGetsAccordion(t *testing.T) {
	body := render(t, jsonviewer.Options{
		Data: decode(t, `{"rfc": "MAMJ900118Y71", "domicilio": {"calle": "Reforma", "cp": "06600"}}`),
	})
	// The nested object renders as a <details> with the key as
	// the summary; clicking it expands without any JS round-trip.
	assert.Contains(t, body, "<details")
	assert.Contains(t, body, "data-jsonviewer-node")
	assert.Contains(t, body, "domicilio")
	// Nested primitives still render (collapsed view is the
	// default; the markup is present either way).
	assert.Contains(t, body, "calle")
	assert.Contains(t, body, "Reforma")
}

func TestView_ArrayOfPrimitivesUsesNumericKeys(t *testing.T) {
	body := render(t, jsonviewer.Options{
		Data: decode(t, `["alpha", "beta", "gamma"]`),
	})
	// Indices appear as 0/1/2 keys; values appear inline.
	assert.Contains(t, body, "alpha")
	assert.Contains(t, body, "beta")
	assert.Contains(t, body, "gamma")
}

func TestView_ArrayOfObjectsEachGetsItsOwnNode(t *testing.T) {
	body := render(t, jsonviewer.Options{
		Data: decode(t, `[{"name":"Alice"},{"name":"Bob"}]`),
	})
	// One <details> per element; indices serve as the labels.
	assert.GreaterOrEqual(t, strings.Count(body, "<details"), 2)
	assert.Contains(t, body, "Alice")
	assert.Contains(t, body, "Bob")
}

func TestView_NumericKeysSortNumerically(t *testing.T) {
	// Mongo-style numeric-keyed maps are common (string keys that
	// happen to look like indices). They must sort numerically so
	// "10" doesn't precede "2".
	body := render(t, jsonviewer.Options{
		Data: decode(t, `{"10":"ten","2":"two","1":"one"}`),
	})
	// Find the index of each rendered value to confirm order.
	idxOne := strings.Index(body, "one")
	idxTwo := strings.Index(body, "two")
	idxTen := strings.Index(body, "ten")
	require.True(t, idxOne >= 0 && idxTwo >= 0 && idxTen >= 0)
	assert.Less(t, idxOne, idxTwo)
	assert.Less(t, idxTwo, idxTen)
}

func TestView_MultipleViewersScopedById(t *testing.T) {
	// When two viewers share a page they must each scope their
	// master toggle to their own wrapper so clicking "Expandir
	// todo" on one doesn't expand the other.
	a := render(t, jsonviewer.Options{
		ID:   "viewer-a",
		Data: decode(t, `{"x":{"y":1}}`),
	})
	b := render(t, jsonviewer.Options{
		ID:   "viewer-b",
		Data: decode(t, `{"x":{"y":1}}`),
	})
	assert.Contains(t, a, `id="viewer-a"`)
	assert.Contains(t, a, `data-jsonviewer-target="#viewer-a"`)
	assert.Contains(t, b, `id="viewer-b"`)
	assert.Contains(t, b, `data-jsonviewer-target="#viewer-b"`)
}

func TestView_StringFallback(t *testing.T) {
	// A top-level primitive (string body that wasn't JSON-decoded)
	// renders as a single value cell, not an accordion.
	body := render(t, jsonviewer.Options{
		Data: "not-actually-json",
	})
	assert.Contains(t, body, "not-actually-json")
	assert.NotContains(t, body, "<details")
}
