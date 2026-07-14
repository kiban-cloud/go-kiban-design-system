# Plan: consola kiban responsive (móvil)

> **Fuente de verdad** de la iniciativa "hacer responsive toda la consola de kiban cloud"
> (tarea `hacer kiban-cloud responsive`, creada 2026-07-09 por alex@kiban.com, asignada a antonio.blancas@kiban.com).
>
> Este doc vive en `go-kiban-design-system` porque el ~80% del trabajo es en la shell compartida.
> Está pensado para ejecutarse **en etapas y/o repartido entre varios agentes**: cada fase es
> autocontenida y tiene su checklist. Al terminar una tarea, marca su casilla y anota el commit/branch.

## Objetivo

Toda la consola htmx (kiban-cloud, rekon, klin, link, workfloo, crm) debe verse y funcionar
en celular (~375px de ancho). Trabajar **a nivel design system** para homogeneizar al máximo:
un arreglo en la shell → los 6 módulos lo heredan.

**Excepción reconocida:** el *editor* de workfloos (canvas de nodos + wizards multi-columna) no
será ideal en celular. Se hará un intento **best-effort** (navegable con scroll, nodos apilados),
sin bloquear al usuario, pero sin garantizar paridad con desktop.

## Estrategia

`design-system-first`. Orden: shell compartida → primitivas compartidas (tablas) → páginas por módulo.
Módulos ordenados de menor a mayor superficie para validar el approach temprano:
**Link (7) + crm (4) → rekon (14) → klin (11) → kiban-cloud (27) → workfloo (56)**.

## Convenciones responsive (aplican a TODO el trabajo)

- **Mobile-first, breakpoints estándar de Tailwind.** Base = móvil; `md:` (768px) = tablet/desktop; `lg:` = desktop amplio.
- **Regla del pulgar:** `< md` es "móvil". El corte de la shell (rails → drawer) es en `md`.
- **Nada de anchos fijos que no colapsen.** `w-56` → `w-full md:w-56`; `min-w-[220px]` → quitar o `md:min-w-[220px]`.
- **Padding fluido:** `p-4 md:p-8`, `px-4 md:px-6`.
- **Grids de formulario:** `grid-cols-1 md:grid-cols-2`.
- **Modales/drawers grandes:** full-screen (`w-full h-full` / `max-w-none`) en móvil, tamaño normal en `md:`.
- **Nunca overflow horizontal del `<body>`.** Contenido ancho (tablas, code, diagramas) scrollea dentro de su propio contenedor `overflow-x-auto`.
- **Verificar** con el viewport a 375px (iPhone SE) y 768px (tablet). El runtime del DS ya trae `<meta viewport>`.

---

## Fase 0 — Shell responsive en el design system  ✅ HECHO (2026-07-10)

Archivos tocados: `view/layout/base.templ`, `view/layout/nav.templ`, `view/layout/admin.templ` (+ `_templ.go` regenerados). **Este es el multiplicador.**
Verificado renderizando el shell a 375px, 753px y 1280px con un harness throwaway (ya eliminado): topbar sin overflow y avatar visible en móvil, drawer off-canvas entra/cierra (backdrop/Escape/link), y el shell desktop (SubNav siempre visible + hamburger togglea el rail 0↔224px) queda idéntico. `go build`/`go vet` limpios.

Estructura actual de la shell (`Layout`):
```
root (min-h-screen flex flex-col, data-sidebar-root)
  Topbar          h-14, fila horizontal fija
  div flex flex-1
    sidebar-rail-slot  (width 0↔14rem toggle CSS) → IconRail (w-56)
    SubNav             (w-56, SIEMPRE visible)   ← problema #1 en móvil
    main               (flex-1, p-8)             ← problema #2 (padding)
```

### Tareas
- [x] **T0.1 — SubNav + IconRail como drawer off-canvas en `< md`.**
      Rail + sub-nav envueltos en `.sidebar-nav-group` (en `Layout`). CSS scoped a `@media (max-width:767px)`:
      `position:fixed`, `translateX(-100%)` → `translateX(0)` cuando `data-sidebar-mobile-open` está en el root.
      Desktop 100% intacto. El hamburger (`data-sidebar-show-tools`) elige rail-toggle (desktop) vs drawer
      (móvil) con `matchMedia('(max-width:767px)')`. En el drawer el rail va siempre expandido y ambos rails
      se apilan a ancho completo.
