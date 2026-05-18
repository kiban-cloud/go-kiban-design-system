# CLAUDE.md — go-kiban-design-system

Sistema de diseño compartido entre los proyectos kiban basados en HTMX + templ (rekon, crm, futuros). Provee componentes templ reusables, helpers de form binding, middlewares HTMX, y la shell de layout (sidebar, topbar, slide animation entre niveles 1/2 estilo GCP).

## Cuándo usar este paquete

**SIEMPRE** que vayas a crear un componente UI (templ, HTML, HTMX, CSS) en rekon o crm:

1. **Busca aquí primero.** Si el componente ya existe, importalo y úsalo.
2. **Si no existe pero tiene sentido reusar entre proyectos**, créalo en este paquete y consúmelo desde el proyecto.
3. **Sólo lo realmente específico** de un proyecto (vistas de páginas individuales, request structs, controllers, lógica de negocio) vive en el repo de ese proyecto.

## Cuándo NO

- Lógica de negocio, controllers HTTP, repositorios, usecases, modelos de dominio.
- Vistas de páginas concretas (ej. la página "lista de clientes" de crm).
- Configuración runtime de un proyecto.
- Cualquier cosa que no sea reutilizable entre proyectos.

## Módulo y dependencias

- **Módulo Go**: `github.com/kiban-cloud/go-kiban-design-system`
- **Versión Go**: 1.24+ (para el `tool` directive en `go.mod`).
- **Stack runtime**: templ + Gin + HTMX (CDN) + Tailwind (CDN, con tokens kiban en `view/layout/base.templ`) + intl-tel-input (CDN). Sin daisyUI por ahora — los componentes son hand-rolled con primitives de Tailwind.
- **Sin assets locales**. Todo lo externo (HTMX, intl-tel-input, Chart.js si fuera necesario) se carga vía CDN, pinned a versión específica.

## Estructura

```
view/
  render.go             Render(c, status, component) — emite text/html, idéntico en todos los consumidores
  layout/               shell HTML completa — dos shells: customer-facing Layout (Topbar + IconRail + SubNav) y AdminLayout (topbar + breadcrumbs, sin multi-tool nav)
    base.templ          <!doctype>, <head> con scripts/CSS, <body>, JS global (sidebar level switching, tooltips, intl-tel-input init, htmx, overlay/menu/nav-loader runtime)
    types.go            Config, AdminConfig, User, Tool, SubItem, NavSection
    nav.templ           Topbar, IconRail (nivel 1), SubNav (nivel 2) — usados por Layout
    admin.templ         AdminLayout — topbar minimal con horizontal nav + user menu, breadcrumbs strip opcional
    icons.templ         (alternativa: ver view/icons/)
  avatar/               Circular profile picture con fallback de iniciales — usado por AdminLayout user menu
  breadcrumbs/          Breadcrumbs(items) — nav trail con separador "/", último item no-clickeable
  icons/                set compartido de iconos SVG (currentColor stroke)
  input/                text, password, number, phone (intl-tel-input wrapper), select, checkbox, checkbox_card, toggle, radio_card, textarea, hidden, date, month, file
  button/               Button(Options) — variant via Options.Variant; renders <button> or <a> (when Href set); Group + RenderGroup for footer/bar rows
  card/                 Card (chrome wrapper) + Section (sub-section divider)
  chip/                 Atomic chip pill — label, optional remove button, variants
  file_chip_input/      File input paired with chips (DataTransfer add/remove, invalid-flag)
  badge/                Variant (generic) + Status (shared code lookup) + VariantForCode helper
  flash/                Banner (generic) + Success / Error / Warning / Info wrappers
  table/                Table (chrome) + Row (helper) + BulkActionBar (Tailwind group-has visibility) + Pagination + EmptyState
  drawer/               SidePanel (slide-in) + Modal (centered) + Confirm (preset). FooterActions reuse button.Group; open/close via window.kibanOpenOverlay / kibanCloseOverlay; Escape closes topmost visible.
  menu/                 Kebab-style action dropdown — icon trigger + panel of <button role="menuitem"> rows. Single-active behaviour (open one closes others); JS in base.templ.
  spinner/              loading indicator (CSS class `.ds-spinner` in base.templ)
  tooltip/              CSS tooltip (`data-tooltip="…"` in base.templ)
  tabs/                 in-page tabs (underline style; distinct from SubNav's pill style for level-2 app navigation)
  comment_input/        textarea + chip-style file uploader + submit, all in one composition (used by klin's delivery comments; designed to be reused by future comment flows)

binding/                form_binding.FieldErrors() — traduce validator errors → map[formField]mensaje en español
middleware/
  authcookie/           cookie auth (kiban_session + kiban_space_id), redirige a /login en caso de fallo, HX-Redirect para HTMX
htmx/                   helpers Go: IsRequest, Redirect, TriggerName
```

> Subpaquetes en **singular** cuando representan un tipo (`input`, `button`, `card`, `drawer`, `flash`, `spinner`, `tooltip`).
> Subpaquetes en **plural** cuando agrupan varios items (`icons`, `tabs`).

## Catálogo de componentes

Estado: `[done]` = implementado, `[planned]` = pendiente. Cuando implementes uno, actualiza esta tabla.

### Layout & navegación (`view/layout/`)

