package api

import (
	"encoding/json"
	"net/http"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/data"
)

// Handler holds dependencies for all API handlers.
type Handler struct {
	store        data.DataStore
	chatSvc      *ChatService
	weatherSvc   *ai.WeatherService
	newsSvc      *ai.NewsService
	trafficSvc   *ai.TrafficService
	authProvider auth.IdentityProvider
	dev          bool // development mode — relaxes Secure cookie flag for HTTP
}

// NewHandler creates a new API handler.
func NewHandler(store data.DataStore, chatSvc *ChatService, weatherSvc *ai.WeatherService, newsSvc *ai.NewsService, trafficSvc *ai.TrafficService, authProvider auth.IdentityProvider, dev bool) *Handler {
	return &Handler{store: store, chatSvc: chatSvc, weatherSvc: weatherSvc, newsSvc: newsSvc, trafficSvc: trafficSvc, authProvider: authProvider, dev: dev}
}

// SetupRouter registers all API routes on the given mux.
func SetupRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Students
	mux.HandleFunc("GET /api/v1/students", h.handleListStudents)
	mux.HandleFunc("GET /api/v1/students/stats", h.handleStudentStats)
	mux.HandleFunc("GET /api/v1/students/at-risk", h.handleAtRisk)
	mux.HandleFunc("GET /api/v1/students/{id}", h.handleGetStudent)

	// Instructors
	mux.HandleFunc("GET /api/v1/instructors", h.handleListInstructors)
	mux.HandleFunc("GET /api/v1/instructors/{id}", h.handleGetInstructor)

	// Qualifications
	mux.HandleFunc("GET /api/v1/qualifications", h.handleListQualifications)
	mux.HandleFunc("GET /api/v1/qual-records", h.handleListQualRecords)
	mux.HandleFunc("GET /api/v1/qual-records/stats", h.handleQualStats)

	// Schedule
	mux.HandleFunc("GET /api/v1/schedule", h.handleListSchedule)

	// Feedback
	mux.HandleFunc("GET /api/v1/feedback", h.handleListFeedback)

	// Chat
	mux.HandleFunc("POST /api/v1/chat", h.handleChat)

	// Tasks
	mux.HandleFunc("GET /api/v1/tasks", h.handleListTasks)
	mux.HandleFunc("GET /api/v1/tasks/{id}", h.handleGetTask)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}", h.handleUpdateTask)

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
