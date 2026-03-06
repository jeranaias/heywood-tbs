package auth

import "net/http"

// UserIdentity represents an authenticated user regardless of auth source.
type UserIdentity struct {
	ID      string // EDIPI for CAC, "demo-xo" for demo mode
	Name    string // "Maj Smith, John" or "Executive Officer"
	Role    string // "xo", "staff", "spc", "student"
	Company string // for SPC/student filtering
	Email   string // from CAC SAN or manual entry
	Source  string // "demo" or "cac"
}

// IdentityProvider abstracts authentication so the same binary supports
// demo mode (role picker cookies) and CAC/PKI login on MCEN.
type IdentityProvider interface {
	// Authenticate inspects the request and returns the user's identity.
	// Implementations must always return a non-nil identity (default to staff).
	Authenticate(r *http.Request) *UserIdentity

	// SupportsSwitch reports whether this provider allows role switching.
	// Demo mode allows it; CAC mode does not.
	SupportsSwitch() bool
}
