package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
)

func TestSuggestedPrompts_ByRole(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		name     string
		role     string
		wantMin  int
		wantWord string // at least one prompt should contain this
	}{
		{"XO gets leadership prompts", auth.RoleXO, 4, "Morning brief"},
		{"Staff gets overview prompts", auth.RoleStaff, 4, "performance"},
		{"SPC gets company prompts", auth.RoleSPC, 3, "company"},
		{"Student gets self-focused prompts", auth.RoleStudent, 3, "doing"},
		{"Unknown role falls back to staff", "badRole", 4, "performance"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/chat/suggestions", nil)
			ctx := context.WithValue(req.Context(), middleware.RoleKey, tc.role)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			h.handleSuggestedPrompts(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", rec.Code)
			}

			var resp struct {
				Prompts []string `json:"prompts"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if len(resp.Prompts) < tc.wantMin {
				t.Errorf("expected at least %d prompts, got %d", tc.wantMin, len(resp.Prompts))
			}

			found := false
			for _, p := range resp.Prompts {
				if containsCI(p, tc.wantWord) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("no prompt contains %q; got %v", tc.wantWord, resp.Prompts)
			}
		})
	}
}

func containsCI(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		len(sub) == 0 ||
		findCI(s, sub))
}

func findCI(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			a, b := s[i+j], sub[j]
			if a >= 'A' && a <= 'Z' {
				a += 32
			}
			if b >= 'A' && b <= 'Z' {
				b += 32
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
