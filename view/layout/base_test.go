package layout_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// renderBase runs layout.Base (the HTML shell shared by every kiban module)
// to a string so we can assert on the global head/body wiring. Base takes
// children; we pass an empty component since the assertions here are about
// the shell itself.
func renderBase(t *testing.T, cfg layout.Config) string {
	t.Helper()
	empty := templ.ComponentFunc(func(_ context.Context, _ io.Writer) error { return nil })
	var buf bytes.Buffer
	if err := layout.Base(cfg).Render(templ.WithChildren(context.Background(), empty), &buf); err != nil {
		t.Fatalf("Base render: %v", err)
	}
	return buf.String()
}

// The global HTMX busy-region rule must be wired into the shared shell so
// every module (link / rekon / crm / klin) gets consistent "blocked +
// something is happening" feedback on a partial swap with no per-button
// wiring. This guards the three moving parts: the CSS class, the spinner
// overlay, and the JS that flags the swap target on htmx:beforeRequest /
// clears it on htmx:afterRequest.
func TestBase_DataHrefRowShowsNavLoader(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	// The data-href row-click delegation must trigger the shared
	// #nav-loader overlay before navigating, so list→detail row clicks
	// aren't a frozen wait.
	assert.Contains(t, body, "getElementById('nav-loader')")
	assert.Contains(t, body, "navLoader.classList.add('is-visible')")
	assert.Contains(t, body, "window.location.href = el.getAttribute('data-href')")
}

// The shared HTMX error-response pattern (see the htmxerror package) needs
// three things wired into the shell so every module gets non-200 errors that
// Cloud Monitoring can see: the response-targets extension loaded and enabled
// on <body>, the #global-error-banner landing zone for 5xx fragments, and the
// beforeRequest cleanup that clears stale feedback so a new error never shows
// next to an old one.
func TestBase_HtmxErrorResponsePattern(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	// Extension script loaded and enabled on the body so hx-target-4xx /
	// hx-target-5xx resolve.
	assert.Contains(t, body, "htmx-ext-response-targets@2.0.4")
	assert.Contains(t, body, `hx-ext="response-targets"`)

	// Landing zone for 5xx fragments routed via hx-target-5xx.
	assert.Contains(t, body, `id="global-error-banner"`)

	// Cleanup: clear the banner and any data-error-type / data-success-type
	// feedback before each request.
	assert.Contains(t, body, "htmx:beforeRequest")
	assert.Contains(t, body, "[data-error-type],[data-success-type]")
}

func TestBase_GlobalHtmxBusyRule(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:       "T",
		ProjectName: "kiban",
	})

	// CSS: the busy region is non-interactive and gets a spinner overlay.
	assert.Contains(t, body, ".kiban-busy")
	assert.Contains(t, body, "pointer-events: none")
	assert.Contains(t, body, ".kiban-busy::after")

	// Trigger gets disabled during its own request (anti double-click),
	// inherited from <body> by every htmx control.
	assert.Contains(t, body, `hx-disabled-elt="this"`)

	// JS: flag the swap target on request start, clear it when the request
	// finishes (success or error).
	assert.Contains(t, body, "htmx:beforeRequest")
	assert.Contains(t, body, "htmx:afterRequest")
	assert.Contains(t, body, "classList.add('kiban-busy')")
	assert.Contains(t, body, "classList.remove('kiban-busy', 'kiban-busy-quiet')")

	// Una petición con su propio hx-indicator (el overlay de carga de una
	// lista) no se queda además con el spinner de la zona: dos spinners al
	// paginar. La zona se sigue bloqueando, sólo pierde el ::after.
	assert.Contains(t, body, "elt.closest('[hx-indicator]')) region.classList.add('kiban-busy-quiet')")
	assert.Contains(t, body, ".kiban-busy.kiban-busy-quiet::after { content: none; }")

	// History cleanup: pagination/filter controls use hx-push-url, so htmx
	// snapshots the page into its history cache — with .kiban-busy still on
	// the swap target. Pressing Back restores that snapshot from cache with no
	// request, so afterRequest never fires and the spinner stays frozen. Strip
	// .kiban-busy before htmx serializes the snapshot (beforeHistorySave) and
	// defensively after any restore (historyRestore).
	assert.Contains(t, body, "htmx:beforeHistorySave")
	assert.Contains(t, body, "htmx:historyRestore")
}

