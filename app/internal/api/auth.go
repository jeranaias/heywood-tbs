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
	case "xo":
		name = "Executive Officer"
	case "spc":
		name = "Staff Platoon Commander"
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

	validRoles := map[string]bool{"staff": true, "spc": true, "student": true, "xo": true}
	if !validRoles[req.Role] {
		writeError(w, 400, "role must be staff, spc, student, or xo")
		return
	}

	expires := time.Now().Add(24 * time.Hour)

	// STIG-compliant cookie settings (Secure disabled in dev mode for HTTP localhost)
	sameSite := http.SameSiteStrictMode
	if h.dev {
		sameSite = http.SameSiteLaxMode
	}
	for _, c := range []http.Cookie{
		{Name: "heywood-role", Value: req.Role},
		{Name: "heywood-company", Value: req.Company},
		{Name: "heywood-student-id", Value: req.StudentID},
	} {
		http.SetCookie(w, &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     "/",
			Expires:  expires,
			HttpOnly: true,
			Secure:   !h.dev,
			SameSite: sameSite,
		})
	}

	name := "TBS Staff"
	switch req.Role {
	case "xo":
		name = "Executive Officer"
	case "spc":
		name = "Staff Platoon Commander"
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
