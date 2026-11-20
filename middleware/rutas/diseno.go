package rutas

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// Aplicar va en lugar de c.Next() al final de un middleware de sesión, una
// vez que las banderas del diseño están en el contexto: lo usan
// middleware/authcookie y el de staff de klin. En un módulo que pasa por
// AceptarNombresNuevos, manda con un 308 al nombre del módulo que
// corresponde al diseño y, con flamboyant, lleva al nombre nuevo los
// enlaces y redirects de la respuesta (writer.go). En cualquier otro caso
// es c.Next().
func Aplicar(c *gin.Context) {
	if destino, ok := destinoDelDiseno(c); ok {
		c.Redirect(http.StatusPermanentRedirect, destino)
		c.Abort()
		return
	}
	ctx := c.Request.Context()
	if !flamboyant.AceptaNombresNuevos(ctx) || !flamboyant.On(ctx) {
		c.Next()
		return
	}
	w := &writer{ResponseWriter: c.Writer}
	c.Writer = w
	c.Next()
	w.terminar()
}

// destinoDelDiseno es a dónde se manda un request cuyo nombre de módulo no
// corresponde a su diseño: flamboyant con el nombre viejo (un enlace de una
// notificación, un marcador, una pantalla de siempre) va al nuevo, y sin
// flamboyant con el nombre nuevo (un enlace compartido) vuelve al viejo. El
// 308 conserva el método y el cuerpo, y htmx lo sigue solo. Sólo en módulos
// que sirven los dos nombres (middleware/rutas).
func destinoDelDiseno(c *gin.Context) (string, bool) {
	ctx := c.Request.Context()
	if !flamboyant.AceptaNombresNuevos(ctx) {
		return "", false
	}
	path := c.Request.URL.Path
	on := flamboyant.On(ctx)
	porNuevo := flamboyant.PorNombreNuevo(ctx)
	switch {
	case on && !porNuevo:
		nuevo, ok := flamboyant.ConNombreNuevo(path)
		if !ok {
			return "", false
		}
		path = nuevo
	case !on && porNuevo:
		// path ya llega traducido al nombre viejo.
	default:
		return "", false
	}
	if q := c.Request.URL.RawQuery; q != "" {
		path += "?" + q
	}
	return path, true
}
