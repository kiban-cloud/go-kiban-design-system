package layout

import (
	"context"
	_ "embed"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/a-h/templ"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// Shell «kiban cloud flamboyant»: el mismo que estrenó la página de link
// (microservices view/flamboyant), ahora en el Layout común para que lo tengan
// todas las herramientas. Con la bandera prendida, Layout cambia la topbar y el
// menú lateral; el contenido de cada pantalla no cambia.
//
//   - La plataforma se llama kiban: lockup isotipo + «kiban» (skill de marca),
//     y un botón de panel como el de la app de Claude que expande y contrae el
//     menú. El espacio pasa a una migaja junto al logo, con buscador.
//   - Topbar blanca con línea azul kiban de 1 px; a la derecha el selector
//     segmentado Producción | Sandbox (si la herramienta tiene sandbox) y el
//     avatar.
//   - Menú lateral en español, una palabra por módulo (Base, CRM, Solicitudes,
//     Pagos, Cumplimiento, Fuentes; Cartera y Monitoreo, por ahora ocultos),
//     con íconos Lucide de trazo 1.5, contraído a íconos por defecto. Abajo Documentación,
//     Referencia API, Desarrolladores y Consultoría.
//   - «Pendientes» (las notificaciones) es el primer renglón del menú, con un
//     panel del 50 % con búsqueda y filtros. La campana sale de la topbar.
//   - El menú del módulo es un pilar azul cloud claro sin líneas verticales.
//   - Todas las fechas como «28 oct 2026 12:02» (flamboyant/fechas.js).
//
// Apagada la bandera, Layout queda exactamente como hoy.
const FlagShellFlamboyant = flamboyant.Flag

// FlagShellFlamboyantSwitch muestra en el menú del usuario la opción de pasar
// de un diseño al otro. Se prende para toda compañía que tiene flamboyant
// disponible, haya elegido el usuario el diseño nuevo o no; FlagShellFlamboyant
// sólo se prende si además lo eligió (middleware/authcookie).
const FlagShellFlamboyantSwitch = "shell-flamboyant-switch"

// FlamboyantFlags son las banderas del shell.
var FlamboyantFlags = []string{FlagShellFlamboyant}

// La elección del usuario vive en una cookie del dominio, path /, así que vale
// para todas las herramientas. La escribe el switch en el navegador (base.templ)
// y la lee middleware/authcookie. Sin cookie se ve el diseño actual.
const (
	UICookie     = "kiban_ui"
	UIFlamboyant = "flamboyant"
	UIClasico    = "clasico"
)

// FlagMonitoreo muestra Monitoreo (folo) en el menú del shell. Entra en una fase
// propia, después de flamboyant: el módulo es nuevo y se revisa aparte.
const FlagMonitoreo = "modulo-monitoreo"

// MonitoreoHref es la entrada del módulo cuando la herramienta no la declara en
// Config.Tools.
const MonitoreoHref = "/monitoreo"

// ShellFlamboyant dice si el request pinta el shell flamboyant. En modo
// Comparar también: el shell envuelve la pantalla entera y no se pueden
// mostrar los dos lado a lado. Las vistas que se pintan dentro del shell (el
// fragmento de Pendientes, por ejemplo) lo usan para elegir su variante.
func ShellFlamboyant(ctx context.Context) bool {
	return flamboyant.On(ctx)
}

var (
	//go:embed flamboyant/shell.css
	flamboyantCSS string
	//go:embed flamboyant/shell.js
	flamboyantJS string
	//go:embed flamboyant/fechas.js
	flamboyantFechasJS string
	//go:embed flamboyant/isotipo-1.png
	isotipo1 []byte
	//go:embed flamboyant/isotipo-2.png
	isotipo2 []byte
	//go:embed flamboyant/isotipo-3.png
	isotipo3 []byte
	//go:embed flamboyant/wordmark.png
	wordmark []byte
)

func pngURI(b []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
}

// El isotipo son tres piezas encimadas, como en el skill de marca.
var (
	isotipoURIs = []string{pngURI(isotipo1), pngURI(isotipo2), pngURI(isotipo3)}
	wordmarkURI = pngURI(wordmark)
)

// Íconos web del brandbook: Lucide a 20 px con trazo 1.5. Sólo el shell los usa
// con este trazo; el resto de la aplicación se queda en trazo 2.
const lucide = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">`

var fsIcons = map[string]string{
	"house":            `<path d="M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8"/><path d="M3 10a2 2 0 0 1 .709-1.528l7-5.999a2 2 0 0 1 2.582 0l7 5.999A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>`,
	"contact":          `<path d="M16 2v2"/><path d="M7 22v-2a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2"/><path d="M8 2v2"/><circle cx="12" cy="11" r="3"/><rect x="3" y="4" width="18" height="18" rx="2"/>`,
	"clipboard-check":  `<rect width="8" height="4" x="8" y="2" rx="1" ry="1"/><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><path d="m9 14 2 2 4-4"/>`,
	"briefcase":        `<rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>`,
	"credit-card":      `<rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/>`,
	"activity":         `<polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>`,
	"shield-check":     `<path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"/><path d="m9 12 2 2 4-4"/>`,
	"blocks":           `<rect width="7" height="7" x="14" y="3" rx="1"/><path d="M10 21V8a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-5a1 1 0 0 0-1-1H3"/>`,
	"book-open":        `<path d="M12 7v14"/><path d="M3 18a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h5a4 4 0 0 1 4 4 4 4 0 0 1 4-4h5a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-6a3 3 0 0 0-3 3 3 3 0 0 0-3-3z"/>`,
	"braces":           `<path d="M8 3H7a2 2 0 0 0-2 2v5a2 2 0 0 1-2 2 2 2 0 0 1 2 2v5c0 1.1.9 2 2 2h1"/><path d="M16 21h1a2 2 0 0 0 2-2v-5c0-1.1.9-2 2-2a2 2 0 0 1-2-2V5a2 2 0 0 0-2-2h-1"/>`,
	"square-terminal":  `<path d="m7 11 2-2-2-2"/><path d="M11 13h4"/><rect width="18" height="18" x="3" y="3" rx="2" ry="2"/>`,
	"handshake":        `<path d="m11 17 2 2a1 1 0 1 0 3-3"/><path d="m14 14 2.5 2.5a1 1 0 1 0 3-3l-3.88-3.88a3 3 0 0 0-4.24 0l-.88.88a1 1 0 1 1-3-3l2.81-2.81a5.79 5.79 0 0 1 7.06-.87l.47.28a2 2 0 0 0 1.42.25L21 4"/><path d="m21 3 1 11h-2"/><path d="M3 3 2 14l6.5 6.5a1 1 0 1 0 3-3"/><path d="M3 4h8"/>`,
	"bell":             `<path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"/><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"/>`,
	"panel-open":       `<rect width="18" height="18" x="3" y="3" rx="2"/><path d="M9 3v18"/><path d="m14 9 3 3-3 3"/>`,
	"panel-close":      `<rect width="18" height="18" x="3" y="3" rx="2"/><path d="M9 3v18"/><path d="m16 15-3-3 3-3"/>`,
	"chevrons-up-down": `<path d="m7 15 5 5 5-5"/><path d="m7 9 5-5 5 5"/>`,
	"search":           `<circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>`,
	"check":            `<path d="M20 6 9 17l-5-5"/>`,
	"check-check":      `<path d="M18 6 7 17l-5-5"/><path d="m22 10-7.5 7.5L13 16"/>`,
	"x":                `<path d="M18 6 6 18"/><path d="m6 6 12 12"/>`,
	"undo-2":           `<path d="M9 14 4 9l5-5"/><path d="M4 9h10.5a5.5 5.5 0 0 1 5.5 5.5a5.5 5.5 0 0 1-5.5 5.5H11"/>`,
	"user":             `<circle cx="12" cy="8" r="5"/><path d="M20 21a8 8 0 0 0-16 0"/>`,
	"log-out":          `<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/>`,
	"arrow-up-right":   `<path d="M7 7h10v10"/><path d="M7 17 17 7"/>`,
	"settings":         `<path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/>`,
	"file-text":        `<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/>`,
	// Los del resumen de rekon flamboyant (ciclo del dinero y tarjetas).
	"play":               `<polygon points="6 3 20 12 6 21 6 3"/>`,
	"circle-check":       `<circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/>`,
	"circle-x":           `<circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/>`,
	"chart-line":         `<path d="M3 3v16a2 2 0 0 0 2 2h16"/><path d="m19 9-5 5-4-4-3 3"/>`,
	"chart-pie":          `<path d="M21.21 15.89A10 10 0 1 1 8 2.83"/><path d="M22 12A10 10 0 0 0 12 2v10z"/>`,
	"workflow":           `<rect width="8" height="8" x="3" y="3" rx="2"/><path d="M7 11v4a2 2 0 0 0 2 2h4"/><rect width="8" height="8" x="13" y="13" rx="2"/>`,
	"tag":                `<path d="M12.586 2.586A2 2 0 0 0 11.172 2H4a2 2 0 0 0-2 2v7.172a2 2 0 0 0 .586 1.414l8.704 8.704a2.426 2.426 0 0 0 3.42 0l6.58-6.58a2.426 2.426 0 0 0 0-3.42z"/><circle cx="7.5" cy="7.5" r=".5" fill="currentColor"/>`,
	"git-branch":         `<path d="M6 3v12"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/>`,
	"arrow-left":         `<path d="m12 19-7-7 7-7"/><path d="M19 12H5"/>`,
	"chevron-down":       `<path d="m6 9 6 6 6-6"/>`,
	"triangle-alert":     `<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/><path d="M12 9v4"/><path d="M12 17h.01"/>`,
	"clipboard":          `<rect width="8" height="4" x="8" y="2" rx="1" ry="1"/><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>`,
	"cpu":                `<rect width="16" height="16" x="4" y="4" rx="2"/><rect width="6" height="6" x="9" y="9" rx="1"/><path d="M15 2v2"/><path d="M15 20v2"/><path d="M2 15h2"/><path d="M2 9h2"/><path d="M20 15h2"/><path d="M20 9h2"/><path d="M9 2v2"/><path d="M9 20v2"/>`,
	"file-up":            `<path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M12 12v6"/><path d="m15 15-3-3-3 3"/>`,
	"git-merge":          `<circle cx="18" cy="18" r="3"/><circle cx="6" cy="6" r="3"/><path d="M6 21V9a9 9 0 0 0 9 9"/>`,
	"plug":               `<path d="M12 22v-5"/><path d="M9 8V2"/><path d="M15 8V2"/><path d="M18 8v5a4 4 0 0 1-4 4h-4a4 4 0 0 1-4-4V8Z"/>`,
	"sliders-horizontal": `<line x1="21" x2="14" y1="4" y2="4"/><line x1="10" x2="3" y1="4" y2="4"/><line x1="21" x2="12" y1="12" y2="12"/><line x1="8" x2="3" y1="12" y2="12"/><line x1="21" x2="16" y1="20" y2="20"/><line x1="12" x2="3" y1="20" y2="20"/><line x1="14" x2="14" y1="2" y2="6"/><line x1="8" x2="8" y1="10" y2="14"/><line x1="16" x2="16" y1="18" y2="22"/>`,
	"history":            `<path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/><path d="M12 7v5l4 2"/>`,
	"clock":              `<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>`,
	"plus":               `<path d="M5 12h14"/><path d="M12 5v14"/>`,
	"copy":               `<rect width="14" height="14" x="8" y="8" rx="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>`,
	"eye":                `<path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"/><circle cx="12" cy="12" r="3"/>`,
	"trash-2":            `<path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/>`,
	"eye-off":            `<path d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49"/><path d="M14.084 14.158a3 3 0 0 1-4.242-4.242"/><path d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143"/><path d="m2 2 20 20"/>`,
	"type":               `<polyline points="4 7 4 4 20 4 20 7"/><line x1="9" x2="15" y1="20" y2="20"/><line x1="12" x2="12" y1="4" y2="20"/>`,
	"calendar":           `<path d="M8 2v4"/><path d="M16 2v4"/><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M3 10h18"/>`,
	"hash":               `<line x1="4" x2="20" y1="9" y2="9"/><line x1="4" x2="20" y1="15" y2="15"/><line x1="10" x2="8" y1="3" y2="21"/><line x1="16" x2="14" y1="3" y2="21"/>`,
	"list":               `<path d="M3 12h.01"/><path d="M3 18h.01"/><path d="M3 6h.01"/><path d="M8 12h13"/><path d="M8 18h13"/><path d="M8 6h13"/>`,
	"folder":             `<path d="M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z"/>`,
	"zap":                `<path d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z"/>`,
	"lock":               `<rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>`,
	"user-check":         `<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><polyline points="16 11 18 13 22 9"/>`,
	"hourglass":          `<path d="M5 22h14"/><path d="M5 2h14"/><path d="M17 22v-4.172a2 2 0 0 0-.586-1.414L12 12l-4.414 4.414A2 2 0 0 0 7 17.828V22"/><path d="M7 2v4.172a2 2 0 0 0 .586 1.414L12 12l4.414-4.414A2 2 0 0 0 17 6.172V2"/>`,
	"layers":             `<path d="m12.83 2.18a2 2 0 0 0-1.66 0L2.6 6.08a1 1 0 0 0 0 1.83l8.58 3.91a2 2 0 0 0 1.66 0l8.58-3.9a1 1 0 0 0 0-1.83Z"/><path d="m22 17.65-9.17 4.16a2 2 0 0 1-1.66 0L2 17.65"/><path d="m22 12.65-9.17 4.16a2 2 0 0 1-1.66 0L2 12.65"/>`,
	"map-pin":            `<path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"/><circle cx="12" cy="10" r="3"/>`,
	"key-round":          `<path d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z"/><circle cx="16.5" cy="7.5" r=".5" fill="currentColor"/>`,
	"flag":               `<path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/><line x1="4" x2="4" y1="22" y2="15"/>`,
	"maximize":           `<path d="M8 3H5a2 2 0 0 0-2 2v3"/><path d="M21 8V5a2 2 0 0 0-2-2h-3"/><path d="M3 16v3a2 2 0 0 0 2 2h3"/><path d="M16 21h3a2 2 0 0 0 2-2v-3"/>`,
	"list-checks":        `<path d="m3 17 2 2 4-4"/><path d="m3 7 2 2 4-4"/><path d="M13 6h8"/><path d="M13 12h8"/><path d="M13 18h8"/>`,
	"grip-vertical":      `<circle cx="9" cy="12" r="1"/><circle cx="9" cy="5" r="1"/><circle cx="9" cy="19" r="1"/><circle cx="15" cy="12" r="1"/><circle cx="15" cy="5" r="1"/><circle cx="15" cy="19" r="1"/>`,
	"arrow-up":           `<path d="m5 12 7-7 7 7"/><path d="M12 19V5"/>`,
	"arrow-down":         `<path d="M12 5v14"/><path d="m19 12-7 7-7-7"/>`,
	"filter":             `<polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>`,
	"alert-circle":       `<circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/>`,
	"corner-down-right":  `<polyline points="15 10 20 15 15 20"/><path d="M4 4v7a4 4 0 0 0 4 4h12"/>`,
	"trash":              `<path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>`,
	"split":              `<path d="M16 3h5v5"/><path d="M8 3H3v5"/><path d="M12 22v-8.3a4 4 0 0 0-1.172-2.872L3 3"/><path d="m15 9 6-6"/>`,
	"trophy":             `<path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/><path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/><path d="M4 22h16"/><path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/><path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/><path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>`,
	"external-link":      `<path d="M15 3h6v6"/><path d="M10 14 21 3"/><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>`,
	"git-fork":           `<circle cx="12" cy="18" r="3"/><circle cx="6" cy="6" r="3"/><circle cx="18" cy="6" r="3"/><path d="M18 9v2c0 .6-.4 1-1 1H7c-.6 0-1-.4-1-1V9"/><path d="M12 12v3"/>`,
	"download":           `<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>`,
	"send":               `<line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/>`,
	"chevron-right":      `<polyline points="9 18 15 12 9 6"/>`,
}

// FlamboyantIcon devuelve el ícono Lucide del shell (trazo 1.5) por su nombre
// en Lucide. Un nombre que el shell no conoce devuelve un componente vacío.
func FlamboyantIcon(name string) templ.Component {
	return templ.Raw(fsIcon(name))
}

func fsIcon(name string) string {
	body, ok := fsIcons[name]
	if !ok {
		return ""
	}
	return lucide + body + `</svg>`
}

// fsModulo es un renglón del menú del shell flamboyant. Sin Href (una
// herramienta que no declara ese enlace) se ve, con tooltip, pero no navega.
type fsModulo struct {
	Key   string
	Label string
	Icon  string // nombre Lucide
	Href  string
	Ext   bool // abre en otra pestaña (documentación)
}

// modulosFlamboyant es el menú, una palabra por módulo (la etiqueta sale de
// shell.module.<Key> en locales/). Key es la misma clave del rail de hoy
// (Config.Tools), de ahí sale el Href, así que cada herramienta sigue mandando
// a sus propias rutas. Cartera (core) y Monitoreo (folo) son módulos nuevos
// sin ruta todavía; por ahora no se muestran (fsOcultos).
var modulosFlamboyant = []fsModulo{
	{Key: "home", Icon: "house"},
	{Key: "customers", Icon: "contact"},
	{Key: "workfloo", Icon: "clipboard-check"},
	{Key: "core", Icon: "briefcase"},
	{Key: "rekon", Icon: "credit-card"},
	{Key: "folo", Icon: "activity"},
	{Key: "klin", Icon: "shield-check"},
	{Key: "link", Icon: "blocks"},
}

// fsColores son el color de cada módulo (el del ciclo del menú de Base en
// kiban-proto): el ícono y su cuadro en la ficha del cliente.
var fsColores = map[string]string{
	"customers": "#5B82F5", "workfloo": "#0047FF", "link": "#0E9FC2",
	"rekon": "#8000FF", "klin": "#00A884", "core": "#8000FF", "folo": "#5B82F5", "home": "#0047FF",
}

// FlamboyantModuloColor es el color del módulo (misma clave que
// FlamboyantModulo); gris si no es ninguno.
func FlamboyantModuloColor(key string) string {
	switch key {
	case "kiban-cloud", "kibancloud":
		key = "home"
	case "crm":
		key = "customers"
	}
	if c, ok := fsColores[key]; ok {
		return c
	}
	return "#606875"
}

// moduleTileStyle pinta el cuadro: color del módulo al 12 % de fondo, o gris
// cuando el módulo todavía no tiene nada que mostrar (dim).
func moduleTileStyle(key string, size int, dim bool) string {
	px := strconv.Itoa(size) + "px"
	bg, fg := "#F5F6F9", "#9096A3"
	if !dim {
		fg = FlamboyantModuloColor(key)
		bg = fg + "1F"
	}
	return "width:" + px + ";height:" + px + ";background:" + bg + ";color:" + fg
}

// iconTileStyle es el tamaño del cuadro de IconTile.
func iconTileStyle(size int) string {
	px := strconv.Itoa(size) + "px"
	return "width:" + px + ";height:" + px
}

func moduleIcon(key string) string {
	_, icon, _ := FlamboyantModulo(context.Background(), key)
	return icon
}

// FlamboyantModulo traduce la clave de una herramienta (la del rail, o la del
// módulo emisor de una notificación: "kiban-cloud", "crm"…) a su nombre y su
// ícono en el menú flamboyant. ok es false si no es ninguno de los ocho.
func FlamboyantModulo(ctx context.Context, key string) (label, icon string, ok bool) {
	switch key {
	case "kiban-cloud", "kibancloud":
		key = "home"
	case "crm":
		key = "customers"
	}
	for _, m := range modulosFlamboyant {
		if m.Key == key {
			return docsCatalog.T(ctx, "shell.module."+m.Key), m.Icon, true
		}
	}
	return "", "", false
}

// fsTools arma el menú de módulos con los Href del rail de la herramienta. Las
// entradas del rail que no son ninguno de los ocho módulos (el «prototipo» del
// proto, por ejemplo) van al final con su ícono de siempre.
// fsOcultos son los módulos del menú que por ahora no se muestran: Cartera y
// Monitoreo todavía no existen. Para volver a mostrarlos, se quitan de aquí.
// Monitoreo se muestra igual cuando su bandera está prendida (FlagMonitoreo).
var fsOcultos = map[string]bool{"core": true, "folo": true}

func fsTools(ctx context.Context, cfg Config) ([]fsModulo, []Tool) {
	hrefs := map[string]string{"folo": MonitoreoHref}
	for _, t := range cfg.Tools {
		if t.Href != "" {
			hrefs[t.Key] = t.Href
		}
	}
	monitoreo := feature.Comparing(ctx) || feature.Enabled(ctx, FlagMonitoreo)
	out := make([]fsModulo, 0, len(modulosFlamboyant))
	for _, m := range modulosFlamboyant {
		if fsOcultos[m.Key] && !(m.Key == "folo" && monitoreo) {
			continue
		}
		m.Label = docsCatalog.T(ctx, "shell.module."+m.Key)
		m.Href = flamboyant.Ruta(ctx, hrefs[m.Key])
		out = append(out, m)
	}
	var extra []Tool
	for _, t := range cfg.Tools {
		if _, _, ok := FlamboyantModulo(ctx, t.Key); !ok {
			t.Href = flamboyant.Ruta(ctx, t.Href)
			extra = append(extra, t)
		}
	}
	return out, extra
}

// Configuración es una sección de la página de Base y Configuración de kiban
// cloud (kibancloud view/flamboyant); se abre por su hash.
const FlamboyantConfigHref = DefaultLogoHref + "/#inicio"

// fsConfig es Configuración, la primera entrada de la parte baja del menú,
// separada de las ayudas por una línea azul (como Pendientes arriba).
func fsConfig(ctx context.Context) fsModulo {
	return fsModulo{Key: "config", Label: docsCatalog.T(ctx, "shell.config"), Icon: "settings", Href: FlamboyantConfigHref}
}

// fsDocs es el resto de la parte baja del menú. Los Href salen de Config.Docs
// (por clave, lo que arma DocsMenu) y caen a los destinos globales de kiban si
// la herramienta no los declara. Desarrolladores sigue la regla de DocsMenu:
// sin Href en Config.Docs no se muestra, porque un enlace a un 404 es peor que
// ningún enlace.
func fsDocs(ctx context.Context, cfg Config) []fsModulo {
	hrefs := map[string]string{
		"docs-kiban":  KnowledgeCenterHref,
		"docs-api":    APIDocsHref,
		"consultancy": ConsultancyHref,
	}
	for _, d := range cfg.Docs {
		if d.Href != "" {
			hrefs[d.Key] = d.Href
		}
	}
	docs := []fsModulo{
		{Key: "docs-kiban", Label: docsCatalog.T(ctx, "shell.docs.knowledge"), Icon: "book-open", Href: hrefs["docs-kiban"], Ext: true},
		{Key: "docs-api", Label: docsCatalog.T(ctx, "shell.docs.api"), Icon: "braces", Href: hrefs["docs-api"], Ext: true},
	}
	if dev := hrefs[DevelopersToolKey]; dev != "" {
		docs = append(docs, fsModulo{Key: DevelopersToolKey, Label: docsCatalog.T(ctx, "nav.docs.developers"), Icon: "square-terminal", Href: dev})
	}
	return append(docs, fsModulo{Key: "consultancy", Label: docsCatalog.T(ctx, "nav.docs.consultancy"), Icon: "handshake", Href: hrefs["consultancy"], Ext: true})
}

// fsModuloActivo es el nombre y el ícono del pilar del módulo. Una herramienta
// que no es ninguno de los ocho módulos usa su SectionLabel.
func fsModuloActivo(ctx context.Context, cfg Config) (string, string) {
	if label, icon, ok := FlamboyantModulo(ctx, cfg.ToolKey); ok {
		return label, icon
	}
	label := cfg.SectionLabel
	if label != "" {
		r := []rune(label)
		r[0] = unicode.ToUpper(r[0])
		label = string(r)
	}
	return label, "file-text"
}

// fsSpaces es la lista del menú de espacios: los de Config.Spaces o, si la
// herramienta no los manda, sólo el actual.
func fsSpaces(cfg Config) []SpaceOption {
	if len(cfg.Spaces) > 0 {
		return cfg.Spaces
	}
	return []SpaceOption{{Id: cfg.SpaceID, Name: spaceChipLabel(cfg)}}
}

// fsSpaceInitials son las dos letras del avatar del espacio.
func fsSpaceInitials(name string) string {
	words := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	var out []rune
	for _, w := range words {
		out = append(out, []rune(w)[0])
		if len(out) == 2 {
			break
		}
	}
	if len(out) == 1 && len([]rune(words[0])) > 1 {
		out = append(out, []rune(words[0])[1])
	}
	if len(out) == 0 {
		return "·"
	}
	return strings.ToUpper(string(out))
}

// Cada espacio tiene su color, por posición en la lista.
var fsSpaceTints = [][2]string{
	{"#E8EEFF", "#0036C4"}, {"#E6FBF6", "#007A63"}, {"#F2E8FF", "#6A00D6"},
	{"#EEF3FF", "#33406B"}, {"#FFF1E6", "#9A4A00"},
}

func fsSpaceStyle(cfg Config, id string, size int) string {
	i := 0
	for k, s := range fsSpaces(cfg) {
		if s.Id == id {
			i = k
		}
	}
	t := fsSpaceTints[i%len(fsSpaceTints)]
	px := strconv.Itoa(size) + "px"
	return "width:" + px + ";height:" + px + ";background:" + t[0] + ";color:" + t[1]
}

func fsSpaceSuffix(id string) string {
	if len(id) > 6 {
		return id[len(id)-6:]
	}
	return id
}

// fsUserItems es el menú del avatar: los de la herramienta o, si no manda
// ninguno, Tu perfil.
func fsUserItems(ctx context.Context, cfg Config) []UserMenuItem {
	if len(cfg.UserMenuItems) > 0 {
		return cfg.UserMenuItems
	}
	return []UserMenuItem{{Label: docsCatalog.T(ctx, "shell.profile"), Href: "/kiban-cloud/profile"}}
}

// fsURL es un enlace del shell con el nombre de módulo del diseño
// (flamboyant.Ruta): /workfloo/… sale como /solicitudes/….
func fsURL(ctx context.Context, path string) templ.SafeURL {
	return templ.URL(flamboyant.Ruta(ctx, path))
}

func fsSubItemClass(cfg Config, item SubItem) string {
	if cfg.ActiveKey == item.Key {
		return "fs-sub fs-sub-on"
	}
	return "fs-sub"
}

// fsMeses son las doce abreviaturas de mes del idioma del request, separadas
// por coma. fechas.js las recibe tal cual en data-fs-months.
func fsMeses(ctx context.Context) string {
	return docsCatalog.T(ctx, "shell.dates.months")
}

// FlamboyantFecha es el formato de fecha de flamboyant: «28 oct 2026 12:02»,
// o «28 oct 2026» sin la hora. El mismo que aplica flamboyant/fechas.js en el
// cliente.
func FlamboyantFecha(ctx context.Context, t time.Time, conHora bool) string {
	meses := strings.Split(fsMeses(ctx), ",")
	s := strconv.Itoa(t.Day()) + " " + meses[t.Month()-1] + " " + strconv.Itoa(t.Year())
	if conHora {
		s += " " + t.Format("15:04")
	}
	return s
}

// InboxToggleAttrs son los atributos que abren y cierran el panel
// «Pendientes» del shell flamboyant desde cualquier botón de la página (el
// «27 por revisar» del hero de un inicio), igual que el botón de la campana
// del menú: lo atiende flamboyant/shell.js y htmx carga la lista.
// notificationsBaseURL vacío usa DefaultNotificationsBaseURL.
func InboxToggleAttrs(notificationsBaseURL string) templ.Attributes {
	base := notificationsBaseURL
	if base == "" {
		base = DefaultNotificationsBaseURL
	}
	return templ.Attributes{
		"data-fs-inbox-toggle": "",
		"aria-controls":        "fs-inbox",
		"hx-get":               base + "/dropdown",
		"hx-target":            "#fs-inbox-list",
		"hx-swap":              "innerHTML",
		"hx-trigger":           "click",
	}
}
