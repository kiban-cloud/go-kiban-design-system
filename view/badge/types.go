// Package badge holds the kiban "pill" / status-marker primitives — the
// small coloured rectangles that flag row state (PAID, PENDING, …) in
// tables and detail headers.
//
// Two entry points:
//
//   - Variant(label, variant, size): the explicit, generic badge. The
//     caller picks the colour scheme (success / warning / danger / info /
//     neutral) and the size; the design system has no opinion about what
//     the label means.
//
//   - Status(code, label, size): a thin convenience that maps kiban-shared
//     backend codes to a Variant. Codes covered today are the ones reused
//     across rekon (payment-order status, payment-method validation
//     status, etc.). Projects with novel codes should call Variant
//     directly — don't extend this lookup with project-only codes.
//
// Sizes:
//
//	"sm"  // px-2 py-1 text-xs    — table-row default
//	"md"  // px-3 py-1 text-sm    — detail-header default
//
// Empty size string falls back to "sm" (the most common case).
//
// Variants follow the same vocabulary used in flash and card:
//
//	"success"   emerald
//	"warning"   amber
//	"danger"    red
//	"info"      kiban-primary tint
//	"neutral"   kiban-surface / kiban-ink3
//
// Empty variant string falls back to "neutral".
package badge

// pillClass es la píldora de estado de kiban cloud flamboyant: redonda, con un
// punto del color del texto delante. Mismas variantes y tamaños que el badge
// de siempre; neutral no aparece en el diseño y toma el gris de los textos
// secundarios.
func pillClass(variant, size string) string {
	c := "inline-flex items-center gap-1.5 whitespace-nowrap rounded-full font-medium before:h-1.5 before:w-1.5 before:shrink-0 before:rounded-full before:bg-current before:content-['']"
	if size == "md" {
		c += " px-3 py-1 text-sm"
	} else {
		c += " px-2.5 py-[3px] text-xs"
	}
	switch variant {
	case "success":
		return c + " bg-[#E3FBF5] text-[#007A63]"
	case "warning":
		return c + " bg-[#FFF4E5] text-[#9A4F00]"
	case "danger":
		return c + " bg-[#FDECEC] text-[#B42318]"
	case "info":
		return c + " bg-[#EAF0FF] text-[#0036C4]"
	default:
		return c + " bg-[#F2F4F8] text-[#474F66]"
	}
}