- [x] **T0.2 — Backdrop + cierre.** `[data-sidebar-backdrop]` (sibling de `main`, `display:none` salvo móvil-open).
      Cierra por: click en backdrop, click en cualquier `.sidebar-nav-group a`, Escape, y al pasar a ≥ md.
      El mobile-open NO se persiste.
- [x] **T0.3 — Topbar compacto en `< md`.** Logo `w-24 md:w-32`; chip de espacio `min-w-0 max-w-[32vw] md:min-w-[220px]`;
      Developers `hidden md:inline-block`; divisor pre-avatar `hidden md:block`; header `gap-2 px-3 md:gap-4 md:px-4`.
      Verificado: avatar visible (right 369 < 375) y sin overflow horizontal.
- [x] **T0.4 — `<main>` padding fluido:** `p-4 md:p-8` en `Layout` y `AdminLayout`.
- [x] **T0.5 — AdminLayout responsive.** Pase ligero: padding fluido, logo `w-24 md:w-32`, gaps reducidos,
      nav horizontal con `overflow-x-auto whitespace-nowrap` para que no desborde. (Sin hamburger dedicado —
      la nav admin tiene pocas secciones y con scroll horizontal basta. Verificación visual pendiente en un
      consumidor admin real, p.ej. klin-internal.)
- [x] **T0.6 — Regenerar y verificar:** `templ generate` + `go build ./...` + `go vet ./...` limpios.
      Verificado visualmente a 375/753/1280px con harness throwaway (eliminado).
- [x] **T0.7 — Actualizar `CLAUDE.md` del DS** — filas "Tools rail toggle" + nueva "Mobile nav drawer (< md)".

**Definición de "hecho" Fase 0:** ✅ a 375px no hay scroll horizontal, la nav es accesible vía hamburger
(drawer), el topbar no desborda, y el contenido usa el ancho completo con padding cómodo. Desktop idéntico.

**Pendiente de arrastre para las fases de módulo:** las *páginas* internas (tablas anchas, forms multi-columna,
modales grandes) siguen sin ajustar — eso es Fase 1 (tabla del DS) + fases por módulo. La tabla de ejemplo del
harness scrolleaba dentro de su card en móvil, lo cual valida que el approach "overflow contenido" funciona.

---

## Fase 1 — Tabla responsive en el design system  ✅ HECHO (2026-07-10)

Archivo: `view/table/` (`table.templ`, `types.go`, `bulk_action_bar.templ`). La mayoría de las vistas de
lista de los 6 módulos usan `table.Table`/`table.Row` → arreglo gratis en todos lados.

### Tareas
- [x] **T1.1 — Patrón elegido:** (a) **scroll horizontal contenido** como baseline universal. Se descartó
      (b) "tarjetas apiladas" como default porque las celdas son `templ` children opacos (sin metadata de
      label por columna), así que requeriría cambiar el API de `Row` y todos los callers. (b) queda como
      opt-in futuro para listas puntuales de alto tráfico, resuelto en el callsite del módulo.
- [x] **T1.2 — Implementado en `view/table`, cero regresión en desktop.** El `<table>` va envuelto en
      `overflow-x-auto` (contención universal: una tabla ancha ya no rompe el layout de la página). Nuevo
      `TableConfig.MinWidth` pone un piso de ancho **solo bajo `md`** vía `min-w-[40rem] md:min-w-0`
      (default 640px; `"0"` = sin piso; valor Tailwind = piso custom). Como el piso lleva `md:min-w-0`,
      **desktop no cambia para ningún caller** (ni tablas dentro de drawers angostos). `BulkActionBar` gana
      `flex-wrap`. `Pagination`/`EmptyState` ya eran mobile-safe.
- [x] **T1.3 — Verificado** con `table.Table` real (7 columnas, MinWidth default) dentro de la shell:
      a 375px `body` no scrollea (layout intacto) y la tabla (726px) scrollea dentro de su caja; a 1280px
      la tabla llena el ancho sin scroll forzado. Regenerado + `go build`/`go vet` limpios.
- [x] **T1.4 — Documentado** en `CLAUDE.md` del DS (fila `table.Table` + contrato `MinWidth`).

