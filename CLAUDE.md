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
  input/                text, phone (intl-tel-input wrapper), select, checkbox, toggle, textarea, hidden, country-code
  button/               primary, secondary, destructive, icon
  card/                 card, section (titled card)
  badge/                status pill, type badge
  table/                table, table row, pagination, bulk-action bar
  drawer/               slide-in side panel (filters), modal
  flash/                success / error / warning banners
  spinner/              loading indicator
  tooltip/              CSS tooltip (data-tooltip="…")
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
| `view.Render(c, status, component)` | done | En `view/render.go`. |
| `layout.Layout(cfg)` | planned | Shell completa: HTML, head, scripts, topbar, sidebar (level 1/2), main content. `cfg` es `LayoutConfig` con Title, SectionName, ToolKey, ActiveKey, User, SpaceID, ShellURL, Tools, SubItems, LogoutAction. |
| `layout.Topbar(cfg)` | planned | Hamburger izquierda + logo kiban + space chip + Developers button + user menu con logout. |
| `layout.IconRail(cfg)` | planned | Nivel 1: iconos de tools + docs en la parte inferior. Tool activo resaltado. Tooltips en hover. |
| `layout.SubNav(cfg)` | planned | Nivel 2: items de la sección activa. Item activo resaltado. |
| Slide animation 1↔2 | planned | CSS-only en `base.templ` con `.sidebar-slot` + `.sidebar-track`. JS persiste en `localStorage[<project>-sidebar-level]`. Hamburger toggla. |
| Tooltip CSS (`data-tooltip`) | planned | Bubble dark-ink, instant on hover/focus, ignora pointer events. |

### Iconos (`view/icons/`)

| Componente | Estado |
|---|---|
| `icons.Home`, `icons.Users`, `icons.ShoppingBag`, `icons.GitMerge`, `icons.CreditCard`, `icons.Code`, `icons.Database`, `icons.BookOpen`, `icons.FileText` | planned |
| `icons.ChevronDown`, `icons.Filter`, `icons.Sort`, `icons.More`, `icons.Close`, `icons.Hamburger`, `icons.Sidebar` | planned |

Convención: 20×20 viewBox=24, `stroke="currentColor"`, sin fill — colorables con clases `text-…` de Tailwind.

### Inputs (`view/input/`)

| Componente | Estado | Notas |
|---|---|---|
| `input.Text(name, label, value, errMsg, required)` | planned | Border `kiban-border`, focus `kiban-primary`, error → `border-red-400`. Helper text rojo bajo el input cuando `errMsg != ""`. |
| `input.Phone(label, ccName, ccValue, phoneName, phoneValue, errMsg, required)` | planned | Wrapper de intl-tel-input. Renderiza un `<input type="tel">` visible (sin `name`) + dos `<input type="hidden">` (countryCode + phone). El init JS en `layout/base.templ` les asocia el widget y sincroniza los hiddens en `countrychange`/`input`. |
| `input.Select(name, label, value, errMsg, options, required)` | planned | `options []SelectOption{Value, Label}`. |
| `input.Checkbox(name, label, value, enabled)` | planned | **`value="true"` obligatorio** — Gin no parsea "on" como bool. |
| `input.Toggle(name, label, value, enabled)` | planned | Checkbox semántico, switch visual. |
| `input.Textarea(name, label, value, errMsg, required, rows)` | planned | |
| `input.Hidden(name, value)` | planned | Helper para hidden context fields (page, filter state, etc.). |
| `input.Date(name, label, value, errMsg, required)` | planned | `<input type="date">`, validación nativa. |

### Botones (`view/button/`)

| Componente | Estado | Notas |
|---|---|---|
| `button.Primary(opts)` | planned | bg `kiban-primary`, text white. |
| `button.Secondary(opts)` | planned | Outline, border `kiban-border`. |
| `button.Destructive(opts)` | planned | bg `red-600`. |
| `button.Icon(opts)` | planned | Solo icono, hover suave. |

`Opts` (struct) lleva: `Label`, `Type` (button/submit), `Disabled`, `HxPost`/`HxGet`/`HxTarget`/`HxSwap`/`HxConfirm`/`HxIndicator`/`HxDisabledElt`, `Class` (extra), `Form` (id of external form).

