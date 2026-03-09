package api

import (
	"encoding/json"
	"net/http"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
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

// handleCreateTrainingEvent creates a new training event.
// POST /api/v1/schedule (staff/XO only)
func (h *Handler) handleCreateTrainingEvent(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent || role == auth.RoleSPC {
		writeError(w, http.StatusForbidden, "only staff and XO can create training events")
		return
	}

	var event models.TrainingEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if event.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	if err := h.store.CreateTrainingEvent(event); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create training event")
		return
	}
	writeJSON(w, http.StatusCreated, event)
}

// handleUpdateTrainingEvent updates an existing training event.
// PUT /api/v1/schedule/{id}
func (h *Handler) handleUpdateTrainingEvent(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent || role == auth.RoleSPC {
		writeError(w, http.StatusForbidden, "only staff and XO can update training events")
		return
	}

	id := r.PathValue("id")
	var event models.TrainingEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.store.UpdateTrainingEvent(id, event); err != nil {
		writeError(w, http.StatusNotFound, "training event not found")
		return
	}

	event.ID = id
	writeJSON(w, http.StatusOK, event)
}

// handleDeleteTrainingEvent deletes a training event.
// DELETE /api/v1/schedule/{id}
func (h *Handler) handleDeleteTrainingEvent(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent || role == auth.RoleSPC {
		writeError(w, http.StatusForbidden, "only staff and XO can delete training events")
		return
	}

	id := r.PathValue("id")
	if err := h.store.DeleteTrainingEvent(id); err != nil {
		writeError(w, http.StatusNotFound, "training event not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