**Nota para las fases de módulo:** las tablas MUY anchas (8+ columnas, IDs/fechas largas) pueden querer un
`MinWidth` mayor (ej. `"56rem"`); las de 2-3 columnas cortas, `MinWidth:"0"`. Eso se ajusta por callsite al
auditar cada módulo — el default (640px) es un punto medio razonable que no rompe nada.

---

## Fases por módulo (2–7)

Para cada módulo el trabajo es el mismo patrón: recorrer sus `*.templ`, aplicar las convenciones
responsive a las páginas propias (forms multi-columna → 1 col en móvil, modales grandes → full-screen,
grids de dashboard, spacing). La shell y las tablas ya vienen del DS.

**Cómo trabajar un módulo (receta para agentes):**
1. `find <repo> -name "*.templ"` para el inventario.
2. Priorizar: dashboards/listas/detalle/forms que se usen desde móvil.
3. Aplicar convenciones. Si algo debería ser reusable → subirlo al DS, no resolverlo local.
4. `go tool templ generate` en el repo (si aplica su comando) + `go build ./...` + `go vet ./...`.
5. Verificar visualmente a 375px las páginas tocadas.
6. Marcar checklist + anotar branch/commit aquí abajo.

### Fase 2 — Link (`microservices`, 7 templ)  ✅ HECHO (2026-07-10)
Vistas: configuration, apikeys, vault, dashboard, testcases (list+detail), layout.

**Hallazgo clave:** Link casi no necesitó cambios propios — sus páginas ya componen primitivas responsive
del DS y usan `flex-wrap`/`w-full`/forms de una columna. La auditoría en cambio destapó **un bug del DS**
(`tabs.Strip` desbordaba en móvil) que se arregló en el DS y beneficia a todos los módulos.

- [x] Dashboard (header simple + tabla de 10 cols) — hereda scroll de Fase 1. OK sin cambios.
- [x] Listas (apikeys, vault, configuration, testcases) — todas usan `@ds_table.Table` → heredan scroll. OK.
- [x] testcases/detail — ya usa `flex-wrap` en meta/acciones, `w-full` en inputs, cards de una columna,
      editor Ace full-width. OK sin cambios.
- [x] Modales (apikeys/testcases): el `drawer.Modal`/`SidePanel` del DS ya es responsive (`w-full mx-4` +
      `max-w-*`, `max-h-[calc(100vh-2rem)]`). OK sin cambios.
- [x] **Fix DS: `tabs.Strip`** (usado por configuration + otros módulos) → `overflow-x-auto` en el contenedor
      + `shrink-0 whitespace-nowrap` en cada tab. Test `TestStrip_ContainerBaselineMatchesActiveBorderWidth`
      actualizado (ya no pinnea el string exacto de clase). Verificado con 8 tabs a 375px: scrollea en su
      caja, sin overflow de página ni scroll vertical.
- [x] Build (`go build` del árbol htmx) + `go vet` + **los 6 paquetes `view/` de Link pasan sus tests**.
      Verificación visual del strip a 375px con harness throwaway (eliminado).

### Fase 3 — crm (`crm-backend`, 4 templ)  ✅ HECHO (2026-07-10)
Vistas: customers/{list, form, bulk_upload}, layout. 3 fixes de página propios (chicos):
- [x] **`form.templ` typeToggle** — `grid grid-cols-2` (forzado) → `grid-cols-1 sm:grid-cols-2` para los RadioCards PF/PM.
- [x] **`bulk_upload.templ`** — mismo fix `grid-cols-2` → `grid-cols-1 sm:grid-cols-2` (RadioCards PF/PM).
- [x] **`list.templ` ActionBar** — fila de 3 botones (`Nuevo cliente`/`Carga masiva`/`Filtros`) `flex justify-end`
      sin wrap → agregado `flex-wrap` para que envuelvan en móvil en vez de cliparse.
- [x] Sin cambios: la tabla de clientes (`dstable.Table` + `dstable.Row`, 5 cols + bulk) hereda scroll de Fase 1;
      el resto del form ya usa `grid-cols-1 md:grid-cols-*`; `ContactRow` responsive; FiltersDrawer = `SidePanel`
      del DS (responsive); headings/actionbar de `bulk_upload` con `flex-wrap`.
- [x] `templ generate` + `go build ./...` + `go vet ./...` limpios. **Nota:** los paquetes `internal/view/customers`
      y `.../layout` **no tienen tests** (gap preexistente del repo, no introducido por este cambio); los fixes son
      presentacionales (solo clases Tailwind), sin ramas/contenido nuevos.

