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

	// Tools are available for all roles
	tools := ai.HeywoodTools

	// Streaming mode
	if req.Stream {
		h.handleStreamingChat(w, r, messages, tools, role)
		return
	}

	// Non-streaming mode
	chatReq := openai.ChatCompletionRequest{
		Model:       h.chatSvc.model,
		Messages:    messages,
		MaxTokens:   4096,
		Temperature: 0.7,
	}
	if len(tools) > 0 {
		chatReq.Tools = tools
	}

	resp, err := h.chatSvc.client.CreateChatCompletion(r.Context(), chatReq)
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

	// Handle tool calls — execute tools and re-send for final answer
	choice := resp.Choices[0]
	if choice.FinishReason == openai.FinishReasonToolCalls && len(choice.Message.ToolCalls) > 0 {
		messages = append(messages, choice.Message)
		for _, tc := range choice.Message.ToolCalls {
			result := h.executeToolCall(tc, role)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}
		// Re-send without tools to get final response
		resp2, err := h.chatSvc.client.CreateChatCompletion(r.Context(), openai.ChatCompletionRequest{
			Model:       h.chatSvc.model,
			Messages:    messages,
			MaxTokens:   4096,
			Temperature: 0.7,
		})
		if err != nil {
			slog.Error("openai tool follow-up error", "error", err)
			writeError(w, 500, "AI follow-up failed")
			return
		}
		if len(resp2.Choices) > 0 {
			writeJSON(w, 200, models.ChatResponse{Response: resp2.Choices[0].Message.Content})
			return
		}
	}

	writeJSON(w, 200, models.ChatResponse{
		Response: choice.Message.Content,
	})
}

// handleStreamingChat sends the response as Server-Sent Events.
// If the model returns tool calls, they are executed synchronously and the
// final response is then streamed.
func (h *Handler) handleStreamingChat(w http.ResponseWriter, r *http.Request, messages []openai.ChatCompletionMessage, tools []openai.Tool, role string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, 500, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	// If tools are available, do a non-streaming call first to check for tool calls.
	// Tool calls are fast (small response) and we need the full response before streaming.
	if len(tools) > 0 {
		resp, err := h.chatSvc.client.CreateChatCompletion(r.Context(), openai.ChatCompletionRequest{
			Model:       h.chatSvc.model,
			Messages:    messages,
			MaxTokens:   4096,
			Temperature: 0.7,
			Tools:       tools,
		})
		if err != nil {
			slog.Error("tool check error", "error", err)
			// Fall through to normal streaming without tools
		} else if len(resp.Choices) > 0 && resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
			// Execute tool calls
			choice := resp.Choices[0]
			messages = append(messages, choice.Message)
			for _, tc := range choice.Message.ToolCalls {
				result := h.executeToolCall(tc, role)
				slog.Info("tool call executed", "tool", tc.Function.Name, "id", tc.ID)
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				})
			}
			// Now stream the final response (no tools on follow-up)
			h.streamMessages(w, r, flusher, messages)
			return
		} else if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
			// No tool calls — model responded directly.
			// Simulate streaming by sending small character chunks (preserves all formatting).
			content := resp.Choices[0].Message.Content
			runes := []rune(content)
			chunkSize := 6
			for i := 0; i < len(runes); i += chunkSize {
				end := i + chunkSize
				if end > len(runes) {
					end = len(runes)
				}
				chunk := string(runes[i:end])
				data, _ := json.Marshal(map[string]string{"content": chunk})
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
				time.Sleep(12 * time.Millisecond)
			}
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}
	}

	// Standard streaming (no tools or tool check failed)
	h.streamMessages(w, r, flusher, messages)
}

// streamMessages opens a streaming connection and sends SSE chunks.
func (h *Handler) streamMessages(w http.ResponseWriter, r *http.Request, flusher http.Flusher, messages []openai.ChatCompletionMessage) {
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
				time.Sleep(8 * time.Millisecond)
			}
		}
	}
}

// executeToolCall dispatches a tool call to the appropriate store method and returns the result as a string.
func (h *Handler) executeToolCall(tc openai.ToolCall, callerRole string) string {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	switch tc.Function.Name {
	case "create_task":
		return h.toolCreateTask(args, callerRole)
	case "send_message":
		return h.toolSendMessage(args, callerRole)
	case "lookup_student":
		return h.toolLookupStudent(args)
	case "lookup_schedule":
		return h.toolLookupSchedule(args)
	case "web_search":
		return h.toolWebSearch(args)
	case "lookup_exam_results":
		return h.toolLookupExamResults(args)
	default:
		return fmt.Sprintf("Unknown tool: %s", tc.Function.Name)
	}
}

