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
  layout/               shell HTML completa
    base.templ          <!doctype>, <head> con scripts/CSS, <body>, JS global (sidebar level switching, tooltips, intl-tel-input init, htmx)
    types.go            PageData, NavConfig, User, Tool, SubItem
    nav.templ           Topbar, IconRail (nivel 1), SubNav (nivel 2)
    icons.templ         (alternativa: ver view/icons/)
  icons/                set compartido de iconos SVG (currentColor stroke)
  input/                text, password, number, phone (intl-tel-input wrapper), select, checkbox, checkbox_card, toggle, radio_card, textarea, hidden, date, file
  button/               primary, secondary, destructive, icon
  card/                 Card (chrome wrapper) + Section (sub-section divider)
  badge/                Variant (generic) + Status (shared code lookup) + VariantForCode helper
  flash/                Banner (generic) + Success / Error / Warning / Info wrappers
  table/                table, table row, pagination, bulk-action bar
  drawer/               SidePanel (slide-in) + Modal (centered) + Confirm (preset). Shared FooterActions / Action structs; open/close via window.kibanOpenOverlay / kibanCloseOverlay; Escape closes topmost visible.
  spinner/              loading indicator (CSS class `.ds-spinner` in base.templ)
  tooltip/              CSS tooltip (`data-tooltip="…"` in base.templ)
  tabs/                 tab strip

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
| `input.File(name, label, accept, hint string, required bool)` | done | `<input type="file">` con `file:` prefix kiban-styleado. `accept` (ej. `.csv,text/csv`) para restringir el picker; pasá `""` para aceptar todo. `required?={required}` se emite (HTML5) para que el browser bloquee submit sin archivo. |
| `input.Hidden(name, value, id string)` | done | `<input type="hidden">`. Pasá `id=""` cuando no haga falta (la mayoría de los casos); poné `id` cuando el JS necesite leer/escribir el valor en runtime (ej. el script de geolocation que rellena lat/lng antes de submit). |
| `input.Textarea(name, label, value, errMsg, hint string, required bool, rows int, placeholder string, mono bool)` | done | `rows <= 0` → default 3. `placeholder` se emite siempre (vacío = no se muestra). `mono=true` swappea a `font-mono` para PEM blocks, listas de emails una-por-línea, etc. |
| `input.Checkbox(name, label string, value, enabled bool, attrs templ.Attributes, disabledHint string)` | done | Checkbox plano + label en línea. **`value="true"` obligatorio** — Gin no parsea "on" como bool. `disabledHint` se renderiza después del label en `text-kiban-ink4` solo cuando `!enabled` (útil para mostrar "— no configurado en …" en toggles disabled). Pasá `""` cuando no haga falta. |
| `input.CheckboxCard(name, title, help string, value, enabled bool, attrs templ.Attributes)` | done | Card-style: borde + padding + título arriba + helper text en muted abajo. Mismo contrato `value="true"`. `enabled=false` deshabilita y baja opacidad. Útil para listas de preferencias / settings cards. |
| `input.Toggle(name, label string, value, enabled bool, attrs templ.Attributes)` | done | Checkbox semántico, switch visual (Tailwind `peer-checked` para animar el tracker). Mismo contrato `value="true"` y `attrs` que Checkbox. |
| `input.RadioCard(name, value, title, subtitle string, checked bool, attrs templ.Attributes)` | done | Card-style radio: borde clickeable, resalta `border-kiban-primary bg-kiban-primary-soft` al checkearse. Múltiples cards con el mismo `name` forman el group. `subtitle` opcional. `attrs` para HTMX cuando la selección dispara un partial re-render del form. |

### Botones (`view/button/`)

| Componente | Estado | Notas |
|---|---|---|
| `button.Primary(opts)` | planned | bg `kiban-primary`, text white. |
| `button.Secondary(opts)` | planned | Outline, border `kiban-border`. |
| `button.Destructive(opts)` | planned | bg `red-600`. |
| `button.Icon(opts)` | planned | Solo icono, hover suave. |

`Opts` (struct) lleva: `Label`, `Type` (button/submit), `Disabled`, `HxPost`/`HxGet`/`HxTarget`/`HxSwap`/`HxConfirm`/`HxIndicator`/`HxDisabledElt`, `Class` (extra), `Form` (id of external form).

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

### Tables (`view/table/`)

| Componente | Estado | Notas |
|---|---|---|
| `table.PaginationConfig` | done | Struct con Page, HasPrev, HasNext, PageURL func(int) string, Target, Indicator. Caller construye con un closure sobre su `pageURL` local para preservar filtros entre páginas. |
| `table.Pagination(cfg)` | done | Anterior/Siguiente botones HTMX. Estados disabled-styled cuando edge. `cfg.PageURL` callback para que el caller maneje filter state. |
| `table.EmptyState(title, hint)` | done | Card centrada con border-top, "Aún no tenemos X qué mostrar". `hint` opcional. |
| `table.Table(headers, rows, opts)` | planned | Componente generico de table chrome. `opts`: `BulkSelect`, `RowHref`. Per-row rendering sigue siendo per-project (las columnas varían). |
| `table.BulkActionBar(opts)` | planned | Barra `:has(input:checked)` que aparece cuando hay selecciones. Botones de acción dentro. |

