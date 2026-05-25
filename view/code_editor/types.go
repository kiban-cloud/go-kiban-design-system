// Package code_editor wraps CodeMirror 5 as a templ component for any
// kiban project that needs a code-editing surface (JS, Python, JSON,
// etc.). The page hosting this component must opt in to the
// CodeMirror CDN bundle via `layout.Config.LoadCodeMirror = true` —
// without that, the JS init in `view/layout/base.templ` no-ops and
// the user just sees the underlying <textarea>.
//
// The init JS lives in base.templ (alongside file_chip_input,
// jsonviewer, canvas, etc.) and re-runs on htmx:afterSwap, so partial
// responses that include an editor get their CodeMirror instance
// without any caller action.
package code_editor

// Options configures a single CodeEditor. The component renders a
// <textarea> wrapped in a div tagged `data-kiban-code-editor`; the
// runtime JS spots it, picks a CodeMirror mode from Language, and
// hides the textarea behind a CodeMirror surface. The textarea's
// value is kept in sync with the editor on every change so a plain
// `new FormData(form)` reads the latest code.
//
// Name is the form field name (mandatory). Label is rendered above
// the editor (omit for an inline editor with no chrome). Hint and
// ErrMsg follow the same contract as view/input components — only
// one of them shows at a time. Required adds the red asterisk.
//
// Language maps to a CodeMirror mode keyword. Today the init JS
// understands the kiban-domain keywords "nodejs20" and "python312"
// plus the generic "javascript" and "python"; anything else falls
// back to JavaScript mode. Add new mappings here if more languages
// become supported.
type Options struct {
	Name     string
	Label    string
	Value    string
	Hint     string
	ErrMsg   string
	Language string
	Required bool
	// Rows sets the initial textarea height (in line-height units).
	// CodeMirror uses `viewportMargin: Infinity` so the editor grows
	// with its content; Rows only affects the underlying textarea's
	// pre-init height, which the user briefly sees before CodeMirror
	// takes over. Defaults to 12 when <= 0.
	Rows int
}

// effectiveRows returns Rows or its default so the templ doesn't need
// to repeat the fallback.
func (o Options) effectiveRows() int {
	if o.Rows <= 0 {
		return 12
	}
	return o.Rows
}

// inputID stamps the textarea (and its <label for>) with a stable id.
// Pattern matches view/input/* helpers.
func (o Options) inputID() string {
	return "f-" + o.Name
}
