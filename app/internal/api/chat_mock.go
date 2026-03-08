package api

import (
	"fmt"
	"strings"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/models"
)

// mockResponse generates a mock response when no OpenAI API key is configured.
func (h *Handler) mockResponse(role, company, studentID, message string) string {
	msg := strings.ToLower(message)

	// XO mock: comprehensive greeting
	if role == auth.RoleXO {
		stats := h.store.StudentStats("")
		qualStats := h.store.QualStats()
		atRisk := h.store.AtRiskStudents("")
		today := nowET().Format("Monday, January 2, 2006")

		if len(message) < 30 || containsAny(msg, "today", "status", "brief", "morning", "what") {
			return fmt.Sprintf("## Good morning, sir. Heywood online.\n\n"+
				"**Morning Brief — %s**\n\n"+
				"### Company Status\n"+
				"- **Active Students:** %d\n"+
				"- **Average Composite:** %.1f\n"+
				"- **At-Risk:** %d (%.1f%%)\n\n"+
				"### Instructor Quals\n"+
				"- **Expired:** %d | **Critical (30d):** %d | **Warning (60d):** %d\n"+
				"- **Coverage Gaps:** %d qualifications below minimum staffing\n\n"+
				"### At-Risk Students Requiring Attention\n"+
				"Top concerns:\n%s\n"+
				"### Recommendations\n"+
				"1. Prioritize counseling for students with declining trends\n"+
				"2. Address the %d expired instructor qualifications before next week's graded events\n"+
				"3. Review coverage gaps to ensure upcoming ranges are adequately staffed\n\n"+
				"Anything else you'd like to drill into, sir?\n\n"+
				"*AI-generated analysis — verify all data before taking action.*",
				today,
				stats.ActiveStudents, stats.AvgComposite,
				stats.AtRiskCount, stats.AtRiskPercent,
				qualStats.ExpiredCount, qualStats.Expiring30, qualStats.Expiring60,
				len(qualStats.CoverageGaps),
				formatTopAtRisk(atRisk, 5),
				qualStats.ExpiredCount)
		}
	}

	// Existing mock logic for other roles
	if containsAny(msg, "at risk", "at-risk", "struggling", "failing") {
		var atRisk []models.Student
		if role == auth.RoleSPC {
			atRisk = h.store.AtRiskStudents(company)
		} else {
			atRisk = h.store.AtRiskStudents("")
		}
		return ai.MockAtRiskResponse(atRisk)
	}

	if containsAny(msg, "counseling", "counsel") {
		sid := extractStudentID(msg)
		if sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				return ai.MockCounselingResponse(st)
			}
		}
		return "I can prepare a counseling outline for any student. Please specify the student ID (e.g., \"Prepare counseling for STU-042\")."
	}

	if containsAny(msg, "scenario", "mett-tc", "tactical") {
		phase := "Phase 2"
		if strings.Contains(msg, "phase 1") || strings.Contains(msg, "phase i") {
			phase = "Phase 1"
		} else if strings.Contains(msg, "phase 3") || strings.Contains(msg, "phase iii") {
			phase = "Phase 3"
		}
		terrain := "wooded/hilly"
		if strings.Contains(msg, "urban") {
			terrain = "urban"
		} else if strings.Contains(msg, "desert") {
			terrain = "desert/open"
		}
		objective := "conduct a deliberate attack"
		if strings.Contains(msg, "defense") || strings.Contains(msg, "defend") {
			objective = "establish a defensive position"
		} else if strings.Contains(msg, "patrol") || strings.Contains(msg, "recon") {
			objective = "conduct a reconnaissance patrol"
		}
		return ai.MockScenarioResponse(phase, objective, terrain)
	}

	if containsAny(msg, "qual", "certification", "expir") && (role == auth.RoleStaff || role == auth.RoleXO) {
		qs := h.store.QualStats()
		var b strings.Builder
		fmt.Fprintf(&b, "**Instructor Qualification Status:**\n\n")
		fmt.Fprintf(&b, "- **Expired:** %d qualifications need immediate renewal\n", qs.ExpiredCount)
		fmt.Fprintf(&b, "- **Critical (30 days):** %d expiring soon\n", qs.Expiring30)
		fmt.Fprintf(&b, "- **Warning (60 days):** %d approaching expiration\n", qs.Expiring60)
		fmt.Fprintf(&b, "- **Caution (90 days):** %d to plan for\n", qs.Expiring90)
		fmt.Fprintf(&b, "- **Current:** %d qualifications in good standing\n\n", qs.CurrentCount)
		if len(qs.CoverageGaps) > 0 {
			b.WriteString("**Coverage Gaps (below minimum required):**\n")
			for _, g := range qs.CoverageGaps {
				fmt.Fprintf(&b, "- %s: %d qualified / %d required (**gap: %d**)\n", g.QualName, g.QualifiedCount, g.RequiredCount, g.Gap)
			}
		}
		b.WriteString("\n*AI-generated analysis. Verify all data before taking action.*")
		return b.String()
	}

	// Exam results — student asks about test performance
	if containsAny(msg, "exam", "test", "quiz", "study", "missed", "wrong", "score") {
		// Figure out which student
		sid := ""
		if role == auth.RoleStudent && studentID != "" {
			sid = studentID
		} else if s := extractStudentID(msg); s != "" {
			sid = s
		}

		if sid != "" {
			// Try to determine exam number from message
			examNum := 0
			for i := 1; i <= 4; i++ {
				if strings.Contains(msg, fmt.Sprintf("exam %d", i)) || strings.Contains(msg, fmt.Sprintf("exam%d", i)) || strings.Contains(msg, fmt.Sprintf("test %d", i)) {
					examNum = i
					break
				}
			}

			st, ok := h.store.GetStudent(sid)
			if ok {
				// If no specific exam mentioned, find their most recent exam with a score
				if examNum == 0 {
					if st.Exam4 > 0 {
						examNum = 4
					} else if st.Exam3 > 0 {
						examNum = 3
					} else if st.Exam2 > 0 {
						examNum = 2
					} else {
						examNum = 1
					}
				}

				results := h.store.GetExamResults(sid, examNum)
				if results != nil {
					var b strings.Builder
					fmt.Fprintf(&b, "## Exam %d Results — %s %s\n\n", examNum, st.Rank, st.LastName)
					fmt.Fprintf(&b, "**Score: %.1f%%** (%d/%d correct)\n\n", results.Score, results.Correct, results.Total)

					// Group by topic
					topicCorrect := make(map[string]int)
					topicTotal := make(map[string]int)
					var weakTopics []string
					for _, q := range results.Questions {
						topicTotal[q.Topic]++
						if q.Correct {
							topicCorrect[q.Topic]++
						}
					}

					b.WriteString("### Performance by Topic\n\n")
					b.WriteString("| Topic | Score | Status |\n|-------|-------|--------|\n")
					for topic, total := range topicTotal {
						correct := topicCorrect[topic]
						pct := float64(correct) / float64(total) * 100
						status := "Strong"
						if pct < 60 {
							status = "**Needs Work**"
							weakTopics = append(weakTopics, topic)
						} else if pct < 80 {
							status = "Fair"
							weakTopics = append(weakTopics, topic)
						}
						fmt.Fprintf(&b, "| %s | %d/%d (%.0f%%) | %s |\n", topic, correct, total, pct, status)
					}

					if len(weakTopics) > 0 {
						b.WriteString("\n### Study Recommendations\n\n")
						b.WriteString("Focus your study time on these areas:\n")
						for _, t := range weakTopics {
							fmt.Fprintf(&b, "- **%s** — review your course materials and reference publications for this topic\n", t)
						}
					} else {
						b.WriteString("\nStrong performance across all topic areas. Keep up the good work.\n")
					}

					b.WriteString("\n*Note: Specific test questions and answers cannot be shared. This analysis shows topic-level performance only.*")
					return b.String()
				}
				return fmt.Sprintf("No detailed exam results on file for %s %s yet. Their exam scores are: Exam 1: %.0f, Exam 2: %.0f, Exam 3: %.0f, Exam 4: %.0f.",
					st.Rank, st.LastName, st.Exam1, st.Exam2, st.Exam3, st.Exam4)
			}
		}

		if role == auth.RoleStudent {
			return "I can look up your exam results and show you which topic areas to focus on. Try asking \"How did I do on Exam 1?\" or \"What should I study for?\""
		}
		return "I can look up exam results for any student. Specify the student ID, e.g. \"How did STU-042 do on Exam 1?\""
	}

	if containsAny(msg, "how", "overall", "status", "summary", "doing") {
		stats := h.store.StudentStats(company)
		return fmt.Sprintf("**Company Status Overview:**\n\n"+
			"- **Active Students:** %d\n"+
			"- **Average Overall Composite:** %.1f\n"+
			"- **At-Risk Students:** %d (%.1f%%)\n\n"+
			"The company is tracking well overall. %d students are flagged at-risk and should be prioritized for counseling.\n\n"+
			"Would you like me to drill into the at-risk students, review specific individuals, or look at something else?",
			stats.ActiveStudents, stats.AvgComposite,
			stats.AtRiskCount, stats.AtRiskPercent, stats.AtRiskCount)
	}

	if sid := extractStudentID(msg); sid != "" {
		if st, ok := h.store.GetStudent(sid); ok {
			return fmt.Sprintf("**%s — %s**\n\n"+
				"Phase: %s | Status: %s\n\n"+
				"| Pillar | Score | Status |\n"+
				"|--------|-------|--------|\n"+
				"| Academic (32%%) | %.1f | %s |\n"+
				"| Mil Skills (32%%) | %.1f | %s |\n"+
				"| Leadership (36%%) | %.1f | %s |\n"+
				"| **Overall** | **%.1f** | **%s** |\n\n"+
				"Trend: %s | At-Risk: %v\n\n"+
				"Would you like me to prepare a counseling outline for this student?",
				st.ID, st.Rank, st.Phase, st.Status,
				st.AcademicComposite, scoreStatus(st.AcademicComposite),
				st.MilSkillsComposite, scoreStatus(st.MilSkillsComposite),
				st.LeadershipComposite, scoreStatus(st.LeadershipComposite),
				st.OverallComposite, scoreStatus(st.OverallComposite),
				st.Trend, st.AtRisk)
		}
	}

	if len(message) < 20 {
		stats := h.store.StudentStats(company)
		return ai.MockGreeting(role, stats)
	}

	return ai.MockGeneralResponse(message)
}

