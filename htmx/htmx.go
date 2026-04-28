// Package htmx holds the small set of helpers every kiban htmx-based
// project reaches for when handling HTMX requests at the controller layer.
//
// HTMX issues AJAX requests via HTML attributes (hx-post, hx-target, etc.)
// and signals its intent through a few request/response headers:
//
//   - Request side: `HX-Request: true` tells the server "this is HTMX,
//     return a partial". `HX-Trigger-Name` carries the `name` attribute
//     of the element that fired the request — useful when one form has
//     multiple submit buttons.
//   - Response side: `HX-Redirect` makes HTMX run a real client-side
//     navigation (window.location.href = url) instead of swapping the
//     redirect page's body into the target div.
//
// The helpers below centralize the header names and the
// HTMX-vs-browser branching so handlers don't open-code
// `c.GetHeader("HX-Request") == "true"` everywhere. They're intentionally
// thin — three functions, ~30 lines of code — but they prevent the most
// common HTMX bug in the codebase (calling `c.Redirect(302, url)` on an
// HTMX request and silently producing no navigation, because HTMX would
// just swap the redirect page's near-empty body).
package htmx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IsRequest reports whether the current request was issued by HTMX. Used
// in handlers to decide between rendering the full Page (with the layout
// shell) or just the inner Content partial that HTMX swaps into the
// target.
//
//	if htmx.IsRequest(c) {
//	    view.Render(c, http.StatusOK, view_page.Content(data))
//	    return
//	}
//	view.Render(c, http.StatusOK, view_page.Page(layout, data))
func IsRequest(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}

// Redirect issues a redirect that works for both HTMX and browser
// requests. HTMX requests get an HTTP 200 + `HX-Redirect: <url>` header
// — HTMX reads that header client-side and runs a real navigation.
// Non-HTMX requests get a regular HTTP 302 + Location header.
//
// Calling `c.Redirect(302, url)` directly on an HTMX request is a
// silent footgun: HTMX swaps the redirect page's body into the target
// div instead of navigating, producing a near-empty render. This helper
// makes the right thing the default.
//
// Status codes:
//   - HTMX path: HTTP 200 (HTMX expects any 2xx with `HX-Redirect`).
//   - Browser path: HTTP 302 (Found) — matches what every existing
//     kiban handler already uses for post-mutation redirects. 303 (See
//     Other) is more semantically correct for "POST then redirect to
//     GET" but the practical difference is negligible and 302 keeps the
//     migration zero-risk.
func Redirect(c *gin.Context, url string) {
	if IsRequest(c) {
		c.Header("HX-Redirect", url)
		c.Status(http.StatusOK)
		return
	}
	c.Redirect(http.StatusFound, url)
}

// TriggerName returns the `HX-Trigger-Name` request header, which HTMX
// sets to the `name` attribute of the element that fired the request.
// Empty string when the header is absent (non-HTMX request, or the
// triggering element had no `name`).
//
// Useful for forms with multiple submit buttons:
//
//	<button name="action" value="draft">Save draft</button>
//	<button name="action" value="publish">Publish</button>
//
// Both POST to the same URL; the handler reads `htmx.TriggerName(c)` and
// branches on `c.PostForm("action")`.
func TriggerName(c *gin.Context) string {
	return c.GetHeader("HX-Trigger-Name")
}
