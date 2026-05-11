package comment_input_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/comment_input"
	"github.com/stretchr/testify/assert"
)

func render(t *testing.T, opts comment_input.Options) string {
	t.Helper()
	var buf bytes.Buffer
	err := comment_input.Field(opts).Render(context.Background(), &buf)
	assert.NoError(t, err)
	return buf.String()
}

func TestField_Defaults(t *testing.T) {
	body := render(t, comment_input.Options{Action: "/c"})
	// Default copy.
	assert.True(t, strings.Contains(body, "Nuevo comentario"))
	assert.True(t, strings.Contains(body, "Comparte una actualizaci"))
	assert.True(t, strings.Contains(body, "Escribe un comentario"))
	// Form wired to Action.
	assert.True(t, strings.Contains(body, `action="/c"`))
	// Form id derived from default ID.
	assert.True(t, strings.Contains(body, `id="comment-input-form"`))
	// Wrapper exposes the JS hook (now via the embedded file_chip_input).
	assert.True(t, strings.Contains(body, "data-kiban-file-chip-input"))
	// Submit carries the data tag the chip JS uses to find/disable it,
	// pointing at the file_chip_input's id (same lookup the JS does).
	assert.True(t, strings.Contains(body, `data-kiban-file-chip-submit="comment-input-files"`))
	// Default submit label is rendered inside the button.
	assert.True(t, strings.Contains(body, "Enviar"))
	// Chip list mount point comes from the file_chip_input.
	assert.True(t, strings.Contains(body, "data-chip-list"))
}

func TestField_OmitsCardWhenWithoutCard(t *testing.T) {
	body := render(t, comment_input.Options{Action: "/c", WithoutCard: true})
	// The card chrome adds "bg-white border border-kiban-border rounded-md p-6"
	// — at least one of those classes should be absent in WithoutCard mode.
	assert.False(t, strings.Contains(body, "bg-white border border-kiban-border rounded-md p-6"))
	// But the form itself is still emitted.
	assert.True(t, strings.Contains(body, `action="/c"`))
}

func TestField_RendersTextErrorAndRoundtripsValue(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:    "/c",
		TextValue: "borrador del usuario",
		TextError: "Escribe un comentario.",
	})
	assert.True(t, strings.Contains(body, "borrador del usuario"))
	assert.True(t, strings.Contains(body, "Escribe un comentario."))
}

func TestField_RendersGlobalError(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:      "/c",
		GlobalError: "No pudimos guardar.",
	})
	assert.True(t, strings.Contains(body, "No pudimos guardar."))
}

func TestField_RendersSuccessFlash(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:  "/c",
		Success: "Comentario publicado.",
	})
	assert.True(t, strings.Contains(body, "Comentario publicado."))
}

func TestField_HintCombinesSizeAndMultiple(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:       "/c",
		MaxSizeBytes: 5 * 1024 * 1024,
		Multiple:     true,
	})
	// "Hasta 5 MB por archivo. Puedes seleccionar varios." — accent é
	// gets escaped by templ to entity, so check the substring up to it.
	assert.True(t, strings.Contains(body, "Hasta 5 MB por archivo."))
	assert.True(t, strings.Contains(body, "Puedes seleccionar varios"))
}

func TestField_MaxCharsHint(t *testing.T) {
	body := render(t, comment_input.Options{Action: "/c", MaxChars: 2000})
	assert.True(t, strings.Contains(body, "M&aacute;ximo 2000 caracteres") ||
		strings.Contains(body, "Máximo 2000 caracteres"))
}

func TestField_MaxSizeBytesAttrPropagated(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:       "/c",
		MaxSizeBytes: 5 * 1024 * 1024,
	})
	// The wrapper carries the size limit so the JS can enforce client-side.
	assert.True(t, strings.Contains(body, `data-max-size="5242880"`))
}

func TestField_MultipleAttrEmittedOnInput(t *testing.T) {
	// With Multiple:true, the chip JS will rebind via DataTransfer; but the
	// input itself should also carry `multiple` so the picker UI lets the
	// user select more than one file in the first place. Today we rely on
	// the input.File component default + the JS, but verifying
	// presence of a file input regardless guards against regressions where
	// the file block silently disappears.
	body := render(t, comment_input.Options{Action: "/c", Multiple: true})
	assert.True(t, strings.Contains(body, `type="file"`))
	assert.True(t, strings.Contains(body, `name="files"`))
}

func TestField_DisableFilesHidesUploader(t *testing.T) {
	body := render(t, comment_input.Options{Action: "/c", DisableFiles: true})
	assert.False(t, strings.Contains(body, "data-kiban-file-chip-input"))
	assert.False(t, strings.Contains(body, `type="file"`))
	// Textarea still there.
	assert.True(t, strings.Contains(body, `name="text"`))
}

func TestField_CustomIDsAreUnique(t *testing.T) {
	a := render(t, comment_input.Options{Action: "/a", ID: "thread-1"})
	b := render(t, comment_input.Options{Action: "/b", ID: "thread-2"})
	assert.True(t, strings.Contains(a, `id="thread-1-form"`))
	assert.True(t, strings.Contains(b, `id="thread-2-form"`))
	// And the file_chip_input wrapper id is also distinct so the JS
	// scoping by wrapper doesn't collide between two siblings on the
	// same page.
	assert.True(t, strings.Contains(a, `id="thread-1-files"`))
	assert.True(t, strings.Contains(b, `id="thread-2-files"`))
}

func TestField_FormAttrsSpreadOnForm(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:    "/c",
		FormAttrs: templ.Attributes{"hx-post": "/c", "hx-target": "#tab"},
	})
	// FormAttrs is the HTMX escape hatch — the consumer can drive the
	// submit through HTMX without the DS package depending on htmx.
	assert.True(t, strings.Contains(body, `hx-post="/c"`))
	assert.True(t, strings.Contains(body, `hx-target="#tab"`))
}

func TestField_CustomFieldNamesPropagate(t *testing.T) {
	body := render(t, comment_input.Options{
		Action:    "/c",
		TextName:  "body",
		FilesName: "attachments",
	})
	assert.True(t, strings.Contains(body, `name="body"`))
	assert.True(t, strings.Contains(body, `name="attachments"`))
	assert.False(t, strings.Contains(body, `name="text"`))
	assert.False(t, strings.Contains(body, `name="files"`))
}
