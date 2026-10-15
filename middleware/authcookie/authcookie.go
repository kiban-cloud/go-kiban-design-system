// Package authcookie authenticates htmx-style page requests from the pair
// of cookies the kiban shell sets at login:
//
//   - kiban_session : the JWT (same token the React shell keeps in
//     localStorage under "@kiban/user")
//   - kiban_space_id: the active space id (same value the React shell keeps
//     under "@kiban/spaceId")
//
// Validation is delegated to go-kiban's IAuthorizationAuthorizeWithSessionUseCase
// — the same use case the header-based SessionAuth middleware already uses
// — so auth semantics are identical between the JSON API and the htmx UI.
//
// Wiring: callers build one instance via `New(uc, loginURL)` (typically in
// their DI container) and apply `m.Middleware()` to the htmx route group.
// Each consumer is free to construct it inline (manual DI, like rekon) or
// wrap it in an inject-tagged adapter struct (reflection DI, like crm) —
// the shared code only needs the two values.
package authcookie

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiban-cloud/go-kiban-fullstack/logger"

	"github.com/kiban-cloud/go-kiban-design-system/view"
	"github.com/kiban-cloud/go-kiban-design-system/view/accountblocked"
	"github.com/kiban-cloud/go-kiban-design-system/view/ipblocked"
	controller_core_middleware "github.com/kiban-cloud/go-kiban/controller_core/middleware"
	controller_core_model "github.com/kiban-cloud/go-kiban/controller_core/model"
	domain_core_authorization_interface "github.com/kiban-cloud/go-kiban/domain/authorization/interface"
	domain_core_authorization_model "github.com/kiban-cloud/go-kiban/domain/authorization/model"
	infrastructure_core_env "github.com/kiban-cloud/go-kiban/infrastructure_core/env"
	utils_http "github.com/kiban-cloud/go-kiban/utils/http"
)

const (
	CookieSession = "kiban_session"
	CookieSpaceID = "kiban_space_id"

	// EnvLocalSession / EnvLocalSpaceID are dev-only env-var fallbacks for the
	// two auth cookies. When the request has no cookie AND the process is NOT
	// running on Cloud Run (!IsCloudRun) — i.e. a developer's local machine —
	// the middleware reads the session token and space id from these env vars
	// instead. This lets a developer paste a real kiban_session JWT (copied
	// from the browser after logging in once) into their local env and reach
	// the htmx UI without a browser cookie. The token still goes to kiban-cloud
	// through the real AuthorizeWithSessionUseCase, so it must be valid there
	// (it will expire like any session). On Cloud Run (dev/hotfix/prod) these
	// vars are ignored entirely.
	EnvLocalSession = "KIBAN_SESSION"
	EnvLocalSpaceID = "KIBAN_SPACE_ID"
)

// Middleware holds the dependencies and configuration the gin handler needs
// to validate the cookie pair on each request.
type Middleware struct {
	authorize domain_core_authorization_interface.IAuthorizationAuthorizeWithSessionUseCase
	loginURL  string
}

// New constructs the middleware. loginURL is where unauthenticated requests
// are redirected (usually `<shell url>/login`).
func New(uc domain_core_authorization_interface.IAuthorizationAuthorizeWithSessionUseCase, loginURL string) *Middleware {
	return &Middleware{authorize: uc, loginURL: loginURL}
}

