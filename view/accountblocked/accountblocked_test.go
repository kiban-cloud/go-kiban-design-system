package accountblocked_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