// El menú de un filtro (filterbar) se corre al abrir para no salirse de la
// pantalla por la derecha, y se abre hacia arriba si abajo no cabe.
// El Dropzone con data-change se pinta elegido y muestra la miniatura.
func TestBase_DropzoneElegido(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.Contains(t, body, "zone.getAttribute('data-change')")
	assert.Contains(t, body, "zone.setAttribute('data-has-file', '')")
	assert.Contains(t, body, "thumb.src = URL.createObjectURL(f)")
	// Pie pendiente con lo que falta (drawer.SummaryMissing).
	assert.Contains(t, body, "function kibanPanelMissing(form, sum)")
	assert.Contains(t, body, "names = src ? (src.getAttribute('data-kiban-missing-names') || '').split(',').filter(Boolean) : [];", "From: los campos a contar salen del form")
	assert.Contains(t, body, "if (el.type === 'radio') { if (!form.querySelector('[name=\"' + name + '\"]:checked')) n++; return; }", "un grupo de radios falta si no hay uno elegido")
	assert.Contains(t, body, "if (done && num.hasAttribute('data-progress'))", "el paso completo pasa a palomita")
	assert.Contains(t, body, "var listas = form.querySelectorAll('[data-kiban-missing-names]');", "requiere-completo también espera a los obligatorios listados")
	assert.Contains(t, body, "(!e.target.name || e.target.hasAttribute('data-nodirty'))) return;", "un buscador sin name no es un cambio")
	assert.Contains(t, body, "kibanPanelSnapshotAll(root) !== root.getAttribute('data-kiban-panel-snapshot')", "es cambio sólo si difiere de la foto al abrir")
	assert.Contains(t, body, "Array.prototype.map.call(root.querySelectorAll('form'), kibanPanelSnapshot).join('|')", "la foto mira todos los forms del panel, no sólo el primero")
	assert.Contains(t, body, "if (save) save.disabled = !!(primary && primary.disabled);", "sin poder guardar, el aviso no ofrece guardar y cerrar")
	assert.Contains(t, body, "var ok = form.checkValidity() && !missing;", "lo que falta gana aunque el navegador dé el form por válido")
	assert.Contains(t, body, "missing || invalid || pending")
	assert.Contains(t, body, "form.querySelectorAll('[data-kiban-require-one]')", "al menos una casilla marcada")
	assert.Contains(t, body, "(total > 1 ? cfg.many : cfg.one) + ' ' + out.join(' ' + cfg.and + ' ')")
	// El submit que espera a los obligatorios no frena a un botón con
	// formnovalidate (guardar la razón social del RCE con su propio POST).
	assert.Contains(t, body, "button[type=submit]:not([formnovalidate]), input[type=submit]:not([formnovalidate])")
}

func TestBase_MenuDeFiltroDentroDeLaPantalla(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.Contains(t, body, "function kibanFilterPlace(d)")
	assert.Contains(t, body, "d.querySelector('[data-kiban-filter-menu]')")
	assert.Contains(t, body, "var shift = r.right - (vw - margin);")
	assert.Contains(t, body, "m.style.bottom = 'calc(100% + 6px)'")
	assert.Contains(t, body, "d.open) kibanFilterPlace(d);")
}

// Non-HTMX navigations (a form submit that posts then redirects, or a
// [data-block-ui] button) must show a full-screen blocking overlay so the
// user gets feedback + can't double-submit during the server round-trip.
// The overlay lives in Base (not just Layout) so login-style bare pages get
// it too. This guards the element, the CSS, and the submit/click handlers —
// including the skip for HTMX forms (those use the per-region block).
func TestBase_GlobalActionOverlay(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:       "T",
		ProjectName: "kiban",
	})

	// The overlay element + its CSS.
	assert.Contains(t, body, `id="kiban-action-overlay"`)
	assert.Contains(t, body, ".kiban-action-overlay")
	assert.Contains(t, body, "position: fixed")

	// JS: show on a real form submit, opt-in via [data-block-ui].
	assert.Contains(t, body, "addEventListener('submit'")
	assert.Contains(t, body, "data-block-ui")

	// HTMX forms / boosted forms are skipped (region block handles them);
	// forms can also opt out with data-no-block-ui.
	assert.Contains(t, body, "hx-boost")
	assert.Contains(t, body, "data-no-block-ui")
}

func TestBase_CustomStylesheetHref(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:                "T",
		ProjectName:          "kiban",
		CustomStylesheetHref: "/portal/theme.css",
	})

	// The per-client stylesheet is emitted as a <link> and must appear after
	// the design-system <style> block so it wins the cascade.
	assert.Contains(t, body, `<link rel="stylesheet" href="/portal/theme.css">`)
	styleEnd := strings.LastIndex(body, "</style>")
	linkAt := strings.Index(body, `href="/portal/theme.css"`)
	assert.Greater(t, linkAt, styleEnd, "custom stylesheet must come after the base <style> block")
}

func TestBase_NoCustomStylesheetByDefault(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.NotContains(t, body, `rel="stylesheet" href="/portal`)
}

