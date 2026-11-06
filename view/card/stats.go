package card

import "github.com/a-h/templ"

// StatOptions es una cifra de Stats: el número (ya formateado), su etiqueta,
// el ícono y el tono ("success" verde, "warning" ámbar, "danger" rojo; vacío,
// azul).
type StatOptions struct {
	Value string
	Label string
	Icon  templ.Component
	Tone  string
}

// statToneClass son el fondo y el color del cuadro del ícono.
func statToneClass(tone string) string {
	switch tone {
	case "success":
		return "bg-[#E3FBF5] text-[#007A63]"
	case "warning":
		return "bg-[#FFF4E5] text-[#9A4F00]"
	case "danger":
		return "bg-[#FDECEC] text-[#B42318]"
	}
	return "bg-kiban-primary-soft text-kiban-primary"
}
