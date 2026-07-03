// Package layout provides the shared HTML shell, topbar, sidebar (icon rail
// + subnav), and the slide animation between sidebar levels — every kiban
// htmx-based project (rekon, crm, …) wraps its pages with the same Layout
// component, passing per-project nav data via Config.
package layout

import (
	"strings"

	"github.com/a-h/templ"

	"github.com/kiban-cloud/go-kiban-design-system/view/breadcrumbs"
)

// DefaultLogoHref is the destination of the topbar kiban logo in BOTH shells
// (Topbar and AdminLayout). It is a global, project-agnostic link: every
// kiban tool's logo points at the same kiban-cloud landing, so no project
// has to wire the logo URL. It is intentionally host-relative so it resolves
// against whatever host served the page — the same link works in
// local/dev/qa/alpha/prod without per-environment configuration. The logo
// always uses this value and ignores Config.ShellURL / AdminConfig.HomeHref
// (those remain for other purposes — see their field docs).
const DefaultLogoHref = "/kiban-cloud"

// logoHref resolves the AdminLayout topbar's ProjectName label link (the app
// name shown next to the logo): the project-supplied HomeHref when set,
// otherwise DefaultLogoHref. The logo image itself does not use this — it is
// always DefaultLogoHref.
func logoHref(explicit string) string {
	if explicit != "" {
		return explicit
	}
	return DefaultLogoHref
}

// DefaultNotificationsBaseURL is the base path the topbar notifications bell
// uses when a Config doesn't set NotificationsBaseURL explicitly. Like
// DefaultLogoHref, it is a global, project-agnostic, host-relative link:
// every kiban tool shares one origin, and the notification endpoints
// (dropdown/badge/page) are served by the kiban-cloud backend under this
// path. Defaulting it here means every tool that wraps its pages with
// Layout gets the bell for free — no per-project wiring, just a DS bump.
const DefaultNotificationsBaseURL = "/kiban-cloud/notifications"

// notificationsBaseURL resolves the bell's base path: the project-supplied
// NotificationsBaseURL when set, otherwise DefaultNotificationsBaseURL. The
// bell is always rendered in the customer-facing Topbar — there is no
// opt-out, since every kiban tool surfaces the same shared notifications.
func notificationsBaseURL(explicit string) string {
	if explicit != "" {
		return explicit
	}
	return DefaultNotificationsBaseURL
}

// User is the minimal identity displayed in the topbar's avatar dropdown.
//
// Picture is optional. The customer-facing Layout doesn't render it
// today (its topbar shows initials only); AdminLayout shows the photo
// when present and falls back to initials when empty. Adding the
// field here keeps a single User source-of-truth across both shells.
type User struct {
	Email   string
	Name    string
	Picture string
}

// Tool is one entry in the icon rail (sidebar level 1). External=true → links
// out to the kiban shell; External=false → links to a route inside the
// hosted project (the project's own self-tool).
type Tool struct {
	Key      string             // matches Config.ToolKey to highlight the active tool
	Icon     func() templ.Component
	Href     string
	External bool
	Label    string             // tooltip text
}

// SubItem is one entry in the sub-nav (sidebar level 2) — items live inside
// the currently-hosted project's tool.
type SubItem struct {
	Key   string                // matches Config.ActiveKey to highlight the active item
	Label string
	Href  string
}

// SandboxToggle drives an optional switch in the topbar — every kiban tool
// has a prod/sandbox concept, so the control lives in the shared shell so
// consumers don't have to re-implement it per project.
//
// Behavior is a link, not a form submit: clicking flips the page to
// ToggleHref, which the project computes (typically the current URL with
// ?sandbox flipped). The shell never persists state — the caller owns it.
// Label defaults to "Modo sandbox" when empty.
type SandboxToggle struct {
	Enabled    bool
	ToggleHref string
	Label      string
}

// LabelOrDefault returns the toggle's display label, defaulting to "Modo
// sandbox" when the caller didn't set one — keeps the contract terse for
// the common case while still allowing overrides for localized or
// project-specific wording.
func (s SandboxToggle) LabelOrDefault() string {
	if strings.TrimSpace(s.Label) != "" {
		return s.Label
	}
	return "Modo sandbox"
}

// SpaceOption is one entry in the topbar's space switcher dropdown. The
// switcher only renders as a clickable dropdown when len(Config.Spaces) >= 2;
// with 0 or 1 spaces the chip falls back to a non-interactive label showing
// the current space's name (or a "—" sentinel when no space is active yet).
type SpaceOption struct {
	Id   string // mongo ObjectId — submitted to SwitchSpaceAction in the spaceId form field
	Name string // user-visible display name (the shell formerly showed the trailing 8 chars of Id)
}

