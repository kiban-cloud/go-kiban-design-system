package canvas_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"

	"github.com/kiban-cloud/go-kiban-design-system/view/canvas"
	"github.com/kiban-cloud/go-kiban-design-system/view/icons"
)

// renderCanvas wraps the children-taking Canvas component with the helper
// child content the test wants emitted into the stage. Mirrors the
// renderPanel pattern in view/tabs.
func renderCanvas(t *testing.T, cfg canvas.CanvasOptions, body string) string {
	t.Helper()
	children := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, body)
		return err
	})
	ctx := templ.WithChildren(context.Background(), children)
	var buf bytes.Buffer
	require.NoError(t, canvas.Canvas(cfg).Render(ctx, &buf))
	return buf.String()
}

// renderComponent is for the children-less components (Node, EdgeButton).
func renderComponent(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

func TestCanvas_RendersWrapperAndEdges(t *testing.T) {
	out := renderCanvas(t, canvas.CanvasOptions{
		ID: "wf-canvas",
		Edges: []canvas.EdgeOptions{
			{From: "node-a", To: "node-b"},
			{From: "node-b", To: "node-c", Variant: "error", Label: "Sí"},
		},
	}, `<span data-test-child="yes">child</span>`)

	if !strings.Contains(out, `id="wf-canvas"`) {
		t.Errorf("expected wrapper id: %s", out)
	}
	if !strings.Contains(out, `data-kiban-canvas`) {
		t.Errorf("expected data-kiban-canvas attribute")
	}
	if !strings.Contains(out, `&#34;from&#34;:&#34;node-a&#34;`) && !strings.Contains(out, `"from":"node-a"`) {
		t.Errorf("expected from in data-edges: %s", out)
	}
	if !strings.Contains(out, `error`) || !strings.Contains(out, `variant`) {
		t.Errorf("expected variant in data-edges JSON: %s", out)
	}
	if !strings.Contains(out, `data-kiban-canvas-edges`) {
		t.Errorf("expected svg overlay")
	}
	if !strings.Contains(out, `data-kiban-canvas-stage`) {
		t.Errorf("expected stage container")
	}
	if !strings.Contains(out, `data-test-child="yes"`) {
		t.Errorf("expected children to render inside stage: %s", out)
	}
}

func TestCanvas_DefaultID(t *testing.T) {
	out := renderCanvas(t, canvas.CanvasOptions{}, "")
	if !strings.Contains(out, `id="kiban-canvas"`) {
		t.Errorf("expected default id when ID empty: %s", out)
	}
}

func TestCanvas_EmptyEdgesAsArray(t *testing.T) {
	out := renderCanvas(t, canvas.CanvasOptions{}, "")
	if !strings.Contains(out, `data-edges="[]"`) {
		t.Errorf("expected data-edges to be []: %s", out)
	}
}

func TestCanvas_EmptyMessage(t *testing.T) {
	out := renderCanvas(t, canvas.CanvasOptions{
		EmptyMessage: "Aún no hay nodos.",
	}, "")
	if !strings.Contains(out, "Aún no hay nodos.") {
		t.Errorf("expected empty message text: %s", out)
	}
	if !strings.Contains(out, `data-kiban-canvas-empty`) {
		t.Errorf("expected empty-state container: %s", out)
	}
}

func TestCanvas_EmptyAction(t *testing.T) {
	marker := templ.Raw(`<button data-empty-action>Crear</button>`)
	out := renderCanvas(t, canvas.CanvasOptions{
		EmptyMessage: "x",
		EmptyAction:  marker,
	}, "")
	if !strings.Contains(out, `data-empty-action`) {
		t.Errorf("expected empty action to render: %s", out)
	}
}

func TestNode_PlainDivWhenNoHref(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:    "node-1",
		Title: "Formulario inicial",
	}))
	if !strings.Contains(out, `id="node-1"`) {
		t.Errorf("expected id: %s", out)
	}
	if strings.Contains(out, "<a ") {
		t.Errorf("expected <div> not <a> when Href empty: %s", out)
	}
	if !strings.Contains(out, "Formulario inicial") {
		t.Errorf("expected title text: %s", out)
	}
	if !strings.Contains(out, `data-status=""`) {
		t.Errorf("expected empty data-status when Status omitted: %s", out)
	}
}

