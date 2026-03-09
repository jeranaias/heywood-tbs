package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/calendar"
	"heywood-tbs/internal/data"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/msgraph"
)

// Handler holds dependencies for all API handlers.
type Handler struct {
	store           data.DataStore
	chatSvc         *ChatService
	weatherSvc      *ai.WeatherService
	newsSvc         *ai.NewsService
	trafficSvc      *ai.TrafficService
	authProvider    auth.IdentityProvider
	chatRateLimiter *middleware.RateLimiter
	dev             bool // development mode — relaxes Secure cookie flag for HTTP

	// Calendar provider (mock or Outlook via Graph)
	calendarProvider calendar.CalendarProvider

	// Microsoft Graph services
	graphClient   *msgraph.Client
	sharePointSvc *msgraph.SharePointService
	teamsSvc      *msgraph.TeamsService

	// Settings file
	settingsPath string
	settingsMu   sync.RWMutex

	// SSE broker for real-time notifications
	sseBroker *SSEBroker
}

// NewHandler creates a new API handler.
func NewHandler(
	store data.DataStore,
	chatSvc *ChatService,
	weatherSvc *ai.WeatherService,
	newsSvc *ai.NewsService,
	trafficSvc *ai.TrafficService,
	authProvider auth.IdentityProvider,
	dev bool,
	calProvider calendar.CalendarProvider,
	graphClient *msgraph.Client,
	sharePointSvc *msgraph.SharePointService,
	teamsSvc *msgraph.TeamsService,
	settingsPath string,
) *Handler {
	if calProvider == nil {
		calProvider = &calendar.MockCalendar{}
	}
	return &Handler{
		store:            store,
		chatSvc:          chatSvc,
		weatherSvc:       weatherSvc,
		newsSvc:          newsSvc,
		trafficSvc:       trafficSvc,
		authProvider:     authProvider,
		chatRateLimiter:  middleware.NewRateLimiter(5, 10), // 5 req/s, burst 10
		dev:              dev,
		calendarProvider: calProvider,
		graphClient:      graphClient,
		sharePointSvc:    sharePointSvc,
		teamsSvc:         teamsSvc,
		settingsPath:     settingsPath,
		sseBroker:        NewSSEBroker(),
	}
}

