package code_editor_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/code_editor"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

func TestCodeEditor_RendersTextareaWithMarkerAndLanguage(t *testing.T) {
	out := render(t, code_editor.CodeEditor(code_editor.Options{
		Name:     "code",
		Label:    "Código",
		Value:    "console.log('hi');",
		Language: "nodejs20",
		Required: true,
	}))
	if !strings.Contains(out, `data-kiban-code-editor`) {
		t.Errorf("expected wrapper data-kiban-code-editor attribute: %s", out)
	}
	if !strings.Contains(out, `data-language="nodejs20"`) {
		t.Errorf("expected data-language attribute on wrapper: %s", out)
	}
	if !strings.Contains(out, `<textarea`) {
		t.Errorf("expected underlying textarea: %s", out)
	}
	if !strings.Contains(out, `name="code"`) {
		t.Errorf("expected textarea name attribute: %s", out)
	}
	if !strings.Contains(out, "Código") {
		t.Errorf("expected label text: %s", out)
	}
	if !strings.Contains(out, "&#42;") {
		t.Errorf("expected required asterisk: %s", out)
	}
	// Value is HTML-escaped inside the textarea.
	if !strings.Contains(out, "console.log") {
		t.Errorf("expected value to round-trip into textarea: %s", out)
	}
}

func TestCodeEditor_ErrorAndHintMutuallyExclusive(t *testing.T) {
	out := render(t, code_editor.CodeEditor(code_editor.Options{
		Name:   "code",
		Hint:   "Escribe código válido.",
		ErrMsg: "Sintaxis inválida.",
	}))
	if !strings.Contains(out, "Sintaxis inválida.") {
		t.Errorf("expected error message: %s", out)
	}
	if strings.Contains(out, "Escribe código válido.") {
		t.Errorf("hint should be suppressed when error is set: %s", out)
	}
}

func TestCodeEditor_DefaultRowsApplied(t *testing.T) {
	out := render(t, code_editor.CodeEditor(code_editor.Options{Name: "code"}))
	if !strings.Contains(out, `rows="12"`) {
		t.Errorf("expected default rows=12: %s", out)
	}
}

func TestCodeEditor_OmitLabel(t *testing.T) {
	out := render(t, code_editor.CodeEditor(code_editor.Options{Name: "code"}))
	if strings.Contains(out, "<label") {
		t.Errorf("expected no label when Label is empty: %s", out)
	}
}
