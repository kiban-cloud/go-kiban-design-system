package htmxerror_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/kiban-cloud/go-kiban-fullstack/pkg/infrastructure/http/htmx/htmxtest"

	"github.com/kiban-cloud/go-kiban-design-system/htmxerror"
	"github.com/kiban-cloud/go-kiban-design-system/view/errormsg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A project-specific sentinel + mapper, mirroring what a consumer (rekon, crm,
// …) wires: an error the canonical categories don't cover, mapped to a 5xx via
// the DS Banner fragment.
var errDependencyDown = errors.New("payments provider unavailable")

func projectMap(err error) (int, htmxerror.Renderable) {
	if errors.Is(err, errDependencyDown) {
		return http.StatusServiceUnavailable, errormsg.Banner("Servicio no disponible. Intenta más tarde.")
	}
	return 0, nil
}

func TestMain(m *testing.M) {
	htmxerror.Setup(projectMap)
	m.Run()
}

// TestCanonicalMapping verifies each canonical sentinel maps to the right
// status AND renders the DS fragment carrying the matching data-error-type —
// the two halves of "GCP sees a non-200, the user sees the error in context".
func TestCanonicalMapping(t *testing.T) {
	cases := []struct {
		name          string
		err           error
		wantStatus    int
		wantErrorType string
	}{
		{"invalid input → 422", htmxerror.ErrInvalidInput, http.StatusUnprocessableEntity, "validation"},
		{"validation → 422", htmxerror.ErrValidation, http.StatusUnprocessableEntity, "validation"},
		{"not found → 404", htmxerror.ErrNotFound, http.StatusNotFound, "not-found"},
		{"unauthorized → 401", htmxerror.ErrUnauthorized, http.StatusUnauthorized, "unauthorized"},
		{"forbidden → 403", htmxerror.ErrForbidden, http.StatusForbidden, "forbidden"},
		{"already exists → 409", htmxerror.ErrAlreadyExists, http.StatusConflict, "validation"},
		{"conflict → 409", htmxerror.ErrConflict, http.StatusConflict, "validation"},
		{"external service → 502", htmxerror.ErrExternalServiceUnavailable, http.StatusBadGateway, "system"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := htmxtest.NewCtx(t)
			htmxerror.Respond(c, tc.err)
			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), `data-error-type="`+tc.wantErrorType+`"`)
			assert.Contains(t, rec.Body.String(), `role="alert"`)
		})
	}
}

// Usecase sentinels that wrap a canonical sentinel must still map correctly —
// errors.Is has to traverse the chain through the DS re-exported var.
func TestWrappedSentinelMaps(t *testing.T) {
	errEmailTaken := fmt.Errorf("%w: el correo ya está registrado", htmxerror.ErrAlreadyExists)

	c, rec := htmxtest.NewCtx(t)
	htmxerror.Respond(c, errEmailTaken)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

// The project mapper resolves errors outside the canonical set.
func TestProjectMapper(t *testing.T) {
	c, rec := htmxtest.NewCtx(t)
	htmxerror.Respond(c, errDependencyDown)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "Servicio no disponible")
}

// An unmapped error falls through to the 500 banner — still a real non-200.
func TestUnmappedFallsThroughTo500(t *testing.T) {
	c, rec := htmxtest.NewCtx(t)
	htmxerror.Respond(c, errors.New("something nobody mapped"))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `data-error-type="system"`)
}

// WithFormFallback re-renders the form on a 4xx (so field values survive)
// instead of the canonical fragment.
func TestFormFallbackOn4xx(t *testing.T) {
	c, rec := htmxtest.NewCtx(t)
	htmxerror.Respond(c, htmxerror.ErrInvalidInput, htmxerror.WithFormFallback(htmxtest.Marker("MY-FORM")))
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "MY-FORM")
}

// The fallback is ignored on 5xx — those route to the global banner, where a
// form would look broken. The canonical banner renders instead.
func TestFormFallbackIgnoredOn5xx(t *testing.T) {
	c, rec := htmxtest.NewCtx(t)
	htmxerror.Respond(c, htmxerror.ErrExternalServiceUnavailable, htmxerror.WithFormFallback(htmxtest.Marker("MY-FORM")))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
	assert.NotContains(t, rec.Body.String(), "MY-FORM")
	assert.Contains(t, rec.Body.String(), `data-error-type="system"`)
}

// StatusForError exposes the mapping without rendering.
func TestStatusForError(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, htmxerror.StatusForError(htmxerror.ErrNotFound))
	assert.Equal(t, http.StatusInternalServerError, htmxerror.StatusForError(errors.New("x")))
}

// TagError stashes the error for the logger middleware without writing a
// response body (used on partial-degradation 200 paths).
func TestTagErrorWritesNoBody(t *testing.T) {
	c, rec := htmxtest.NewCtx(t)
	htmxerror.TagError(c, errors.New("degraded section"))
	assert.Empty(t, rec.Body.String())

	v, ok := c.Get("ERROR_MESSAGE")
	require.True(t, ok)
	assert.Error(t, v.(error))
}