| Componente | Estado | Notas |
|---|---|---|
| `view.Render(c, status, component)` | done | En `view/render.go`. Consumido por crm + rekon. |
| `layout.Config` | done | Struct con Title, ProjectName, SectionLabel, User, SpaceID, ShellURL, LogoutAction, ToolKey, ActiveKey, Tools, Docs, SubItems. Cada proyecto crea un type alias local (`view_layout.PageData = layout.Config`) y un `enrich(data)` que añade los campos fijos (Tools, etc.) antes de pasar a `Layout`. |
| `layout.Layout(cfg)` | done | Shell completa: HTML, head, scripts, topbar, sub-nav (siempre visible) + tools rail (toggleable), main content. |
| `layout.Topbar(cfg)` | done | Hamburger izquierda + logo kiban + space chip + Developers button + user menu con logout. |
| `layout.IconRail(cfg)` | done | Tools rail con `icon + label` por entrada (kiban tools + docs en la parte inferior). w-56. Hidden por defecto, slide-in cuando el hamburger está activo. Tool activo resaltado. |
| `layout.SubNav(cfg)` | done | Sub-nav siempre visible (no se oculta nunca). Items de la sección activa. Item activo resaltado. |
| `layout.AdminConfig` | done | Struct para AdminLayout: `Title`, `ProjectName`, `HomeHref`, `User`, `LogoutAction`, `NavSections`, `ActiveSection`, `Breadcrumbs`. Distinto de `Config` — no carga el modelo multi-tool (Tools/Docs/SubItems). |
| `layout.NavSection` | done | Una entrada de la nav horizontal del topbar admin: `Key`, `Label`, `Href`. `Key` matchea `AdminConfig.ActiveSection` para resaltar. |
| `layout.AdminLayout(cfg)` | done | Shell minimal para apps admin/internal: topbar (logo + horizontal nav + user menu) + breadcrumbs strip opcional + main. Reusa `Base()` para el chrome HTML; los blocks JS de Base que dependen de elementos del shell customer-facing (sidebar, sub-nav) son no-ops cuando esos elementos no existen. Para apps con 2-5 secciones top-level y autenticación propia, sin multi-tool switcher. |
| Admin user menu | done | Dentro de `AdminLayout` topbar: avatar + nombre como trigger de un dropdown que muestra email + form de logout. Reusa los handlers JS de `view/menu` (kibanToggleMenu/kibanCloseMenu) sin código nuevo. ID fijo `admin-user-menu`. |
| Tools rail toggle (slide) | done | CSS-only en `base.templ` con `.sidebar-rail-slot` (width 0 ↔ 14rem). El hamburger toggla el atributo booleano `data-sidebar-rail-open` en el root; JS persiste en `localStorage[<ProjectName>-sidebar-rail-open]` (key namespaceada por proyecto). |
| Navigation loader (`#nav-loader`) | done | Overlay full-screen mostrado al hacer click en cualquier `<a data-nav-loader>` (los items del tool rail lo llevan automáticamente). Cubre el wait del navegador mientras la siguiente página carga — útil cuando se cambia entre tools de backends distintos. Skip para clicks con modificadores, target=_blank, anchors `#…`, y links HTMX. Auto-hide via `pageshow` para no quedarse pegado en restores de bfcache. |
| Tooltip CSS (`data-tooltip`) | done | Bubble dark-ink, instant on hover/focus, ignora pointer events. CSS en `base.templ`. |
| intl-tel-input init | done | JS en `base.templ` escanea `[data-phone-input]` al DOMContentLoaded y tras cada htmx swap; sincroniza hiddens (`data-tel-cc`, `data-tel-national`) con el widget. |
| Spinner CSS (`.ds-spinner`) | done | Spinner kiban-primary, 32×32, en `base.templ`. |

### Iconos (`view/icons/`)

| Componente | Estado |
|---|---|
| `icons.Home`, `icons.Users`, `icons.ShoppingBag`, `icons.GitMerge`, `icons.CreditCard`, `icons.Code`, `icons.Database`, `icons.BookOpen`, `icons.FileText` | done |
| `icons.ChevronDown`, `icons.Filter`, `icons.Sort`, `icons.More`, `icons.Close`, `icons.Hamburger`, `icons.Sidebar` | done |

Convención: 20×20 viewBox=24, `stroke="currentColor"`, sin fill — colorables con clases `text-…` de Tailwind. Nombres sin prefijo `Icon` (el package ya lo aporta — `icons.Home` no `icons.IconHome`).

### Avatar (`view/avatar/`)

Foto de perfil circular con fallback de iniciales para cuando la URL es vacía, 404 o lenta. Hoy lo usa `AdminLayout` para el user menu del topbar.

**Patrón visual sin JS**: las iniciales siempre se renderizan en el background del span; cuando `Options.Src` está seteado, el `<img>` se posiciona absoluto encima con `object-cover` full-bleed. Si la imagen 404 o tarda, las iniciales se ven naturalmente — no hay handler `onerror`, no hay fallback dinámico.

| Componente | Estado | Notas |
|---|---|---|
| `avatar.Options` | done | `Src string` (URL imagen, vacío = solo iniciales), `Name string` (deriva iniciales + aria-label default), `Size string` (`sm` 24px / `md` 32px default / `lg` 40px), `AltText string` (override opcional del aria-label, default = Name). |
| `avatar.Avatar(o Options)` | done | Renderiza `<span>` circular con `bg-kiban-primary-soft text-kiban-primary`, iniciales centradas, `<img>` opcional encima. `overflow-hidden` corta la imagen al círculo. `shrink-0` protege el tamaño dentro de flex layouts (típico en topbars/listas). |
| `avatar.Initials(name)` | done | Helper público — devuelve hasta 2 letras mayúsculas: `"Antonio Blancas"` → `"AB"`, `"Antonio"` → `"A"`, `""` → `"?"`. Útil cuando un consumer quiere las iniciales fuera del templ (ej. server-rendered notification list que muestra el icono inicialmente). |

### Breadcrumbs (`view/breadcrumbs/`)

Trail de navegación: lista de `{Label, Href?}` con separador "/" entre items. Items con `Href` no-vacío son links; items con `Href = ""` se renderizan como texto plano (típicamente el último, página actual).

| Componente | Estado | Notas |
|---|---|---|
| `breadcrumbs.Item` | done | `Label string`, `Href string`. Convención: el último Item tiene `Href = ""` para que se vea como página actual (`text-kiban-ink font-medium aria-current="page"`); los anteriores con `Href` actúan como links de retroceso. Items intermedios sin `Href` se permiten — útil para segmentos no-navegables. |
| `breadcrumbs.Breadcrumbs(items)` | done | `<nav aria-label="Breadcrumb">` con `<ol>` flex-wrap. Lista vacía/`nil` no emite nada (útil en páginas sin jerarquía). Texto base muted (`kiban-ink3`); active en `kiban-ink + font-medium`; separador "/" en `kiban-ink4 select-none`. |

### Inputs (`view/input/`)

Convenciones del paquete (ver godoc de `types.go` para el detalle):
- Cada control emite `id="f-<name>"` y un `<label for="f-<name>">` correspondiente. Pasá `label=""` para omitir el `<label>` (útil cuando una sección o card más prominente ya cumple ese rol — el componente no emite el tag vacío).
- `errMsg != ""` → bordes rojos + mensaje en `text-red-600` debajo del control. `errMsg == "" && hint != ""` → mensaje en `text-kiban-ink3`. Nunca aparecen los dos a la vez.
- `required` es marcador visual (asterisco rojo). HTML5 `required` se omite a propósito — la validación vive server-side.
- HTMX se inyecta vía `attrs templ.Attributes` en componentes que lo soportan (Checkbox, CheckboxCard, RadioCard, Toggle); el design system se queda HTMX-agnóstico.

