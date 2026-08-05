// Package view holds the shared rendering helpers used across every kiban
// HTMX-based project. Components themselves live in feature-named subpackages
// (view/layout, view/input, view/button, …) so callers can import only what
// they use.
package view

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// HxRequestHeader is the request header htmx sets on every call it makes, and
// the one handlers branch on to decide full page vs fragment. Declared here
// (rather than imported from the htmx package) to keep view dependency-free.
const HxRequestHeader = "HX-Request"

// Render writes a templ component as the HTTP response body with the given
// status and Content-Type: text/html. Identical signature in every
// kiban-design-system consumer so handlers don't have to be aware of which
// project they're in.
//
// Every page in these apps answers the SAME url with two different bodies: the
// full page for a browser navigation, and just the inner fragment when htmx
// asks (handlers branch on htmx.IsRequest). Without `Vary: HX-Request` nothing
// tells the browser those are different resources, so a fragment cached under
// that url gets replayed as a whole document — the page comes back with no
// <head>, hence no styles and no scripts.
//
// It bites after a control that pushes the url (pagination, filters) and a
// later back-navigation: htmx fetched the fragment and pushed the url, so the
// cache now holds a fragment for it. Vary makes the cache key include the
// header, so a navigation without it misses and refetches the full page.
func Render(c *gin.Context, status int, component templ.Component) {
	c.Status(status)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Writer.Header().Add("Vary", HxRequestHeader)
	_ = component.Render(c.Request.Context(), c.Writer)
}
