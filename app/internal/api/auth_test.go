package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

// injectIdentity creates a request context with the given identity and legacy keys set.
func injectIdentity(r *http.Request, identity *auth.UserIdentity, studentID string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, middleware.RoleKey, identity.Role)
	ctx = context.WithValue(ctx, middleware.CompanyKey, identity.Company)
	ctx = context.WithValue(ctx, middleware.StudentIDKey, studentID)
	// Use the unexported identityKey via a round-trip through middleware.GetIdentity
	// Instead, directly embed using the context key that middleware expects.
	// We need to match the identityKey used in middleware.GetIdentity.
	// Since identityKey is unexported from middleware, we set the identity via
	// a wrapping approach: build a handler that sets it.
	return r.WithContext(ctx)
}

// withIdentityContext returns a context with identity values set using middleware's
// exported context keys plus the identity embedded via a test middleware pass-through.
func withIdentityContext(identity *auth.UserIdentity, studentID string) context.Context {
	// Create a dummy request, run it through AuthWithProvider, and capture the context.
	req := httptest.NewRequest("GET", "/", nil)
	if identity.Source == "demo" {
		req.AddCookie(&http.Cookie{Name: "heywood-role", Value: identity.Role})
		req.AddCookie(&http.Cookie{Name: "heywood-company", Value: identity.Company})
		req.AddCookie(&http.Cookie{Name: "heywood-student-id", Value: studentID})
	}

	var captured context.Context
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Context()
	})
	handler := middleware.AuthWithProvider(&auth.DemoProvider{})(inner)
	handler.ServeHTTP(httptest.NewRecorder(), req)
	return captured
}

func TestHandleAuthMe(t *testing.T) {
	h := newTestHandler(t)

	// Set up an XO identity via demo mode cookies
	ctx := withIdentityContext(&auth.UserIdentity{
		ID:     "demo-xo",
		Role:   auth.RoleXO,
		Source: "demo",
	}, "")

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.handleAuthMe(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var info models.AuthInfo
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Role != auth.RoleXO {
		t.Errorf("expected role %q, got %q", auth.RoleXO, info.Role)
	}
	if info.Name == "" {
		t.Error("expected non-empty name in response")
	}
}

func TestHandleAuthSwitch_DemoMode(t *testing.T) {
	h := newTestHandler(t)
	h.authProvider = &auth.DemoProvider{} // explicitly demo

	body := `{"role":"spc","company":"Alpha","studentId":""}`
	req := httptest.NewRequest("POST", "/api/v1/auth/switch", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.handleAuthSwitch(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var info models.AuthInfo
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Role != auth.RoleSPC {
		t.Errorf("expected role %q, got %q", auth.RoleSPC, info.Role)
	}
	if info.Company != "Alpha" {
		t.Errorf("expected company Alpha, got %q", info.Company)
	}

	// Verify role-switch cookie was set
	cookies := rec.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "heywood-role" && c.Value == "spc" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected heywood-role cookie to be set to 'spc'")
	}
}

func TestHandleAuthSwitch_CACMode(t *testing.T) {
	h := newTestHandler(t)
	// Use CAC provider which does NOT support switching
	h.authProvider = auth.NewCACProvider("nonexistent-roster.json")

	body := `{"role":"xo","company":"","studentId":""}`
	req := httptest.NewRequest("POST", "/api/v1/auth/switch", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.handleAuthSwitch(rec, req)

	if rec.Code != 403 {
		t.Fatalf("expected 403 for CAC mode, got %d: %s", rec.Code, rec.Body.String())
	}

	var errResp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if !strings.Contains(errResp["error"], "CAC") && !strings.Contains(errResp["error"], "not available") {
		t.Errorf("expected error mentioning CAC restriction, got: %q", errResp["error"])
	}
}
