package api

import (
	"fmt"
	"log/slog"
	"strings"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/models"
)

// buildChatContext creates the system prompt and relevant data context for the chat.
func (h *Handler) buildChatContext(role, company, studentID, message string) (systemPrompt, userContext string) {
	msg := strings.ToLower(message)

	switch role {
	case auth.RoleXO:
		// XO gets EVERYTHING in the system prompt — full brief mode
		today := nowET().Format("2006-01-02")

		// Fetch live weather
		var weatherStr string
		var weatherData *ai.WeatherData
		if h.weatherSvc != nil {
			if wd, err := h.weatherSvc.Get(); err == nil {
				weatherData = wd
				weatherStr = ai.FormatWeatherForPrompt(wd)
			} else {
				slog.Error("weather fetch failed", "error", err)
				weatherStr = "Weather data temporarily unavailable."
			}
		}

		// Fetch live news headlines
		var newsStr string
		if h.newsSvc != nil {
			if items, err := h.newsSvc.Get(); err == nil {
				newsStr = ai.FormatNewsForPrompt(items)
			} else {
				slog.Error("news fetch failed", "error", err)
			}
		}

		// Calculate real traffic/routes for off-base appointments
		xoSchedule := h.store.XOScheduleForDate(today)
		var trafficStr string
		if h.trafficSvc != nil {
			routes := h.trafficSvc.CalculateRoutes(xoSchedule, weatherData)
			trafficStr = ai.FormatTrafficForPrompt(routes)
		}

		stats := h.store.StudentStats("")
		qualStats := h.store.QualStats()
		atRisk := h.store.AtRiskStudents("")
		todayEvents := h.store.TodaySchedule(today)
		weekEvents := h.store.ThisWeekSchedule(today)
		recentFeedback := h.store.RecentFeedback(10)
		instructors := h.store.ListInstructors("")

		systemPrompt = ai.XOSystemPrompt(
			today, weatherStr, newsStr, trafficStr, stats, qualStats,
			atRisk, todayEvents, weekEvents,
			recentFeedback, instructors, xoSchedule,
		)
		// No userContext needed — everything is in the system prompt
		return systemPrompt, ""

	case auth.RoleStaff:
		stats := h.store.StudentStats("")
		systemPrompt = ai.StaffSystemPrompt(stats)

		// Always inject today's schedule + at-risk summary so Heywood can answer proactively
		today := nowET().Format("2006-01-02")
		var ctxParts []string

		// Today's schedule — always relevant
		todayEvents := h.store.TodaySchedule(today)
		if len(todayEvents) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("Today's date: %s\nToday's schedule:", today))
			for _, e := range todayEvents {
				graded := ""
				if e.IsGraded {
					graded = " [GRADED]"
				}
				ctxParts = append(ctxParts, fmt.Sprintf("- %s–%s: %s (%s)%s at %s | Lead: %s", e.StartTime, e.EndTime, e.Title, e.Code, graded, e.Location, e.LeadInstructor))
			}
		} else {
			ctxParts = append(ctxParts, fmt.Sprintf("Today's date: %s\nNo training events scheduled for today.", today))
		}

		// This week's schedule
		weekEvents := h.store.ThisWeekSchedule(today)
		if len(weekEvents) > 0 {
			ctxParts = append(ctxParts, "\nThis week's remaining events:")
			shown := 0
			for _, e := range weekEvents {
				if e.StartDate == today {
					continue
				}
				graded := ""
				if e.IsGraded {
					graded = " [GRADED]"
				}
				ctxParts = append(ctxParts, fmt.Sprintf("- %s %s–%s: %s (%s)%s at %s | Lead: %s",
					e.StartDate, e.StartTime, e.EndTime, e.Title, e.Code, graded, e.Location, e.LeadInstructor))
				shown++
				if shown >= 15 {
					ctxParts = append(ctxParts, fmt.Sprintf("  ...and %d more", len(weekEvents)-shown))
					break
				}
			}
		}

		// At-risk summary
		atRisk := h.store.AtRiskStudents("")
		if len(atRisk) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("\nAt-risk students: %d total (%.1f%% of class)", len(atRisk), float64(len(atRisk))/float64(stats.ActiveStudents)*100))
			ctxParts = append(ctxParts, "Most critical (lowest overall):")
			show := 10
			if show > len(atRisk) {
				show = len(atRisk)
			}
			for _, s := range atRisk[:show] {
				flags := ""
				if len(s.RiskFlags) > 0 {
					flags = " — " + strings.Join(s.RiskFlags, ", ")
				}
				ctxParts = append(ctxParts, fmt.Sprintf("  %s %s, %s (%s): Acad %.1f | MilSk %.1f | Ldr %.1f | Overall %.1f | Trend: %s%s",
					s.Rank, s.LastName, s.FirstName, s.ID,
					s.AcademicComposite, s.MilSkillsComposite, s.LeadershipComposite,
					s.OverallComposite, s.Trend, flags))
			}
		}

		// Qual alerts — with specifics
		qs := h.store.QualStats()
		if qs.ExpiredCount > 0 || qs.Expiring30 > 0 || len(qs.CoverageGaps) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("\nQual alerts: %d expired, %d critical (30d), %d warning (60d), %d coverage gaps",
				qs.ExpiredCount, qs.Expiring30, qs.Expiring60, len(qs.CoverageGaps)))
			if len(qs.CoverageGaps) > 0 {
				ctxParts = append(ctxParts, "Coverage gaps:")
				for _, g := range qs.CoverageGaps {
					ctxParts = append(ctxParts, fmt.Sprintf("  %s: %d qualified / %d required (gap: %d)", g.QualName, g.QualifiedCount, g.RequiredCount, g.Gap))
				}
			}
		}

		// Specific student lookup if mentioned
		if sid := extractStudentID(msg); sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				ctxParts = append(ctxParts, "\nRequested student detail:\n"+ai.AnonymizeStudent(st))
			}
		}

		userContext = strings.Join(ctxParts, "\n")

	case auth.RoleSPC:
		stats := h.store.StudentStats(company)
		systemPrompt = ai.SPCSystemPrompt(stats, company)

		// Always inject today's schedule + company at-risk
		today := nowET().Format("2006-01-02")
		var ctxParts []string

		todayEvents := h.store.TodaySchedule(today)
		if len(todayEvents) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("Today's date: %s\nToday's schedule:", today))
			for _, e := range todayEvents {
				graded := ""
				if e.IsGraded {
					graded = " [GRADED]"
				}
				ctxParts = append(ctxParts, fmt.Sprintf("- %s–%s: %s (%s)%s at %s | Lead: %s", e.StartTime, e.EndTime, e.Title, e.Code, graded, e.Location, e.LeadInstructor))
			}
		} else {
			ctxParts = append(ctxParts, fmt.Sprintf("Today's date: %s\nNo training events scheduled for today.", today))
		}

		atRisk := h.store.AtRiskStudents(company)
		if len(atRisk) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("\nAt-risk students in %s Company: %d", company, len(atRisk)))
			show := 10
			if show > len(atRisk) {
				show = len(atRisk)
			}
			for _, s := range atRisk[:show] {
				flags := ""
				if len(s.RiskFlags) > 0 {
					flags = " — " + strings.Join(s.RiskFlags, ", ")
				}
				ctxParts = append(ctxParts, fmt.Sprintf("- %s (%s): Overall %.1f, Trend: %s%s", fmt.Sprintf("%s %s, %s", s.Rank, s.LastName, s.FirstName), s.ID, s.OverallComposite, s.Trend, flags))
			}
		}

		if sid := extractStudentID(msg); sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				ctxParts = append(ctxParts, "\nRequested student detail:\n"+ai.AnonymizeStudent(st))
			}
		}

		userContext = strings.Join(ctxParts, "\n")

	case auth.RoleStudent:
		var student *models.Student
		if studentID != "" {
			student, _ = h.store.GetStudent(studentID)
		}
		systemPrompt = ai.StudentSystemPrompt(student)

		// Inject today's schedule + the student's own data
		today := nowET().Format("2006-01-02")
		var ctxParts []string

		todayEvents := h.store.TodaySchedule(today)
		if len(todayEvents) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("Today's date: %s\nToday's training schedule:", today))
			for _, e := range todayEvents {
				graded := ""
				if e.IsGraded {
					graded = " [GRADED]"
				}
				ctxParts = append(ctxParts, fmt.Sprintf("- %s–%s: %s%s at %s", e.StartTime, e.EndTime, e.Title, graded, e.Location))
			}
		} else {
			ctxParts = append(ctxParts, fmt.Sprintf("Today's date: %s\nNo training events scheduled for today.", today))
		}

		if student != nil {
			ctxParts = append(ctxParts, fmt.Sprintf("\nYour performance data (use these numbers in your response):\n"+
				"- Academic: %.1f (Exam1: %.0f, Exam2: %.0f, Exam3: %.0f, Exam4: %.0f | Quiz Avg: %.1f)\n"+
				"- Mil Skills: %.1f (PFT: %d, CFT: %d, Rifle: %s, Pistol: %s, LandNav Day: %s, Night: %s, Written: %.1f)\n"+
				"- Leadership: %.1f (Week12: %.1f, Week22: %.1f | Peer Week12: %.1f, Peer Week22: %.1f)\n"+
				"- Overall: %.1f | Standing: %s | Company Rank: %d\n"+
				"- Trend: %s | At-Risk: %v",
				student.AcademicComposite, student.Exam1, student.Exam2, student.Exam3, student.Exam4, student.QuizAvg,
				student.MilSkillsComposite, student.PFTScore, student.CFTScore, student.RifleQual, student.PistolQual,
				student.LandNavDay, student.LandNavNight, student.LandNavWritten,
				student.LeadershipComposite, student.LeadershipWeek12, student.LeadershipWeek22,
				student.PeerEvalWeek12, student.PeerEvalWeek22,
				student.OverallComposite, student.ClassStandingThird, student.CompanyRank,
				student.Trend, student.AtRisk))
		}

		userContext = strings.Join(ctxParts, "\n")
	}

	return systemPrompt, userContext
}

// extractStudentID attempts to parse a student ID (e.g., "STU-042") from a message string.
func extractStudentID(msg string) string {
	msg = strings.ToUpper(msg)
	idx := strings.Index(msg, "STU-")
	if idx >= 0 && idx+7 <= len(msg) {
		return msg[idx : idx+7]
	}
	for _, prefix := range []string{"STUDENT #", "STUDENT "} {
		idx = strings.Index(msg, prefix)
		if idx >= 0 {
			numStart := idx + len(prefix)
			numEnd := numStart
			for numEnd < len(msg) && msg[numEnd] >= '0' && msg[numEnd] <= '9' {
				numEnd++
			}
			if numEnd > numStart {
				num := msg[numStart:numEnd]
				for len(num) < 3 {
					num = "0" + num
				}
				return "STU-" + num
			}
		}
	}
	return ""
}
