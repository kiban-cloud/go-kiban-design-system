package input_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/input"
)

func TestDropzone_InputDeArchivoYTextos(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, input.Dropzone(input.DropzoneOptions{Name: "file", Accept: ".xlsx", Title: "Suelta tu archivo", Hint: "Máximo 5 MB", Required: true}).Render(context.Background(), &buf))
	body := buf.String()

	assert.Contains(t, body, "data-kiban-dropzone")
	assert.Contains(t, body, `type="file" name="file" accept=".xlsx" required`)
	assert.Contains(t, body, ">Suelta tu archivo</span>")
	assert.Contains(t, body, ">Máximo 5 MB</span>")
	assert.Contains(t, body, "data-kiban-dropzone-file")
}
