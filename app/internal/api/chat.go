package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"

	openai "github.com/sashabaranov/go-openai"
)

// ChatService manages OpenAI API interactions.
type ChatService struct {
	client *openai.Client
}

// NewChatService creates a chat service. If apiKey is empty, returns nil (mock mode).
func NewChatService(apiKey string) *ChatService {
	if apiKey == "" {
		return nil
	}
	client := openai.NewClient(apiKey)
	return &ChatService{client: client}
}

func (h *Handler) handleChat(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		writeError(w, 400, "message is required")
		return
	}

	role := middleware.GetRole(r.Context())
	company := middleware.GetCompany(r.Context())
	studentID := middleware.GetStudentID(r.Context())

	// Build system prompt and context based on role
	systemPrompt, userContext := h.buildChatContext(role, company, studentID, req.Message)

	// If no OpenAI client, use mock responses
	if h.chatSvc == nil {
		response := h.mockResponse(role, company, studentID, req.Message)
		writeJSON(w, 200, models.ChatResponse{Response: response})
		return
	}

	// Build messages for OpenAI
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Add history (limited to last 10 exchanges)
	historyStart := 0
	if len(req.History) > 20 {
		historyStart = len(req.History) - 20
	}
	for _, msg := range req.History[historyStart:] {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add current message with injected context
	userMessage := req.Message
	if userContext != "" {
		userMessage = req.Message + "\n\n---\n[Relevant data context]\n" + userContext
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: userMessage,
	})

	// Call OpenAI
	resp, err := h.chatSvc.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT4o,
			Messages:    messages,
			MaxTokens:   1500,
			Temperature: 0.7,
		},
	)
	if err != nil {
		slog.Error("openai error", "error", err)
		// Fall back to mock on API error
		response := h.mockResponse(role, company, studentID, req.Message)
		writeJSON(w, 200, models.ChatResponse{Response: response + "\n\n*(Note: AI service temporarily unavailable — showing cached response)*"})
		return
	}

	if len(resp.Choices) == 0 {
		writeError(w, 500, "no response from AI")
		return
	}

	writeJSON(w, 200, models.ChatResponse{
		Response: resp.Choices[0].Message.Content,
	})
}

// buildChatContext creates the system prompt and relevant data context for the chat.
func (h *Handler) buildChatContext(role, company, studentID, message string) (systemPrompt, userContext string) {
	msg := strings.ToLower(message)

	switch role {
	case "staff":
		stats := h.store.StudentStats("")
		systemPrompt = ai.StaffSystemPrompt(stats)

		// Inject relevant data based on query intent
		if containsAny(msg, "at risk", "at-risk", "struggling", "failing", "concern") {
			atRisk := h.store.AtRiskStudents("")
			userContext = ai.AnonymizeStudentList(atRisk)
		} else if containsAny(msg, "counseling", "counsel") {
			sid := extractStudentID(msg)
			if sid != "" {
				if st, ok := h.store.GetStudent(sid); ok {
					userContext = ai.AnonymizeStudent(st)
				}
			}
		} else if containsAny(msg, "qual", "certification", "expir") {
			qs := h.store.QualStats()
			userContext = fmt.Sprintf("Qualification Status:\nTotal records: %d\nExpired: %d\nExpiring 30 days: %d\nExpiring 60 days: %d\nExpiring 90 days: %d\nCurrent: %d\n",
				qs.TotalRecords, qs.ExpiredCount, qs.Expiring30, qs.Expiring60, qs.Expiring90, qs.CurrentCount)
			if len(qs.CoverageGaps) > 0 {
				userContext += "\nCoverage Gaps:\n"
				for _, g := range qs.CoverageGaps {
					userContext += fmt.Sprintf("- %s: %d qualified / %d required (gap: %d)\n", g.QualName, g.QualifiedCount, g.RequiredCount, g.Gap)
				}
			}
		} else if containsAny(msg, "schedule", "training", "event", "calendar") {
			events := h.store.ListSchedule("")
			userContext = formatScheduleSummary(events)
		} else if containsAny(msg, "how", "overall", "status", "summary", "company") {
			// General overview — include stats and at-risk summary
			atRisk := h.store.AtRiskStudents("")
			if len(atRisk) > 5 {
				atRisk = atRisk[:5]
			}
			userContext = ai.AnonymizeStudentList(atRisk)
		} else if sid := extractStudentID(msg); sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				userContext = ai.AnonymizeStudent(st)
			}
		}

	case "spc":
		stats := h.store.StudentStats(company)
		systemPrompt = ai.SPCSystemPrompt(stats, company)

		if containsAny(msg, "at risk", "at-risk", "struggling") {
			atRisk := h.store.AtRiskStudents(company)
			userContext = ai.AnonymizeStudentList(atRisk)
		} else if containsAny(msg, "counseling", "counsel") {
			sid := extractStudentID(msg)
			if sid != "" {
				if st, ok := h.store.GetStudent(sid); ok {
					userContext = ai.AnonymizeStudent(st)
				}
			}
		} else if sid := extractStudentID(msg); sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				userContext = ai.AnonymizeStudent(st)
			}
		}

	case "student":
		var student *models.Student
		if studentID != "" {
			student, _ = h.store.GetStudent(studentID)
		}
		systemPrompt = ai.StudentSystemPrompt(student)
		// Student context is already in the system prompt
	}

	return systemPrompt, userContext
}

// mockResponse generates a mock response when no OpenAI API key is configured.
func (h *Handler) mockResponse(role, company, studentID, message string) string {
	msg := strings.ToLower(message)

	// Check for specific intents
	if containsAny(msg, "at risk", "at-risk", "struggling", "failing") {
		var atRisk []models.Student
		if role == "spc" {
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

	if containsAny(msg, "qual", "certification", "expir") && role == "staff" {
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
		b.WriteString("\n*This is AI-generated analysis. Verify all data before taking action.*")
		return b.String()
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

	// Check for specific student ID
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

	// Short messages are likely greetings
	if len(message) < 20 {
		stats := h.store.StudentStats(company)
		return ai.MockGreeting(role, stats)
	}

	return ai.MockGeneralResponse(message)
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func extractStudentID(msg string) string {
	// Look for STU-XXX pattern
	msg = strings.ToUpper(msg)
	idx := strings.Index(msg, "STU-")
	if idx >= 0 && idx+7 <= len(msg) {
		return msg[idx : idx+7]
	}
	// Look for "student #XX" or "student XX" pattern
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

