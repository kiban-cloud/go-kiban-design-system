// Package flash holds the kiban "banner" primitive — the full-width
// status messages rendered above forms / above list pages after a
// mutation, an authorization failure, or any other actionable
// notification the user needs to see before continuing.
//
// Banners are server-rendered and re-rendered on the next request /
// HTMX swap. They are intentionally NOT dismissable from the client —
// no JS, no per-user state, no localStorage. If a banner needs to
// disappear, the next render decides; that keeps the contract simple
// and matches how the existing controllers already think about
// FlashSuccess / FlashError fields on their views.
//
// Variants follow the same vocabulary used in card and badge:
//
//	"success"   emerald — completed actions, positive confirmations
//	"error"     red     — operations that failed, hard blocks
//	"warning"   amber   — soft warnings, attention-required states
//	"info"      blue / kiban-primary — neutral informative messages
//
// Empty variant falls back to "info".
//
// Two API shapes:
//
//   - Banner(variant, msg): generic, takes a variant string. Use this
//     when the variant is computed at runtime.
//   - Success(msg) / Error(msg) / Warning(msg) / Info(msg): typed
//     wrappers, the more idiomatic call when the variant is known at
//     compile time. They delegate to Banner.
package flash
