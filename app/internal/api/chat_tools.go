package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"heywood-tbs/internal/models"

	openai "github.com/sashabaranov/go-openai"
)

// executeToolCall dispatches a tool call to the appropriate store method and returns the result as a string.
func (h *Handler) executeToolCall(tc openai.ToolCall, callerRole, callerCompany string) string {
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
	case "lookup_calendar":
		return h.toolLookupCalendar(args, callerRole, callerCompany)
	case "schedule_event":
		return h.toolScheduleEvent(args, callerRole, callerCompany)
	case "setup_guide":
		return h.toolSetupGuide(args)
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

func (h *Handler) toolLookupCalendar(args map[string]interface{}, role, company string) string {
	dateStr, _ := args["date"].(string)
	query, _ := args["query"].(string)

	now := time.Now()
	var start, end time.Time

	switch dateStr {
	case "today", "":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)
	case "tomorrow":
		start = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)
	case "this week":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// Go to end of week (Sunday)
		daysUntilSunday := 7 - int(start.Weekday())
		end = start.AddDate(0, 0, daysUntilSunday+1)
	default:
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return "Error: invalid date format. Use YYYY-MM-DD, 'today', 'tomorrow', or 'this week'."
		}
		start = parsed
		end = parsed.AddDate(0, 0, 1)
	}

	events := h.calendarProvider.GetEvents(role, company, start, end)

	// Also merge TBS schedule
	scheduleEvents := h.scheduleToCalendarEvents(role, company, start, end)
	events = append(events, scheduleEvents...)

	// Filter by query if provided
	if query != "" {
		queryLower := strings.ToLower(query)
		var filtered []models.CalendarEvent
		for _, e := range events {
			if strings.Contains(strings.ToLower(e.Title), queryLower) ||
				strings.Contains(strings.ToLower(e.Category), queryLower) ||
				strings.Contains(strings.ToLower(e.Description), queryLower) {
				filtered = append(filtered, e)
			}
		}
		events = filtered
	}

	if len(events) == 0 {
		return fmt.Sprintf("No calendar events found for %s.", dateStr)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Calendar Events (%d found):\n\n", len(events))
	for _, e := range events {
		loc := ""
		if e.Location != "" {
			loc = " at " + e.Location
		}
		source := ""
		if e.Source == "outlook" {
			source = " [Outlook]"
		}
		fmt.Fprintf(&b, "- %s: %s%s%s\n", e.Start, e.Title, loc, source)
	}
	return b.String()
}

func (h *Handler) toolScheduleEvent(args map[string]interface{}, role, company string) string {
	title, _ := args["title"].(string)
	dateStr, _ := args["date"].(string)
	startTime, _ := args["start_time"].(string)
	endTime, _ := args["end_time"].(string)
	location, _ := args["location"].(string)
	description, _ := args["description"].(string)
	category, _ := args["category"].(string)

	if title == "" || startTime == "" || endTime == "" {
		return "Error: title, start_time, and end_time are required."
	}

	// Resolve date
	now := time.Now()
	var eventDate string
	switch strings.ToLower(dateStr) {
	case "today", "":
		eventDate = now.Format("2006-01-02")
	case "tomorrow":
		eventDate = now.AddDate(0, 0, 1).Format("2006-01-02")
	case "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday":
		// Find next occurrence of that day
		target := dayNameToWeekday(dateStr)
		d := now
		for i := 1; i <= 7; i++ {
			d = d.AddDate(0, 0, 1)
			if d.Weekday() == target {
				eventDate = d.Format("2006-01-02")
				break
			}
		}
		if eventDate == "" {
			eventDate = now.AddDate(0, 0, 1).Format("2006-01-02")
		}
	default:
		// Try parsing as ISO date
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return "Error: invalid date. Use YYYY-MM-DD, 'today', 'tomorrow', or a day name like 'friday'."
		}
		eventDate = dateStr
	}

	if category == "" {
		category = "admin"
	}

	event := models.CalendarEvent{
		Title:       title,
		Start:       eventDate + "T" + startTime + ":00",
		End:         eventDate + "T" + endTime + ":00",
		Location:    location,
		Description: description,
		Source:      "outlook",
		Category:    category,
		Company:     company,
	}

	created, err := h.calendarProvider.CreateEvent(event)
	if err != nil {
		return fmt.Sprintf("Error creating event: %s", err)
	}

	loc := ""
	if location != "" {
		loc = " at " + location
	}
	return fmt.Sprintf("Event created: \"%s\" on %s from %s to %s%s (ID: %s)",
		created.Title, eventDate, startTime, endTime, loc, created.ID)
}

func dayNameToWeekday(name string) time.Weekday {
	switch strings.ToLower(name) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return time.Monday
	}
}

