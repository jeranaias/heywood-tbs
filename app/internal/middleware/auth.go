package middleware

import (
	"context"
	"net/http"
	"strings"

	"heywood-tbs/internal/auth"
)

type contextKey string

const (
	identityKey contextKey = "identity"
	// Legacy keys kept for backward compatibility during transition
	RoleKey      contextKey = "role"
	CompanyKey   contextKey = "company"
	StudentIDKey contextKey = "studentId"
)

// AuthWithProvider returns middleware that uses the given IdentityProvider
// to authenticate requests and inject UserIdentity into the context.
func AuthWithProvider(provider auth.IdentityProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity := provider.Authenticate(r)

			ctx := context.WithValue(r.Context(), identityKey, identity)
			// Also set legacy keys so existing handlers work without changes
			ctx = context.WithValue(ctx, RoleKey, identity.Role)
			ctx = context.WithValue(ctx, CompanyKey, identity.Company)
			// For demo mode, StudentID comes from cookies (stored in Identity.ID as "demo-{studentId}")
			// We need to preserve the raw student ID for data lookups
			studentID := ""
			if identity.Source == "demo" {
				// In demo mode the DemoProvider doesn't store studentID separately,
				// so we read it from the cookie in the identity's ID field
				if c, err := r.Cookie("heywood-student-id"); err == nil && c.Value != "" {
					studentID = c.Value
				}
			} else {
				// In CAC mode, ID is the EDIPI itself
				studentID = identity.ID
			}
			ctx = context.WithValue(ctx, StudentIDKey, studentID)

			// Block unauthorized users from all API routes except /auth/me and /healthz
			if identity.Role == auth.RoleUnauthorized &&
				strings.HasPrefix(r.URL.Path, "/api/") &&
				r.URL.Path != "/api/v1/auth/me" &&
				r.URL.Path != "/api/v1/healthz" {
				http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Auth is the legacy middleware — uses DemoProvider for backward compatibility.
// Prefer AuthWithProvider for new code.
var Auth = AuthWithProvider(&auth.DemoProvider{})

// GetIdentity extracts the full UserIdentity from request context.
func GetIdentity(ctx context.Context) *auth.UserIdentity {
	if v, ok := ctx.Value(identityKey).(*auth.UserIdentity); ok {
		return v
	}
	return &auth.UserIdentity{ID: "unknown", Role: auth.RoleUnauthorized, Source: "unknown"}
}

// GetRole extracts the role from request context.
func GetRole(ctx context.Context) string {
	if v, ok := ctx.Value(RoleKey).(string); ok {
		return v
	}
	return auth.RoleUnauthorized
}

// GetCompany extracts the company filter from request context.
func GetCompany(ctx context.Context) string {
	if v, ok := ctx.Value(CompanyKey).(string); ok {
		return v
	}
	return ""
}

// GetStudentID extracts the student ID from request context.
func GetStudentID(ctx context.Context) string {
	if v, ok := ctx.Value(StudentIDKey).(string); ok {
		return v
	}
	return ""
}
