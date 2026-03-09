package api

import (
	"encoding/csv"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
)

func TestExportStudentsCSV(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/export/students", nil)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleExportStudents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/csv" {
		t.Errorf("expected Content-Type text/csv, got %s", ct)
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, "students.csv") {
		t.Errorf("expected Content-Disposition with students.csv, got %s", cd)
	}

	// Parse CSV
	reader := csv.NewReader(rec.Body)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	// Header + at least 1 data row
	if len(records) < 2 {
		t.Errorf("expected at least 2 CSV rows (header + data), got %d", len(records))
	}
	// Verify header
	if records[0][0] != "ID" {
		t.Errorf("expected first column header 'ID', got %s", records[0][0])
	}
}

func TestExportAtRiskCSV(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/export/at-risk", nil)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleExportAtRisk(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/csv" {
		t.Errorf("expected text/csv, got %s", ct)
	}

	reader := csv.NewReader(rec.Body)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	// Should have header row
	if len(records) < 1 {
		t.Error("expected at least a header row")
	}
}

func TestExportQualRecordsCSV(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/export/qual-records", nil)
	req = withFullContext(req, auth.RoleXO, "", "")
	rec := httptest.NewRecorder()

	h.handleExportQualRecords(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/csv" {
		t.Errorf("expected text/csv, got %s", ct)
	}

	reader := csv.NewReader(rec.Body)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	if len(records) < 2 {
		t.Errorf("expected at least 2 rows (header + data), got %d", len(records))
	}
}

func TestCompanyPerformanceSummary(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/reports/company-summary", nil)
	req = withFullContext(req, auth.RoleXO, "", "")
	rec := httptest.NewRecorder()

	h.handleCompanyPerformanceSummary(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Companies []struct {
			Company      string  `json:"company"`
			StudentCount int     `json:"studentCount"`
			AvgOverall   float64 `json:"avgOverall"`
			AtRiskCount  int     `json:"atRiskCount"`
			AtRiskPct    float64 `json:"atRiskPct"`
		} `json:"companies"`
		Total int `json:"total"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Total == 0 {
		t.Error("expected total > 0")
	}
	if len(resp.Companies) == 0 {
		t.Error("expected at least one company")
	}
	for _, co := range resp.Companies {
		if co.AvgOverall <= 0 {
			t.Errorf("company %s has non-positive avg overall: %.1f", co.Company, co.AvgOverall)
		}
	}
}