func (h *Handler) toolSetupGuide(args map[string]interface{}) string {
	topic, _ := args["topic"].(string)

	// Load current settings for context
	settings := h.loadCurrentSettings()
	aiConfigured := os.Getenv("OPENAI_API_KEY") != "" || os.Getenv("AZURE_OPENAI_KEY") != ""
	studentCount := len(h.store.ListStudents("", "", "", false))

	var b strings.Builder

	switch topic {
	case "overview":
		b.WriteString("HEYWOOD SETUP STATUS:\n\n")
		b.WriteString(fmt.Sprintf("1. Student Data: %s\n", statusLabel(studentCount > 0)))
		if studentCount > 0 {
			b.WriteString(fmt.Sprintf("   - %d students loaded from %s\n", studentCount, settings.DataSource.Type))
		} else {
			b.WriteString("   - No students loaded. Upload an Excel roster or use demo data.\n")
		}
		b.WriteString(fmt.Sprintf("2. AI Assistant: %s\n", statusLabel(aiConfigured)))
		if aiConfigured {
			b.WriteString("   - AI is active and ready for chat.\n")
		} else {
			b.WriteString("   - Not configured. Chat uses placeholder responses. Ask your S-6 to set the OPENAI_API_KEY environment variable on the server.\n")
		}
		b.WriteString(fmt.Sprintf("3. Outlook: %s\n", statusLabel(settings.Outlook.Enabled && settings.Outlook.TenantID != "")))
		if settings.Outlook.Enabled && settings.Outlook.TenantID != "" {
			b.WriteString("   - Outlook calendar and mail sync is connected.\n")
		} else {
			b.WriteString("   - Not connected. Calendar uses demo events. Connect Outlook to see real calendar and mail.\n")
		}
		b.WriteString("\nTo configure any of these, go to the Settings page (gear icon in sidebar) or ask me about a specific topic: 'How do I upload my roster?', 'How do I connect Outlook?', etc.")

	case "student-data":
		b.WriteString("HOW TO LOAD YOUR STUDENT ROSTER:\n\n")
		b.WriteString(fmt.Sprintf("Current: %d students from %s source.\n\n", studentCount, settings.DataSource.Type))
		b.WriteString("OPTION A — Upload Excel (Recommended for most units):\n")
		b.WriteString("1. Go to Settings (gear icon in sidebar)\n")
		b.WriteString("2. Under 'Student Data', click 'Upload Excel Roster'\n")
		b.WriteString("3. Upload your .xlsx or .csv file with student data\n")
		b.WriteString("4. Heywood automatically maps your column headers (Last Name, EDIPI, Platoon, etc.)\n")
		b.WriteString("5. Review the preview, then click Save\n\n")
		b.WriteString("Your spreadsheet should have columns like: Last Name, First Name, Rank, EDIPI, Company, Platoon, Academics Score, etc.\n")
		b.WriteString("Heywood recognizes 50+ common column name variations used by Marine units.\n\n")
		b.WriteString("OPTION B — SharePoint (for units with IT support):\n")
		b.WriteString("Ask your S-6 to create an Azure App Registration with Sites.Selected permission, then enter the credentials in Settings > Advanced > SharePoint.\n\n")
		b.WriteString("OPTION C — Database (for production):\n")
		b.WriteString("Enter a Cosmos DB, PostgreSQL, or SQL Server connection string in Settings > Advanced > Database.")

	case "ai":
		b.WriteString("HOW TO SET UP THE AI ASSISTANT:\n\n")
		b.WriteString(fmt.Sprintf("Current status: %s\n\n", statusLabel(aiConfigured)))
		if aiConfigured {
			b.WriteString("AI is already configured and active. The chat assistant is fully operational.\n\n")
		}
		b.WriteString("The AI assistant requires an API key set as an environment variable on the server.\n\n")
		b.WriteString("FOR OPENAI (most common):\n")
		b.WriteString("1. Get an API key from platform.openai.com\n")
		b.WriteString("2. Set the environment variable: OPENAI_API_KEY=sk-your-key-here\n")
		b.WriteString("3. Restart Heywood\n\n")
		b.WriteString("FOR AZURE OPENAI (MCEN/government):\n")
		b.WriteString("1. Get credentials from your Azure OpenAI resource\n")
		b.WriteString("2. Set environment variables:\n")
		b.WriteString("   AZURE_OPENAI_KEY=your-key\n")
		b.WriteString("   AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com\n")
		b.WriteString("   AZURE_OPENAI_DEPLOYMENT=your-deployment-name\n")
		b.WriteString("3. Restart Heywood\n\n")
		b.WriteString("NOTE: These are server-side environment variables. Your S-6 or IT admin needs to set them where Heywood runs (Azure Container App, VM, etc.). They cannot be set from the web UI for security.")

	case "outlook":
		b.WriteString("HOW TO CONNECT OUTLOOK:\n\n")
		b.WriteString(fmt.Sprintf("Current status: %s\n\n", statusLabel(settings.Outlook.Enabled && settings.Outlook.TenantID != "")))
		b.WriteString("Connecting Outlook lets Heywood show your calendar events and recent mail alongside the TBS training schedule.\n\n")
		b.WriteString("WHAT YOUR S-6/IT ADMIN NEEDS TO DO:\n")
		b.WriteString("1. Go to Azure Portal > Azure Active Directory > App Registrations\n")
		b.WriteString("2. Create a new registration (name: 'Heywood TBS')\n")
		b.WriteString("3. Add API permissions: Microsoft Graph > Application:\n")
		b.WriteString("   - Calendars.Read (read calendar events)\n")
		b.WriteString("   - Mail.Read (read mail summaries)\n")
		b.WriteString("4. Grant admin consent for the permissions\n")
		b.WriteString("5. Go to Certificates & Secrets > New client secret\n")
		b.WriteString("6. Copy these three values:\n")
		b.WriteString("   - Tenant ID (from Overview page)\n")
		b.WriteString("   - Client ID / Application ID (from Overview page)\n")
		b.WriteString("   - Client Secret (the value you just created)\n\n")
		b.WriteString("WHAT YOU DO:\n")
		b.WriteString("1. Go to Settings (gear icon in sidebar)\n")
		b.WriteString("2. Under 'Outlook Mail & Calendar', toggle it ON\n")
		b.WriteString("3. Click 'Enter connection details'\n")
		b.WriteString("4. Paste the Tenant ID, Client ID, and Client Secret\n")
		b.WriteString("5. Select your network (Commercial for most, GCC High for MCEN)\n")
		b.WriteString("6. Click Save Changes\n\n")
		b.WriteString("Until connected, the Calendar page uses realistic demo events so you can see how it will look.")

	case "sharepoint":
		b.WriteString("HOW TO CONNECT SHAREPOINT:\n\n")
		b.WriteString("SharePoint integration lets Heywood read data directly from your unit's SharePoint lists instead of uploaded files.\n\n")
		b.WriteString("PREREQUISITES:\n")
		b.WriteString("- Your unit has a SharePoint site with student/roster data in Lists\n")
		b.WriteString("- Your S-6 can create an Azure App Registration\n\n")
		b.WriteString("SETUP STEPS:\n")
		b.WriteString("1. S-6 creates Azure App Registration with Sites.Selected permission\n")
		b.WriteString("2. S-6 grants the app access to your specific SharePoint site\n")
		b.WriteString("3. Go to Settings > Student Data > Advanced > SharePoint\n")
		b.WriteString("4. Enter Tenant ID, Client ID, Client Secret, and Site URL\n")
		b.WriteString("5. Select network (Commercial or GCC High for MCEN)\n")
		b.WriteString("6. Click Test Connection to verify\n")
		b.WriteString("7. Save Changes\n\n")
		b.WriteString("Most units should start with Excel upload instead — it's simpler and doesn't need IT support.")

	case "database":
		b.WriteString("HOW TO CONNECT A DATABASE:\n\n")
		b.WriteString("Database connections are for production deployments where you need persistent, shared data storage.\n\n")
		b.WriteString("SUPPORTED DATABASES:\n")
		b.WriteString("- Azure Cosmos DB (recommended for Azure deployments)\n")
		b.WriteString("- PostgreSQL (self-hosted or Azure Database for PostgreSQL)\n")
		b.WriteString("- Azure SQL / SQL Server\n\n")
		b.WriteString("SETUP:\n")
		b.WriteString("1. Your DBA or cloud admin provisions the database\n")
		b.WriteString("2. Get the connection string from them\n")
		b.WriteString("3. Go to Settings > Student Data > Advanced > Database\n")
		b.WriteString("4. Select database type and paste connection string\n")
		b.WriteString("5. Click Test Connection\n")
		b.WriteString("6. Save Changes\n\n")
		b.WriteString("NOTE: Most units don't need a database. Excel upload or demo data works great for evaluation and smaller deployments.")

	default:
		b.WriteString("I can help with setup for: student data, AI assistant, Outlook, SharePoint, or database. What would you like to configure?")
	}

	return b.String()
}

func statusLabel(ok bool) string {
	if ok {
		return "CONFIGURED"
	}
	return "NOT CONFIGURED"
}

func (h *Handler) loadCurrentSettings() AppSettings {
	data, err := os.ReadFile(h.settingsPath)
	if err != nil {
		return AppSettings{}
	}
	var s AppSettings
	json.Unmarshal(data, &s)
	return s
}
