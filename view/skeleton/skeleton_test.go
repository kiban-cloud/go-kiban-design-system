package skeleton_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/skeleton"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTable_MirrorsRealTable(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"ID de pago", "Monto", "Fecha"}
	require.NoError(t, skeleton.Table(headers, 7).Render(context.Background(), &buf))
	out := buf.String()
	assert.Contains(t, out, "animate-pulse")
	assert.Contains(t, out, "data-kiban-skeleton")
	assert.Contains(t, out, `aria-hidden="true"`)
	// Real headers render in the header row (mirrors React's tableHeadings).
	assert.Contains(t, out, "ID de pago")
	assert.Contains(t, out, "Monto")
	// 7 body rows × 3 columns = 21 shimmer cells.
	assert.Equal(t, 21, strings.Count(out, "rounded-md bg-kiban-border/70"))
}

func TestTable_ClampsRowsToAtLeastOne(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, skeleton.Table([]string{"A", "B"}, 0).Render(context.Background(), &buf))
	// 1 row × 2 columns.
	assert.Equal(t, 2, strings.Count(buf.String(), "rounded-md bg-kiban-border/70"))
}

func TestSkeleton_GenericFourColumnTable(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, skeleton.Skeleton().Render(context.Background(), &buf))
	out := buf.String()
	assert.Contains(t, out, "data-kiban-skeleton")
	// 7 rows × 4 blank columns.
	assert.Equal(t, skeleton.DefaultRows*4, strings.Count(out, "rounded-md bg-kiban-border/70"))
}
