package card

// PreviewOptions configura PreviewBody. Eyebrow es el antetítulo chico en
// mayúsculas (el negocio que cobra); vacío, no se pinta. Amount es el monto
// ya formateado («$1,500.00») y Currency la moneda que va chica a su lado.
type PreviewOptions struct {
	Eyebrow  string
	Amount   string
	Currency string
	Rows     []PreviewRow
}

// PreviewRow es un renglón de PreviewBody: la etiqueta gris a la izquierda y
// el valor a la derecha. Mono pinta el valor en monoespaciado (una CLABE).
type PreviewRow struct {
	Label, Value string
	Mono         bool
}
