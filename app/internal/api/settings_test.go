package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
)

func TestSettings_PrivilegedAllowed(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	req := httptest.NewRequest("GET", "/api/v1/settings", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleGetSettings(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSettings_NonPrivilegedBlocked(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	req := httptest.NewRequest("GET", "/api/v1/settings", nil)
	req = withRoleContext(req, auth.RoleSPC)
	rec := httptest.NewRecorder()

	h.handleGetSettings(rec, req)

	if rec.Code != 403 {
		t.Fatalf("expected 403 for non-privileged role, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSettings_SecretsMasked(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	s := defaultSettings()
	s.Outlook.ClientSecret = "super-secret-value"
	if err := h.saveSettings(s); err != nil {
		t.Fatalf("failed to save settings: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/settings", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleGetSettings(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var got AppSettings
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Outlook.ClientSecret != "••••••••" {
		t.Errorf("expected masked secret, got %q", got.Outlook.ClientSecret)
	}
}

func TestSettings_MaskedSecretsPreserved(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	// Save settings with a real secret
	s := defaultSettings()
	s.Outlook.ClientSecret = "super-secret-value"
	if err := h.saveSettings(s); err != nil {
		t.Fatalf("failed to save settings: %v", err)
	}

	// PUT with masked value — should preserve the original
	update := *s
	update.Outlook.ClientSecret = "••••••••"
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("failed to marshal update: %v", err)
	}

	req := httptest.NewRequest("PUT", "/api/v1/settings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleUpdateSettings(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Reload and verify the original secret is preserved
	loaded, err := h.loadSettings()
	if err != nil {
		t.Fatalf("failed to load settings: %v", err)
	}
	if loaded.Outlook.ClientSecret != "super-secret-value" {
		t.Errorf("expected preserved secret %q, got %q", "super-secret-value", loaded.Outlook.ClientSecret)
	}
}

func TestTestConnection_JSONAlwaysOK(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	body := `{"type":"json"}`
	req := httptest.NewRequest("POST", "/api/v1/settings/test-connection", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleTestConnection(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status %q, got %q", "ok", resp["status"])
	}
}

func TestTestConnection_PostgresNeedsDSN(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	body := `{"type":"postgres"}`
	req := httptest.NewRequest("POST", "/api/v1/settings/test-connection", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleTestConnection(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "error" {
		t.Errorf("expected status %q, got %q", "error", resp["status"])
	}
	msg, _ := resp["message"].(string)
	if !strings.Contains(msg, "Connection string is required") {
		t.Errorf("expected message about connection string, got %q", msg)
	}
}

func TestTestConnection_NonPrivileged(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	body := `{"type":"json"}`
	req := httptest.NewRequest("POST", "/api/v1/settings/test-connection", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleSPC)
	rec := httptest.NewRecorder()

	h.handleTestConnection(rec, req)

	if rec.Code != 403 {
		t.Fatalf("expected 403 for non-privileged role, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUpload_RejectsInvalidType(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("dummy content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/settings/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleUpload(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for .txt upload, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !strings.Contains(resp["error"], ".xlsx") || !strings.Contains(resp["error"], ".csv") {
		t.Errorf("expected error mentioning allowed types, got %q", resp["error"])
	}
}

func TestSystemInfo_Shape(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	req := httptest.NewRequest("GET", "/api/v1/settings/system-info", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleSystemInfo(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var info map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for _, key := range []string{"version", "studentCount", "authMode"} {
		if _, ok := info[key]; !ok {
			t.Errorf("expected key %q in system info response", key)
		}
	}
}

func TestSystemInfo_NonPrivileged(t *testing.T) {
	h := newTestHandler(t)
	h.settingsPath = filepath.Join(t.TempDir(), "settings.json")

	req := httptest.NewRequest("GET", "/api/v1/settings/system-info", nil)
	req = withRoleContext(req, auth.RoleSPC)
	rec := httptest.NewRecorder()

	h.handleSystemInfo(rec, req)

	if rec.Code != 403 {
		t.Fatalf("expected 403 for non-privileged role, got %d: %s", rec.Code, rec.Body.String())
	}
}