| Componente | Estado | Notas |
|---|---|---|
| `input.Text(name, label, value, errMsg, hint string, required bool)` | done | Texto plano. Mismo contrato error/hint descrito arriba. |
| `input.Password(name, label, value, errMsg, hint string, required bool, placeholder string)` | done | Como Text pero `type="password"` + `autocomplete="off"`. `placeholder` para renderizar el sentinel "••••••••" cuando ya hay un secreto guardado y el usuario puede dejar el campo en blanco para conservarlo. |
| `input.Number(name, label, value, errMsg, hint string, required bool, min, max, step string)` | done | `type="number"`. `min`/`max`/`step` como strings para soportar decimales ("0.01") y permitir suprimir el atributo pasando `""`. `value` también string para round-trippar lo que el usuario tipeó tras un error de validación. |
| `input.Phone(label, ccName, ccValue, phoneName, phoneValue, errMsg string, required bool)` | done | Wrapper de intl-tel-input. Renderiza `<input type="tel" data-tel-visible>` (sin `name`) + dos `<input type="hidden">` (`data-tel-cc` / `data-tel-national`). El init JS de `layout/base.templ` escanea `[data-phone-input]`. `data-tel-initial` se pre-popula con `+<cc><national>` cuando ambos lados están seteados (edit). |
| `input.Select(name, label, value, errMsg, hint string, options []SelectOption, required bool)` | done | `options []SelectOption{Value, Label}`. Si necesitas placeholder, prepéndelo como una `SelectOption{Value: "", Label: "Selecciona…"}`. |
| `input.Date(name, label, value, errMsg, hint string, required bool)` | done | `<input type="date">`, mismo contrato que Text. |
| `input.Month(name, label, value, errMsg, hint string, required bool)` | done | `<input type="month">`. El browser muestra un picker mes/año (sin día); el value que sube en el form es `"YYYY-MM"`. Soporte ubicuo en Chrome/Edge/Safari/Firefox; browsers viejos caen a un text input libre, así que el server siempre debe validar la shape `YYYY-MM`. |
| `input.File(name, label, accept, hint string, required, multiple bool)` | done | `<input type="file">` con `file:` prefix kiban-styleado. `accept` (ej. `.csv,text/csv`) para restringir el picker; pasá `""` para aceptar todo. `required?={required}` se emite (HTML5) para que el browser bloquee submit sin archivo. `multiple?={multiple}` permite que el cuadro de diálogo del SO seleccione varios archivos en una sola apertura — pareá con `view/file_chip_input` cuando además querés chips + remove + agregar más en sucesivos picks. |
| `input.Hidden(name, value, id string)` | done | `<input type="hidden">`. Pasá `id=""` cuando no haga falta (la mayoría de los casos); poné `id` cuando el JS necesite leer/escribir el valor en runtime (ej. el script de geolocation que rellena lat/lng antes de submit). |
| `input.Textarea(name, label, value, errMsg, hint string, required bool, rows int, placeholder string, mono bool)` | done | `rows <= 0` → default 3. `placeholder` se emite siempre (vacío = no se muestra). `mono=true` swappea a `font-mono` para PEM blocks, listas de emails una-por-línea, etc. |
| `input.Checkbox(name, label string, value, enabled bool, attrs templ.Attributes, disabledHint string)` | done | Checkbox plano + label en línea. **`value="true"` obligatorio** — Gin no parsea "on" como bool. `disabledHint` se renderiza después del label en `text-kiban-ink4` solo cuando `!enabled` (útil para mostrar "— no configurado en …" en toggles disabled). Pasá `""` cuando no haga falta. |
| `input.CheckboxCard(name, title, help string, value, enabled bool, attrs templ.Attributes)` | done | Card-style: borde + padding + título arriba + helper text en muted abajo. Mismo contrato `value="true"`. `enabled=false` deshabilita y baja opacidad. Útil para listas de preferencias / settings cards. |
| `input.Toggle(name, label string, value, enabled bool, attrs templ.Attributes)` | done | Checkbox semántico, switch visual (Tailwind `peer-checked` para animar el tracker). Mismo contrato `value="true"` y `attrs` que Checkbox. |
| `input.RadioCard(name, value, title, subtitle string, checked bool, attrs templ.Attributes)` | done | Card-style radio: borde clickeable, resalta `border-kiban-primary bg-kiban-primary-soft` al checkearse. Múltiples cards con el mismo `name` forman el group. `subtitle` opcional. `attrs` para HTMX cuando la selección dispara un partial re-render del form. |

### Botones (`view/button/`)

Primitiva única para todo control "botón o link" del kit: variantes `primary` / `secondary` / `danger` / `icon`, soporte de iconos (string registry o `templ.Component` arbitrario), modo anchor cuando `Options.Href` está seteado, y `Group` para componer rows de acciones (drawer footers, bulk bars, action sheets) con defaults de variante por slot.

| Componente | Estado | Notas |
|---|---|---|
| `button.Button(p Options)` | done | Único template. `Options.Variant`: `primary` (default), `secondary`, `danger`, `icon`. (`destructive` se normaliza a `danger` por compatibilidad.) Cuando `Options.Href != ""`, renderiza `<a>` y silenciosamente ignora `IsSubmit`/`IsReset`/`Form`/`Disabled` (no aplican a anchors). Atajo `Icon` + `IconPosition` para glifos del registro interno; para iconos nuevos u otros SVGs usar `IconComponent: icons.Algo()` (u otro `templ.Component`). Variante `icon`: usar `AriaLabel` (o `Title`). |
| `button.Group` (struct) | done | `PrimaryAction *Options` (rightmost, default variant `primary`) + `SecondaryActions []Options` (a la izquierda en orden, default variant `secondary`). Método `IsEmpty()` para que renderers skipeen el chrome cuando no hay nada que mostrar. Reusado por `drawer.SidePanelConfig.FooterActions`, `drawer.ModalConfig.FooterActions`, `drawer.ConfirmConfig.PrimaryAction`, `table.BulkActionBarConfig.Actions`. |
| `button.RenderGroup(g Group)` | done | Renderiza un Group como `flex gap-3 justify-end` (secondaries izquierda, primary derecha). Emite nada cuando `g.IsEmpty()`. Per-slot defaults aplicados automáticamente vía `WithDefaultVariant`. |
| `button.WithDefaultVariant(o, def)` | done | Helper público que devuelve `o` con `Variant` reemplazado por `def` cuando el original viene vacío. Útil cuando un componente del DS quiere consumir `Options` y aplicar sus propios defaults sin forzar al caller a repetir el variant. |

