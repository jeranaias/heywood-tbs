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

// withFullContext returns a request with role, company, and studentID injected into context.
func withFullContext(req *http.Request, role, company, studentID string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, role)
	ctx = context.WithValue(ctx, middleware.CompanyKey, company)
	ctx = context.WithValue(ctx, middleware.StudentIDKey, studentID)
	return req.WithContext(ctx)
}

// --- handleListStudents ---

func TestListStudents_StudentWithID(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students", nil)
	req = withFullContext(req, auth.RoleStudent, "", "STU-001")
	rec := httptest.NewRecorder()

	h.handleListStudents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Students []json.RawMessage `json:"students"`
		Total    int               `json:"total"`
		Filtered int               `json:"filtered"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Students) != 1 {
		t.Fatalf("expected exactly 1 student, got %d", len(resp.Students))
	}
	if resp.Total != 1 {
		t.Errorf("expected total=1, got %d", resp.Total)
	}
	if resp.Filtered != 1 {
		t.Errorf("expected filtered=1, got %d", resp.Filtered)
	}

	// Verify the returned student is STU-001
	var stu struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp.Students[0], &stu); err != nil {
		t.Fatalf("failed to decode student: %v", err)
	}
	if stu.ID != "STU-001" {
		t.Errorf("expected student ID STU-001, got %q", stu.ID)
	}
}

func TestListStudents_StudentEmptyID(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students", nil)
	req = withFullContext(req, auth.RoleStudent, "", "")
	rec := httptest.NewRecorder()

	h.handleListStudents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Students []json.RawMessage `json:"students"`
		Total    int               `json:"total"`
		Filtered int               `json:"filtered"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Students) != 0 {
		t.Errorf("expected 0 students, got %d", len(resp.Students))
	}
	if resp.Total != 0 {
		t.Errorf("expected total=0, got %d", resp.Total)
	}
	if resp.Filtered != 0 {
		t.Errorf("expected filtered=0, got %d", resp.Filtered)
	}
}

func TestListStudents_SPCCompanyFilter(t *testing.T) {
	h := newTestHandler(t)

	// SPC with company="Alpha" — even if query param says "Bravo", the SPC's
	// own company must override the query param. Since all test data students
	// are in Alpha, the SPC still sees them. The key assertion is that the
	// query param "Bravo" is ignored: we get Alpha students, not zero.
	req := httptest.NewRequest("GET", "/api/v1/students?company=Bravo", nil)
	req = withFullContext(req, auth.RoleSPC, "Alpha", "")
	rec := httptest.NewRecorder()

	h.handleListStudents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Students []struct {
			ID      string `json:"id"`
			Company string `json:"company"`
		} `json:"students"`
		Total    int `json:"total"`
		Filtered int `json:"filtered"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// SPC Alpha with ?company=Bravo must still see Alpha students (override).
	// If the query param were honored we would get 0 results.
	if len(resp.Students) == 0 {
		t.Fatal("SPC Alpha should see students (query param Bravo must be overridden)")
	}
	for _, s := range resp.Students {
		if s.Company != "Alpha" {
			t.Errorf("SPC Alpha saw student %s in company %q — expected Alpha only", s.ID, s.Company)
		}
	}
}

func TestListStudents_StaffFullList(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students", nil)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleListStudents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Students []json.RawMessage `json:"students"`
		Total    int               `json:"total"`
		Filtered int               `json:"filtered"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Total != 200 {
		t.Errorf("expected total=200, got %d", resp.Total)
	}
	if resp.Filtered != 200 {
		t.Errorf("expected filtered=200, got %d", resp.Filtered)
	}
	if len(resp.Students) != 200 {
		t.Errorf("expected 200 students in list, got %d", len(resp.Students))
	}
}

// --- handleGetStudent ---

func TestGetStudent_StaffAccessAny(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/STU-001", nil)
	req.SetPathValue("id", "STU-001")
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleGetStudent(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var stu struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&stu); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if stu.ID != "STU-001" {
		t.Errorf("expected student ID STU-001, got %q", stu.ID)
	}
}

