package ai

import (
	"strings"
	"testing"

	"heywood-tbs/internal/models"
)

func TestAnonymizeStudent(t *testing.T) {
	s := &models.Student{
		ID:                  "STU-042",
		Rank:                "2ndLt",
		LastName:            "Thompson",
		FirstName:           "James",
		Phase:               "Phase 2",
		Status:              "Active",
		Exam1:               85.0,
		Exam2:               78.0,
		Exam3:               92.0,
		Exam4:               88.0,
		QuizAvg:             84.5,
		AcademicComposite:   85.5,
		PFTScore:            275,
		CFTScore:            290,
		RifleQual:           "Expert",
		PistolQual:          "Sharpshooter",
		LandNavDay:          "GO",
		LandNavNight:        "GO",
		LandNavWritten:      90.0,
		ObstacleCourse:      "GO",
		EnduranceCourse:     "GO",
		MilSkillsComposite:  82.3,
		LeadershipWeek12:    80.0,
		LeadershipWeek22:    85.0,
		PeerEvalWeek12:      78.0,
		PeerEvalWeek22:      82.0,
		LeadershipComposite: 81.2,
		OverallComposite:    83.0,
		Trend:               "up",
		AtRisk:              false,
		RiskFlags:           []string{"low_quiz_avg"},
		Notes:               "Improving steadily",
	}

	result := AnonymizeStudent(s)

	// Should contain ID, rank, scores
	if !strings.Contains(result, "STU-042") {
		t.Error("expected result to contain student ID STU-042")
	}
	if !strings.Contains(result, "2ndLt") {
		t.Error("expected result to contain rank 2ndLt")
	}
	if !strings.Contains(result, "85.0") {
		t.Errorf("expected result to contain Exam 1 score 85.0")
	}
	if !strings.Contains(result, "275") {
		t.Error("expected result to contain PFT score 275")
	}
	if !strings.Contains(result, "low_quiz_avg") {
		t.Error("expected result to contain risk flag low_quiz_avg")
	}

	// Should NOT contain the student's name
	if strings.Contains(result, "Thompson") {
		t.Error("expected result to NOT contain last name Thompson (PII)")
	}
	if strings.Contains(result, "James") {
		t.Error("expected result to NOT contain first name James (PII)")
	}
}

func TestStripPII_EDIPI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "EDIPI in sentence",
			input: "Contact 1234567890 for info",
			want:  "Contact [REDACTED] for info",
		},
		{
			name:  "multiple EDIPIs",
			input: "Users 1234567890 and 9876543210 are active",
			want:  "Users [REDACTED] and [REDACTED] are active",
		},
		{
			name:  "9 digits not redacted",
			input: "Phone 123456789 is not an EDIPI",
			want:  "Phone 123456789 is not an EDIPI",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := StripPII(tc.input)
			if got != tc.want {
				t.Errorf("StripPII(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestStripPII_NoMatch(t *testing.T) {
	input := "Hello world"
	got := StripPII(input)
	if got != input {
		t.Errorf("StripPII(%q) = %q, want unchanged", input, got)
	}
}

func TestStripNames(t *testing.T) {
	students := []models.Student{
		{ID: "STU-001", FirstName: "James", LastName: "Thompson"},
		{ID: "STU-002", FirstName: "Maria", LastName: "Garcia"},
	}

	text := "Thompson scored higher than Garcia. James Thompson leads the platoon."
	got := StripNames(text, students)

	if strings.Contains(got, "Thompson") {
		t.Errorf("expected Thompson to be replaced, got: %s", got)
	}
	if strings.Contains(got, "Garcia") {
		t.Errorf("expected Garcia to be replaced, got: %s", got)
	}
	if !strings.Contains(got, "STU-001") {
		t.Errorf("expected STU-001 in result, got: %s", got)
	}
	if !strings.Contains(got, "STU-002") {
		t.Errorf("expected STU-002 in result, got: %s", got)
	}
}