`Options`: `Label`, `Icon`, `IconPosition`, `IconComponent templ.Component` (opcional; si no es nil, sustituye al registro de `Icon` y permite cualquier icono de `view/icons` u otro fragmento templ), `Variant`, `Href` (cuando ≠ "" → renderiza `<a>` en vez de `<button>`), `IsSubmit`, `IsReset`, `OnClick` (raw inline JS — patrón típico `kibanCloseOverlay('id')`), `Disabled`, `ExtraClass`, `Form`, `AriaLabel`, `Title`, `Attrs templ.Attributes`. HTMX y cualquier atributo extra van en `Attrs` (mismo patrón que `view/input`). Helpers públicos: `button.BuildClass`, `button.NonEmptyAttrs`, `button.WithDefaultVariant`.

### Display (`view/card/`, `view/badge/`, `view/flash/`, `view/spinner/`, `view/tooltip/`)

Convenciones del grupo:
- **Variantes** comparten vocabulario en card / badge / flash: `success` (emerald), `warning` (amber), `danger` (red), `info` (kiban-primary tint), `neutral` (kiban-surface). En `flash` el destructivo se llama `error` por consistencia con el sentimiento de los banners; el resto de los componentes usan `danger`. Strings vacíos caen al default razonable (`neutral` para badge, `info` para flash, sin variant para card).
- **HTMX-agnóstico**: ninguno de estos componentes inyecta atributos hx-*. Si un caso de uso requiere HTMX (swap parcial de un Card después de una mutación, etc.), envuelve el componente en un `<div hx-*>` o un `<section id>` al lado del callsite.

#### `view/card/`

| Componente | Estado | Notas |
|---|---|---|
| `card.Card(variant string)` | done | Chrome blanca: `bg-white border border-kiban-border rounded-md p-6 space-y-5`. `variant` swappa el color del borde (`success`/`warning`/`danger`/`info`); el body se queda neutral — para mensajes de estado loud usá `flash.*`. Acepta cualquier contenido como templ children. Si necesitás `id` (HTMX target, anchor link, CSS selector), envolvé el `Card` en un `<section id="…">`. |
| `card.Section(title, subtitle string)` | done | Sub-sección DENTRO de un `Card`. Header (heading + subtitle opcional) + body via templ children. Cuando hay 2+ Sections como hijos directos del mismo Card, las que no son first-child reciben `pt-4 border-t border-kiban-border` automáticamente vía Tailwind `first:`. Para card de UNA sola sección no uses Section — `Card { contenido inline }` ya alcanza. Subtitle es plain text; para subtítulos con HTML rico (links, asteriscos coloreados) renderizalos inline en el body. |

#### `view/badge/`

| Componente | Estado | Notas |
|---|---|---|
| `badge.Variant(label, variant, size string)` | done | Pill genérico. `variant`: `success`/`warning`/`danger`/`info`/`neutral` (vacío = `neutral`). `size`: `sm` (`px-2 py-1 text-xs`, default para celdas de tabla) o `md` (`px-3 py-1 text-sm`, headers de detail). Siempre que el código de status sea project-only, llamá `Variant` directo y mapeá el código a un variant en el callsite. |
| `badge.Status(code, label, size string)` | done | Wrapper de conveniencia para los códigos compartidos por más de un proyecto kiban. Llama `VariantForCode(code)` y delega a `Variant`. Códigos cubiertos hoy: `PAID/PENDING/EXPIRED/CANCELLED` (lifecycle de pagos), `VALIDATED/TO_VALIDATE` (validación de payment-methods), `ACTIVE/COMPLETED/PROCESSING/FAILED/REJECTED/DRAFT` (varios). Códigos desconocidos caen a `neutral` (no rompen el render). |
| `badge.VariantForCode(code) string` | done | Lookup compartido `code → variant`. Públicamente accesible para que un proyecto pueda combinar el lookup con su propio override. **Sólo agregar códigos que se reutilicen en >1 proyecto** — los códigos project-only viven en mappings locales en cada proyecto. |

#### `view/flash/`

| Componente | Estado | Notas |
|---|---|---|
| `flash.Banner(variant, msg string)` | done | Banner full-width genérico. `variant`: `success`/`error`/`warning`/`info` (vacío = `info`). `msg` es plain text — si necesitás HTML (links, formato), renderizá inline el div con las clases del variant correspondiente. |
| `flash.Success(msg string)` | done | Wrapper tipado: confirmaciones post-mutación ("Cliente guardado", "Pago registrado", …). |
| `flash.Error(msg string)` | done | Wrapper tipado: operaciones fallidas y bloqueos ("No tienes permiso", "No pudimos guardar", …). |
| `flash.Warning(msg string)` | done | Wrapper tipado: warnings suaves donde el usuario puede continuar pero debería leer algo ("Banco Donde no está configurado todavía", …). |
| `flash.Info(msg string)` | done | Wrapper tipado: mensajes informativos neutrales — el catch-all cuando nada salió mal pero hay algo que comunicar. |

Banners son **estáticos**: no hay JS de dismiss, no hay localStorage. Re-renderizan con la siguiente request / HTMX swap. Esto coincide con cómo los controllers ya piensan en `FlashSuccess`/`FlashError` en sus views.

#### Spinner / Tooltip (CSS-only, no templ helpers)

| Componente | Estado | Notas |
|---|---|---|
| `spinner` (CSS class `.ds-spinner`) | done | Definida en `view/layout/base.templ`. Uso: `<div class="ds-spinner" role="status"></div>`. No hay templ helper — la clase CSS es la API. |
| `tooltip` | done | Patrón CSS-only: añade `data-tooltip="texto"` a cualquier elemento. CSS en `base.templ`. |

### Chip (`view/chip/`)

Atomic pill primitive — a small label-shaped tag with an optional close "×" button. Use it anywhere a list of removable labels makes sense: file uploaders (via `view/file_chip_input/`), filter pills, multi-select readouts, tag inputs.