### Fase 4 — rekon (`rekon`, 14 templ)  ✅ HECHO (2026-07-10)
Vistas: dashboard, developers, dispersion(create+list), notifications, paymentmethods(detail+list),
paymentorders(create+detail+list), payments(detail+list), withdrawals, layout.

- [x] **3 fixes idénticos:** el filtro de fecha "Desde/Hasta" en los drawers de `payments/list`,
      `paymentorders/list` y `withdrawals/list` usaba `grid grid-cols-2` (forzado) → `grid-cols-1 sm:grid-cols-2`.
- [x] Sin cambios: dashboard ya usa `grid-cols-1 md:grid-cols-3` (stat cards) + `lg:grid-cols-3` (charts);
      tablas (payments 6 cols, paymentorders 9 cols, withdrawals 6 cols) usan `dstable.Table` → scroll heredado,
      MinWidth default OK; headers de cards (`justify-between`) ya usan `min-w-0`+`shrink-0`; footers de forms
      son de 2 botones; no hay `w-[…]` fijos ni `tabs.Strip`.
- [x] `templ generate` + `go build ./...` limpios; **view tests verdes** (payments/dispersion/paymentmethods
      con tests pasan; paymentorders/withdrawals/dashboard/etc. sin tests — gap preexistente). Cambios class-only.
- [x] ⚠️ **Bug preexistente (fuera de alcance, NO tocado):** `go vet` marca `internal/repository/client_payment_methods/ClientPaymentMethodsRepository.go`
      (bson.E unkeyed ×2 + unreachable code). No relacionado con responsive; reportado para arreglo aparte.

### Fase 5 — klin (`klin-backend`, 11 templ)  ✅ HECHO (2026-07-10, con caveat de entorno)
Vistas: apikeys, dashboard, deliveries(create+detail+list), staff(auth/login, deliveries detail+list, layout, tools), layout.

- [x] **`staff/deliveries/list.templ`** — filtro fecha Desde/Hasta `grid-cols-2` → `grid-cols-1 sm:grid-cols-2`.
- [x] **`deliveries/list.templ` header** — tenía título + 2 botones (Filtros + Nueva entrega) en
      `flex items-center justify-between` que apretaban en móvil → `flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between`
      (stack en móvil, fila en `sm+`).
- [x] Sin cambios: deliveries/create + detail usan `grid-cols-1 md:grid-cols-2/3`; dashboard `md:grid-cols-2 lg:grid-cols-4`;
      staff/deliveries/detail ya tiene `flex-wrap`; `tabs.Strip` (staff/tools, deliveries/detail) responsive desde Fase 2;
      tablas heredan scroll; headers de 1 botón (apikeys, staff/deliveries) funcionan; login simple. Sin `w-[…]` fijos.
- [x] `templ generate` corrido; `_templ.go` regenerados.
- [x] **Fix de dependencia (bloqueaba el build de TODO el repo, ahora resuelto):** klin-backend fijaba
      `bitbucket.org/alexandregrin/go-kiban v0.0.297`, una tag **borrada/retagueada** en el remoto (la secuencia
      de tags salta 296→298; el `go.sum` aún tenía el hash de 297 de cuando existía). Bumpeado en `go.mod` a
      **v0.0.298** (sucesora inmediata; su `go.mod` hash es idéntico al de 295/297 → mismas deps, sin transitivas
      nuevas) + `go mod download` + `go mod tidy` (limpió las líneas 297 stale del go.sum). **Nota:** este cambio
      de `go.mod`/`go.sum` es parte del working tree de klin-backend.
- [x] **Verificado:** `go build ./...` limpio, `go vet ./...` sin issues nuevos (2 warnings PREEXISTENTES no
      relacionados: `internal/domain/sepomex/model/model.go` json tag duplicado + `internal/repository/deliveries/DeliveriesRepository.go:215`
      context leak), **10 paquetes de unit tests OK / 0 FAILs**, y **todos los view tests verdes** (apikeys,
      dashboard, deliveries, staff/deliveries, staff/tools).

### Fase 6 — kiban-cloud (`kibancloud`, 27 templ)  ✅ HECHO (2026-07-10)
Vistas: admin(billings×4, home, profile, users×2), auth(×10), notifications, spaces(×7), workfloo(host).

