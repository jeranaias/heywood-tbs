package api

import (
	"net/http"

	"heywood-tbs/internal/middleware"
)

// suggestedPrompts defines role-specific prompt suggestions for the chat UI.
var suggestedPrompts = map[string][]string{
	"xo": {
		"Morning brief",
		"At-risk students",
		"Qual expiration alerts",
		"This week's schedule",
		"Instructor workload summary",
		"Company performance comparison",
	},
	"staff": {
		"Company performance summary",
		"Schedule this week",
		"At-risk students",
		"Instructor workload",
		"Qual coverage gaps",
		"Recent feedback trends",
	},
	"spc": {
		"My company at-risk",
		"Today's schedule",
		"Company performance summary",
		"Counseling outline for...",
	},
	"student": {
		"How am I doing?",
		"What should I study?",
		"My schedule today",
		"What are my weakest areas?",
	},
}

// handleSuggestedPrompts returns role-specific prompt suggestions.
func (h *Handler) handleSuggestedPrompts(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	prompts, ok := suggestedPrompts[role]
	if !ok {
		prompts = suggestedPrompts["staff"]
	}
	writeJSON(w, 200, map[string]interface{}{"prompts": prompts})
}
