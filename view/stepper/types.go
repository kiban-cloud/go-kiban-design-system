// Package stepper renders a horizontal multi-stage progress
// indicator: numbered dots connected by lines, each dot picking a
// colour from its status. Generic across the kiban backends that
// surface a sequential multi-step flow (NIP authentication phases,
// onboarding wizards, KYC progress, payment settlement steps …).
//
// Visual contract:
//   - dot variant by status:
//       "complete"   → emerald background, checkmark icon
//       "active"     → kiban-primary tinted background, kiban-primary border + number
//       "incomplete" → surface background, ink-4 border + number
//   - connector colour matches the *preceding* stage's status, so
//     "complete → incomplete" draws the line in emerald and
//     "incomplete → anything" draws it in border-grey.
//
// Status strings outside the three above fall through to the
// "incomplete" appearance so an unexpected value never explodes
// the layout.
package stepper

// Stage is one slot in the stepper. Status drives the dot colour
// and whether the dot shows a checkmark (complete) or the position
// number (active / incomplete).
type Stage struct {
	Label  string
	Status string // "complete" / "active" / "incomplete"
}

// Status constants — exported so callers can avoid stringly-typed
// values at the call site.
const (
	StatusComplete   = "complete"
	StatusActive     = "active"
	StatusIncomplete = "incomplete"
)
