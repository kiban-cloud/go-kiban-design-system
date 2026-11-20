package flamboyant

import (
	"context"
	"strings"
)

// Los módulos cambiaron de nombre (workfloo → Solicitudes, rekon → Pagos,
// link → Fuentes, klin → Cumplimiento). Las páginas flamboyant usan el nombre
// nuevo en la URL; las de siempre, el viejo. Cada módulo sirve los dos
// prefijos (middleware/rutas), middleware/authcookie lleva a cada quien al
// que corresponde a su diseño, y las vistas flamboyant arman sus enlaces con
// Ruta. crm y kiban-cloud no cambian.

// Renombre es un prefijo de módulo con su nombre viejo y el nuevo.
type Renombre struct {
	Viejo string
	Nuevo string
}

// Renombres son los prefijos que cambian de nombre con flamboyant.
var Renombres = []Renombre{
	{Viejo: "/workfloo", Nuevo: "/solicitudes"},
	{Viejo: "/rekon", Nuevo: "/pagos"},
	{Viejo: "/link", Nuevo: "/fuentes"},
	{Viejo: "/klin-internal", Nuevo: "/cumplimiento-interno"},
	{Viejo: "/klin", Nuevo: "/cumplimiento"},
}

// cambiaPrefijo cambia de por a al inicio de path si de es un segmento
// completo: «/klin» no cambia «/klin-internal» ni «/klinx».
func cambiaPrefijo(path, de, a string) (string, bool) {
	if !strings.HasPrefix(path, de) {
		return path, false
	}
	rest := path[len(de):]
	if rest != "" && !strings.ContainsRune("/?#", rune(rest[0])) {
		return path, false
	}
	return a + rest, true
}

// ConNombreNuevo traduce un path con el nombre viejo del módulo al nuevo. ok
// es false si el path no es de un módulo que cambió de nombre.
func ConNombreNuevo(path string) (string, bool) {
	for _, r := range Renombres {
		if p, ok := cambiaPrefijo(path, r.Viejo, r.Nuevo); ok {
			return p, true
		}
	}
	return path, false
}

// ConNombreViejo traduce un path con el nombre nuevo del módulo al viejo. ok
// es false si el path no lleva un nombre nuevo.
func ConNombreViejo(path string) (string, bool) {
	for _, r := range Renombres {
		if p, ok := cambiaPrefijo(path, r.Nuevo, r.Viejo); ok {
			return p, true
		}
	}
	return path, false
}

// Ruta es el path de un enlace con el nombre del módulo que corresponde al
// diseño del request: con flamboyant, el nuevo; si no, el path tal cual. Las
// vistas flamboyant arman sus enlaces y redirecciones con ella; un enlace que
// se escape igual funciona, porque authcookie lleva el nombre viejo al nuevo.
func Ruta(ctx context.Context, path string) string {
	if !On(ctx) {
		return path
	}
	p, _ := ConNombreNuevo(path)
	return p
}

type rutasKey struct{}

type rutasInfo struct {
	porNombreNuevo bool
}

// ConRutas marca el contexto de un request que pasó por middleware/rutas: el
// módulo sirve los dos nombres, y porNombreNuevo dice si el request llegó con
// el nuevo (el path ya va traducido al viejo).
func ConRutas(ctx context.Context, porNombreNuevo bool) context.Context {
	return context.WithValue(ctx, rutasKey{}, rutasInfo{porNombreNuevo: porNombreNuevo})
}

// AceptaNombresNuevos dice si el módulo sirve los dos nombres (pasó por
// middleware/rutas). Sin eso authcookie no redirige: mandaría al nombre
// nuevo a un módulo que todavía no lo entiende.
func AceptaNombresNuevos(ctx context.Context) bool {
	_, ok := ctx.Value(rutasKey{}).(rutasInfo)
	return ok
}

// PorNombreNuevo dice si el request llegó con el nombre nuevo del módulo.
func PorNombreNuevo(ctx context.Context) bool {
	info, _ := ctx.Value(rutasKey{}).(rutasInfo)
	return info.porNombreNuevo
}