func TestNode_AnchorWhenHrefSet(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:    "node-1",
		Title: "Click me",
		Href:  "/edit/abc",
	}))
	if !strings.Contains(out, `<a `) {
		t.Errorf("expected <a> when Href set: %s", out)
	}
	if !strings.Contains(out, `href="/edit/abc"`) {
		t.Errorf("expected href attribute: %s", out)
	}
}

func TestNode_Subtitle(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:       "n",
		Title:    "T",
		Subtitle: "Formulario",
	}))
	if !strings.Contains(out, "Formulario") {
		t.Errorf("expected subtitle to render: %s", out)
	}
}

func TestNode_IconRenders(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:    "n",
		Title: "T",
		Icon:  icons.GitMerge(),
	}))
	if !strings.Contains(out, "<svg") {
		t.Errorf("expected svg icon: %s", out)
	}
}

func TestNode_StatusError(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:     "n",
		Title:  "T",
		Status: canvas.StatusError,
	}))
	if !strings.Contains(out, "Error") {
		t.Errorf("expected Error pill: %s", out)
	}
	if !strings.Contains(out, "border-red-300") {
		t.Errorf("expected red border class: %s", out)
	}
}

func TestNode_StatusNotConfigured(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:     "n",
		Title:  "T",
		Status: canvas.StatusNotConfigured,
	}))
	if !strings.Contains(out, "Pendiente") {
		t.Errorf("expected Pendiente pill: %s", out)
	}
	if !strings.Contains(out, "border-amber-300") {
		t.Errorf("expected amber border class: %s", out)
	}
}

func TestNode_StatusOK_NoPill(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:     "n",
		Title:  "T",
		Status: canvas.StatusOK,
	}))
	if strings.Contains(out, "Pendiente") || strings.Contains(out, "Error") {
		t.Errorf("ok status should not render a status pill: %s", out)
	}
}

func TestNode_ActionMenuRenders(t *testing.T) {
	menuMarker := templ.Raw(`<div data-test-action-menu="yes"></div>`)
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:         "n",
		Title:      "T",
		ActionMenu: menuMarker,
	}))
	if !strings.Contains(out, `data-test-action-menu="yes"`) {
		t.Errorf("expected action menu slot to render its content: %s", out)
	}
}

func TestNode_AttrsSpread(t *testing.T) {
	out := renderComponent(t, canvas.Node(canvas.NodeOptions{
		ID:    "n",
		Title: "T",
		Attrs: templ.Attributes{"hx-get": "/test", "data-foo": "bar"},
	}))
	if !strings.Contains(out, `hx-get="/test"`) {
		t.Errorf("expected hx-get attribute: %s", out)
	}
	if !strings.Contains(out, `data-foo="bar"`) {
		t.Errorf("expected data-foo attribute: %s", out)
	}
}

func TestEdgeButton_RendersWithAriaLabel(t *testing.T) {
	out := renderComponent(t, canvas.EdgeButton(canvas.EdgeButtonOptions{
		AriaLabel: "Agregar nodo después de X",
	}))
	if !strings.Contains(out, `Agregar nodo`) {
		t.Errorf("expected aria-label text: %s", out)
	}
	if !strings.Contains(out, `data-kiban-canvas-edge-button`) {
		t.Errorf("expected data-kiban-canvas-edge-button marker: %s", out)
	}
	if !strings.Contains(out, "<svg") {
		t.Errorf("expected + icon: %s", out)
	}
}

func TestEdgeButton_AttrsSpread(t *testing.T) {
	out := renderComponent(t, canvas.EdgeButton(canvas.EdgeButtonOptions{
		AriaLabel: "x",
		Attrs:     templ.Attributes{"hx-get": "/add", "hx-target": "#wf-canvas"},
	}))
	if !strings.Contains(out, `hx-get="/add"`) {
		t.Errorf("expected hx-get on button: %s", out)
	}
	if !strings.Contains(out, `hx-target="#wf-canvas"`) {
		t.Errorf("expected hx-target on button: %s", out)
	}
}
