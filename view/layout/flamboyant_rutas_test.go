package layout_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/icons"
	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

// workflooCfg es el shell de workfloo con enlaces a los demás módulos, como
// lo arma cada herramienta con su lista Tools de siempre.
func workflooCfg() layout.Config {
	return layout.Config{
		Title:             "Flujos",
		ProjectName:       "workfloo",
		SectionLabel:      "workfloo",
		User:              layout.User{Name: "Ana Pérez", Email: "ana@kiban.com"},
		SpaceID:           "space-123456",
		LogoutAction:      "/workfloo/auth/logout",
		SwitchSpaceAction: "/workfloo/spaces/switch",
		Spaces:            []layout.SpaceOption{{Id: "space-123456", Name: "Uno"}, {Id: "space-2", Name: "Dos"}},
		ToolKey:           "workfloo",
		ActiveKey:         "workfloos",
		Tools: []layout.Tool{
			{Key: "customers", Label: "CRM", Href: "/crm", Icon: icons.Users},
			{Key: "workfloo", Label: "workfloo", Href: "/workfloo", Icon: icons.GitMerge},
			{Key: "rekon", Label: "rekon", Href: "/rekon", Icon: icons.CreditCard},
			{Key: "link", Label: "link", Href: "/link", Icon: icons.Code},
			{Key: "klin", Label: "klin", Href: "/klin", Icon: icons.Database},
		},
		SubItems:  []layout.SubItem{{Key: "workfloos", Label: "Flujos", Href: "/workfloo/workfloos"}},
		SubAction: &layout.SubAction{Label: "Nueva ejecución", Href: "/workfloo/execute"},
	}
}

func TestFlamboyant_ElMenuUsaLosNombresNuevos(t *testing.T) {
	body := shell(t, renderLayout(t, flags(layout.FlagShellFlamboyant), workflooCfg()))

	for _, href := range []string{
		`href="/solicitudes"`, `href="/pagos"`, `href="/fuentes"`, `href="/cumplimiento"`,
		`href="/solicitudes/workfloos"`, `href="/solicitudes/execute"`,
		`action="/solicitudes/auth/logout"`, `action="/solicitudes/spaces/switch"`,
		`href="/crm"`,
	} {
		assert.Contains(t, body, href)
	}
	for _, viejo := range []string{`href="/workfloo"`, `href="/rekon"`, `href="/link"`, `href="/klin"`, `"/workfloo/`} {
		assert.NotContains(t, body, viejo)
	}
}

func TestLayout_ElShellDeSiempreConservaLosNombresViejos(t *testing.T) {
	body := renderLayout(t, context.Background(), workflooCfg())

	assert.Contains(t, body, `href="/workfloo/workfloos"`)
	assert.Contains(t, body, `href="/rekon"`)
	assert.NotContains(t, body, `/solicitudes`)
	assert.NotContains(t, body, `/pagos`)
}
