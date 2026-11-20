package rutas

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kiban-cloud/go-kiban-design-system/view/flamboyant"
)

// headersDeRuta son los headers con los que una respuesta manda al navegador
// a otra página: el redirect de siempre y los de htmx.
var headersDeRuta = []string{"Location", "HX-Redirect", "HX-Location", "HX-Push-Url", "HX-Replace-Url"}

// atributoDeRuta encuentra en el HTML los atributos que llevan una URL
// relativa: enlaces, forms, htmx y los data-*-url/-href/-action/-path/-src
// que lee el JS. Grupo 1 es el atributo con la comilla, grupo 2 el path.
var atributoDeRuta = regexp.MustCompile(`(\s(?:href|action|formaction|src|hx-(?:get|post|put|patch|delete|push-url|replace-url)|data-[a-z0-9-]*?(?:href|url|action|path|src))=["'])(/[^"'\s>]*)`)

// writer lleva al nombre nuevo del módulo (flamboyant.ConNombreNuevo) lo
// que una respuesta flamboyant manda al navegador: los headers de redirección
// y, si es HTML, los atributos con URL. Así las vistas y los controllers no
// cambian: un enlace o un redirect hacia /workfloo/… sale hacia
// /solicitudes/…. El HTML se junta en memoria y se reescribe al final; lo que
// no es HTML pasa tal cual. Las URLs dentro de JS no se tocan (llegan por el
// 308), ni las absolutas, ni un HX-Location en JSON.
type writer struct {
	gin.ResponseWriter
	html     bool
	buf      bytes.Buffer
	enviado  bool
	decidido bool
}

// fixHeaders corre en cada escritura mientras los headers no se hayan
// mandado: un header puesto después de c.Status también se corrige.
func (w *writer) fixHeaders() {
	if w.ResponseWriter.Written() {
		return
	}
	h := w.Header()
	for _, k := range headersDeRuta {
		if v := h.Get(k); v != "" {
			if nuevo, ok := flamboyant.ConNombreNuevo(v); ok {
				h.Set(k, nuevo)
			}
		}
	}
}

func (w *writer) esHTML() bool {
	if !w.decidido {
		w.decidido = true
		w.html = strings.HasPrefix(w.Header().Get("Content-Type"), "text/html")
	}
	return w.html
}

func (w *writer) WriteHeader(code int) {
	w.fixHeaders()
	w.ResponseWriter.WriteHeader(code)
}

func (w *writer) WriteHeaderNow() {
	w.fixHeaders()
	w.ResponseWriter.WriteHeaderNow()
}

func (w *writer) Write(b []byte) (int, error) {
	if !w.enviado && w.esHTML() {
		return w.buf.Write(b)
	}
	w.fixHeaders()
	return w.ResponseWriter.Write(b)
}

func (w *writer) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *writer) Written() bool {
	return w.buf.Len() > 0 || w.ResponseWriter.Written()
}

// Flush manda lo juntado antes de vaciar: quien hace streaming de HTML lo ve
// salir, ya reescrito, y lo que siga pasa tal cual.
func (w *writer) Flush() {
	w.terminar()
	w.ResponseWriter.Flush()
}

// terminar escribe el HTML juntado con las URLs reescritas y corrige los
// headers de una respuesta sin cuerpo (gin los escribe después de los
// middlewares, sin pasar por aquí).
func (w *writer) terminar() {
	if w.enviado {
		return
	}
	w.enviado = true
	w.fixHeaders()
	if w.buf.Len() == 0 {
		return
	}
	_, _ = w.ResponseWriter.Write(reescribirHTML(w.buf.Bytes()))
	w.buf.Reset()
}

func reescribirHTML(b []byte) []byte {
	return atributoDeRuta.ReplaceAllFunc(b, func(m []byte) []byte {
		sub := atributoDeRuta.FindSubmatch(m)
		nuevo, ok := flamboyant.ConNombreNuevo(string(sub[2]))
		if !ok {
			return m
		}
		return append(append([]byte{}, sub[1]...), nuevo...)
	})
}
