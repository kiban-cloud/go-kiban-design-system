package layout_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiban-cloud/go-kiban-design-system/view/layout"
)

func TestNovedades_SinBanderaNoHayMenuNiModal(t *testing.T) {
	for _, nombre := range []string{"flamboyant", "clasico"} {
		ctx := flags()
		if nombre == "flamboyant" {
			ctx = flags(layout.FlagShellFlamboyant)
		}
		body := renderLayout(t, ctx, flamboyantCfg())
		assert.NotContains(t, body, layout.NovedadesModalPath, nombre)
		assert.NotContains(t, body, `href="`+layout.NovedadesPath+`"`, nombre)
	}
}

func TestNovedades_EnLosDosShellsMenuYModalConSuHerramienta(t *testing.T) {
	for _, nombre := range []string{"flamboyant", "clasico"} {
		ctx := flags(layout.FlagNovedades)
		cfg := flamboyantCfg()
		if nombre == "flamboyant" {
			ctx = flags(layout.FlagShellFlamboyant, layout.FlagNovedades)
		} else {
			cfg.UserMenuItems = []layout.UserMenuItem{{Label: "Perfil", Href: "/kiban-cloud/profile"}}
		}
		body := renderLayout(t, ctx, cfg)
		assert.Contains(t, body, `hx-get="/kiban-cloud/novedades/modal?desde=CRM"`, nombre+": el modal sabe que está en CRM")
		assert.Contains(t, body, `href="`+layout.NovedadesPath+`"`, nombre+": Novedades en el menú del avatar")
	}
}

func TestNovedades_EnElChangelogNoSePideElModal(t *testing.T) {
	cfg := flamboyantCfg()
	cfg.ActiveKey = layout.NovedadesActiveKey
	body := renderLayout(t, flags(layout.FlagShellFlamboyant, layout.FlagNovedades), cfg)
	assert.NotContains(t, body, layout.NovedadesModalPath)
}

func TestNovedades_NoDuplicaElRenglonSiLaHerramientaYaLoTrae(t *testing.T) {
	cfg := flamboyantCfg()
	cfg.UserMenuItems = []layout.UserMenuItem{{Label: "Novedades", Href: layout.NovedadesPath}}
	body := renderLayout(t, flags(layout.FlagShellFlamboyant, layout.FlagNovedades), cfg)
	assert.Equal(t, 1, strings.Count(body, `href="`+layout.NovedadesPath+`"`))
}

func TestNovedades_HerramientaPorToolKey(t *testing.T) {
	casos := map[string]string{
		"home": "KIBAN_CLOUD", "config": "KIBAN_CLOUD", "customers": "CRM", "workfloo": "WORKFLOO",
		"rekon": "REKON", "link": "LINK", "klin": "KLIN", "folo": "MONITOREO", "monitoreo": "MONITOREO",
	}
	for toolKey, want := range casos {
		cfg := flamboyantCfg()
		cfg.ToolKey = toolKey
		body := renderLayout(t, flags(layout.FlagShellFlamboyant, layout.FlagNovedades), cfg)
		assert.Contains(t, body, "/kiban-cloud/novedades/modal?desde="+want, toolKey)
	}
}
