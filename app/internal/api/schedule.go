package api

import (
	"net/http"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
)

func (h *Handler) handleListSchedule(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent {
		writeError(w, 403, "access denied")
		return
	}
	phase := r.URL.Query().Get("phase")
	events := h.store.ListSchedule(phase)
	writeJSON(w, 200, map[string]interface{}{
		"events": events,
		"total":  len(events),
	})
}
