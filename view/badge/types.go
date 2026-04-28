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
