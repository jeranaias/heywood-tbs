package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/models"
)

func TestCreateCounseling(t *testing.T) {
	h := newTestHandler(t)

	body := `{"studentId":"STU-001","type":"initial","notes":"Test counseling"}`
	req := httptest.NewRequest("POST", "/api/v1/counselings", strings.NewReader(body))
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleCreateCounseling(rec, req)

	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var session models.CounselingSession
	if err := json.NewDecoder(rec.Body).Decode(&session); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if session.StudentID != "STU-001" {
		t.Errorf("expected studentId=STU-001, got %s", session.StudentID)
	}
	if session.Type != "initial" {
		t.Errorf("expected type=initial, got %s", session.Type)
	}
	if session.Status != "draft" {
		t.Errorf("expected status=draft, got %s", session.Status)
	}
	if session.StudentName == "" {
		t.Error("expected studentName to be populated from store")
	}
}

func TestListCounselings(t *testing.T) {
	h := newTestHandler(t)

	// Create two sessions
	s1 := models.CounselingSession{StudentID: "STU-001", Type: "initial", Status: "draft"}
	s2 := models.CounselingSession{StudentID: "STU-002", Type: "progress", Status: "draft"}
	h.store.CreateCounseling(s1)
	h.store.CreateCounseling(s2)

	// List all
	req := httptest.NewRequest("GET", "/api/v1/counselings", nil)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()
	h.handleListCounselings(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Sessions []models.CounselingSession `json:"sessions"`
		Total    int                         `json:"total"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Total < 2 {
		t.Errorf("expected at least 2 sessions, got %d", resp.Total)
	}

	// List filtered by student
	req2 := httptest.NewRequest("GET", "/api/v1/counselings?studentId=STU-001", nil)
	req2 = withFullContext(req2, auth.RoleStaff, "", "")
	rec2 := httptest.NewRecorder()
	h.handleListCounselings(rec2, req2)

	var resp2 struct {
		Sessions []models.CounselingSession `json:"sessions"`
		Total    int                         `json:"total"`
	}
	json.NewDecoder(rec2.Body).Decode(&resp2)
	if resp2.Total != 1 {
		t.Errorf("expected 1 session for STU-001, got %d", resp2.Total)
	}
}

func TestUpdateCounseling_AddFollowUp(t *testing.T) {
	h := newTestHandler(t)

	// Create a session first
	s := models.CounselingSession{StudentID: "STU-001", Type: "progress", Status: "draft"}
	h.store.CreateCounseling(s)
	sessions := h.store.ListCounselings("STU-001")
	if len(sessions) == 0 {
		t.Fatal("no sessions created")
	}
	id := sessions[0].ID

	// Update with follow-ups
	body := `{"notes":"Updated notes","status":"conducted","followUps":[{"description":"Retake land nav written","dueDate":"2026-04-01","status":"pending"}]}`
	req := httptest.NewRequest("PUT", "/api/v1/counselings/"+id, strings.NewReader(body))
	req.SetPathValue("id", id)
	req = withFullContext(req, auth.RoleSPC, "Alpha", "")
	rec := httptest.NewRecorder()

	h.handleUpdateCounseling(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.CounselingSession
	json.NewDecoder(rec.Body).Decode(&updated)
	if updated.Status != "conducted" {
		t.Errorf("expected status=conducted, got %s", updated.Status)
	}
	if updated.Notes != "Updated notes" {
		t.Errorf("expected updated notes, got %q", updated.Notes)
	}
	if len(updated.FollowUps) != 1 {
		t.Errorf("expected 1 follow-up, got %d", len(updated.FollowUps))
	}
}

func TestGenerateOutline(t *testing.T) {
	h := newTestHandler(t)

	body := `{"studentId":"STU-001","type":"progress"}`
	req := httptest.NewRequest("POST", "/api/v1/counselings/generate-outline", strings.NewReader(body))
	req = withFullContext(req, auth.RoleSPC, "Alpha", "")
	rec := httptest.NewRecorder()

	h.handleGenerateCounselingOutline(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Outline string `json:"outline"`
	}
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Outline == "" {
		t.Error("expected non-empty outline")
	}
	// Should contain student data references
	if !strings.Contains(resp.Outline, "Counseling") && !strings.Contains(resp.Outline, "counseling") {
		t.Error("expected outline to contain counseling-related content")
	}
}

func TestCounselingRoleGating(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		name    string
		method  string
		path    string
		handler http.HandlerFunc
		body    string
	}{
		{"create", "POST", "/api/v1/counselings", h.handleCreateCounseling, `{"studentId":"STU-001"}`},
		{"list", "GET", "/api/v1/counselings", h.handleListCounselings, ""},
		{"get", "GET", "/api/v1/counselings/COUNS-001", h.handleGetCounseling, ""},
		{"update", "PUT", "/api/v1/counselings/COUNS-001", h.handleUpdateCounseling, `{"notes":"test"}`},
		{"generate", "POST", "/api/v1/counselings/generate-outline", h.handleGenerateCounselingOutline, `{"studentId":"STU-001"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			req = withFullContext(req, auth.RoleStudent, "", "STU-001")
			rec := httptest.NewRecorder()

			tc.handler(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Errorf("expected 403 for student role on %s, got %d", tc.name, rec.Code)
			}
		})
	}
}
