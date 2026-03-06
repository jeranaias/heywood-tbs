package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"

	openai "github.com/sashabaranov/go-openai"
)

// eastern is the US Eastern timezone for Quantico, VA.
var eastern *time.Location

func init() {
	var err error
	eastern, err = time.LoadLocation("America/New_York")
	if err != nil {
		eastern = time.FixedZone("EST", -5*3600)
	}
}

// nowET returns the current time in Eastern timezone.
func nowET() time.Time { return time.Now().In(eastern) }

// ChatService manages OpenAI/Azure OpenAI API interactions.
type ChatService struct {
	client     *openai.Client
	model      string // "gpt-4o" or Azure deployment name
	isAzure    bool
}

// NewChatService creates a chat service with automatic Azure/OpenAI detection.
//
// Env vars checked:
//   - AZURE_OPENAI_ENDPOINT + OPENAI_API_KEY → Azure OpenAI
//   - OPENAI_API_KEY alone → Public OpenAI
//   - Neither → nil (mock mode)
func NewChatService() *ChatService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	azureEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	azureDeployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT")

	if azureEndpoint != "" && apiKey != "" {
		// Azure OpenAI (IL5-ready path)
		config := openai.DefaultAzureConfig(apiKey, azureEndpoint)
		config.APIVersion = "2024-12-01-preview"
		model := "gpt-4o"
		if azureDeployment != "" {
			model = azureDeployment
		}
		slog.Info("Azure OpenAI configured", "endpoint", azureEndpoint, "deployment", model)
		return &ChatService{
			client:  openai.NewClientWithConfig(config),
			model:   model,
			isAzure: true,
		}
	}

	if apiKey != "" {
		// Public OpenAI
		slog.Info("OpenAI configured (public API)")
		return &ChatService{
			client: openai.NewClient(apiKey),
			model:  string(openai.GPT4o),
		}
	}

	return nil // mock mode
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

	// Add history (limited to last 20 messages)
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

	// Streaming mode
	if req.Stream {
		h.handleStreamingChat(w, r, messages)
		return
	}

	// Non-streaming mode
	resp, err := h.chatSvc.client.CreateChatCompletion(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:       h.chatSvc.model,
			Messages:    messages,
			MaxTokens:   4096,
			Temperature: 0.7,
		},
	)
	if err != nil {
		slog.Error("openai error", "error", err)
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

// handleStreamingChat sends the response as Server-Sent Events.
func (h *Handler) handleStreamingChat(w http.ResponseWriter, r *http.Request, messages []openai.ChatCompletionMessage) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, 500, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	stream, err := h.chatSvc.client.CreateChatCompletionStream(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:       h.chatSvc.model,
			Messages:    messages,
			MaxTokens:   4096,
			Temperature: 0.7,
			Stream:      true,
		},
	)
	if err != nil {
		slog.Error("stream create error", "error", err)
		data, _ := json.Marshal(map[string]string{"error": "stream failed"})
		fmt.Fprintf(w, "data: %s\n\n", data)
		fmt.Fprintf(w, "data: [DONE]\n\n")
		flusher.Flush()
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}
		if err != nil {
			slog.Error("stream recv error", "error", err)
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}

		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			if chunk != "" {
				data, _ := json.Marshal(map[string]string{"content": chunk})
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}

// buildChatContext creates the system prompt and relevant data context for the chat.
func (h *Handler) buildChatContext(role, company, studentID, message string) (systemPrompt, userContext string) {
	msg := strings.ToLower(message)

	switch role {
	case "xo":
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

	case "staff":
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

		// At-risk summary
		atRisk := h.store.AtRiskStudents("")
		if len(atRisk) > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("\nAt-risk students: %d total", len(atRisk)))
			show := 10
			if show > len(atRisk) {
				show = len(atRisk)
			}
			for _, s := range atRisk[:show] {
				flags := ""
				if len(s.RiskFlags) > 0 {
					flags = " — " + strings.Join(s.RiskFlags, ", ")
				}
				ctxParts = append(ctxParts, fmt.Sprintf("- %s (%s): Overall %.1f, Trend: %s%s", s.ID, s.Rank, s.OverallComposite, s.Trend, flags))
			}
		}

		// Qual alerts
		qs := h.store.QualStats()
		if qs.ExpiredCount > 0 || qs.Expiring30 > 0 {
			ctxParts = append(ctxParts, fmt.Sprintf("\nQual alerts: %d expired, %d critical (30d), %d coverage gaps", qs.ExpiredCount, qs.Expiring30, len(qs.CoverageGaps)))
		}

		// Specific student lookup if mentioned
		if sid := extractStudentID(msg); sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				ctxParts = append(ctxParts, "\nRequested student detail:\n"+ai.AnonymizeStudent(st))
			}
		}

		userContext = strings.Join(ctxParts, "\n")

	case "spc":
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
				ctxParts = append(ctxParts, fmt.Sprintf("- %s (%s): Overall %.1f, Trend: %s%s", s.ID, s.Rank, s.OverallComposite, s.Trend, flags))
			}
		}

		if sid := extractStudentID(msg); sid != "" {
			if st, ok := h.store.GetStudent(sid); ok {
				ctxParts = append(ctxParts, "\nRequested student detail:\n"+ai.AnonymizeStudent(st))
			}
		}

		userContext = strings.Join(ctxParts, "\n")

	case "student":
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
			ctxParts = append(ctxParts, fmt.Sprintf("\nYour current scores:\n- Academic: %.1f (Exams: %.0f, %.0f, %.0f, %.0f | Quiz Avg: %.1f)\n- Mil Skills: %.1f (PFT: %d, CFT: %d)\n- Leadership: %.1f\n- Overall: %.1f\n- Trend: %s\n- At-Risk: %v",
				student.AcademicComposite, student.Exam1, student.Exam2, student.Exam3, student.Exam4, student.QuizAvg,
				student.MilSkillsComposite, student.PFTScore, student.CFTScore,
				student.LeadershipComposite,
				student.OverallComposite, student.Trend, student.AtRisk))
		}

		userContext = strings.Join(ctxParts, "\n")
	}

	return systemPrompt, userContext
}

// mockResponse generates a mock response when no OpenAI API key is configured.
func (h *Handler) mockResponse(role, company, studentID, message string) string {
	msg := strings.ToLower(message)

	// XO mock: comprehensive greeting
	if role == "xo" {
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

	if containsAny(msg, "qual", "certification", "expir") && (role == "staff" || role == "xo") {
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
		fmt.Fprintf(&b, "- **%s** (%s): Overall %.1f, Trend: %s — %s\n", s.ID, s.Rank, s.OverallComposite, s.Trend, flags)
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
