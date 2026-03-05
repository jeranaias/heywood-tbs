package api

import (
	"encoding/json"
	"net/http"
	"time"

	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

func (h *Handler) handleAuthMe(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	company := middleware.GetCompany(r.Context())
	studentID := middleware.GetStudentID(r.Context())

	name := "TBS Staff"
	switch role {
	case "spc":
		name = "Capt Harris (SPC)"
	case "student":
		if st, ok := h.store.GetStudent(studentID); ok {
			name = st.Rank + " " + st.LastName
		} else {
			name = "Student"
		}
	}

	writeJSON(w, 200, models.AuthInfo{
		Role:      role,
		Company:   company,
		StudentID: studentID,
		Name:      name,
	})
}

type switchRequest struct {
	Role      string `json:"role"`
	Company   string `json:"company"`
	StudentID string `json:"studentId"`
}

func (h *Handler) handleAuthSwitch(w http.ResponseWriter, r *http.Request) {
	var req switchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if req.Role != "staff" && req.Role != "spc" && req.Role != "student" {
		writeError(w, 400, "role must be staff, spc, or student")
		return
	}

	expires := time.Now().Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "heywood-role",
		Value:   req.Role,
		Path:    "/",
		Expires: expires,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "heywood-company",
		Value:   req.Company,
		Path:    "/",
		Expires: expires,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "heywood-student-id",
		Value:   req.StudentID,
		Path:    "/",
		Expires: expires,
	})

	name := "TBS Staff"
	switch req.Role {
	case "spc":
		name = "Capt Harris (SPC)"
	case "student":
		if st, ok := h.store.GetStudent(req.StudentID); ok {
			name = st.Rank + " " + st.LastName
		}
	}

	writeJSON(w, 200, models.AuthInfo{
		Role:      req.Role,
		Company:   req.Company,
		StudentID: req.StudentID,
		Name:      name,
	})
}
