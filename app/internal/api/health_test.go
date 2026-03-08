package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz_ReturnsOK(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	h.handleHealthz(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for _, key := range []string{"status", "version", "uptime", "timestamp", "checks"} {
		if _, ok := resp[key]; !ok {
			t.Errorf("response missing key %q", key)
		}
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}

	checks, ok := resp["checks"].(map[string]interface{})
	if !ok {
		t.Fatal("checks is not an object")
	}
	if checks["dataStore"] != "ok" {
		t.Errorf("expected dataStore 'ok', got %q", checks["dataStore"])
	}
}

func TestHealthz_NilStoreReturns503(t *testing.T) {
	h := newTestHandler(t)
	h.store = nil

	req := httptest.NewRequest("GET", "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	h.handleHealthz(rec, req)

	if rec.Code != 503 {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["status"] != "degraded" {
		t.Errorf("expected status 'degraded', got %q", resp["status"])
	}
}

func TestHealthz_NoAuthRequired(t *testing.T) {
	h := newTestHandler(t)

	// Send request with no role context at all — simulates unauthenticated probe
	req := httptest.NewRequest("GET", "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	h.handleHealthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 without auth, got %d", rec.Code)
	}
}
