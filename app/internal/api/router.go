package api

import (
	"encoding/json"
	"net/http"

	"heywood-tbs/internal/data"
)

// Handler holds dependencies for all API handlers.
type Handler struct {
	store   *data.Store
	chatSvc *ChatService
}

// NewHandler creates a new API handler with the given data store and optional chat service.
func NewHandler(store *data.Store, chatSvc *ChatService) *Handler {
	return &Handler{store: store, chatSvc: chatSvc}
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
