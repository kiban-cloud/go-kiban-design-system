package file_chip_input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/kiban-cloud/go-kiban-design-system/view/file_chip_input"
	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, opts file_chip_input.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := file_chip_input.Field(opts).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestField_DefaultsRenderInputAndChipList(t *testing.T) {
	body := render(t, file_chip_input.Options{Name: "files"})
	// File input emitted with the right name attribute.
	assert.True(t, strings.Contains(body, `type="file"`))
	assert.True(t, strings.Contains(body, `name="files"`))
	// Wrapper carries the JS-binding attribute (default id).
	assert.True(t, strings.Contains(body, `data-kiban-file-chip-input="kiban-file-chip-input"`))
	// Chip list is the JS contract for chip injection.
	assert.True(t, strings.Contains(body, "data-chip-list"))
}

func TestField_CustomIDGoesOnWrapperAndDataAttr(t *testing.T) {
	body := render(t, file_chip_input.Options{Name: "files", ID: "create-files"})
	// Wrapper id matches.
	assert.True(t, strings.Contains(body, `id="create-files"`))
	// And the JS data attribute carries the same id so the submit button
	// (data-kiban-file-chip-submit="create-files") binds correctly.
	assert.True(t, strings.Contains(body, `data-kiban-file-chip-input="create-files"`))
}

func TestField_MaxSizeBytesAttrPropagated(t *testing.T) {
	body := render(t, file_chip_input.Options{
		Name:         "files",
		MaxSizeBytes: 5 * 1024 * 1024,
	})
	assert.True(t, strings.Contains(body, `data-max-size="5242880"`))
}

func TestField_FileVariantPropagated(t *testing.T) {
	body := render(t, file_chip_input.Options{
		Name:        "files",
		FileVariant: "info",
	})
	assert.True(t, strings.Contains(body, `data-file-variant="info"`))
}

func TestField_HintAutoBuildsFromSizeAndMultiple(t *testing.T) {
	body := render(t, file_chip_input.Options{
		Name:         "files",
		MaxSizeBytes: 5 * 1024 * 1024,
		Multiple:     true,
	})
	assert.True(t, strings.Contains(body, "Hasta 5 MB por archivo."))
	assert.True(t, strings.Contains(body, "Puedes seleccionar varios"))
}

func TestField_HintRespectsExplicitOverride(t *testing.T) {
	body := render(t, file_chip_input.Options{
		Name: "files",
		Hint: "Selecciona los archivos a procesar.",
	})
	assert.True(t, strings.Contains(body, "Selecciona los archivos a procesar."))
}

func TestField_AcceptPropagatedToInput(t *testing.T) {
	body := render(t, file_chip_input.Options{
		Name:   "files",
		Accept: ".csv,text/csv",
	})
	assert.True(t, strings.Contains(body, `accept=".csv,text/csv"`))
}

func TestField_RequiredEmitsAttrOnInput(t *testing.T) {
	body := render(t, file_chip_input.Options{Name: "files", Required: true})
	assert.True(t, strings.Contains(body, "required"))
}