func (h *Handler) toolCreateTask(args map[string]interface{}, callerRole string) string {
	title, _ := args["title"].(string)
	desc, _ := args["description"].(string)
	assignedTo, _ := args["assigned_to"].(string)
	priority, _ := args["priority"].(string)
	dueDate, _ := args["due_date"].(string)
	relatedID, _ := args["related_id"].(string)

	if priority == "" {
		priority = "medium"
	}

	task := models.Task{
		Title:       title,
		Description: desc,
		AssignedTo:  assignedTo,
		CreatedBy:   "heywood",
		Priority:    priority,
		DueDate:     dueDate,
		RelatedID:   relatedID,
	}

	if err := h.store.CreateTask(task); err != nil {
		return fmt.Sprintf("Failed to create task: %v", err)
	}

	// Also create a notification for the assignee
	_ = h.store.CreateNotification(models.Notification{
		UserRole:  assignedTo,
		Type:      "task",
		Title:     "New Task: " + title,
		Body:      fmt.Sprintf("Heywood has assigned you a %s-priority task: %s", priority, title),
		ActionURL: "/tasks",
	})

	return fmt.Sprintf("Task created successfully. Assigned to %s with %s priority. Notification sent.", assignedTo, priority)
}

func (h *Handler) toolSendMessage(args map[string]interface{}, callerRole string) string {
	to, _ := args["to"].(string)
	subject, _ := args["subject"].(string)
	body, _ := args["body"].(string)
	relatedID, _ := args["related_id"].(string)

	msg := models.Message{
		From:      "heywood (on behalf of " + callerRole + ")",
		To:        to,
		Subject:   subject,
		Body:      body,
		RelatedID: relatedID,
	}

	if err := h.store.CreateMessage(msg); err != nil {
		return fmt.Sprintf("Failed to send message: %v", err)
	}

	_ = h.store.CreateNotification(models.Notification{
		UserRole: to,
		Type:     "message",
		Title:    "Message: " + subject,
		Body:     "From: Heywood (XO) — " + subject,
	})

	return fmt.Sprintf("Message sent to %s. Subject: %s. Notification delivered.", to, subject)
}

func (h *Handler) toolLookupStudent(args map[string]interface{}) string {
	query, _ := args["query"].(string)
	if query == "" {
		return "No query provided"
	}

	// Try exact ID first
	if st, ok := h.store.GetStudent(strings.ToUpper(query)); ok {
		return formatStudentForTool(st)
	}

	// Search by name
	students := h.store.ListStudents("", "", query, false)
	if len(students) == 0 {
		return fmt.Sprintf("No students found matching '%s'", query)
	}
	if len(students) == 1 {
		return formatStudentForTool(&students[0])
	}

	// Multiple matches — return summary
	var b strings.Builder
	fmt.Fprintf(&b, "Found %d students matching '%s':\n", len(students), query)
	for _, s := range students {
		fmt.Fprintf(&b, "- %s %s, %s (%s): Overall %.1f, %s\n", s.Rank, s.LastName, s.FirstName, s.ID, s.OverallComposite, s.Trend)
	}
	return b.String()
}

func formatStudentForTool(st *models.Student) string {
	flags := strings.Join(st.RiskFlags, ", ")
	if flags == "" {
		flags = "none"
	}
	return fmt.Sprintf("Student: %s %s, %s (%s)\n"+
		"Company: %s | Platoon: %s | Phase: %s | SPC: %s\n"+
		"Academic: %.1f (Exams: %.0f, %.0f, %.0f, %.0f | Quiz: %.1f)\n"+
		"Mil Skills: %.1f (PFT: %d, CFT: %d, Rifle: %s, Pistol: %s)\n"+
		"Leadership: %.1f (Wk12: %.1f, Wk22: %.1f, PeerWk12: %.1f, PeerWk22: %.1f)\n"+
		"Overall: %.1f | Trend: %s | At-Risk: %v | Flags: %s",
		st.Rank, st.LastName, st.FirstName, st.ID,
		st.Company, st.Platoon, st.Phase, st.SPC,
		st.AcademicComposite, st.Exam1, st.Exam2, st.Exam3, st.Exam4, st.QuizAvg,
		st.MilSkillsComposite, st.PFTScore, st.CFTScore, st.RifleQual, st.PistolQual,
		st.LeadershipComposite, st.LeadershipWeek12, st.LeadershipWeek22, st.PeerEvalWeek12, st.PeerEvalWeek22,
		st.OverallComposite, st.Trend, st.AtRisk, flags)
}

