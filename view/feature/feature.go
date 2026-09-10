// Package feature marca en un templ lo que depende de una funcionalidad que
// todavía no está (o no siempre está) prendida: un bloque nuevo, un reemplazo
// de algo existente. La decisión de qué está prendido no vive en la vista sino
// en el contexto del request, que llena quien renderiza:
//
//   - en producción, un middleware del módulo, desde su configuración: es un
//     feature flag común, y una funcionalidad se despliega apagada y se prende
//     sin tocar el templ;
//   - en kiban-proto, desde la cookie del selector de fase de la topbar, con el
//     modo Comparar que muestra todo y resalta lo que pertenece a cada bandera.
//
// Las banderas se nombran por funcionalidad ("cuotas", "dispersion-masiva"),
// nunca por número de fase: las fases cambian, las funcionalidades no.
package feature

import "context"

type ctxKey struct{}

// Flags es lo que el request sabe de las banderas.
type Flags struct {
	// Enabled son las banderas prendidas. Una bandera ausente está apagada.
	Enabled map[string]bool
	// Compare hace que Block y Replace rendericen todo, prendido o no, y
	// resalten cada bloque con su etiqueta. Es el modo «Comparar» del proto;
	// en producción nadie lo prende.
	Compare bool
	// Labels es la etiqueta que se muestra al resaltar, por bandera ("Fase 2").
	// Sin etiqueta se muestra el nombre de la bandera.
	Labels map[string]string
}

// With deja las banderas en el contexto. Para que lleguen al templ hay que
// ponerlas en el contexto del request antes de renderizar, porque view.Render
// renderiza con c.Request.Context():
//
//	c.Request = c.Request.WithContext(feature.With(c.Request.Context(), flags))
func With(ctx context.Context, f Flags) context.Context {
	return context.WithValue(ctx, ctxKey{}, f)
}

func flags(ctx context.Context) Flags {
	f, _ := ctx.Value(ctxKey{}).(Flags)
	return f
}

// Enabled dice si la bandera está prendida. Sin banderas en el contexto, todo
// está apagado: una vista renderizada sin middleware muestra sólo lo que ya
// existe.
func Enabled(ctx context.Context, name string) bool {
	return flags(ctx).Enabled[name]
}

// Comparing dice si el request está en modo Comparar.
func Comparing(ctx context.Context) bool {
	return flags(ctx).Compare
}

// Label es la etiqueta con la que se resalta una bandera en modo Comparar.
func Label(ctx context.Context, name string) string {
	if l := flags(ctx).Labels[name]; l != "" {
		return l
	}
	return name
}
