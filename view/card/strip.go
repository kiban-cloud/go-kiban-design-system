package card

import "github.com/a-h/templ"

// StripItem es una cifra de Strip: el número en negritas, su etiqueta y, si
// hace falta, algo más al final (un enlace «¿por qué?»). Tone pinta el número:
// "warn" ámbar, "ok" verde; vacío, el color del texto.
type StripItem struct {
	Value string
	Label string
	Tone  string
	Extra templ.Component
}

func stripValueClass(tone string) string {
	switch tone {
	case "warn":
		return "font-semibold text-[#B26B00]"
	case "ok":
		return "font-semibold text-[#007A63]"
	}
	return "font-semibold text-kiban-ink"
}
