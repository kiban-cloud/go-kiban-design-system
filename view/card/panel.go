package card

import "github.com/a-h/templ"

// noteTone es el tono de un Note para data-kiban-note: warning o info.
func noteTone(tone string) string {
	if tone == "warning" {
		return "warning"
	}
	return "info"
}

// CalloutOptions configura Callout. Tone es "ok" (verde), "warn" (ámbar),
// "bad" (rojo) o, cualquier otro, "info" (azul). Check cambia el ícono de
// alerta por una palomita: para un resultado que salió bien aunque el tono
// sea informativo. Attrs se agregan al contenedor (p. ej. role="status"
// cuando el aviso es el estado de lo que se está viendo).
type CalloutOptions struct {
	Tone  string
	Text  string
	Check bool
	Attrs templ.Attributes
}

func calloutTone(tone string) string {
	switch tone {
	case "ok", "warn", "bad":
		return tone
	}
	return "info"
}

func calloutToneClass(tone string) string {
	switch calloutTone(tone) {
	case "ok":
		return "bg-[#E3FBF5] text-[#007A63]"
	case "warn":
		return "bg-[#FFF6E5] text-[#8A5A00]"
	case "bad":
		return "bg-[#FDECEC] text-[#B42318]"
	}
	return "bg-[#EAF0FF] text-[#0036C4]"
}