// UserMenuItem is one row in the topbar avatar dropdown. The DS
// renders these as plain anchor links between the user-info header
// and the logout form. Permission checks (e.g. "only owners see
// Usuarios") are the consuming project's responsibility — the slice
// it builds is what the menu shows verbatim.
type UserMenuItem struct {
	Label string // user-visible label, e.g. "Perfil"
	Href  string // destination path or full URL
}

// Config carries everything the Layout templ needs. Each consuming project
// builds this once per render (typically in a project-specific PageData
// helper) and hands it to Layout.
//
// The split between project-specific and shared concerns:
//   - Per project: ProjectName, SectionLabel, LogoutAction, Tools, Docs,
//     SubItems — these define the nav data and identity of the hosting app.
//   - Per page: Title, ToolKey, ActiveKey — what's currently active.
//   - Per request: User, SpaceID, ShellURL — pulled from the auth context.
type Config struct {
	// Page metadata
	Title string

	// Project identity. Used for:
	//   - <title> suffix: "<Title> · <ProjectName>"
	//   - localStorage namespace for sidebar level persistence
	//     (key: "<ProjectName>-sidebar-level")
	ProjectName string
	// SectionLabel shown above the sub-nav items (e.g. "rekon", "crm").
	SectionLabel string

	// User context
	User    User
	SpaceID string
	// CurrentSpaceName is the human-readable name of SpaceID. Used by the
	// topbar's space chip so the user sees "Mi espacio" instead of the last
	// 8 chars of the ObjectId. Falls back to a "—" sentinel when empty.
	CurrentSpaceName string
	// Spaces is the list shown in the topbar's space-switcher dropdown.
	// When len(Spaces) < 2 the chip is non-interactive (no menu, no
	// chevron). When >= 2 the chip becomes a clickable dropdown whose
	// items each POST to SwitchSpaceAction with a hidden spaceId field.
	Spaces []SpaceOption
	// SwitchSpaceAction is the POST URL the space-switcher submits to.
	// Owned by the hosting project (kibancloud sets it to
	// "/cloud-htmx/spaces/switch"). Empty disables the dropdown even if
	// Spaces has items — defensive default so a misconfigured shell
	// can't issue submits to "".
	SwitchSpaceAction string

	// Routing
	ShellURL     string // kiban shell base URL — external tool prefix (icon rail builds links as ShellURL + "/crm/customers", etc.). NOT used by the logo (the logo is always DefaultLogoHref).
	LogoutAction string // POST URL that clears auth cookies
	// UserMenuItems are the links rendered inside the topbar avatar
	// dropdown, between the user-info header and the "Salir" button.
	// Each project decides what goes there (kibancloud surfaces
	// Perfil / Facturación / Usuarios; other tools may surface
	// different shortcuts or none).
	//
	// Order is preserved. Empty slice / nil → menu renders just the
	// header + logout. Pages that need permission gating do the
	// filtering before building the slice; the DS itself just
	// renders.
	UserMenuItems []UserMenuItem
	// DevelopersURL is the topbar "Developers" link target. The button is
	// opt-in: leave this empty to hide it entirely. Set it (typically to an
	// in-project route or to the kiban shell's developers hub) to render the
	// button pointing there.
	DevelopersURL string

	// Active state
	ToolKey   string // matches one of Tools[i].Key — highlights the active tool icon
	ActiveKey string // matches one of SubItems[i].Key — highlights the active sub-nav item

	// Nav data
	Tools    []Tool // top of the icon rail — kiban tools
	Docs     []Tool // bottom of the icon rail — docs links
	SubItems []SubItem

	// SandboxToggle, when non-nil, renders a switch in the topbar (next to
	// the user menu). Leave nil to omit the control entirely on pages that
	// don't have a sandbox dimension.
	SandboxToggle *SandboxToggle

	// Optional CDN-pinned client libraries. Each flag emits the matching
	// <script> / <link> in base.templ's <head>. Consumers opt in per page
	// (or per layout helper) so projects that don't need a library don't
	// pay the byte cost. Versions are pinned in base.templ — see the
	// "Externals" section of the design-system CLAUDE.md for the
	// current set.
	LoadChartJS    bool // Chart.js — line/bar/donut charts (view/chart)
	LoadCytoscape  bool // Cytoscape.js + cytoscape-dagre — flow graphs (view/flow_graph)
	LoadCodeMirror bool // CodeMirror 6 (JS mode) — code editor (view/code_editor)
	LoadSortable   bool // SortableJS — drag-to-reorder lists (view/sortable_list)
	LoadFlatpickr  bool // flatpickr + es locale — date range picker (view/date_range)
	LoadMarked     bool // marked.js — markdown rendering

	// NotificationsBaseURL overrides the base path the topbar
	// notifications bell uses. Leave it empty in normal use: the bell is
	// always rendered and falls back to DefaultNotificationsBaseURL
	// ("/kiban-cloud/notifications"), so every tool gets it with no
	// wiring. Set it only to point the bell at a non-default path. The
	// bell expects three routes under the resolved base:
	//   - GET  <base>/dropdown → the panel list fragment (lazy-loaded on open)
	//   - GET  <base>/badge    → an OOB <span id="kiban-notif-badge"> (polled)
	//   - GET  <base>          → the full notifications page (the "ver todas" link)
	NotificationsBaseURL string
	// NotificationsUnread is an optional seed for the bell badge on first
	// paint, avoiding a flash before the count loads. It is purely a
	// nicety: the badge poller fires on page load (not just every 30s),
	// so the count self-corrects over HTTP from <base>/badge immediately
	// — tools that can't compute the count locally (everything except
	// kiban-cloud, which owns the count usecase) just leave this zero.
	NotificationsUnread int

	// CustomStylesheetHref, when non-empty, appends a <link rel="stylesheet">
	// as the LAST entry in <head> so it overrides the design-system tokens
	// and component styles. Purpose: per-client theming for customer-facing
	// shells (e.g. workfloo's portal cliente loads the tenant's uploaded
	// CSS). It's a plain href — the consumer owns serving the stylesheet
	// (typically a route that streams the stored CSS). Empty = no extra
	// stylesheet, so tools that don't theme pay nothing.
	CustomStylesheetHref string
}