- [x] **`notifications/notifications.templ`** — único `<table>` hand-rolled del repo (los demás usan `ds_table.Table`);
      estaba en un card `overflow-hidden` → se aplastaba/clipaba en móvil. Envuelto en `<div class="overflow-x-auto">`
      + `min-w-[36rem] md:min-w-0` (patrón Fase 1). El resto de tablas hereda scroll de `ds_table.Table`.
- [x] Sin cambios: sin `grid-cols-N` sin prefijo (todos `md:`/`sm:`); auth es card centrada `max-w-md px-4`
      (responsive por naturaleza, cubre las 10 vistas auth); `tabs.Strip` responsive de Fase 2; sin `w-[…]` fijos;
      headers/filas flex restantes son label+control seguros.
- [x] `templ generate` + `go build ./...` limpio (rc=0) + `go vet ./view/...` limpio. (notifications sin test files —
      gap preexistente; cambio class/wrapper-only.)

---

## ⚠️ PASO OBLIGATORIO PENDIENTE — Publicar el DS y bumpear consumidores

**Descubrimiento (2026-07-10):** NINGÚN repo consumidor tiene activo el `replace => ../go-kiban-design-system`.
Todos usan una versión **publicada** (tag) del DS, y **distinta por repo**:

| Repo | DS pin actual |
|---|---|
| crm-backend | v0.0.22 |
| rekon | v0.0.29 |
| klin-backend | v0.0.35 |
| kibancloud | v0.0.38 |
| microservices (Link) | v0.0.42 |
| workfloo | v0.0.52 |

**Implicaciones:**
1. Los cambios de Fases 0/1/2 en el DS (shell drawer, tabla scroll, tabs scroll) **viven solo en el working tree
   local del DS** — hoy **no llegan a ningún módulo**. Se verificaron vía harness (que usa el DS local), no vía el
   build de cada consumidor.
2. Los fixes de página por módulo (Fases 2–6) **sí** están en cada repo y toman efecto en su build.
3. Que los 6 corran shells de 6 versiones distintas del DS es, de hecho, la causa de la inconsistencia que motiva
   esta tarea — bumpearlos todos a UNA versión nueva entrega lo responsive **y homogeneíza** (lo que pide Alex).

**Para shippear las Fases 0/1/2 a los módulos (cuando el usuario dé el OK a commitear/publicar):**
- [ ] Commit + tag de una versión nueva del DS (ej. `v0.0.53`) con los cambios de Fases 0/1/2.
- [ ] Bumpear los 6 consumidores a ese tag y `go mod tidy` — usar la skill **`ds-bump`** (`go get @tag` + tidy +
      build por repo, sin commitear). Nota: `ds-bump` hoy lista rekon/crm/link/klin; workfloo y kiban-cloud "más
      tarde" → confirmar que cubra los 6.
- [ ] Verificación visual de la shell responsive en cada módulo real (a 375px) tras el bump.
- [ ] Regenerar `_templ.go` de cada consumidor si su versión de templ lo requiere; `go build`/`go vet`/tests verdes.

### Fase 7 — workfloo (`workfloo`, 56 templ)  ⬜ PENDIENTE
El módulo más grande. Separar en dos frentes:

**7A — workfloo "normal" (responsive completo):** ✅ HECHO (2026-07-11)
- [x] **`workfloo_definition/list.templ`** — filtro fecha "Creado desde/hasta" `grid-cols-2` → `grid-cols-1 sm:grid-cols-2`.
- [x] **`history/detail.templ`** — grid de detalle API (Duración/Status/…) `grid-cols-2` → `grid-cols-1 sm:grid-cols-2`.
- [x] **Ejecución/Portal (PRIORIDAD) — ya mobile-first, sin cambios.** `portal/execution.templ` envuelve en
      `min-h-screen ... px-4 py-6` + `w-full max-w-2xl p-6 sm:p-8`; fue construido para celular. Las 10 vistas
      portal + las 7 de execute no tienen grids/anchos/tablas problemáticos; execute hereda la shell (Fase 0).
- [x] Sin cambios: dashboard KPI tiles (`grid-cols-3`) = 3 tiles compactos side-by-side intencional; las 2 tablas
      raw (dashboard pivot, ab_testing) YA están envueltas en `overflow-x-auto`; headers de listas (api_keys,
      backtesting, scenery, ab_testing, fronts) son de 1 botón (seguros); `tabs.Strip` responsive (Fase 2).

