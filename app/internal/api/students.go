package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

func (h *Handler) handleListStudents(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	company := middleware.GetCompany(r.Context())
	studentID := middleware.GetStudentID(r.Context())

	// Student role: can only see their own record
	if role == auth.RoleStudent {
		if studentID == "" {
			writeJSON(w, 200, map[string]interface{}{"students": []interface{}{}, "total": 0, "filtered": 0})
			return
		}
		st, ok := h.store.GetStudent(studentID)
		if !ok {
			writeJSON(w, 200, map[string]interface{}{"students": []interface{}{}, "total": 0, "filtered": 0})
			return
		}
		writeJSON(w, 200, map[string]interface{}{"students": []interface{}{st}, "total": 1, "filtered": 1})
		return
	}

	// SPC role: force company filter
	qCompany := r.URL.Query().Get("company")
	if role == auth.RoleSPC && company != "" {
		qCompany = company
	}

	phase := r.URL.Query().Get("phase")
	search := r.URL.Query().Get("search")
	atRiskOnly := r.URL.Query().Get("atRisk") == "true"

	students := h.store.ListStudents(qCompany, phase, search, atRiskOnly)
	if students == nil {
		writeJSON(w, 200, map[string]interface{}{"students": []models.Student{}, "total": 0, "filtered": 0})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"students": students,
		"total":    h.store.TotalStudentCount(),
		"filtered": len(students),
	})
}

func (h *Handler) handleGetStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	role := middleware.GetRole(r.Context())
	studentID := middleware.GetStudentID(r.Context())

	// Student role: can only see their own record
	if role == auth.RoleStudent && id != studentID {
		writeError(w, 403, "access denied")
		return
	}

	st, ok := h.store.GetStudent(id)
	if !ok {
		writeError(w, 404, "student not found")
		return
	}

	// SPC role: can only see students in their company
	if role == auth.RoleSPC {
		company := middleware.GetCompany(r.Context())
		if company != "" && st.Company != company {
			writeError(w, 403, "access denied")
			return
		}
	}

	writeJSON(w, 200, st)
}

func (h *Handler) handleStudentStats(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	company := ""
	if role == auth.RoleSPC {
		company = middleware.GetCompany(r.Context())
	}
	if qc := r.URL.Query().Get("company"); qc != "" && role == auth.RoleStaff {
		company = qc
	}

	stats := h.store.StudentStats(company)
	writeJSON(w, 200, stats)
}

func (h *Handler) handleAtRisk(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	company := ""
	if role == auth.RoleSPC {
		company = middleware.GetCompany(r.Context())
	}

	students := h.store.AtRiskStudents(company)
	if students == nil {
		writeJSON(w, 200, map[string]interface{}{"students": []interface{}{}, "total": 0})
		return
	}
	writeJSON(w, 200, map[string]interface{}{"students": students, "total": len(students)})
}

func (h *Handler) handleUpdateStudent(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	// Only staff/XO/SPC can update students
	if role == auth.RoleStudent {
		writeError(w, 403, "students cannot modify their own records")
		return
	}

	id := r.PathValue("id")
	var req models.StudentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if err := h.store.UpdateStudent(id, req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, 404, err.Error())
		} else {
			writeError(w, 500, err.Error())
		}
		return
	}

	st, _ := h.store.GetStudent(id)

	// Broadcast SSE if risk status changed
	if h.sseBroker != nil && st != nil && req.AtRisk != nil {
		h.sseBroker.Broadcast("xo", SSEEvent{
			Type: "at-risk-alert",
			Data: map[string]interface{}{
				"studentId":   st.ID,
				"studentName": st.Rank + " " + st.LastName,
				"atRisk":      st.AtRisk,
				"riskFlags":   st.RiskFlags,
				"message":     st.Rank + " " + st.LastName + " risk status updated",
			},
		})
		h.sseBroker.Broadcast("staff", SSEEvent{
			Type: "at-risk-alert",
			Data: map[string]interface{}{
				"studentId":   st.ID,
				"studentName": st.Rank + " " + st.LastName,
				"atRisk":      st.AtRisk,
				"riskFlags":   st.RiskFlags,
				"message":     st.Rank + " " + st.LastName + " risk status updated",
			},
		})
	}

	writeJSON(w, 200, st)
}

func (h *Handler) handleCreateStudentNote(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent {
		writeError(w, 403, "students cannot create notes")
		return
	}

	studentID := r.PathValue("id")
	// Verify student exists
	if _, ok := h.store.GetStudent(studentID); !ok {
		writeError(w, 404, "student not found")
		return
	}

	var body struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if strings.TrimSpace(body.Content) == "" {
		writeError(w, 400, "content is required")
		return
	}

	noteType := body.Type
	if noteType == "" {
		noteType = "note"
	}

	identity := middleware.GetIdentity(r.Context())
	note := models.StudentNote{
		StudentID:  studentID,
		AuthorRole: role,
		AuthorName: identity.Name,
		Content:    body.Content,
		Type:       noteType,
	}

	if err := h.store.CreateStudentNote(note); err != nil {
		writeError(w, 500, "failed to create note")
		return
	}

	// Return the new list of notes
	notes := h.store.ListStudentNotes(studentID)
	writeJSON(w, 201, map[string]interface{}{"notes": notes})
}

func (h *Handler) handleListStudentNotes(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role == auth.RoleStudent {
		writeError(w, 403, "access denied")
		return
	}

	studentID := r.PathValue("id")
	notes := h.store.ListStudentNotes(studentID)
	if notes == nil {
		notes = []models.StudentNote{}
	}
	writeJSON(w, 200, map[string]interface{}{"notes": notes})
}