// NavSection is one entry in AdminLayout's horizontal top-level nav.
// Section nav is a single-level row of links (no icon rail, no
// sub-nav) — admin/internal apps with 2-5 sections fit here cleanly.
type NavSection struct {
	Key   string // matches AdminConfig.ActiveSection to highlight
	Label string
	Href  string
}

// AdminConfig drives AdminLayout — a minimal admin/internal shell:
// topbar with logo + horizontal nav + user menu, optional breadcrumbs
// strip below, and the page content. No icon rail, no sub-nav,
// no multi-tool switcher. Suitable for staff/internal apps that
// aren't part of the multi-tool kiban shell.
//
// Reuses Base() under the hood for the HTML chrome, scripts, and
// runtime helpers (overlay, menu, nav-loader, …).
type AdminConfig struct {
	// Page metadata
	Title string

	// Identity. Used for the <title> suffix and the topbar logo link.
	ProjectName string
	HomeHref    string // where the topbar ProjectName label links to (the app's own home, e.g. /klin-internal-htmx/). Empty → falls back to DefaultLogoHref. The logo itself always links to DefaultLogoHref regardless.

	// User context (avatar + dropdown in the topbar)
	User User

	// LogoutAction is the POST URL the "Cerrar sesión" item submits to
	// (the user-menu wraps the item in a <form method="POST">). Leave
	// empty to omit the logout entry from the menu entirely.
	LogoutAction string

	// Top-level horizontal nav rendered in the topbar.
	NavSections []NavSection
	// ActiveSection highlights the matching NavSection.Key. Empty =
	// nothing highlighted (e.g. on the login page or a 404).
	ActiveSection string

	// Breadcrumbs is the trail rendered in a strip below the topbar.
	// Empty/nil collapses the strip entirely (zero height).
	Breadcrumbs []breadcrumbs.Item
}

// adminToBaseConfig narrows an AdminConfig down to the few Config
// fields Base actually uses (title chrome, project namespacing for
// sidebar localStorage). The customer-facing fields (Tools, SubItems,
// etc.) stay zero-valued — Base.templ tolerates them.
func adminToBaseConfig(cfg AdminConfig) Config {
	return Config{
		Title:       cfg.Title,
		ProjectName: cfg.ProjectName,
	}
}

// adminNavSectionClass returns the class set for a horizontal nav
// item, branching on whether the item's Key matches the active one.
// Active gets the kiban-primary tint; inactive renders muted ink with
// hover affordance.
func adminNavSectionClass(itemKey, activeKey string) string {
	base := "px-3 py-1.5 rounded-md text-sm transition-colors"
	if itemKey != "" && itemKey == activeKey {
		return base + " bg-kiban-primary-soft text-kiban-primary font-medium"
	}
	return base + " text-kiban-ink3 hover:text-kiban-ink hover:bg-kiban-surface"
}

// adminUserMenuLabel returns the visible name in the topbar's user
// menu trigger. Falls back to email when name is empty (some staff
// users may register with email-only).
func adminUserMenuLabel(u User) string {
	if strings.TrimSpace(u.Name) != "" {
		return u.Name
	}
	return u.Email
}

