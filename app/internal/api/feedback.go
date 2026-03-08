package api

import (
	"net/http"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
)

func (h *Handler) handleListFeedback(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent {
		writeError(w, 403, "access denied")
		return
	}
	eventCode := r.URL.Query().Get("eventCode")
	feedback := h.store.ListFeedback(eventCode)
	writeJSON(w, 200, map[string]interface{}{
		"feedback": feedback,
		"total":    len(feedback),
	})
}