// SetupRouter registers all API routes on the given mux.
func SetupRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check (bypasses auth — suitable for k8s/Docker probes)
	mux.HandleFunc("GET /api/v1/healthz", h.handleHealthz)

	// Students
	mux.HandleFunc("GET /api/v1/students", h.handleListStudents)
	mux.HandleFunc("GET /api/v1/students/stats", h.handleStudentStats)
	mux.HandleFunc("GET /api/v1/students/at-risk", h.handleAtRisk)
	mux.HandleFunc("GET /api/v1/students/{id}", h.handleGetStudent)
	mux.HandleFunc("PATCH /api/v1/students/{id}", h.handleUpdateStudent)
	mux.HandleFunc("GET /api/v1/students/{id}/notes", h.handleListStudentNotes)
	mux.HandleFunc("POST /api/v1/students/{id}/notes", h.handleCreateStudentNote)

	// Counseling
	mux.HandleFunc("POST /api/v1/counselings", h.handleCreateCounseling)
	mux.HandleFunc("GET /api/v1/counselings", h.handleListCounselings)
	mux.HandleFunc("GET /api/v1/counselings/{id}", h.handleGetCounseling)
	mux.HandleFunc("PUT /api/v1/counselings/{id}", h.handleUpdateCounseling)
	mux.HandleFunc("POST /api/v1/counselings/generate-outline", h.handleGenerateCounselingOutline)

	// Instructors
	mux.HandleFunc("GET /api/v1/instructors", h.handleListInstructors)
	mux.HandleFunc("GET /api/v1/instructors/{id}", h.handleGetInstructor)

	// Qualifications
	mux.HandleFunc("GET /api/v1/qualifications", h.handleListQualifications)
	mux.HandleFunc("GET /api/v1/qual-records", h.handleListQualRecords)
	mux.HandleFunc("GET /api/v1/qual-records/stats", h.handleQualStats)

	// Schedule
	mux.HandleFunc("GET /api/v1/schedule", h.handleListSchedule)
	mux.HandleFunc("POST /api/v1/schedule", h.handleCreateTrainingEvent)
	mux.HandleFunc("PUT /api/v1/schedule/{id}", h.handleUpdateTrainingEvent)
	mux.HandleFunc("DELETE /api/v1/schedule/{id}", h.handleDeleteTrainingEvent)

	// Export / Reports
	mux.HandleFunc("GET /api/v1/export/students", h.handleExportStudents)
	mux.HandleFunc("GET /api/v1/export/at-risk", h.handleExportAtRisk)
	mux.HandleFunc("GET /api/v1/export/qual-records", h.handleExportQualRecords)
	mux.HandleFunc("GET /api/v1/export/counselings", h.handleExportCounselings)
	mux.HandleFunc("GET /api/v1/reports/company-summary", h.handleCompanyPerformanceSummary)

	// Feedback
	mux.HandleFunc("GET /api/v1/feedback", h.handleListFeedback)

	// Chat (stricter rate limit: 5 req/s per IP)
	chatHandler := h.chatRateLimiter.Middleware(http.HandlerFunc(h.handleChat))
	mux.Handle("POST /api/v1/chat", chatHandler)

	// Chat history (requires SQL-backed store)
	mux.HandleFunc("GET /api/v1/chat/sessions", h.handleListChatSessions)
	mux.HandleFunc("GET /api/v1/chat/sessions/{id}/messages", h.handleGetChatMessages)
	mux.HandleFunc("GET /api/v1/chat/sessions/{id}/export", h.handleExportChatSession)
	mux.HandleFunc("DELETE /api/v1/chat/sessions/{id}", h.handleDeleteChatSession)

	// Chat suggestions (role-specific prompts)
	mux.HandleFunc("GET /api/v1/chat/suggestions", h.handleSuggestedPrompts)

	// Tasks
	mux.HandleFunc("GET /api/v1/tasks", h.handleListTasks)
	mux.HandleFunc("POST /api/v1/tasks", h.handleCreateTask)
	mux.HandleFunc("GET /api/v1/tasks/{id}", h.handleGetTask)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}", h.handleUpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", h.handleDeleteTask)

	// Messages
	mux.HandleFunc("GET /api/v1/messages", h.handleListMessages)
	mux.HandleFunc("POST /api/v1/messages/{id}/read", h.handleMarkMessageRead)

	// Notifications
	mux.HandleFunc("GET /api/v1/notifications", h.handleListNotifications)
	mux.HandleFunc("GET /api/v1/notifications/count", h.handleNotificationCount)
	mux.HandleFunc("POST /api/v1/notifications/{id}/read", h.handleMarkNotificationRead)

	// Auth
	mux.HandleFunc("GET /api/v1/auth/me", h.handleAuthMe)
	mux.HandleFunc("POST /api/v1/auth/switch", h.handleAuthSwitch)

	// Calendar
	mux.HandleFunc("GET /api/v1/calendar/events", h.handleCalendarEvents)
	mux.HandleFunc("POST /api/v1/calendar/events", h.handleCreateCalendarEvent)
	mux.HandleFunc("GET /api/v1/calendar/today", h.handleCalendarToday)
	mux.HandleFunc("GET /api/v1/mail/summary", h.handleMailSummary)
	mux.HandleFunc("GET /api/v1/mail/unread-count", h.handleMailUnreadCount)
	mux.HandleFunc("POST /api/v1/mail/send", h.handleSendMail)
	mux.HandleFunc("POST /api/v1/mail/{id}/reply", h.handleReplyMail)
	mux.HandleFunc("POST /api/v1/calendar/events/{id}/respond", h.handleRespondEvent)

	// Server-Sent Events (real-time notifications)
	mux.HandleFunc("GET /api/v1/events/stream", h.handleSSE)

	// Settings (XO/Staff only — enforced in handlers)
	mux.HandleFunc("GET /api/v1/settings", h.handleGetSettings)
	mux.HandleFunc("PUT /api/v1/settings", h.handleUpdateSettings)
	mux.HandleFunc("POST /api/v1/settings/test-connection", h.handleTestConnection)
	mux.HandleFunc("POST /api/v1/settings/upload", h.handleUpload)
	mux.HandleFunc("POST /api/v1/settings/column-map", h.handleColumnMap)
	mux.HandleFunc("POST /api/v1/settings/upload/preview", h.handleUploadPreview)
	mux.HandleFunc("GET /api/v1/settings/system-info", h.handleSystemInfo)

	// Microsoft Graph integrations (XO/Staff only)
	mux.HandleFunc("POST /api/v1/graph/test", h.handleGraphTest)

	// SharePoint
	mux.HandleFunc("GET /api/v1/sharepoint/site", h.handleSharePointSite)
	mux.HandleFunc("GET /api/v1/sharepoint/lists", h.handleSharePointLists)
	mux.HandleFunc("GET /api/v1/sharepoint/list-items", h.handleSharePointListItems)
	mux.HandleFunc("GET /api/v1/sharepoint/drives", h.handleSharePointDrives)
	mux.HandleFunc("GET /api/v1/sharepoint/files", h.handleSharePointFiles)

	// Teams
	mux.HandleFunc("GET /api/v1/teams", h.handleTeamsList)
	mux.HandleFunc("GET /api/v1/teams/channels", h.handleTeamsChannels)
	mux.HandleFunc("GET /api/v1/teams/files", h.handleTeamsFiles)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
