package lazy_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/lazy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContent_RendersSelfReplacingLazyBoxWithPlaceholder(t *testing.T) {
	placeholder := templ.Raw(`<div data-test-placeholder>loading</div>`)
	var buf bytes.Buffer
	require.NoError(t, lazy.Content("/rekon/payments/content?page=1").Render(templ.WithChildren(context.Background(), placeholder), &buf))
	out := buf.String()
	// Fetches the content URL on load and replaces itself.
	assert.Contains(t, out, `hx-get="/rekon/payments/content?page=1"`)
	assert.Contains(t, out, `hx-trigger="load"`)
	assert.Contains(t, out, `hx-target="this"`)
	assert.Contains(t, out, `hx-swap="outerHTML"`)
	assert.Contains(t, out, "data-kiban-lazy")
	// The caller-supplied placeholder renders inside until content arrives.
	assert.Contains(t, out, "data-test-placeholder")
}