// Middleware returns the gin handler. Place it on the htmx route group
// (or apply per-route inline for projects that don't have a dedicated
// htmx group).
func (m *Middleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(CookieSession)
		if err != nil || token == "" {
			token = localDevOverride(EnvLocalSession)
		}
		if token == "" {
			m.unauthorized(c)
			return
		}
		spaceID, err := c.Cookie(CookieSpaceID)
		if err != nil || spaceID == "" {
			spaceID = localDevOverride(EnvLocalSpaceID)
		}
		if spaceID == "" {
			m.unauthorized(c)
			return
		}

		ip := utils_http.ClientIP(c.Request)
		authorization, err := m.authorize.AuthorizeWithSessionUseCase(logger.TraceIDFromContext(c.Request.Context()), spaceID, token, false, ip)

		// Acceso por IP, lista Consola: es el firewall, lo primero. Lo decide
		// kiban-cloud al validar la sesión (con la IP del navegador y esta
		// herramienta), con la misma regla que la API key; aquí solo se
		// responde. Cambiar de espacio y salir quedan fuera: son la salida de
		// la pantalla de bloqueo.
		if errors.Is(err, domain_core_authorization_model.ErrorIpNotAllowed) {
			if !ipAccessExempt(c) {
				ipBlocked(c, infrastructure_core_env.Get(infrastructure_core_env.TOOL), ip)
				return
			}
			err = nil
		}
		if err != nil {
			m.unauthorized(c)
			return
		}
		authorization = authorization.WithCorrelation(c.Request.Context())

		// Cuenta de facturación suspendida o inactiva: bloqueo total de la
		// herramienta, sandbox incluido. A diferencia del gate de API de
		// go-kiban (AbortIfBillingAccountBlocked), acá sandbox no se exime:
		// es la misma consola con un toggle, y dejarla abierta en sandbox
		// sería mostrar la herramienta a medias a una cuenta cortada.
		if kind, blocked := accountblocked.KindFor(authorization.BillingAccountStatus); blocked {
			accountBlocked(c, kind)
			return
		}

		c.Set(controller_core_model.CONTEXT_KEY_AUTHORIZATION_OBJECT, authorization)
		c.Next()
	}
}

// localDevOverride returns the value of a cookie-fallback env var, but ONLY
// off Cloud Run. On Cloud Run (any deployed env: dev/hotfix/prod) it always
// returns "" so the env vars can never stand in for a real session cookie on a
// deployed instance — the guard is independent of whether the vars are set.
func localDevOverride(envKey string) string {
	if infrastructure_core_env.IsCloudRun() {
		return ""
	}
	return strings.TrimSpace(os.Getenv(envKey))
}

// unauthorized redirects to the kiban shell login. HTMX requests get an
// HX-Redirect header so the browser does a full-page navigation instead of
// trying to swap the login page into a partial.
func (m *Middleware) unauthorized(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", m.loginURL)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Redirect(http.StatusFound, m.loginURL)
	c.Abort()
}

// accountBlocked answers a request from a cut account with the shared
// full-page screen, always as 403 so monitors see it for what it is.
//
// An htmx request can't take the page as a swap — it would land inside
// whatever region asked (a lazy box, a table) — so it gets HX-Refresh
// instead: htmx reads that header before deciding whether to swap a 4xx, the
// browser reloads, and the full-page request renders the screen.
func accountBlocked(c *gin.Context, kind accountblocked.Kind) {
	c.Abort()
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Refresh", "true")
		c.Status(http.StatusForbidden)
		return
	}
	c.Request = c.Request.WithContext(accountblocked.WithRequestLanguage(c))
	view.Render(c, http.StatusForbidden, accountblocked.Page(kind))
}

// ipAccessExempt are the routes the IP gate never blocks: switching space (the
// user may be allowed in another one) and logging out.
func ipAccessExempt(c *gin.Context) bool {
	path := c.Request.URL.Path
	return strings.HasSuffix(path, "/spaces/switch") || strings.HasSuffix(path, "/logout")
}

// ipBlocked answers a request from an IP the Console list doesn't allow, with
// the shared full-page screen as 403 — same htmx handling as accountBlocked:
// HX-Refresh so the browser reloads and the full page shows the screen.
func ipBlocked(c *gin.Context, tool, ip string) {
	c.Abort()
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Refresh", "true")
		c.Status(http.StatusForbidden)
		return
	}
	c.Request = c.Request.WithContext(ipblocked.WithRequestLanguage(c))
	view.Render(c, http.StatusForbidden, ipblocked.Page(tool, ip))
}

// GetAuthorization is a thin re-export so htmx controllers don't have to
// import the go-kiban middleware package directly.
var GetAuthorization = controller_core_middleware.GetAuthorization
