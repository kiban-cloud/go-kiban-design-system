# go-kiban-design-system

Shared **templ + Tailwind + HTMX** primitives for every kiban htmx-based
project. Layout shell, inputs, buttons (planned), tables, drawers,
flash banners, badges, form-binding helpers, HTMX helpers, and cookie
auth middleware all live here so the look-and-feel — and the wording
of error messages — stays consistent across kiban services.

> **Working in this repo with an AI agent?** See [`CLAUDE.md`](./CLAUDE.md)
> for the per-component catalog, conventions, and naming rules. The
> README is the human-facing landing page; CLAUDE.md is the detailed
> implementation reference.

## Who uses it

- [`rekon-backend`](../rekon-backend/) — payments / payment orders /
  withdrawals / dispersion backoffice.
- [`crm-backend`](../crm-backend/) — customers / contacts /
  payment-method configuration.
- Future kiban htmx-based projects.

## How to consume

```go
require github.com/kiban-cloud/go-kiban-design-system v0.x.y
```

For local development, point Go at a sibling checkout via a `replace`
directive in your project's `go.mod`:

```go
replace github.com/kiban-cloud/go-kiban-design-system => ../go-kiban-design-system
```

When you publish a new version, drop the replace and bump the version
in `go get`.

## Repo layout

```
view/
  layout/      Layout shell: Topbar / IconRail / SubNav / global JS helpers
  icons/       Shared SVG icon set
  action/      Action + Group structs reused by drawer footers and table.BulkActionBar
  input/       Text, Password, Number, Phone, Select, Date, File, Hidden, Textarea, Checkbox, CheckboxCard, Toggle, RadioCard
  button/      Primary / Secondary / Destructive / Icon (planned — thin wrappers around action.Render)
  card/        Card + Section
  badge/       Variant + Status + VariantForCode
  flash/       Banner + Success / Error / Warning / Info
  table/       Table + Row + BulkActionBar + Pagination + EmptyState
  drawer/      SidePanel + Modal + Confirm
  tabs/        Strip (in-page tabs)
binding/       Form-error translator (Spanish messages for validator tags)
htmx/          IsRequest / Redirect / TriggerName helpers
middleware/    Cookie-auth middleware
```

## Conventions all kiban projects share

A few cross-cutting conventions to keep services consistent — these are
the ones that span multiple components, not the per-component details
(those live in [`CLAUDE.md`](./CLAUDE.md)).

### Form-error messages → always go through `binding/`

When a Gin handler calls `c.ShouldBind(&req)` and gets a validation
error, the **only** way to turn it into Spanish messages for the
template is `binding.FieldErrors(err)`. Every kiban project does this —
do not copy-paste the message-translation switch into your project.

When a project registers a project-specific validator tag (e.g.
`regexRFC` for Mexican fiscal IDs, `regexCURP`, `regexCLABE`), register
its Spanish message in this package's registry at startup, **not in a
local helper**:

```go
// crm-backend/cmd/api/main.go (or similar startup path)
binding.RegisterMessage("regexRFC", func(_ string) string {
    return "RFC inválido"
})
```

This way:

- The wording stays consistent across kiban projects (if rekon ever
  also validates RFC, both projects' Spanish messages agree).
- New kiban projects pick up the same translations by default — no
  forking the dispatcher.
- The translation table for your project is a small, discoverable
  startup hook, not a switch buried in a helper file.

The full list of built-in tags + custom-registration rules is in
[`CLAUDE.md`'s `binding/` section](./CLAUDE.md#form-binding-binding).

### Open / close overlays → `kibanOpenOverlay` / `kibanCloseOverlay`

Drawers, modals, and confirmations all live in `view/drawer/` and
share a global JS helper (defined in `view/layout/base.templ`):

```html
<button onclick="kibanOpenOverlay('the-drawer-id')">Filtros</button>
```

Close happens via the close button + backdrop click + the global
Escape-key listener (closes only the topmost visible overlay). Don't
ship per-project open/close JS.

### HTMX-vs-browser branching → `htmx.IsRequest` and `htmx.Redirect`

Handlers that render a page differently for HTMX swaps vs full browser
navigation use `htmx.IsRequest(c)`. Post-mutation redirects use
`htmx.Redirect(c, url)` — it sends `HX-Redirect` for HTMX requests and
falls back to `c.Redirect(302, url)` for browser POSTs. Calling
`c.Redirect` directly in an HTMX request is a silent footgun; the
helper is the canonical way.

### Status badges → reuse `badge.VariantForCode` for shared codes

`badge.Status(code, label, size)` translates kiban-shared backend
status codes (`PAID`, `PENDING`, `EXPIRED`, `CANCELLED`, `VALIDATED`,
`TO_VALIDATE`, etc.) to a colour variant. Project-specific codes
(SPEI/Banxico lifecycle, reconciliation states, etc.) get a local
`<feature>Variant(code) string` mapper and use `badge.Variant(label,
variant, size)` directly. Don't extend `VariantForCode` with codes
that only one project uses.

## Build & generate

```bash
go tool templ generate    # regenerate *_templ.go for every templ file
go build ./...            # compile
go vet ./...              # static check
go mod tidy               # tidy go.mod / go.sum
```

`*_templ.go` files are committed (regenerated by `templ generate`).
Don't hand-edit them.

## Versioning

Semver. Breaking changes bump the major. Pre-1.0.0 may include
breaking changes in minor releases, so consumers should pin to an
exact version (and use the `replace` directive locally).