**7B — editor de workfloos (best-effort):** ✅ HECHO (2026-07-11)
- [x] **`edit/decision.templ`** — grid de condiciones del nodo decisión `grid-cols-2` → `grid-cols-1 sm:grid-cols-2`
      (ayuda en el drawer de config en móvil).
- [x] **Canvas ya no rompe el layout:** el editor tiene `[data-kiban-canvas-viewport] absolute inset-0 overflow-auto`
      + controles de zoom (out/readout/in). En móvil el canvas panea/zoomea dentro de su marco — no causa overflow
      de página. El editor hereda la shell responsive (Fase 0), y la config de nodos usa `drawer.Modal/SidePanel`
      del DS (w-full en móvil). **Cumple el bar "navegable y no roto".**
- [x] Dejado como está (best-effort, orientado a desktop, según lo acordado): wizards densos multi-columna
      (RULESET 4 cols), branch columns con `min-w-[160px]`, y `edit/variables.templ` grid `dt/dd` 1:2 (layout
      intencional). No serán ideales en celular pero no rompen la navegación.
- [x] `templ generate` + `go build ./...` (rc=0) + `go vet ./internal/view/...` limpios + **view tests verdes**
      (workfloo_definition, edit, history, portal, execute).

---

## Tablero de progreso

| Fase | Módulo/Área | Estado | Branch / commit | Notas |
|---|---|---|---|---|
| 0 | DS shell | ✅ Hecho (2026-07-10) | sin commit (working tree) | Drawer móvil + topbar compacto + padding fluido. base/nav/admin.templ. Verificado 375/1280px. |
| 1 | DS tabla | ✅ Hecho (2026-07-10) | sin commit (working tree) | overflow-x-auto + MinWidth (piso solo < md). Cero regresión desktop. Verificado 375/1280px. |
| 2 | Link (microservices) | ✅ Hecho (2026-07-10) | sin commit (working tree) | Sin cambios de página (ya responsive). Fix DS: tabs.Strip scrollable. Tests de Link verdes. |
| 3 | crm-backend | ✅ Hecho (2026-07-10) | sin commit (working tree) | 3 fixes: grid-cols-2→1 sm:2 (×2) + ActionBar flex-wrap. Sin tests preexistentes en las vistas. |
| 4 | rekon | ✅ Hecho (2026-07-10) | sin commit (working tree) | 3 fixes date-range grid-cols-2→1 sm:2. Vet: bug preexistente en repo (fuera de alcance). |
| 5 | klin-backend | ✅ Hecho (2026-07-10) | sin commit (working tree) | 2 fixes: grid date-range + header stack. + fix dep: go-kiban 297 (tag borrada) → 298. build/vet/tests verdes. |
| 6 | kibancloud | ✅ Hecho (2026-07-10) | sin commit (working tree) | Fix tabla notifications (overflow-x-auto). build/vet limpios. |
| — | **DS publish + bump 6 consumidores** | ⬜ **PENDIENTE (obligatorio para shippear Fases 0/1/2)** | — | Ningún repo usa replace local; publicar DS tag + ds-bump. Ver sección arriba. |
| 7A | workfloo (normal + ejecución) | ✅ Hecho (2026-07-11) | sin commit (working tree) | 2 grids (list, history). Portal ya mobile-first. build/vet/tests verdes. |
| 7B | workfloo (editor, best-effort) | ✅ Hecho (2026-07-11) | sin commit (working tree) | 1 grid (decision). Canvas ya panea/zoomea en su viewport (no rompe). Wizards densos = desktop. |

Estados: ⬜ Pendiente · 🟡 En progreso · ✅ Hecho · ⏸️ Bloqueado

## Notas de arquitectura (contexto para agentes)

- Los 6 backends Go consumen el DS vía `replace ... => ../go-kiban-design-system` en su `go.mod` (dev local).
  Un cambio en el DS se ve inmediatamente en los consumidores sin publicar versión.
- Runtime local: Caddy en `:8080` rutea a cada backend (ver memoria "Local platform setup").
- Overlay/menu/drawer runtime ya existe en `view/layout/base.templ` (`kibanOpenOverlay`, `kibanToggleMenu`,
  backdrop, Escape). **Reusar, no reinventar.**
- El DS carga Tailwind por CDN con tokens kiban inline en `base.templ`. Breakpoints `md:`/`lg:` disponibles.
