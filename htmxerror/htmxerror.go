// Package htmxerror wires the shared HTMX error-response pattern for every
// kiban design-system consumer (rekon, crm, klin, workfloo, …) and re-exports
// the call-site surface so controllers never import go-kiban-fullstack
// directly.
//
// # Why
//
// A handler that renders an error into a 200 response is invisible to GCP
// Cloud Monitoring, which alerts on non-2xx status codes. This package makes
// error paths return the real 4xx/5xx status AND an HTMX fragment the client
// swaps into the page, so the user sees the error in context and the monitors
// see the failure. The original error is stashed in the gin context so the
// go-kiban-fullstack logger middleware records it (structured, with the HTMX
// headers) even on the rare partial-degradation path that still returns 200.
//
// # Boot
//
// Call Setup once, before serving, passing an optional ProjectMapper for the
// project's own error sentinels (rate limits, payment-required, external
// dependencies unique to that tool). Pass nil when the project only uses the
// canonical categories:
//
//	func main() {
//	    htmxerror.Setup(nil) // or htmxerror.Setup(projectMap)
//	    // …
//	}
//
// # Call sites
//
// Handlers map errors to responses with a single call. The status is derived
// from the error — no manual status codes:
//
//	user, err := ctrl.uc.Execute(ctx, id)
//	if err != nil {
//	    htmxerror.Respond(c, err, htmxerror.WithFormFallback(
//	        view_user.Form(view_user.FormData{Values: req.values(), Error: msg}),
//	    ))
//	    return
//	}
//
// # Wiring errors on the domain side
//
// Usecase sentinels should wrap a canonical sentinel with fmt.Errorf("%w: …")
// so errors.Is traverses to the canonical match and Respond can map it:
//
//	var ErrEmailTaken = fmt.Errorf("%w: el correo ya está registrado", htmxerror.ErrAlreadyExists)
package htmxerror

import (
	errors_domain "github.com/kiban-cloud/go-kiban-fullstack/pkg/domain/errors"
	fs_htmx "github.com/kiban-cloud/go-kiban-fullstack/pkg/infrastructure/http/htmx"

	"github.com/kiban-cloud/go-kiban-design-system/view/errormsg"

	"github.com/gin-gonic/gin"
)

// Renderable is any value with the templ.Component shape. Re-exported so
// consumers (and ProjectMapper implementations) don't import go-kiban-fullstack.
type Renderable = fs_htmx.Renderable

// ProjectMapper resolves errors specific to a project. It runs after the
// canonical mapping fails; return (0, nil) to fall through to the default 500
// banner. Re-exported so a consumer can declare its mapper against a DS type.
type ProjectMapper = fs_htmx.ProjectMapper

// Option mutates a single Respond call (see WithFormFallback).
type Option = fs_htmx.Option

// Canonical HTTP-aligned error sentinels, re-exported from go-kiban-fullstack.
// Usecases wrap one of these so Respond can map the error to a status; because
// these are the same underlying values, errors.Is matches across the re-export.
var (
	ErrInvalidInput               = errors_domain.ErrInvalidInput               // → 422
	ErrValidation                 = errors_domain.ErrValidation                 // → 422
	ErrNotFound                   = errors_domain.ErrNotFound                   // → 404
	ErrUnauthorized               = errors_domain.ErrUnauthorized               // → 401
	ErrForbidden                  = errors_domain.ErrForbidden                  // → 403
	ErrAlreadyExists              = errors_domain.ErrAlreadyExists              // → 409
	ErrConflict                   = errors_domain.ErrConflict                   // → 409
	ErrExternalServiceUnavailable = errors_domain.ErrExternalServiceUnavailable // → 502
)

// Setup wires the package-level Config in go-kiban-fullstack with the DS error
// fragments, the default text/html renderer, and the project's own mapper
// (nil for none). Call once at boot, before serving. Safe to call again to
// swap the mapper (e.g. in tests).
func Setup(projectMapper ProjectMapper) {
	fs_htmx.Default = fs_htmx.Config{
		Fragments: fs_htmx.CanonicalFragments{
			Validation:   func(msg string) Renderable { return errormsg.Validation(msg) },
			NotFound:     func(msg string) Renderable { return errormsg.NotFound(msg) },
			Unauthorized: func() Renderable { return errormsg.Unauthorized() },
			Forbidden:    func() Renderable { return errormsg.Forbidden() },
			Banner:       func(msg string) Renderable { return errormsg.Banner(msg) },
		},
		ProjectMapper: projectMapper,
		Render:        fs_htmx.DefaultRender,
	}
}

// Respond writes an HTMX error response: the status derived from err plus the
// mapped fragment. Use in error paths only. Delegates to the Config wired by
// Setup — calling it before Setup renders with nil fragments and will panic on
// the matching category, so Setup must run at boot.
func Respond(c *gin.Context, err error, opts ...Option) {
	fs_htmx.RespondHTMX(c, err, opts...)
}

// WithFormFallback re-renders form with the mapped 4xx status instead of the
// canonical fragment, so the handler can preserve form state (values, field
// errors). Honored only for 4xx — the client routes 4xx to the form and 5xx
// to the global banner, so a form in the banner would look broken.
func WithFormFallback(form Renderable) Option {
	return fs_htmx.WithFormFallback(form)
}

// TagError stashes err in the gin context so the logger middleware records it,
// without writing a response. Use in partial-degradation paths that still
// return 200 (a dashboard section that renders a degraded state instead of
// failing the whole page). No-op when err is nil.
func TagError(c *gin.Context, err error) {
	fs_htmx.TagError(c, err)
}

// StatusForError returns the HTTP status Respond would use for err, without
// rendering. Useful when a handler builds its own response but wants the
// mapping to stay consistent.
func StatusForError(err error) int {
	return fs_htmx.StatusForError(err)
}
