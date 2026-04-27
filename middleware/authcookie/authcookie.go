// Package authcookie authenticates htmx-style page requests from the pair
// of cookies the kiban shell sets at login:
//
//   - kiban_session : the JWT (same token the React shell keeps in
//                     localStorage under "@kiban/user")
//   - kiban_space_id: the active space id (same value the React shell keeps
//                     under "@kiban/spaceId")
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
	"net/http"

	controller_core_middleware "bitbucket.org/alexandregrin/go-kiban/controller_core/middleware"
	controller_core_model "bitbucket.org/alexandregrin/go-kiban/controller_core/model"
	domain_core_authorization_interface "bitbucket.org/alexandregrin/go-kiban/domain/authorization/interface"
	utils_http "bitbucket.org/alexandregrin/go-kiban/utils/http"
	"github.com/gin-gonic/gin"
)

const (
	CookieSession = "kiban_session"
	CookieSpaceID = "kiban_space_id"
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
			m.unauthorized(c)
			return
		}
		spaceID, err := c.Cookie(CookieSpaceID)
		if err != nil || spaceID == "" {
			m.unauthorized(c)
			return
		}

		authorization, err := m.authorize.AuthorizeWithSessionUseCase(spaceID, token, false)
		if err != nil {
			m.unauthorized(c)
			return
		}
		authorization.Ip = utils_http.GetClientIpv4(c)

		c.Set(controller_core_model.CONTEXT_KEY_AUTHORIZATION_OBJECT, authorization)
		c.Next()
	}
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

// GetAuthorization is a thin re-export so htmx controllers don't have to
// import the go-kiban middleware package directly.
var GetAuthorization = controller_core_middleware.GetAuthorization