**What it is NOT:**
- Not a [badge](#viewbadge): badges are status displays, chips are interactive/removable.
- Not a [button](#viewbutton): the chip body is non-interactive; only the optional remove button is clickable.
- Not a standalone widget: chips are always rendered as items inside a parent that owns the list state. The chip itself emits markup only; the consumer wires the remove click to whatever JS / HTMX contract removes the chip.

The chip's remove button (when `Removable=true`) emits `<button type="button" data-chip-remove>` so a parent's delegated click listener can find it without per-chip wiring. Consumers can attach extra attributes via `Options.RemoveAttrs` (e.g., `data-remove="<key>"` to identify which chip).

| Componente | Estado | Notas |
|---|---|---|
| `chip.Options` (struct) | done | `Label string` (required), `Subtext string` (optional muted secondary line, e.g. file size), `Title string` (native tooltip), `Variant string` (`""`/`"default"`, `"danger"`, `"info"`, `"success"`, `"warning"`; unknown values fall back to default), `Removable bool` (renders the × button), `RemoveAttrs templ.Attributes` (spread onto the remove `<button>`), `RemoveAriaLabel string` (override default Spanish "Quitar"), `Attrs templ.Attributes` (spread onto the outer `<span>`). |
| `chip.Chip(opts)` | done | Renders an `<span>` pill. Variant tokens come from `helpers.go` so adding a new variant is one map entry, not a templ rewrite. Used today by `view/file_chip_input` for the file selection list; the markup is mirrored 1:1 by JS in `view/layout/base.templ` (when chips are added client-side, the JS injects equivalent markup — keep them in lockstep when changing styling). |

### File chip input (`view/file_chip_input/`)

A `<input type="file">` paired with a chip-style readout of the selected files (one `chip.Chip` per file). Selecting more files appends instead of replacing (DataTransfer trick), each chip has a × that drops just that file, and oversize files render as `danger` chips with a tooltip + the consumer's submit button auto-disabled until they're removed.

**Use when** you need a multi-file (or single-file) picker with per-file remove + cumulative add behavior. **Don't use** for inline-image fields, single-shot upload-and-go flows, or anywhere the chip readout would be visual noise (a single `[input.File]` is simpler).

**Architecture:**
- The component emits markup only: an `[input.File]` plus a sibling `<ul data-chip-list>`. The chip list starts empty server-side; JS populates it from the user's pick.
- The DataTransfer / dedup / per-chip-remove / disable-on-invalid JS lives in `view/layout/base.templ`, scoped to `[data-kiban-file-chip-input]`. Re-runs on `DOMContentLoaded` and `htmx:afterSwap`, mirroring the rest of the design system's init pattern.
- Each instance is identified by a unique `[data-kiban-file-chip-input]` value (defaulting to `Options.ID`). Multiple instances on the same page work without collisions.

**Caller integration:**
- Wrap the field in a `<form>` (or HTMX equivalent). The component does not render the form; it's a sub-widget.
- For "disable submit while invalid file present", the consumer's submit button must carry `data-kiban-file-chip-submit="<id>"` (where `<id>` matches `Options.ID`). Without that attribute the JS still renders chips correctly; only the disable-while-invalid auto-wire is skipped.
- File inputs can't be pre-populated server-side (browser security), so there's no "round-trip on validation error" affordance — the user re-picks on form re-render.

| Componente | Estado | Notas |
|---|---|---|
| `file_chip_input.Options` (struct) | done | `Name string` (required, the form field), `Label string` (above the picker; empty omits), `ID string` (wrapper id + JS lookup key; defaults to `kiban-file-chip-input`), `Hint string` (under the picker; empty auto-builds from MaxSize + Multiple), `Accept string` (native `accept` filter), `Multiple bool`, `Required bool`, `MaxSizeBytes int64` (per-file cap; over-cap chips render `danger` and disable submit), `FileVariant string` (override the chip variant for valid files; default neutral). |
| `file_chip_input.Field(opts)` | done | Renders the wrapper + `[input.File]` + empty chip list. Composed by `view/comment_input` internally — comment_input.Field is the higher-level "comment with attachments" widget, file_chip_input.Field is the picker primitive used anywhere else. |

### Tables (`view/table/`)

| Componente | Estado | Notas |
|---|---|---|
| `table.PaginationConfig` | done | Struct con Page, HasPrev, HasNext, PageURL func(int) string, Target, Indicator, `NextVariant` (`""`/`"secondary"` = ambos outlined, default project-wide; `"primary"` = legacy solid kiban-primary "Siguiente", opt-in). Caller construye con un closure sobre su `pageURL` local para preservar filtros entre páginas. |
| `table.Pagination(cfg)` | done | Anterior/Siguiente botones HTMX. Estados disabled-styled cuando edge. `cfg.PageURL` callback para que el caller maneje filter state. Default visual: ambos botones outlined secondary; `cfg.NextVariant = "primary"` para opt-in al look solid azul. |
| `table.EmptyState(title, hint)` | done | Card centrada con border-top, "Aún no tenemos X qué mostrar". `hint` opcional. |
| `table.Table(cfg TableConfig)` | done | Chrome estándar (border-t + table). `cfg.Headers []string` (plain text por ahora; sortable headers con chevron-down se agregan después). `cfg.HeaderAlignRight []bool` opcional — slice paralelo a `Headers`; cuando el índice está marcado en `true` el `<th>` correspondiente usa `text-right` en vez del default `text-left`. Pensado para columnas de acciones cuya celda body ya está alineada a la derecha (kebab menus, button rows). `nil` o slice corto = todo `text-left` (los callers existentes no se rompen). `cfg.HeaderNoWrap bool` opcional — cuando es `true` cada `<th>` lleva `whitespace-nowrap`, evitando que títulos multi-palabra como "Fecha de creación" wrappen en columnas estrechas. Útil cuando las celdas body son `whitespace-nowrap` y la asimetría visual con headers wrappeando se vuelve obvia. `cfg.BulkSelect=true` prepende un `<th>` con checkbox "select all" que togglea todos los `input[name=ids]` siblings via inline JS. Body via templ children: el caller renderiza sus `<tr>` o, mejor, usa `@table.Row(href, bulkValue)` para el chrome estándar. |
| `table.Row(href, bulkValue string)` | done | `<tr>` con `border-b + hover:bg-kiban-primary-soft`. Cuando `href != ""` agrega `data-href` + `cursor-pointer` (la JS de `view/layout/base.templ` intercepta clicks y navega; los clicks en `a/button/input/textarea/select/label` no disparan navigation, así no compite con anchors anidados). Cuando `bulkValue != ""` prepende un `<td><input type="checkbox" name="ids" value={bulkValue}>` — combinar con `Table.BulkSelect=true` para el toggle de header. |
| `table.BulkActionBar(cfg BulkActionBarConfig)` | done | Barra de acciones que aparece sobre la tabla cuando hay selecciones. **Visibilidad puramente Tailwind** (sin JS, sin `<style>` por instancia): la barra usa `hidden group-has-[input[name=ids]:checked]/bulk:flex`. **Contrato del caller**: envolver el form que contiene tabla + barra en `<form class="group/bulk">` para que el variant nombrado resuelva. `cfg.Message` muestra texto muted a la izquierda; `cfg.Actions button.Group` renderiza primary + secondaries a la derecha (mismo `button.Group` que drawer). |

### Drawer / overlay (`view/drawer/`)

Convenciones del paquete:
- **Open / close** vía `window.kibanOpenOverlay(id)` / `kibanCloseOverlay(id)` (definidos en `view/layout/base.templ`). El caller cablea su trigger con `onclick="kibanOpenOverlay('the-id')"`; el componente se ocupa del close (close button + backdrop click + Escape global).
- **Escape** cierra solo el overlay topmost visible (por orden de `[data-kiban-overlay]:not(.hidden)` en el DOM) — multiple overlays apilados se cierran uno a la vez, igual que el comportamiento nativo del browser.
- **Sizes** compartidos: `sm` (max-w-sm) / `md` (max-w-md, default) / `lg` (max-w-lg) / `xl` (max-w-xl). String vacío cae a `md` (excepto Confirm, que cae a `sm`).
- **Z-index**: SidePanel = 40, Modal/Confirm = 50.
- **Footer actions** vía `button.Group{PrimaryAction *Options, SecondaryActions []Options}` (definido en `view/button/`, ver sección de arriba). PrimaryAction renderiza a la derecha (default variant `primary`); SecondaryActions a la izquierda en orden, default variant `secondary`. Si no hay acciones (`Group.IsEmpty()`), no se renderiza el footer.
- **Variantes** del row alineadas con `button.Button`: `primary` (kiban-primary), `secondary` (outline border-kiban-border), `danger` (red-600 — usar para confirms destructivos).

#### Componentes

| Componente | Estado | Notas |
|---|---|---|
| `drawer.SidePanel(cfg SidePanelConfig)` | done | Slide-in desde la derecha. `cfg.ID` + `Title` + `Size` + `FooterActions button.Group`. Body via templ children, padding `px-6 py-4` aplicado al body (caller no se preocupa por el gutter). Patrón típico para filter-drawers: el caller renderiza un `<form id="filter-form">` en el body y la PrimaryAction usa `IsSubmit:true` + `Form:"filter-form"` + `OnClick:"kibanCloseOverlay('id')"` para submit-and-dismiss. |
| `drawer.Modal(cfg ModalConfig)` | done | Centrado, backdrop oscuro (`bg-black/40`). Mismos campos que SidePanel (incluido `FooterActions button.Group`) + `Icon templ.Component` opcional renderizado a la izquierda del título (caller pasa el block fully styled — patrón típico kiban: `<div class="w-8 h-8 rounded-full bg-kiban-primary-soft text-kiban-primary flex items-center justify-center"><svg…/></div>`). Para flujos HTMX-form-bound (modal con submit): wrapeá el `@drawer.Modal(...)` entero en un `<form hx-post=… hx-on::after-request="if(event.detail.successful){ kibanCloseOverlay('id'); }">`. |
| `drawer.Confirm(cfg ConfirmConfig)` | done | Preset de Modal con shape fijo: title (opcional) + message + cancel/confirm buttons. Default size `sm`. `cfg.PrimaryAction button.Options` es el botón de confirm; setear `Variant:"danger"` para deletes. Cancel se cablea automáticamente al `kibanCloseOverlay(id)`; label default "Cancelar". Para "are you sure?"-level simple usar `hx-confirm` nativo de HTMX; usar Confirm cuando se necesita styling kiban + HTMX wiring custom + título largo. |

### Menu (`view/menu/`)

Kebab-style action menu para tablas y filas con varias acciones. El trigger es un botón con icono (reusa `button.Button` con `Variant:"icon"` + `IconComponent: icons.More()`); el panel es un dropdown `absolute right-0 top-full` con items renderizados como `<button role="menuitem">`. JS toggle/close + outside-click + Escape viven en `view/layout/base.templ` (`window.kibanToggleMenu`, `window.kibanCloseMenu`).

**Single-active**: abrir un menú cierra cualquier otro menú abierto en la página (matching native OS menu behaviour). Click fuera de cualquier `[data-kiban-menu]` cierra todos los abiertos. Escape cierra los menús primero, antes que cualquier overlay (para que cerrar el menú no dismisse el modal subyacente).

**Cierre automático tras click**: cada item agrega `window.kibanCloseMenu(id)` después del `OnClick` del consumer, así no hay que recordar cerrarlo manualmente.

**ID por instancia**: `cfg.ID` debe ser único en la página (ej. `apikey-menu-<rowID>`). Las funciones JS keyean por ese ID para encontrar trigger + panel.

| Componente | Estado | Notas |
|---|---|---|
| `menu.Config` (struct) | done | `ID string` (required, único en la página), `AriaLabel string` (aria-label del trigger, ej. "Acciones para {nombre}"), `Items []MenuItem`. Si `Items` está vacío, `Menu` no renderiza nada (consumer puede pasar lista vacía sin chequear). |
| `menu.MenuItem` (struct) | done | `Label string`, `OnClick string` (raw inline JS — el menú agrega `kibanCloseMenu` al final), `Variant string` (`""`/`"default"` → `text-kiban-ink`; `"danger"` → `text-red-600` para destructivos), `Icon templ.Component` opcional a la izquierda, `Attrs templ.Attributes` para HTMX/data-* (mismo patrón que `button.Options.Attrs`). |
| `menu.Menu(cfg Config)` | done | Renderiza `<div class="relative inline-block" data-kiban-menu={id}>` con el trigger button + panel oculto. Panel se posiciona `absolute right-0 top-full mt-1 z-30`; asegurate de que el call-site tenga espacio a la derecha o abajo. Si `cfg.Items` está vacío, no emite nada. |

### Tabs (`view/tabs/`)

In-page tabs primitive — distinto de `layout.SubNav` (level-2 nav del shell que usa el estilo pill). Tabs usan underline-style para "switch view inside this page", pills mean "switch tools/sections in the app". Mantener visualmente diferenciados evita confundir state in-page con navegación global.

**Visual contract (polished SaaS look)**: el strip está envuelto en un `<div class="border-b border-kiban-border">` que dibuja un raíl gris de **1 px** a lo ancho del contenedor. El tab activo tiene un subrayado de **1 px** color `kiban-primary` (mismo grosor que el raíl) ajustado al ancho del label vía un `<span class="inline-block">` interno. Tabs inactivos preservan `border-transparent` 1 px para que la fila no se desplace verticalmente al cambiar el active. Gap entre tabs: `gap-3` (12 px). Active label en `kiban-primary` + `font-medium`; inactive en `kiban-ink3` con `hover:text-kiban-ink`. La regla de "mismo grosor de borde entre raíl y subrayado" está blindada por `TestStrip_ContainerBaselineMatchesActiveBorderWidth` — no bumpear uno sin el otro.

| Componente | Estado | Notas |
|---|---|---|
| `tabs.TabItem` (struct) | done | Una entrada del strip. Required: `Key string` (matched against activeKey), `Label string`, `Href string` (URL canónica para fallback de browser nav). Affordances opcionales: `Icon templ.Component` (slot a la izquierda del label, mismo patrón que `button.Options.IconComponent`), `Count int` + `HasCount bool` (badge pill a la derecha — el toggle explícito permite que un `0` real se renderice como "Inbox (0)" sin ambigüedad sentinel-vs-zero), `Disabled bool` (saca el `href`, mete `aria-disabled="true"` + `pointer-events-none` + `opacity-50`), `Title string` (atributo `title=""` nativo — útil para explicar el por qué cuando `Disabled=true`). Escape hatch HTMX: `Attrs templ.Attributes` (spread sobre el `<a>` para flujos `hx-get` / `hx-target` / `hx-push-url`). |
| `tabs.Strip(items []TabItem, activeKey string)` | done | Renderiza el strip con el visual contract descrito arriba. Active state es 100% caller-driven: si ningún `Key` matchea `activeKey`, no se resalta nada. Múltiples instancias en la misma página son seguras (cada `<a>` lleva su propio `href`/Attrs). |

### Comment input (`view/comment_input/`)

Composición de alto nivel: textarea + `[file_chip_input.Field]` + botón submit, todo bajo un único `<form>`. Encapsula el look-and-feel para que cualquier proyecto que necesite un "post a comment" UI lo consuma con un solo callsite. Hoy lo usa klin para los comentarios de entrega; fue diseñado para que rekon (notas en órdenes de pago) / crm (actividad por cliente) / futuros tools tengan la misma UX sin re-implementar.

**Reglas que respeta**:
- HTMX-agnóstico: el form es POST plano al `Options.Action`. Los consumers que quieran HTMX inyectan los atributos vía `Options.FormAttrs` (mismo patrón que `button.Options.Attrs`, `input.Toggle.attrs`, etc.).
- File handling delega 100% en `view/file_chip_input` — comment_input ya no carga JS propio, sólo lo compone. Esto significa que cualquier mejora en el chip-uploader (drag&drop, previews, etc.) se hace en `file_chip_input` y comment_input la hereda gratis. El submit del comment_input lleva `data-kiban-file-chip-submit="<filesID>"` para que la JS de file_chip_input lo deshabilite mientras haya un archivo inválido.
- Reusa `card.Card`, `flash.Success`/`Error`, `input.Textarea`, `file_chip_input.Field`, `button.Button` — cero primitivas nuevas a este nivel.

| Componente | Estado | Notas |
|---|---|---|
| `comment_input.Options` (struct) | done | Configura el `Field`. Sólo `Action` es required. Defaults razonables para todo lo demás (Title="Nuevo comentario", SubmitLabel="Enviar", Placeholder="Escribe un comentario…", textName="text", filesName="files", id="comment-input"). Toggles importantes: `Multiple bool` para multi-archivo, `MaxSizeBytes int64` para el cap por archivo client-side, `Accept string` para `<input accept>`, `DisableFiles bool` para text-only, `WithoutCard bool` para skipear el `card.Card` chrome (consumer renderiza su propio container), `FormAttrs templ.Attributes` para escape hatch HTMX. Round-trip de error: `TextValue` + `TextError` para draft + per-field error; `GlobalError` para flash banner; `Success` para flash post-submit. `MaxChars int` es informativo (renderiza "Máximo N caracteres" como hint del textarea); el server sigue siendo el de la verdad. |
| `comment_input.Field(opts Options)` | done | El templ. Emite (en orden): título + subtítulo opcionales, banners de Success/GlobalError si aplica, el `<form>` con `enctype="multipart/form-data"` apuntando a `Action`, textarea, `[file_chip_input.Field]` con id `<opts.ID>-files`, botón submit con `data-kiban-file-chip-submit="<id>-files"` para que la JS de `file_chip_input` lo deshabilite mientras hay archivos inválidos. Múltiples instancias en la misma página funcionan: cada `Options.ID` produce un par único `<id>-form` / `<id>-files`. |

### Form binding (`binding/`)

**Convención cross-project**: este paquete es la única fuente de traducciones de errores de validator across todos los proyectos kiban. Cuando un proyecto registra un validator tag custom (`regexRFC`, `regexCURP`, `regexCLABE`, etc.), también debe registrar acá el mensaje en español via `RegisterMessage` — así el wording se mantiene consistente entre rekon, crm, y futuros proyectos. **No copy-pastear el switch de mensajes en cada proyecto**: el dispatcher vive acá y se extiende via la función pública `RegisterMessage`.

| Componente | Estado | Notas |
|---|---|---|
| `binding.FieldErrors(err) map[string]string` | done | Mapea `validator.ValidationErrors` → claves = nombre del tag `form:"…"`, valores = mensaje en español. Internamente delega a `MessageFor` (incluyendo el fallback a "Valor inválido" para tags desconocidos). |
| `init()` en el package | done | Registra `TagNameFunc` en el validator de Gin para que use `form:"…"` como key. Corre automáticamente al importar el package. |
| `binding.MessageFor(tag, param string) string` | done | Lookup público. Orden: (1) registrados via `RegisterMessage`, (2) tags built-in (`required`, `email`, `min`, `max`, `len`, `url`, `oneof`, `gt`, `gte`, `eqfield`), (3) "Valor inválido" como fallback. Útil cuando el handler necesita el mensaje de un tag fuera del flujo de `FieldErrors` (ej. building API JSON responses). |
| `binding.RegisterMessage(tag string, formatter func(param string) string)` | done | Registra el mensaje en español para un tag custom. El formatter recibe el `param` del tag (lo que va después del `=` en `binding:"tagName=param"`); ignorar cuando el tag no toma parámetro. **Override**: las registraciones custom ganan sobre los built-ins, así un proyecto puede tweakear el wording de un tag estándar (`required`, etc.) si producto lo pide. **Panics** si `tag == ""` o `formatter == nil` — ambos son bugs de programación que serían silently broken sin el panic. Patrón típico: llamar desde un `init()` o setup en el `cmd/api/main.go` del consumidor, antes de servir requests. |

### HTMX helpers (`htmx/`)

Helpers de controller para trabajar con headers HTMX. Centralizan los nombres de header y la rama HTMX-vs-browser para que ningún handler open-codee `c.GetHeader("HX-Request") == "true"` ni se olvide del branching cuando hace un redirect.

Import:
```go
ds_htmx "github.com/kiban-cloud/go-kiban-design-system/htmx"
```

| Componente | Estado | Notas |
|---|---|---|
| `htmx.IsRequest(c *gin.Context) bool` | done | `c.GetHeader("HX-Request") == "true"`. Usado en handlers para decidir entre `Page(layout, view)` (full page, browser nav) y `Content(view)` (partial, HTMX swap). |
| `htmx.Redirect(c *gin.Context, url string)` | done | Redirect que funciona para ambos casos. HTMX request → `HX-Redirect: url` header + HTTP 200 (HTMX hace la navegación client-side). Browser request → `c.Redirect(302, url)`. **Importante**: llamar `c.Redirect(302, url)` directo en una request HTMX falla silenciosamente — HTMX swap-ea el body del redirect dentro del target div en lugar de navegar. Este helper hace que the right thing sea el default. Se usa 302 (`StatusFound`) en el path de browser, matcheando el patrón dominante en rekon + crm; 303 (StatusSeeOther) sería más correcto semánticamente para POST-then-GET pero la diferencia práctica es nula. |
| `htmx.TriggerName(c *gin.Context) string` | done | Atajo a `c.GetHeader("HX-Trigger-Name")` — el `name` attribute del element que disparó la request HTMX. Útil cuando un form tiene multiple submit buttons y el handler necesita saber cuál se clickeó. Empty string cuando el header no está. |

### Middleware (`middleware/authcookie/`)

| Componente | Estado | Notas |
|---|---|---|
| `authcookie.New(uc, loginURL) *Middleware` | done | Constructor con dependencias explícitas. Cada proyecto lo wirea según su DI: rekon llama directo en el container; crm lo envuelve en una struct con tag `inject:""` que delega tras `SetupAfterInjection`. |
| `(*Middleware).Middleware() gin.HandlerFunc` | done | El gin handler. Lee cookies `kiban_session` + `kiban_space_id`, valida vía `IAuthorizationAuthorizeWithSessionUseCase` de go-kiban, setea `controller_core_model.CONTEXT_KEY_AUTHORIZATION_OBJECT`. En falla emite `HX-Redirect` (HTMX) o `c.Redirect` (browser). |
| `authcookie.GetAuthorization(c)` | done | Re-exporta `controller_core_middleware.GetAuthorization` para que los htmx controllers no tengan que importar go-kiban directamente. |
| `authcookie.CookieSession`, `CookieSpaceID` | done | Constantes con los nombres de las cookies. |

## Convenciones

### Templ

- Funcs templ exportadas (`PascalCase`) cuando los consumidores las usan. Lowercase para helpers internos.
- Cada componente vive en su subpackage. Un subpackage puede tener múltiples templ funcs relacionadas (ej. `input/text.templ`, `input/phone.templ`, todas en package `input`).
- Helpers Go (mapeo de datos, structs, funciones puras) van en `.go` adyacentes al `.templ`, mismo package.
- `_templ.go` se commitea (regenera con `go tool templ generate`).

### Tailwind & tokens

- Tokens del design system kiban (colores `kiban.primary`, `kiban.ink/ink2/ink3/ink4`, `kiban.surface`, `kiban.border`, fonts Inter/Manrope) configurados inline en `view/layout/base.templ` vía `tailwind.config = {…}`.
- No editar los tokens a menos que cambien en el design system fuente.
- Mobile-first, breakpoints estándar de Tailwind (`md:`, `lg:`).

### HTMX patterns

- **Form binding**: handlers reciben request structs con tags `form:"…"`, llaman `c.ShouldBind(&req)`, mapean errores con `binding.FieldErrors`.
- **Partial swap**: vistas de página exponen `Page(layout, view)` (con shell completa) y `Content(view)` (sólo el bloque interno). Handlers eligen cuál renderizar según `htmx.IsRequest(c)`.
- **HX-Redirect**: tras mutaciones exitosas, `htmx.Redirect(c, url)` decide entre header HTMX vs `c.Redirect`.
- **Loading overlays**: cualquier elemento con `class="htmx-indicator"` se torna visible cuando un `[hx-indicator="#…"]` apunta a él. Reglas CSS en `view/layout/base.templ`.
- **Tooltips**: `data-tooltip="…"` en cualquier elemento. CSS en `base.templ`.
- **Sidebar level switching**: persistencia en `localStorage["<project>-sidebar-level"]` (la clave debe ser específica por proyecto para no colisionar entre rekon y crm en el mismo origen).

### Naming

- Subpaquetes singulares para tipos de elementos (`input`, `button`, `card`).
- Subpaquetes plurales para colecciones (`icons`, `tabs`).
- Nombres de funciones templ describen el componente sin redundancia con el package: `input.Text` (no `input.TextInput`).

### Externals

- Toda librería JS/CSS externa se carga vía CDN, **pinned a versión exacta** (no `@latest`). Lista actual:
  - HTMX: `https://unpkg.com/htmx.org@2.0.4`
  - Tailwind: `https://cdn.tailwindcss.com` (CDN sin pinning — accept upstream)
  - intl-tel-input: `https://cdn.jsdelivr.net/npm/intl-tel-input@25.3.1/build/...`
  - Inter font: `https://rsms.me/inter/inter.css`
  - Manrope font: Google Fonts `https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;700`
- Al añadir una nueva externa, agregarla al `<head>` de `view/layout/base.templ` y documentarla aquí.

## Cómo agregar un componente

1. Decide el subpackage correcto (`input`, `card`, `button`, etc.). Si ninguno encaja, crea uno (nombre singular si representa un tipo de elemento).
2. Crea el `.templ` en el subpackage. Templ exportada (PascalCase).
3. Si el componente necesita helpers Go (mapeo de datos, configuración), añádelos en un `.go` adyacente (mismo package).
4. **Actualiza el catálogo de componentes** en este CLAUDE.md (cambia `[planned]` → `[done]` o añade una fila nueva).
5. Si requiere una librería externa nueva, añade su `<link>`/`<script>` al `<head>` de `view/layout/base.templ` y documéntala en la sección "Externals".
6. Corre `go tool templ generate` desde el root del paquete.
7. Verifica: `go build ./...` y `go vet ./...` limpios.
8. Commit (preferentemente con un mensaje del estilo `feat(input): add Phone component`).

## Comandos

```bash
go tool templ generate                # regenera *_templ.go en todos los paquetes view/
go build ./...                        # compila todo
go vet ./...                          # vet estático
go mod tidy                           # limpia go.mod / go.sum
```

## Versionado y consumo

Los proyectos consumidores agregan el paquete a su `go.mod`:

```go
require github.com/kiban-cloud/go-kiban-design-system v0.x.y
```

**Durante desarrollo local**, usar un `replace` directive para evitar tagged releases por cada cambio:

```go
// go.mod del proyecto consumidor (rekon, crm)
replace github.com/kiban-cloud/go-kiban-design-system => ../go-kiban-design-system
```

Cuando se publica una versión nueva en GitHub:

```bash
# en el repo del consumidor
go get -u github.com/kiban-cloud/go-kiban-design-system@latest
go mod tidy
```

Versionado: semver. Cambios breaking incrementan major. Antes de v1.0.0 puede haber cambios breaking en minor.