func formatTopAtRisk(students []models.Student, n int) string {
	if len(students) == 0 {
		return "No students currently at-risk.\n"
	}
	var b strings.Builder
	show := n
	if show > len(students) {
		show = len(students)
	}
	for _, s := range students[:show] {
		flags := strings.Join(s.RiskFlags, ", ")
		if flags == "" {
			flags = "composite/trend"
		}
		fmt.Fprintf(&b, "- **%s** (%s): Overall %.1f, Trend: %s — %s\n", fmt.Sprintf("%s %s, %s", s.Rank, s.LastName, s.FirstName), s.ID, s.OverallComposite, s.Trend, flags)
	}
	if len(students) > n {
		fmt.Fprintf(&b, "- ...and %d more\n", len(students)-n)
	}
	return b.String()
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func scoreStatus(score float64) string {
	if score >= 85 {
		return "Strong"
	}
	if score >= 75 {
		return "Satisfactory"
	}
	return "Below Standard"
}

func formatScheduleSummary(events []models.TrainingEvent) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Training Schedule (%d events):\n\n", len(events))
	for _, e := range events {
		graded := ""
		if e.IsGraded {
			graded = " [GRADED]"
		}
		fmt.Fprintf(&b, "- %s (%s): %s %s-%s at %s%s\n",
			e.Title, e.Code, e.StartDate, e.StartTime, e.EndTime, e.Location, graded)
	}
	return b.String()
}
