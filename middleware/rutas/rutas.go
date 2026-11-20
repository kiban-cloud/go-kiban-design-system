// Package rutas deja que un módulo sirva sus páginas con el nombre nuevo
// (/solicitudes, /pagos, /fuentes, /cumplimiento) sin registrar cada ruta dos
// veces: envuelve el handler del servidor, traduce el prefijo nuevo al viejo
// antes del ruteo de gin y marca el request. middleware/authcookie usa la
// marca para llevar a cada quien al nombre de su diseño (flamboyant.Ruta).
package rutas

import (
	"net/http"

	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// AceptarNombresNuevos envuelve el handler del módulo (el engine de gin):
//
//	http.ListenAndServe(":"+port, rutas.AceptarNombresNuevos(engine.Handler()))
//
// Un request con el nombre nuevo llega a las rutas de siempre con el path
// traducido al viejo; cualquier otro pasa tal cual. Los dos quedan marcados
// en el contexto (flamboyant.AceptaNombresNuevos).
func AceptarNombresNuevos(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		viejo, porNuevo := flamboyant.ConNombreViejo(r.URL.Path)
		ctx := flamboyant.ConRutas(r.Context(), porNuevo)
		if !porNuevo {
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		r2 := r.Clone(ctx)
		r2.URL.Path = viejo
		if r2.URL.RawPath != "" {
			if raw, ok := flamboyant.ConNombreViejo(r2.URL.RawPath); ok {
				r2.URL.RawPath = raw
			} else {
				r2.URL.RawPath = ""
			}
		}
		h.ServeHTTP(w, r2)
	})
}
