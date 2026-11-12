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

// Las capturas de la certificación del proto: invitación con «elígela» en
// azul, ícono de imagen con miniatura, alto mínimo y, al elegir, la zona en
// verde con el nombre y «Cambiar archivo» (data-change, lo lee base.templ).
func TestDropzone_ComoLasCapturasDelProto(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, input.Dropzone(input.DropzoneOptions{
		Name: "screenshotForm", Accept: "image/png", Title: "Arrastra la captura", Or: "o", Pick: "elígela",
		Hint: "PNG · hasta 5 MB", Image: true, Tall: true, Change: "Cambiar archivo",
	}).Render(context.Background(), &buf))
	body := buf.String()

	assert.Contains(t, body, `data-change="Cambiar archivo"`)
	assert.Contains(t, body, "min-h-[150px]")
	assert.Contains(t, body, "data-[has-file]:border-[#00B894]", "elegido, en verde")
	assert.Contains(t, body, `<b class="font-medium">Arrastra la captura</b> o <span class="font-medium text-kiban-primary">elígela</span>`)
	assert.Contains(t, body, "data-kiban-dropzone-thumb")
	assert.Contains(t, body, `<rect x="3" y="3" width="18" height="18" rx="2">`, "ícono de imagen")
	assert.Contains(t, body, "data-kiban-dropzone-hint>PNG · hasta 5 MB</span>")
}

func TestDropzone_SinOpcionesNuevasQuedaComoSiempre(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, input.Dropzone(input.DropzoneOptions{Name: "file", Title: "Suelta tu archivo"}).Render(context.Background(), &buf))
	body := buf.String()

	assert.NotContains(t, body, "data-change")
	assert.NotContains(t, body, "min-h-[150px]")
	assert.NotContains(t, body, "data-kiban-dropzone-thumb")
	assert.Contains(t, body, `<span class="text-sm font-medium text-kiban-ink" data-kiban-dropzone-title>Suelta tu archivo</span>`)
}