### Drawer / overlay (`view/drawer/`)

Convenciones del paquete:
- **Open / close** vía `window.kibanOpenOverlay(id)` / `kibanCloseOverlay(id)` (definidos en `view/layout/base.templ`). El caller cablea su trigger con `onclick="kibanOpenOverlay('the-id')"`; el componente se ocupa del close (close button + backdrop click + Escape global).
- **Escape** cierra solo el overlay topmost visible (por orden de `[data-kiban-overlay]:not(.hidden)` en el DOM) — multiple overlays apilados se cierran uno a la vez, igual que el comportamiento nativo del browser.
- **Sizes** compartidos: `sm` (max-w-sm) / `md` (max-w-md, default) / `lg` (max-w-lg) / `xl` (max-w-xl). String vacío cae a `md` (excepto Confirm, que cae a `sm`).
- **Z-index**: SidePanel = 40, Modal/Confirm = 50.
- **Footer actions** vía `FooterActions{PrimaryAction *Action, SecondaryActions []Action}`. PrimaryAction renderiza a la derecha (default variant `primary`); SecondaryActions a la izquierda en orden, default variant `secondary`. Si no hay acciones, no se renderiza el footer.
- **Action variants** alineados con la categoría buttons (futura): `primary` (kiban-primary), `secondary` (outline border-kiban-border, default para anchors), `danger` (red-600 — usar para confirms destructivos).

#### `Action` (struct compartido)

```go
type Action struct {
    Label    string             // visible text
    Variant  string             // "primary"|"secondary"|"danger"; default per slot
    Href     string             // when non-empty: <a href=Href>; otherwise <button>
    Type     string             // "button"|"submit" (default "button"); only for button mode
    Form     string             // for Type="submit": submits external form by id
    OnClick  string             // raw inline JS (e.g. "kibanCloseOverlay('id')" para dismiss tras submit)
    Attrs    templ.Attributes   // HTMX escape hatch (hx-post / hx-target / …)
    Disabled bool
}
```

#### Componentes

| Componente | Estado | Notas |
|---|---|---|
| `drawer.SidePanel(cfg SidePanelConfig)` | done | Slide-in desde la derecha. `cfg.ID` + `Title` + `Size` + `FooterActions`. Body via templ children, padding `px-6 py-4` aplicado al body (caller no se preocupa por el gutter). Patrón típico para filter-drawers: el caller renderiza un `<form id="filter-form">` en el body y la PrimaryAction usa `Type:"submit"` + `Form:"filter-form"` + `OnClick:"kibanCloseOverlay('id')"` para submit-and-dismiss. |
| `drawer.Modal(cfg ModalConfig)` | done | Centrado, backdrop oscuro (`bg-black/40`). Mismos campos que SidePanel + `Icon templ.Component` opcional renderizado a la izquierda del título (caller pasa el block fully styled — patrón típico kiban: `<div class="w-8 h-8 rounded-full bg-kiban-primary-soft text-kiban-primary flex items-center justify-center"><svg…/></div>`). Para flujos HTMX-form-bound (modal con submit): wrapeá el `@drawer.Modal(...)` entero en un `<form hx-post=… hx-on::after-request="if(event.detail.successful){ kibanCloseOverlay('id'); }">`. |
| `drawer.Confirm(cfg ConfirmConfig)` | done | Preset de Modal con shape fijo: title (opcional) + message + cancel/confirm buttons. Default size `sm`. `cfg.PrimaryAction` (full Action) es el botón de confirm; setear `Variant:"danger"` para deletes. Cancel se cablea automáticamente al `kibanCloseOverlay(id)`; label default "Cancelar". Para "are you sure?"-level simple usar `hx-confirm` nativo de HTMX; usar Confirm cuando se necesita styling kiban + HTMX wiring custom + título largo. |

### Tabs (`view/tabs/`)

| Componente | Estado | Notas |
|---|---|---|
| `tabs.Strip(items, activeKey)` | planned | Lista horizontal de tabs. `items []TabItem{Key, Label, Href}`. Item activo resaltado. |

### Form binding (`binding/`)

| Componente | Estado | Notas |
|---|---|---|
| `binding.FieldErrors(err) map[string]string` | done | Mapea `validator.ValidationErrors` → claves = nombre del tag `form:"…"`, valores = mensaje en español. Soporta tags `required`, `email`, `min`, `max`, `len`, `url`, `oneof`, `gt`, `gte`, `eqfield`. |
| `init()` en el package | done | Registra `TagNameFunc` en el validator de Gin para que use `form:"…"` como key. Corre automáticamente al importar el package. |
| `binding.MessageFor(tag, param)` (extensible) | planned | Helper para extender messages a tags custom (ej. `regexCURP`, `regexRFC`). Hoy `messageFor` es interno. |

### HTMX helpers (`htmx/`)

| Componente | Estado | Notas |
|---|---|---|
| `htmx.IsRequest(c) bool` | planned | `c.GetHeader("HX-Request") == "true"`. |
| `htmx.Redirect(c, url string)` | planned | Si HTMX request → header `HX-Redirect`; si no → `c.Redirect(302, url)`. |
| `htmx.TriggerName(c) string` | planned | Atajo a `c.GetHeader("HX-Trigger-Name")`. |

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
