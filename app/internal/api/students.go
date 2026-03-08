package api

import (
	"net/http"

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
