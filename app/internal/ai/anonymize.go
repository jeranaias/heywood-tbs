package ai

import (
	"fmt"
	"regexp"
	"strings"

	"heywood-tbs/internal/models"
)

var edipiPattern = regexp.MustCompile(`\b\d{10}\b`)

// AnonymizeStudent returns a text summary of a student's data with PII stripped.
func AnonymizeStudent(s *models.Student) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Student Data (anonymized — ID: %s):\n", s.ID)
	fmt.Fprintf(&b, "Rank: %s | Phase: %s | Status: %s\n", s.Rank, s.Phase, s.Status)
	fmt.Fprintf(&b, "\nAcademic Pillar (32%% weight):\n")
	fmt.Fprintf(&b, "  Exam 1: %.1f | Exam 2: %.1f | Exam 3: %.1f | Exam 4: %.1f\n", s.Exam1, s.Exam2, s.Exam3, s.Exam4)
	fmt.Fprintf(&b, "  Quiz Average: %.1f | Academic Composite: %.1f\n", s.QuizAvg, s.AcademicComposite)
	fmt.Fprintf(&b, "\nMilitary Skills Pillar (32%% weight):\n")
	fmt.Fprintf(&b, "  PFT: %d | CFT: %d\n", s.PFTScore, s.CFTScore)
	fmt.Fprintf(&b, "  Rifle: %s | Pistol: %s\n", s.RifleQual, s.PistolQual)
	fmt.Fprintf(&b, "  Land Nav Day: %s | Night: %s | Written: %.1f\n", s.LandNavDay, s.LandNavNight, s.LandNavWritten)
	fmt.Fprintf(&b, "  Obstacle Course: %s | Endurance Course: %s\n", s.ObstacleCourse, s.EnduranceCourse)
	fmt.Fprintf(&b, "  Mil Skills Composite: %.1f\n", s.MilSkillsComposite)
	fmt.Fprintf(&b, "\nLeadership Pillar (36%% weight):\n")
	fmt.Fprintf(&b, "  Week 12 SPC Eval: %.1f | Peer Eval: %.1f\n", s.LeadershipWeek12, s.PeerEvalWeek12)
	fmt.Fprintf(&b, "  Week 22 SPC Eval: %.1f | Peer Eval: %.1f\n", s.LeadershipWeek22, s.PeerEvalWeek22)
	fmt.Fprintf(&b, "  Leadership Composite: %.1f\n", s.LeadershipComposite)
	fmt.Fprintf(&b, "\nOverall: %.1f | Trend: %s | At-Risk: %v\n", s.OverallComposite, s.Trend, s.AtRisk)
	if len(s.RiskFlags) > 0 {
		fmt.Fprintf(&b, "Risk Flags: %s\n", strings.Join(s.RiskFlags, ", "))
	}
	if s.Notes != "" {
		fmt.Fprintf(&b, "Notes: %s\n", s.Notes)
	}
	return b.String()
}

// AnonymizeStudentList returns a summary table of students with PII stripped.
func AnonymizeStudentList(students []models.Student) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Student Roster (%d students):\n\n", len(students))
	fmt.Fprintf(&b, "%-8s %-6s %-8s %6s %6s %6s %6s %-5s %s\n",
		"ID", "Rank", "Phase", "Acad", "MilSk", "Ldr", "OvAll", "Trend", "Flags")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 80))
	for _, s := range students {
		flags := "-"
		if len(s.RiskFlags) > 0 {
			flags = strings.Join(s.RiskFlags, ",")
		}
		phase := s.Phase
		if len(phase) > 8 {
			phase = phase[:8]
		}
		fmt.Fprintf(&b, "%-8s %-6s %-8s %6.1f %6.1f %6.1f %6.1f %-5s %s\n",
			s.ID, s.Rank, phase, s.AcademicComposite, s.MilSkillsComposite,
			s.LeadershipComposite, s.OverallComposite, s.Trend, flags)
	}
	return b.String()
}

// StripPII removes any EDIPIs or names from arbitrary text.
func StripPII(text string) string {
	return edipiPattern.ReplaceAllString(text, "[REDACTED]")
}
