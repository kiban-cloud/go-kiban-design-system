package badge

// VariantForCode maps a kiban-shared backend status code to the visual
// variant the badge should use. Only codes that are reused across more
// than one project belong here — keeping this lookup small avoids the
// design system becoming a kitchen-sink of every domain vocabulary.
//
// Today the lookup covers:
//
//   - Payment-order lifecycle (PAID / PENDING / EXPIRED / CANCELLED) used
//     by rekon paymentorders + payments lists.
//   - Payment-method validation status (VALIDATED / TO_VALIDATE) used by
//     rekon paymentmethods.
//
// Codes are matched case-sensitively; pass the raw backend value and let
// the localized label travel separately. Unknown codes fall through to
// "neutral" so badges never crash a render — the worst outcome is a grey
// pill, which the user can still read.
//
// To override or extend per project, projects should call `Variant`
// directly with the variant they want. Adding cases here is reserved for
// codes that are genuinely shared across kiban projects.
func VariantForCode(code string) string {
	switch code {
	case "PAID", "VALIDATED", "ACTIVE", "COMPLETED":
		return "success"
	case "PENDING", "TO_VALIDATE", "PROCESSING":
		return "warning"
	case "EXPIRED", "FAILED", "REJECTED":
		return "danger"
	case "CANCELLED", "DRAFT":
		return "neutral"
	default:
		return "neutral"
	}
}