// The code editor must hide its source <textarea> from CSS, never from the
// inline style CodeMirror.fromTextArea() writes.
//
// htmx's settle phase calls cloneAttributes() on every element it matches by
// id across a swap, dropping any attribute the server did not send. `style`
// is never sent, so the display:none CodeMirror sets right after
// htmx:afterSwap is stripped ~20ms later and the raw textarea reappears next
// to its own widget — two editors stacked on the same field. The wrapper
// carries no id, so a rule keyed on its ready flag survives the settle.
func TestBase_CodeEditorHidesTextareaViaCSS(t *testing.T) {
	body := renderBase(t, layout.Config{
		Title:          "T",
		ProjectName:    "kiban",
		LoadCodeMirror: true,
	})

	assert.Contains(t, body, `[data-kiban-code-editor-ready="1"] > textarea`)
	// The flag the rule hangs off must actually be set once a widget mounts.
	assert.Contains(t, body, `setAttribute('data-kiban-code-editor-ready', '1')`)
}

// El guard de "sólo guardar si cambió algo" vive en el shell, así que cada
// módulo lo hereda marcando el form con data-kiban-dirty-form. Este test
// protege las tres piezas: la foto de los valores, el re-armado tras un swap
// de htmx, y que el estado inicial lo ponga el JS.
func TestBase_DirtyFormGuard(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "data-kiban-dirty-form")
	assert.Contains(t, body, "kibanFormSnapshot")
	assert.Contains(t, body, "data-kiban-dirty-snapshot")
	// Un form intercambiado por htmx es otro nodo y con los valores ya
	// guardados: sin re-armarlo, el botón quedaría habilitado comparando
	// contra la foto vieja.
	assert.Contains(t, body, "kibanArmDirtyForm")
	assert.True(t, strings.Contains(body, "htmx:afterSwap"))
}

// El botón se deshabilita por JS y no desde el servidor: si el script no
// corriera, el form tiene que seguir siendo usable en vez de quedar bloqueado.
func TestBase_DirtyFormEnablesByDefault(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "btn.disabled = !dirty")
	assert.NotContains(t, body, `<button type="submit" disabled`)
}

// input.Money con tope: el runtime marca inválido el monto arriba de
// data-max, y un radio con data-kiban-money-max-for le pone o le quita el
// tope según la opción elegida (OXXO Pay, hasta $10,000).
func TestBase_TopeDelMonto(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "function kibanMoneyCheck(el)")
	assert.Contains(t, body, "el.setCustomValidity(!isNaN(max) && n > max ? (el.getAttribute('data-max-error') || ' ') : '');")
	assert.Contains(t, body, "String(el.value).replace(/,/g, '')", "acepta comas de miles")
	assert.Contains(t, body, "function kibanMoneySync(scope)")
	assert.Contains(t, body, "input[type=radio][data-kiban-money-max-for]")
	assert.Contains(t, body, "on.getAttribute('data-kiban-money-max')")
	assert.Contains(t, body, "el._kibanMoneyBase.max", "otra opción vuelve al tope propio del campo")
	// El paso cuenta el monto capturado aunque pase el tope; el envío no.
	assert.Contains(t, body, "function kibanOnlyOverMax(el)")
	assert.Contains(t, body, "todos[m].validity.customError) return true;")
}

// input.Picker: el buscador filtra sin acentos ni mayúsculas, nunca esconde
// la elegida y Enter no envía el form.
func TestBase_BuscadorDelPicker(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "function kibanPickerNorm(s)")
	assert.Contains(t, body, ".normalize('NFD')")
	assert.Contains(t, body, "function kibanPickerFilter(box)")
	assert.Contains(t, body, "|| (radio && radio.checked);", "la elegida no se esconde")
	assert.Contains(t, body, "if (empty) empty.hidden = any;")
	assert.Contains(t, body, "if (e.key === 'Enter' && e.target.closest && e.target.closest('[data-kiban-picker-q]')) e.preventDefault();")
	assert.Contains(t, body, "document.addEventListener('htmx:afterSettle', kibanPickerInit);")
}

// input.Chips: agrega con Enter, coma o punto y coma, sincroniza el hidden
// (uno por renglón) y avisa al form con un input.
func TestBase_Chips(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "// ---- input.Chips ([data-kiban-chips])")
	assert.Contains(t, body, "hidden.value = vals.join('\\n');")
	assert.Contains(t, body, "hidden.dispatchEvent(new Event('input', { bubbles: true }));")
	assert.Contains(t, body, "if (e.key === 'Enter' || e.key === ',' || e.key === ';')")
	assert.Contains(t, body, "box.getAttribute('data-kiban-chips') === 'email' && !reEmail.test(v)")
	assert.Contains(t, body, "have.indexOf(v.toLowerCase()) !== -1", "no repite valores")
}