### Display (`view/card/`, `view/badge/`, `view/flash/`, `view/spinner/`, `view/tooltip/`)

| Componente | Estado | Notas |
|---|---|---|
| `card.Card` | planned | White, border `kiban-border`, rounded-md, padding 6. |
| `card.Section(title, subtitle)` | planned | Card con header (heading + texto secundario). |
| `badge.Status(code, label)` | planned | Pill con colores derivados del código (`PAID`→emerald, `PENDING`→amber, `EXPIRED`→red, etc.). Override por proyecto si necesario. |
| `badge.Type(label, variant)` | planned | Variant: `success`, `info`, `warning`, `error`, `neutral`. |
| `flash.Success(msg)`, `flash.Error(msg)`, `flash.Warning(msg)` | planned | Banners post-mutación. |
| `spinner.Default()` | planned | Spinner kiban-primary. CSS class `ds-spinner` definida en `base.templ`. |
| `tooltip` | planned | No es un componente sino un patrón: añade `data-tooltip="texto"` a cualquier elemento. CSS in `base.templ`. |

### Tables (`view/table/`)

| Componente | Estado | Notas |
|---|---|---|
| `table.Table(headers, rows, opts)` | planned | `opts`: `BulkSelect`, `RowHref` (data-href para click delegation). |
| `table.Pagination(view, baseURL)` | planned | Anterior/Siguiente botones HTMX. `view`: estado de paginación (Page, HasPrev, HasNext, ItemsPerPage). |
| `table.BulkActionBar(opts)` | planned | Barra `:has(input:checked)` que aparece cuando hay selecciones. Botones de acción dentro. |

### Drawer / overlay (`view/drawer/`)

| Componente | Estado | Notas |
|---|---|---|
| `drawer.SidePanel(id, title, body)` | planned | Slide-in desde la derecha. Backdrop click cierra. |
| `drawer.Modal(id, title, body, footer)` | planned | Centrado, backdrop oscuro. |
| `drawer.Confirm(id, message, confirmLabel)` | planned | Wrapper sobre Modal con dos botones. (Alternativa simple: usar `hx-confirm` nativo de HTMX.) |

### Tabs (`view/tabs/`)

| Componente | Estado | Notas |
|---|---|---|
| `tabs.Strip(items, activeKey)` | planned | Lista horizontal de tabs. `items []TabItem{Key, Label, Href}`. Item activo resaltado. |

### Form binding (`binding/`)

| Componente | Estado | Notas |
|---|---|---|
| `binding.FieldErrors(err) map[string]string` | planned | Mapea `validator.ValidationErrors` → claves = nombre del tag `form:"…"`, valores = mensaje en español. |
| `init()` en el package | planned | Registra `TagNameFunc` en el validator de Gin para que use `form:"…"` como key. Corre automáticamente al importar. |
| `binding.MessageFor(tag, param)` | planned | Helper para extender messages a tags custom (ej. `regexCURP`, `regexRFC`). |

### HTMX helpers (`htmx/`)

| Componente | Estado | Notas |
|---|---|---|
| `htmx.IsRequest(c) bool` | planned | `c.GetHeader("HX-Request") == "true"`. |
| `htmx.Redirect(c, url string)` | planned | Si HTMX request → header `HX-Redirect`; si no → `c.Redirect(302, url)`. |
| `htmx.TriggerName(c) string` | planned | Atajo a `c.GetHeader("HX-Trigger-Name")`. |

### Middleware (`middleware/authcookie/`)

| Componente | Estado | Notas |
|---|---|---|
| `authcookie.Middleware(authorizeUC, loginURL string) gin.HandlerFunc` | planned | Lee cookies `kiban_session` + `kiban_space_id`, valida vía `IAuthorizationAuthorizeWithSessionUseCase` de go-kiban, setea `controller_core_model.CONTEXT_KEY_AUTHORIZATION_OBJECT` en el contexto. En falla emite `HX-Redirect` (HTMX) o `c.Redirect` (browser). |
| `authcookie.GetAuthorization(c)` | planned | Atajo al lookup en contexto. |

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
