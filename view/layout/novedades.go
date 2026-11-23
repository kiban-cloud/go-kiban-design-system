package layout

import (
	"context"

	"github.com/kiban-cloud/go-kiban-design-system/view/feature"
)

// Novedades de kiban en el shell de TODAS las herramientas: «Novedades» en el
// menú del avatar y el modal de novedades al entrar. Los dos los sirve kiban
// cloud (mismo dominio, como la campana de notificaciones): el changelog en
// NovedadesPath y el modal en NovedadesModalPath, que responde vacío si al
// usuario no le toca (ya lo cerró en esta sesión o lo descartó).
//
// FlagNovedades la prende el middleware de cada herramienta cuando la sesión
// trae el feature flag NOVEDADES de la compañía (authorization.HasFeature),
// como SHELL_FLAMBOYANT; en kiban-proto, la tabla de fases.
const FlagNovedades = "novedades"

const (
	NovedadesPath      = "/kiban-cloud/novedades"
	NovedadesModalPath = "/kiban-cloud/novedades/modal"
	// NovedadesActiveKey es el ActiveKey del changelog: ahí no se pide el
	// modal (ya se están viendo todas las novedades).
	NovedadesActiveKey = "novedades"
)

func novedadesOn(ctx context.Context) bool { return feature.Enabled(ctx, FlagNovedades) }

// conNovedades agrega «Novedades» al final del menú del avatar si la bandera
// está prendida y la herramienta no la trae ya.
func conNovedades(ctx context.Context, items []UserMenuItem) []UserMenuItem {
	if !novedadesOn(ctx) {
		return items
	}
	for _, it := range items {
		if it.Href == NovedadesPath {
			return items
		}
	}
	out := make([]UserMenuItem, 0, len(items)+1)
	out = append(out, items...)
	return append(out, UserMenuItem{Label: docsCatalog.T(ctx, "shell.novedades"), Href: NovedadesPath})
}

// novedadesHerramienta es la clave de herramienta (HERRAMIENTA_* de los
// feature flags) del ToolKey del shell: el modal pone primero las entradas de
// la herramienta donde está el usuario.
func novedadesHerramienta(toolKey string) string {
	switch toolKey {
	case "customers":
		return "CRM"
	case "workfloo":
		return "WORKFLOO"
	case "rekon":
		return "REKON"
	case "link":
		return "LINK"
	case "klin":
		return "KLIN"
	case "folo", "monitoreo":
		return "MONITOREO"
	}
	return "KIBAN_CLOUD"
}

// novedadesModalURL es lo que pide el cargador; "" si no se pide (sin bandera,
// o en el propio changelog).
func novedadesModalURL(ctx context.Context, cfg Config) string {
	if !novedadesOn(ctx) || cfg.ActiveKey == NovedadesActiveKey {
		return ""
	}
	return NovedadesModalPath + "?desde=" + novedadesHerramienta(cfg.ToolKey)
}
