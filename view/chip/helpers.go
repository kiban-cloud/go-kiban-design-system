package chip

// chipClass builds the outer span's class string for a given variant.
// All variants share the base layout (inline-flex pill with gap +
// rounded border + xs text); only the color tokens differ.
func chipClass(variant string) string {
	base := "inline-flex items-center gap-2 px-2 py-1 rounded-md border text-xs"
	switch variant {
	case "danger":
		return base + " border-red-400 bg-red-50 text-red-700"
	case "info":
		return base + " border-kiban-primary/40 bg-kiban-primary-soft text-kiban-primary"
	case "success":
		return base + " border-emerald-400 bg-emerald-50 text-emerald-700"
	case "warning":
		return base + " border-amber-400 bg-amber-50 text-amber-700"
	}
	return base + " border-kiban-border bg-kiban-surface text-kiban-ink2"
}

// subtextClass picks the color for the muted secondary line. The
// muted-on-color contrast varies by variant — keep them in sync with
// chipClass so the pair always reads.
func subtextClass(variant string) string {
	switch variant {
	case "danger":
		return "text-red-500"
	case "info":
		return "text-kiban-primary/70"
	case "success":
		return "text-emerald-600"
	case "warning":
		return "text-amber-600"
	}
	return "text-kiban-ink3"
}

// removeClass picks the color for the "×" button. Hover state slides
// to the variant's strong color; default-variant hover slides toward
// the ink token to match the rest of the kit.
func removeClass(variant string) string {
	base := "ml-0.5"
	switch variant {
	case "danger":
		return base + " text-red-500 hover:text-red-700"
	case "info":
		return base + " text-kiban-primary/70 hover:text-kiban-primary"
	case "success":
		return base + " text-emerald-600 hover:text-emerald-800"
	case "warning":
		return base + " text-amber-600 hover:text-amber-800"
	}
	return base + " text-kiban-ink3 hover:text-kiban-ink"
}
