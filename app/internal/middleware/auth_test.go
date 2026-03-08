package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"heywood-tbs/internal/auth"
)

// mockProvider implements auth.IdentityProvider for testing.
type mockProvider struct {
	identity *auth.UserIdentity
}

func (m *mockProvider) Authenticate(_ *http.Request) *auth.UserIdentity {
	return m.identity
}

func (m *mockProvider) SupportsSwitch() bool { return false }

// okHandler is a simple handler that returns 200 OK.
func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func TestAuthMiddleware_InjectsIdentity(t *testing.T) {
	want := &auth.UserIdentity{
		ID:      "1234567890",
		Name:    "Maj Smith, John",
		Role:    auth.RoleXO,
		Company: "Alpha",
		Email:   "john.smith@usmc.mil",
		Source:  "cac",
	}

	provider := &mockProvider{identity: want}
	mw := AuthWithProvider(provider)

	var gotIdentity *auth.UserIdentity
	var gotRole string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotIdentity = GetIdentity(r.Context())
		gotRole = GetRole(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	w := httptest.NewRecorder()

	mw(inner).ServeHTTP(w, r)

	if gotIdentity == nil {
		t.Fatal("GetIdentity returned nil")
	}
	assertField(t, "ID", want.ID, gotIdentity.ID)
	assertField(t, "Role", want.Role, gotIdentity.Role)
	assertField(t, "Company", want.Company, gotIdentity.Company)
	assertField(t, "ContextRole", want.Role, gotRole)
}

func TestAuthMiddleware_UnauthorizedBlocked(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID:   "unknown",
		Role: auth.RoleUnauthorized,
	}}
	mw := AuthWithProvider(provider)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	w := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_UnauthorizedAllowsAuthMe(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID:   "unknown",
		Role: auth.RoleUnauthorized,
	}}
	mw := AuthWithProvider(provider)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d for /api/v1/auth/me", w.Code, http.StatusOK)
	}
}

func TestAuthMiddleware_NonAPIPassesThrough(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID:   "unknown",
		Role: auth.RoleUnauthorized,
	}}
	mw := AuthWithProvider(provider)

	paths := []string{"/index.html", "/static/app.js", "/", "/favicon.ico"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()

			mw(okHandler()).ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want %d for %s", w.Code, http.StatusOK, path)
			}
		})
	}
}

func TestGetIdentity_NoContext(t *testing.T) {
	ctx := context.Background()
	id := GetIdentity(ctx)

	if id == nil {
		t.Fatal("GetIdentity returned nil on empty context")
	}
	assertField(t, "Role", auth.RoleUnauthorized, id.Role)
	assertField(t, "ID", "unknown", id.ID)
}

func TestGetRole_NoContext(t *testing.T) {
	ctx := context.Background()
	role := GetRole(ctx)

	assertField(t, "Role", auth.RoleUnauthorized, role)
}

// assertField is a test helper that compares two strings.
func assertField(t *testing.T, field, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}