func (h *Handler) toolLookupSchedule(args map[string]interface{}) string {
	date, _ := args["date"].(string)
	scope, _ := args["scope"].(string)

	if date == "" {
		date = nowET().Format("2006-01-02")
	}
	if scope == "" {
		scope = "day"
	}

	var events []models.TrainingEvent
	if scope == "week" {
		events = h.store.ThisWeekSchedule(date)
	} else {
		events = h.store.TodaySchedule(date)
	}

	if len(events) == 0 {
		return fmt.Sprintf("No training events found for %s (%s)", date, scope)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Training schedule for %s (%s) — %d events:\n", date, scope, len(events))
	for _, e := range events {
		graded := ""
		if e.IsGraded {
			graded = " [GRADED]"
		}
		fmt.Fprintf(&b, "- %s %s–%s: %s (%s)%s at %s | Lead: %s\n",
			e.StartDate, e.StartTime, e.EndTime, e.Title, e.Code, graded, e.Location, e.LeadInstructor)
	}
	return b.String()
}

func (h *Handler) toolWebSearch(args map[string]interface{}) string {
	query, _ := args["query"].(string)
	if query == "" {
		return "No search query provided"
	}

	// SearXNG instance — sidecar or env-configured URL
	searxURL := os.Getenv("SEARXNG_URL")
	if searxURL == "" {
		searxURL = "http://localhost:8888"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", searxURL+"/search", nil)
	if err != nil {
		return fmt.Sprintf("Search error: %v", err)
	}

	q := req.URL.Query()
	q.Set("q", query)
	q.Set("format", "json")
	q.Set("categories", "general")
	q.Set("language", "en")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Search failed (SearXNG unreachable): %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Sprintf("Search returned %d: %s", resp.StatusCode, string(body[:min(200, len(body))]))
	}

	var result struct {
		Results []struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			URL     string `json:"url"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Sprintf("Failed to parse search results: %v", err)
	}

	if len(result.Results) == 0 {
		return fmt.Sprintf("No results found for '%s'", query)
	}

	// Cap at 5 results
	show := 5
	if show > len(result.Results) {
		show = len(result.Results)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Search results for '%s':\n\n", query)
	for i, r := range result.Results[:show] {
		fmt.Fprintf(&b, "%d. **%s**\n   %s\n   Source: %s\n\n", i+1, r.Title, r.Content, r.URL)
	}
	return b.String()
}

func (h *Handler) toolLookupExamResults(args map[string]interface{}) string {
	studentID, _ := args["student_id"].(string)
	examNumF, _ := args["exam_number"].(float64)
	examNum := int(examNumF)

	if studentID == "" {
		return "No student_id provided"
	}
	if examNum < 1 || examNum > 4 {
		return "exam_number must be 1-4"
	}

	st, ok := h.store.GetStudent(studentID)
	if !ok {
		return fmt.Sprintf("Student %s not found", studentID)
	}

	results := h.store.GetExamResults(studentID, examNum)
	if results == nil {
		return fmt.Sprintf("No Exam %d results on file for %s %s", examNum, st.Rank, st.LastName)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Exam %d Results for %s %s, %s\n", examNum, st.Rank, st.LastName, st.FirstName)
	fmt.Fprintf(&b, "Score: %.1f%% (%d/%d correct)\n\n", results.Score, results.Correct, results.Total)
	fmt.Fprintf(&b, "IMPORTANT: Do NOT reveal specific test questions or correct answers to the student.\n")
	fmt.Fprintf(&b, "Instead, identify topic areas where they struggled and provide study guidance.\n\n")

	// Group by topic
	topicCorrect := make(map[string]int)
	topicTotal := make(map[string]int)
	for _, q := range results.Questions {
		topicTotal[q.Topic]++
		if q.Correct {
			topicCorrect[q.Topic]++
		}
	}

	fmt.Fprintf(&b, "Performance by Topic Area:\n")
	for topic, total := range topicTotal {
		correct := topicCorrect[topic]
		pct := float64(correct) / float64(total) * 100
		status := "STRONG"
		if pct < 60 {
			status = "NEEDS WORK"
		} else if pct < 80 {
			status = "FAIR"
		}
		fmt.Fprintf(&b, "- %s: %d/%d (%.0f%%) — %s\n", topic, correct, total, pct, status)
	}

	return b.String()
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
				ctxParts = append(ctxParts, fmt.Sprintf("- %s (%s): Overall %.1f, Trend: %s%s", fmt.Sprintf("%s %s, %s", s.Rank, s.LastName, s.FirstName), s.ID, s.OverallComposite, s.Trend, flags))
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
				ctxParts = append(ctxParts, fmt.Sprintf("- %s (%s): Overall %.1f, Trend: %s%s", fmt.Sprintf("%s %s, %s", s.Rank, s.LastName, s.FirstName), s.ID, s.OverallComposite, s.Trend, flags))
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
