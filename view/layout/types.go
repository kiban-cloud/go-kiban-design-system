// Package layout provides the shared HTML shell, topbar, sidebar (icon rail
// + subnav), and the slide animation between sidebar levels — every kiban
// htmx-based project (rekon, crm, …) wraps its pages with the same Layout
// component, passing per-project nav data via Config.
package layout

import "github.com/a-h/templ"

// User is the minimal identity displayed in the topbar's avatar dropdown.
type User struct {
	Email string
	Name  string
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

	// Routing
	ShellURL     string // kiban shell base URL (logo link + external tool prefix)
	LogoutAction string // POST URL that clears auth cookies
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
}
