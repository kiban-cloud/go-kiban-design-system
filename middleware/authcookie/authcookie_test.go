package authcookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fsi18n "github.com/kiban-cloud/go-kiban-fullstack/pkg/i18n"
	domain_core_authorization_model "github.com/kiban-cloud/go-kiban/domain/authorization/model"

	"github.com/kiban-cloud/go-kiban-design-system/middleware/authcookie"
)

type fakeAuthorize struct {
	status string
}

func (f fakeAuthorize) AuthorizeWithSessionUseCase(_, spaceID, token string, sandbox bool) (domain_core_authorization_model.Authorization, error) {
	return domain_core_authorization_model.Authorization{
		SpaceId:              spaceID,
		Token:                token,
		Sandbox:              sandbox,
		BillingAccountStatus: f.status,
	}, nil
}

type result struct {
	rec        *httptest.ResponseRecorder
	reachedApp bool
}

func serve(t *testing.T, status string, prepare func(*http.Request)) result {
	t.Helper()
	gin.SetMode(gin.TestMode)

	res := result{rec: httptest.NewRecorder()}
	engine := gin.New()
	m := authcookie.New(fakeAuthorize{status: status}, "/login")
	engine.GET("/tool", m.Middleware(), func(c *gin.Context) {
		res.reachedApp = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/tool", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSession, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: authcookie.CookieSpaceID, Value: "space"})
	if prepare != nil {
		prepare(req)
	}
	engine.ServeHTTP(res.rec, req)
	return res
}

func TestMiddleware_ActiveAccountReachesTheTool(t *testing.T) {
	res := serve(t, "Cuenta Activa", nil)

	assert.True(t, res.reachedApp)
	assert.Equal(t, http.StatusOK, res.rec.Code)
}

func TestMiddleware_EmptyStatusIsNotBlocked(t *testing.T) {
	// kiban-cloud deploys older than the billingAccountStatus field send
	// nothing; that must not lock every user out.
	res := serve(t, "", nil)

	assert.True(t, res.reachedApp)
}

func TestMiddleware_SuspendedAccountGetsTheFullPage(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, nil)

	require.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	body := res.rec.Body.String()
	assert.Contains(t, body, `data-account-blocked="suspended"`)
	assert.Contains(t, body, "Tu cuenta está suspendida")
	assert.Contains(t, body, `href="/kiban-cloud/billing"`)
	assert.Contains(t, body, "<!doctype html>")
}

func TestMiddleware_InactiveAccountPointsToSupport(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_INACTIVE, nil)

	require.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	body := res.rec.Body.String()
	assert.Contains(t, body, `data-account-blocked="inactive"`)
	assert.Contains(t, body, `href="mailto:help@kiban.com"`)
	assert.NotContains(t, body, `href="/kiban-cloud/billing"`)
}

func TestMiddleware_BlockedHtmxRequestRefreshesInsteadOfSwapping(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, func(r *http.Request) {
		r.Header.Set("HX-Request", "true")
	})

	require.False(t, res.reachedApp)
	assert.Equal(t, http.StatusForbidden, res.rec.Code)
	assert.Equal(t, "true", res.rec.Header().Get("HX-Refresh"))
	assert.Empty(t, res.rec.Body.String())
}

func TestMiddleware_BlockedPageFollowsTheLanguageCookie(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: fsi18n.CookieLang, Value: "en"})
	})

	assert.Contains(t, res.rec.Body.String(), "Your account is suspended")
}

func TestMiddleware_BlockedPageFallsBackToAcceptLanguage(t *testing.T) {
	res := serve(t, domain_core_authorization_model.BILLING_ACCOUNT_STATUS_SUSPENDED, func(r *http.Request) {
		r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	})

	assert.Contains(t, res.rec.Body.String(), "Your account is suspended")
}
