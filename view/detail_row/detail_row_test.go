package detail_row_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/detail_row"
)

func renderRow(t *testing.T, label, value string) string {
	t.Helper()
	var buf bytes.Buffer
	err := detail_row.Row(label, value).Render(context.Background(), &buf)
	require.NoError(t, err)
	return buf.String()
}

func renderList(t *testing.T, items []detail_row.Item) string {
	t.Helper()
	var buf bytes.Buffer
	err := detail_row.List(items).Render(context.Background(), &buf)
	require.NoError(t, err)
	return buf.String()
}

func TestRow_RendersLabelAndValue(t *testing.T) {
	body := renderRow(t, "Fecha de ejecución", "15/05/2026 15:18")
	assert.True(t, strings.Contains(body, "Fecha de ejecuci"))
	assert.True(t, strings.Contains(body, "15/05/2026 15:18"))
	assert.True(t, strings.Contains(body, "<dt"))
	assert.True(t, strings.Contains(body, "<dd"))
}

// Empty value collapses to "-" — the dash placeholder keeps the
// row from rendering visually broken next to populated rows. The
// "hide row entirely when value is empty" rule lives at the
// caller; the component takes a value and always renders it.
func TestRow_EmptyValueRendersDashPlaceholder(t *testing.T) {
	body := renderRow(t, "Creado por", "")
	assert.True(t, strings.Contains(body, "Creado por"))
	assert.True(t, strings.Contains(body, ">-<"))
}

// List wraps rows in a <dl> with vertical spacing. The order
// caller passes is preserved.
func TestList_PreservesOrderAndWrapsInDL(t *testing.T) {
	body := renderList(t, []detail_row.Item{
		{Label: "Uno", Value: "1"},
		{Label: "Dos", Value: "2"},
		{Label: "Tres", Value: "3"},
	})
	assert.True(t, strings.Contains(body, "<dl"))
	// Position check: "1" appears before "2" before "3".
	idx1 := strings.Index(body, ">1<")
	idx2 := strings.Index(body, ">2<")
	idx3 := strings.Index(body, ">3<")
	assert.True(t, idx1 >= 0 && idx2 > idx1 && idx3 > idx2, "expected values to render in order")
}

func TestList_EmptyItemsRendersEmptyDL(t *testing.T) {
	body := renderList(t, nil)
	assert.True(t, strings.Contains(body, "<dl"))
	assert.False(t, strings.Contains(body, "<dt"))
}