// El diálogo de confirmación con tono info: ícono de aviso, botón azul y el
// foco en confirmar; data-kiban-confirm-strong pone en negritas un pedazo.
func TestBase_ConfirmacionConTono(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "var info = btn.getAttribute('data-kiban-confirm-tone') === 'info';")
	assert.Contains(t, body, "wrap.querySelector(info ? '[data-yes]' : '[data-no]').focus();")
	assert.Contains(t, body, "btn.getAttribute('data-kiban-confirm-strong')")
	assert.Contains(t, body, "pel.append(ptext.slice(0, at), sb, ptext.slice(at + strong.length));")
}

// El pie del FormPanel dice en rojo el mensaje de un campo con error propio
// (setCustomValidity) cuando no falta nada.
func TestBase_PieConMensajeDeCampoInvalido(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})

	assert.Contains(t, body, "function kibanPanelInvalidMessage(form)")
	assert.Contains(t, body, "else if (!missing) invalid = kibanPanelInvalidMessage(form);")
	assert.Contains(t, body, "if (el.closest && el.closest('[data-kiban-require-one]')) continue;", "la casilla de require-one no sale en rojo")
	assert.Contains(t, body, "sum.toggleAttribute('data-error', !!invalid);")
	// La flecha de un button.Split sigue al botón principal.
	assert.Contains(t, body, "form.querySelectorAll('[data-kiban-split-trigger]')")
}

// El texto de varias opciones de un filtro sale de data-many.
func TestBase_FiltroVariasOpciones(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.Contains(t, body, "(d.getAttribute('data-many') || '{n} seleccionados').split('{n}').join(vals.length)")
}

// hero.Cycle: animaciones del ciclo, quietas con prefers-reduced-motion.
func TestBase_AnimacionesDelHero(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.Contains(t, body, "@keyframes kiban-cy-in")
	assert.Contains(t, body, ".kiban-cy-flow {")
	assert.Contains(t, body, "@media (prefers-reduced-motion: reduce) { .kiban-cy-in { animation: none; } .kiban-cy-flow { display: none; } }")
}

// Flamboyant carga además el peso 600 de Be Vietnam Pro.
func TestBase_PesoSemiboldSoloEnFlamboyant(t *testing.T) {
	assert.Contains(t, renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"}), "Be+Vietnam+Pro:wght@400;500;700&")

	empty := templ.ComponentFunc(func(_ context.Context, _ io.Writer) error { return nil })
	var buf bytes.Buffer
	require.NoError(t, layout.Base(layout.Config{Title: "T", ProjectName: "kiban"}).Render(templ.WithChildren(flags(layout.FlagShellFlamboyant), empty), &buf))
	assert.Contains(t, buf.String(), "Be+Vietnam+Pro:wght@400;500;600;700&")
}

// data-kiban-geo: la ubicación del navegador viaja en el request htmx.
func TestBase_Ubicacion(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.Contains(t, body, "elt.closest('[data-kiban-geo]')")
	assert.Contains(t, body, "geo.getAttribute('data-kiban-geo-lat') || 'latitude'")
	assert.Contains(t, body, "navigator.geolocation.getCurrentPosition")
}

// Un panel que abre sobre otro no oscurece el fondo y es 5vw más angosto
// que el de abajo; uno de 50 % sobre uno ancho conserva su ancho.
func TestBase_PanelesApilados(t *testing.T) {
	body := renderBase(t, layout.Config{Title: "T", ProjectName: "kiban"})
	assert.Contains(t, body, "[data-kiban-panel] ~ [data-kiban-panel] > [data-kiban-panel-close]:not(button)", "el de arriba sin velo")
	assert.Contains(t, body, "[data-kiban-panel] ~ .kiban-panel-ph > .kiban-panel-ph-bg", "la hoja de carga apilada tampoco oscurece")
	assert.Contains(t, body, "[data-kiban-panel] ~ [data-kiban-panel] > [role=dialog],")
	assert.Contains(t, body, "> .kiban-panel-ph-sheet { width: 90vw; }")
	assert.Contains(t, body, "> .kiban-panel-ph-sheet { width: 85vw; }", "un tercer nivel, otros 5vw")
	assert.Contains(t, body, "[data-kiban-panel]:not([data-half]) ~ [data-kiban-panel][data-half] > [role=dialog]", "el de 50 % sobre uno ancho no crece")
	assert.Contains(t, body, "> .kiban-panel-ph-sheet { width: 45vw; }", "sobre uno de 50 %, 45vw")
}
