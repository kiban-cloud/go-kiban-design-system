package accountblocked_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban/domain"
	domain_core_authorization_model "github.com/kiban-cloud/go-kiban/domain/authorization/model"

	"github.com/kiban-cloud/go-kiban-design-system/view/accountblocked"
)

func TestKindFor(t *testing.T) {
	cases := []struct {
		status string
		kind   accountblocked.Kind
		ok     bool
	}{
		{"Cuenta Suspendida", accountblocked.Suspended, true},
		{"Cuenta Inactiva", accountblocked.Inactive, true},
		{"Cuenta Activa", "", false},
		{"", "", false},
		// Same permissive reading as go-kiban's BillingAccountBlocked: only
		// the literals kiban-cloud writes cut the service.
		{"SUSPENDED", "", false},
	}
	for _, tc := range cases {
		kind, ok := accountblocked.KindFor(tc.status)
		assert.Equal(t, tc.ok, ok, tc.status)
		assert.Equal(t, tc.kind, kind, tc.status)
	}
}

// Los cortes que kiban-cloud decide y responde con 403 llegan como error de la
// sesión; cualquier otro error no es un corte.
func TestKindForError(t *testing.T) {
	kind, ok := accountblocked.KindForError(domain_core_authorization_model.ContractExpiredError())
	assert.True(t, ok)
	assert.Equal(t, accountblocked.ContractExpired, kind)

	kind, ok = accountblocked.KindForError(domain_core_authorization_model.AccessBlockedError(domain_core_authorization_model.ERROR_ACCOUNT_SUSPENDED))
	assert.True(t, ok)
	assert.Equal(t, accountblocked.Suspended, kind)

	_, ok = accountblocked.KindForError(domain.ErrorUnauthorized)
	assert.False(t, ok)
	_, ok = accountblocked.KindForError(nil)
	assert.False(t, ok)
}
