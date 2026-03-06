package api

import (
	"net/http"

	"heywood-tbs/internal/middleware"
)

func (h *Handler) handleListInstructors(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == "student" {
		writeError(w, 403, "access denied")
		return
	}

	company := r.URL.Query().Get("company")
	instructors := h.store.ListInstructors(company)
	writeJSON(w, 200, map[string]interface{}{
		"instructors": instructors,
		"total":       len(instructors),
	})
}

func (h *Handler) handleGetInstructor(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "staff" && role != "xo" {
		writeError(w, 403, "access denied")
		return
	}

	id := r.PathValue("id")
	inst, ok := h.store.GetInstructor(id)
	if !ok {
		writeError(w, 404, "instructor not found")
		return
	}
	writeJSON(w, 200, inst)
}

func (h *Handler) handleListQualifications(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "staff" && role != "xo" {
		writeError(w, 403, "access denied")
		return
	}
	writeJSON(w, 200, h.store.ListQualifications())
}

func (h *Handler) handleListQualRecords(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "staff" && role != "xo" {
		writeError(w, 403, "access denied")
		return
	}
	writeJSON(w, 200, h.store.ListQualRecords())
}

func (h *Handler) handleQualStats(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "staff" && role != "xo" {
		writeError(w, 403, "access denied")
		return
	}
	stats := h.store.QualStats()
	writeJSON(w, 200, stats)
}
