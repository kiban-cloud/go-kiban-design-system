// Package view holds the shared rendering helpers used across every kiban
// HTMX-based project. Components themselves live in feature-named subpackages
// (view/layout, view/input, view/button, …) so callers can import only what
// they use.
package view

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Render writes a templ component as the HTTP response body with the given
// status and Content-Type: text/html. Identical signature in every
// kiban-design-system consumer so handlers don't have to be aware of which
// project they're in.
func Render(c *gin.Context, status int, component templ.Component) {
	c.Status(status)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = component.Render(c.Request.Context(), c.Writer)
}
