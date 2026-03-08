package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     10,
		burst:    5,
		cleanup:  5 * time.Minute,
	}

	for i := 0; i < 5; i++ {
		if !rl.allow("192.168.1.1") {
			t.Fatalf("request %d was blocked, expected allowed (burst=5)", i+1)
		}
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     0.1, // very slow refill
		burst:    3,
		cleanup:  5 * time.Minute,
	}

	// Exhaust all tokens
	for i := 0; i < 3; i++ {
		if !rl.allow("10.0.0.1") {
			t.Fatalf("request %d was blocked unexpectedly", i+1)
		}
	}

	// Next request should be blocked
	if rl.allow("10.0.0.1") {
		t.Error("request after burst exhaustion was allowed, expected blocked")
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     100, // 100 tokens/sec — refills quickly
		burst:    2,
		cleanup:  5 * time.Minute,
	}

	// Exhaust tokens
	for i := 0; i < 2; i++ {
		rl.allow("10.0.0.2")
	}

	if rl.allow("10.0.0.2") {
		t.Fatal("should be blocked after exhaustion")
	}

	// Wait long enough for at least 1 token to refill (100/sec = 10ms per token)
	time.Sleep(50 * time.Millisecond)

	if !rl.allow("10.0.0.2") {
		t.Error("request should be allowed after token refill")
	}
}

func TestRateLimiter_IndependentPerIP(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     0.01, // effectively no refill during test
		burst:    2,
		cleanup:  5 * time.Minute,
	}

	// Exhaust IP A
	rl.allow("10.0.0.10")
	rl.allow("10.0.0.10")

	if rl.allow("10.0.0.10") {
		t.Error("IP A should be blocked after exhaustion")
	}

	// IP B should still be allowed
	if !rl.allow("10.0.0.20") {
		t.Error("IP B should be allowed (independent limit)")
	}
	if !rl.allow("10.0.0.20") {
		t.Error("IP B second request should be allowed")
	}
}

func TestRateLimiter_Middleware429(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*bucket),
		rate:     0.01,
		burst:    1,
		cleanup:  5 * time.Minute,
	}

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request succeeds
	r1 := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	r1.RemoteAddr = "192.168.0.1:12345"
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, r1)

	if w1.Code != http.StatusOK {
		t.Fatalf("first request: got %d, want %d", w1.Code, http.StatusOK)
	}

	// Second request gets rate-limited
	r2 := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	r2.RemoteAddr = "192.168.0.1:12346"
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r2)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: got %d, want %d", w2.Code, http.StatusTooManyRequests)
	}

	// Verify JSON response body
	var body map[string]string
	if err := json.Unmarshal(w2.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["error"] != "rate limit exceeded" {
		t.Errorf("error message: got %q, want %q", body["error"], "rate limit exceeded")
	}

	// Verify Retry-After header
	if ra := w2.Header().Get("Retry-After"); ra != "1" {
		t.Errorf("Retry-After: got %q, want %q", ra, "1")
	}

	// Verify Content-Type
	if ct := w2.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: got %q, want %q", ct, "application/json")
	}
}
