package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"heywood-tbs/internal/auth"
)

// securityHeaderChecks asserts the standard STIG security headers are present.
func securityHeaderChecks(t *testing.T, h http.Header) {
	t.Helper()

	assertField(t, "Strict-Transport-Security",
		"max-age=31536000; includeSubDomains",
		h.Get("Strict-Transport-Security"))

	assertField(t, "X-Frame-Options", "DENY", h.Get("X-Frame-Options"))
	assertField(t, "X-Content-Type-Options", "nosniff", h.Get("X-Content-Type-Options"))
	assertField(t, "Referrer-Policy", "strict-origin-when-cross-origin", h.Get("Referrer-Policy"))
	assertField(t, "Permissions-Policy", "camera=(), microphone=(), geolocation=()", h.Get("Permissions-Policy"))

	if csp := h.Get("Content-Security-Policy"); csp == "" {
		t.Error("Content-Security-Policy header is empty, expected non-empty")
	}
}

func TestChain_200WithSecurityHeaders(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID: "test", Role: auth.RoleXO, Source: "test",
	}}
	limiter := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     100,
		burst:    200,
		cleanup:  5 * time.Minute,
	}
	chain := Chain(
		Recovery,
		MaxBodySize(1<<20),
		limiter.Middleware,
		SecurityHeaders,
		CORS(false),
		AuthWithProvider(provider),
	)
	handler := chain(okHandler())

	r := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	r.RemoteAddr = "10.0.0.1:9999"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	securityHeaderChecks(t, w.Header())
}

func TestChain_Unauthorized401WithHeaders(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID: "unknown", Role: auth.RoleUnauthorized, Source: "test",
	}}
	limiter := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     100,
		burst:    200,
		cleanup:  5 * time.Minute,
	}
	chain := Chain(
		Recovery,
		MaxBodySize(1<<20),
		limiter.Middleware,
		SecurityHeaders,
		CORS(false),
		AuthWithProvider(provider),
	)
	handler := chain(okHandler())

	r := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	r.RemoteAddr = "10.0.0.2:9999"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	// SecurityHeaders runs before Auth in the chain, so security headers
	// must still be present even on 401 responses.
	securityHeaderChecks(t, w.Header())
}

func TestChain_PanicRecovery500WithHeaders(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID: "test", Role: auth.RoleXO, Source: "test",
	}}
	limiter := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     100,
		burst:    200,
		cleanup:  5 * time.Minute,
	}

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	chain := Chain(
		Recovery,
		MaxBodySize(1<<20),
		limiter.Middleware,
		SecurityHeaders,
		CORS(false),
		AuthWithProvider(provider),
	)
	handler := chain(panicHandler)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	r.RemoteAddr = "10.0.0.3:9999"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	body := w.Body.String()
	if !strings.Contains(body, "internal server error") {
		t.Errorf("body = %q, want it to contain %q", body, "internal server error")
	}

	// Recovery catches the panic after SecurityHeaders has already set
	// headers on the ResponseWriter, so they should still be present.
	securityHeaderChecks(t, w.Header())
}

func TestChain_RateLimit429(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID: "test", Role: auth.RoleXO, Source: "test",
	}}
	limiter := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     0.01, // effectively no refill during test
		burst:    1,
		cleanup:  5 * time.Minute,
	}
	chain := Chain(
		Recovery,
		MaxBodySize(1<<20),
		limiter.Middleware,
		SecurityHeaders,
		CORS(false),
		AuthWithProvider(provider),
	)
	handler := chain(okHandler())

	// First request — should be allowed.
	r1 := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	r1.RemoteAddr = "10.0.0.4:9999"
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)

	if w1.Code != http.StatusOK {
		t.Fatalf("first request: status = %d, want %d", w1.Code, http.StatusOK)
	}

	// Second request — should be rate-limited.
	r2 := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	r2.RemoteAddr = "10.0.0.4:9998"
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: status = %d, want %d", w2.Code, http.StatusTooManyRequests)
	}

	if ra := w2.Header().Get("Retry-After"); ra != "1" {
		t.Errorf("Retry-After: got %q, want %q", ra, "1")
	}
}

func TestChain_CORSInDevMode(t *testing.T) {
	provider := &mockProvider{identity: &auth.UserIdentity{
		ID: "test", Role: auth.RoleXO, Source: "test",
	}}
	limiter := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     100,
		burst:    200,
		cleanup:  5 * time.Minute,
	}
	chain := Chain(
		Recovery,
		MaxBodySize(1<<20),
		limiter.Middleware,
		SecurityHeaders,
		CORS(true), // dev mode enabled
		AuthWithProvider(provider),
	)
	handler := chain(okHandler())

	r := httptest.NewRequest(http.MethodOptions, "/api/v1/students", nil)
	r.RemoteAddr = "10.0.0.5:9999"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	assertField(t, "Access-Control-Allow-Origin",
		"http://localhost:5173",
		w.Header().Get("Access-Control-Allow-Origin"))

	assertField(t, "Access-Control-Allow-Credentials",
		"true",
		w.Header().Get("Access-Control-Allow-Credentials"))
}