func TestGetStudent_StudentOwnID(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/STU-001", nil)
	req.SetPathValue("id", "STU-001")
	req = withFullContext(req, auth.RoleStudent, "", "STU-001")
	rec := httptest.NewRecorder()

	h.handleGetStudent(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestGetStudent_StudentDifferentID(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/STU-001", nil)
	req.SetPathValue("id", "STU-001")
	req = withFullContext(req, auth.RoleStudent, "", "STU-050")
	rec := httptest.NewRecorder()

	h.handleGetStudent(rec, req)

	if rec.Code != 403 {
		t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != "access denied" {
		t.Errorf("expected error 'access denied', got %q", resp.Error)
	}
}

func TestGetStudent_SPCCrossCompany(t *testing.T) {
	h := newTestHandler(t)

	// All test data students are in Alpha company.
	// SPC with company="Bravo" tries to access STU-001 (Alpha) — should be 403.
	req := httptest.NewRequest("GET", "/api/v1/students/STU-001", nil)
	req.SetPathValue("id", "STU-001")
	req = withFullContext(req, auth.RoleSPC, "Bravo", "")
	rec := httptest.NewRecorder()

	h.handleGetStudent(rec, req)

	if rec.Code != 403 {
		t.Fatalf("expected 403 for SPC Bravo accessing Alpha student, got %d: %s", rec.Code, rec.Body.String())
	}

	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error: %v", err)
	}
	if errResp.Error != "access denied" {
		t.Errorf("expected error 'access denied', got %q", errResp.Error)
	}
}

func TestGetStudent_NotFound(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/STU-999", nil)
	req.SetPathValue("id", "STU-999")
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleGetStudent(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != "student not found" {
		t.Errorf("expected error 'student not found', got %q", resp.Error)
	}
}

// --- handleStudentStats ---

func TestStudentStats_SPCAlpha(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/stats", nil)
	req = withFullContext(req, auth.RoleSPC, "Alpha", "")
	rec := httptest.NewRecorder()

	h.handleStudentStats(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var stats struct {
		ActiveStudents int     `json:"activeStudents"`
		AvgComposite   float64 `json:"avgComposite"`
		AtRiskCount    int     `json:"atRiskCount"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode stats: %v", err)
	}

	// SPC Alpha sees only Alpha company students. All 200 test students are
	// in Alpha, so activeStudents should equal 200 (scoped correctly).
	if stats.ActiveStudents != 200 {
		t.Errorf("SPC Alpha activeStudents=%d, expected 200 (all students are Alpha)", stats.ActiveStudents)
	}
	if stats.AvgComposite <= 0 {
		t.Errorf("expected positive avgComposite, got %f", stats.AvgComposite)
	}
}

func TestStudentStats_StaffFull(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/stats", nil)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleStudentStats(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var stats struct {
		ActiveStudents int `json:"activeStudents"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode stats: %v", err)
	}

	if stats.ActiveStudents != 200 {
		t.Errorf("staff activeStudents=%d, expected 200", stats.ActiveStudents)
	}
}

// --- handleAtRisk ---

func TestAtRisk_SPCAlpha(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/at-risk", nil)
	req = withFullContext(req, auth.RoleSPC, "Alpha", "")
	rec := httptest.NewRecorder()

	h.handleAtRisk(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Students []struct {
			ID      string `json:"id"`
			Company string `json:"company"`
			AtRisk  bool   `json:"atRisk"`
		} `json:"students"`
		Total int `json:"total"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode at-risk: %v", err)
	}

	// All returned students must be Alpha and at-risk
	for _, s := range resp.Students {
		if s.Company != "Alpha" {
			t.Errorf("SPC Alpha at-risk list contains student %s in company %q", s.ID, s.Company)
		}
		if !s.AtRisk {
			t.Errorf("at-risk list contains student %s with atRisk=false", s.ID)
		}
	}

	if resp.Total != len(resp.Students) {
		t.Errorf("total=%d does not match students length=%d", resp.Total, len(resp.Students))
	}
}

func TestAtRisk_StaffFull(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/students/at-risk", nil)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleAtRisk(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Students []struct {
			ID     string `json:"id"`
			AtRisk bool   `json:"atRisk"`
		} `json:"students"`
		Total int `json:"total"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode at-risk: %v", err)
	}

	// Staff sees all at-risk students across all companies
	for _, s := range resp.Students {
		if !s.AtRisk {
			t.Errorf("at-risk list contains student %s with atRisk=false", s.ID)
		}
	}

	if resp.Total != len(resp.Students) {
		t.Errorf("total=%d does not match students length=%d", resp.Total, len(resp.Students))
	}

	if resp.Total == 0 {
		t.Error("expected at least some at-risk students for staff view")
	}
}
