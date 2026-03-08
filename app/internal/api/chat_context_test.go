package api

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/data"
)

func TestExtractStudentID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"exact STU format", "STU-042", "STU-042"},
		{"lowercase stu format", "stu-042", "STU-042"},
		{"in sentence", "Tell me about STU-001 scores", "STU-001"},
		{"student hash", "student #42", "STU-042"},
		{"student word", "student 7", "STU-007"},
		{"no match", "no match here", ""},
		{"empty string", "", ""},
		{"student without number", "student abc", ""},
		{"just number", "42", ""},
		{"STU prefix too short", "STU-0", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractStudentID(tc.input)
			if got != tc.want {
				t.Errorf("extractStudentID(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// testDataDir returns the absolute path to the project data directory.
func testDataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file location")
	}
	// From app/internal/api/ -> app/data/
	return filepath.Join(filepath.Dir(filename), "..", "..", "data")
}

// newTestHandler creates a Handler with a real data store for testing.
func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	dataDir := testDataDir(t)
	store, err := data.NewStore(dataDir)
	if err != nil {
		t.Fatalf("failed to create store from %s: %v", dataDir, err)
	}
	return &Handler{
		store:        store,
		weatherSvc:   &ai.WeatherService{},
		newsSvc:      &ai.NewsService{},
		trafficSvc:   &ai.TrafficService{},
		authProvider: &auth.DemoProvider{},
		dev:          true,
	}
}

func TestBuildChatContext_XO(t *testing.T) {
	h := newTestHandler(t)

	systemPrompt, _ := h.buildChatContext(auth.RoleXO, "", "", "good morning")

	// XO system prompt should contain key elements
	checks := []struct {
		label    string
		expected string
	}{
		{"persona reference", "Executive Officer"},
		{"weather section", "WEATHER"},
		{"schedule section", "SCHEDULE"},
		{"student overview", "STUDENT OVERVIEW"},
		{"at-risk section", "AT-RISK"},
		{"qualification section", "QUALIFICATION"},
		{"instructor section", "INSTRUCTOR"},
	}

	for _, c := range checks {
		if !strings.Contains(systemPrompt, c.expected) {
			t.Errorf("XO system prompt missing %s (expected substring %q)", c.label, c.expected)
		}
	}

	// System prompt should not be empty
	if len(systemPrompt) < 200 {
		t.Errorf("XO system prompt unexpectedly short: %d chars", len(systemPrompt))
	}
}

func TestBuildChatContext_Student(t *testing.T) {
	h := newTestHandler(t)

	systemPrompt, userContext := h.buildChatContext(auth.RoleStudent, "", "STU-001", "how am i doing")

	// Student system prompt should reference student-appropriate content
	if !strings.Contains(systemPrompt, "student") && !strings.Contains(systemPrompt, "Student") {
		t.Error("student system prompt should mention student role")
	}

	// Student context should include their own schedule or scores
	if userContext == "" {
		t.Log("userContext is empty — may be expected if STU-001 is not in data set")
	}

	// Student prompt should NOT contain other students' data indicators
	// (XO-specific sections like "INSTRUCTOR WORKLOAD" should be absent)
	if strings.Contains(systemPrompt, "INSTRUCTOR WORKLOAD") {
		t.Error("student system prompt should not contain INSTRUCTOR WORKLOAD section")
	}
	if strings.Contains(systemPrompt, "AT-RISK STUDENTS") {
		t.Error("student system prompt should not contain AT-RISK STUDENTS overview section")
	}
}
