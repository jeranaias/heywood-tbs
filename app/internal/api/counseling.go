package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

// handleCreateCounseling creates a new counseling session.
// POST /api/v1/counselings
func (h *Handler) handleCreateCounseling(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot create counseling sessions")
		return
	}

	var session models.CounselingSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if session.StudentID == "" {
		writeError(w, http.StatusBadRequest, "studentId is required")
		return
	}

	// Verify student exists and populate name
	student, ok := h.store.GetStudent(session.StudentID)
	if !ok {
		writeError(w, http.StatusNotFound, "student not found")
		return
	}
	session.StudentName = fmt.Sprintf("%s %s, %s", student.Rank, student.LastName, student.FirstName)

	// Set counselor info from auth context
	session.CounselorRole = role
	company, _ := r.Context().Value(middleware.CompanyKey).(string)
	session.CounselorName = strings.Title(role)
	if company != "" {
		session.CounselorName = company + " " + session.CounselorName
	}

	if session.Type == "" {
		session.Type = "initial"
	}
	if session.Status == "" {
		session.Status = "draft"
	}

	if err := h.store.CreateCounseling(session); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create counseling session")
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

// handleListCounselings returns counseling sessions, optionally filtered by studentId.
// GET /api/v1/counselings?studentId=...
func (h *Handler) handleListCounselings(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot view counseling sessions")
		return
	}

	studentID := r.URL.Query().Get("studentId")
	sessions := h.store.ListCounselings(studentID)
	if sessions == nil {
		sessions = []models.CounselingSession{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sessions": sessions,
		"total":    len(sessions),
	})
}

// handleGetCounseling returns a single counseling session.
// GET /api/v1/counselings/{id}
func (h *Handler) handleGetCounseling(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot view counseling sessions")
		return
	}

	id := r.PathValue("id")
	session, ok := h.store.GetCounseling(id)
	if !ok {
		writeError(w, http.StatusNotFound, "counseling session not found")
		return
	}
	writeJSON(w, http.StatusOK, session)
}

// handleUpdateCounseling updates an existing counseling session.
// PUT /api/v1/counselings/{id}
func (h *Handler) handleUpdateCounseling(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot modify counseling sessions")
		return
	}

	id := r.PathValue("id")
	existing, ok := h.store.GetCounseling(id)
	if !ok {
		writeError(w, http.StatusNotFound, "counseling session not found")
		return
	}

	var update models.CounselingSession
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Preserve immutable fields
	update.StudentID = existing.StudentID
	update.StudentName = existing.StudentName

	if err := h.store.UpdateCounseling(id, update); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update counseling session")
		return
	}

	updated, _ := h.store.GetCounseling(id)
	writeJSON(w, http.StatusOK, updated)
}

// handleGenerateCounselingOutline generates an AI counseling outline from student data.
// POST /api/v1/counselings/generate-outline
func (h *Handler) handleGenerateCounselingOutline(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot generate counseling outlines")
		return
	}

	var req struct {
		StudentID string `json:"studentId"`
		Type      string `json:"type"` // counseling type for context
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.StudentID == "" {
		writeError(w, http.StatusBadRequest, "studentId is required")
		return
	}

	student, ok := h.store.GetStudent(req.StudentID)
	if !ok {
		writeError(w, http.StatusNotFound, "student not found")
		return
	}

	// If no AI service, return a template outline
	if h.chatSvc == nil {
		outline := generateTemplateOutline(student, req.Type)
		writeJSON(w, http.StatusOK, map[string]string{"outline": outline})
		return
	}

	// Build student data prompt + counseling suffix
	studentPrompt := ai.StudentSystemPrompt(student)
	counselingType := req.Type
	if counselingType == "" {
		counselingType = "progress"
	}
	fullPrompt := fmt.Sprintf("Generate a %s counseling outline for this student.%s\n\n%s",
		counselingType, ai.CounselingPromptSuffix, studentPrompt)

	writeJSON(w, http.StatusOK, map[string]string{"outline": fullPrompt})
}

// generateTemplateOutline creates a data-driven outline without AI.
func generateTemplateOutline(s *models.Student, counselingType string) string {
	if counselingType == "" {
		counselingType = "Progress"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s Counseling — %s %s, %s\n\n", strings.Title(counselingType), s.Rank, s.LastName, s.FirstName)
	fmt.Fprintf(&b, "**Date:** ____________________\n\n")

	// Opening
	fmt.Fprintf(&b, "## 1. Opening Statement\n")
	fmt.Fprintf(&b, "This counseling is being conducted to review your performance at The Basic School. ")
	fmt.Fprintf(&b, "Your current overall composite is **%.1f** and you are trending **%s**.\n\n", s.OverallComposite, s.Trend)

	// Strengths
	fmt.Fprintf(&b, "## 2. Strengths Observed\n")
	pillars := []struct {
		name  string
		score float64
	}{
		{"Academic", s.AcademicComposite},
		{"Military Skills", s.MilSkillsComposite},
		{"Leadership", s.LeadershipComposite},
	}
	for _, p := range pillars {
		if p.score >= 80 {
			fmt.Fprintf(&b, "- **%s (%.1f):** Performing well in this pillar.\n", p.name, p.score)
		}
	}
	fmt.Fprintf(&b, "\n")

	// Areas for improvement
	fmt.Fprintf(&b, "## 3. Areas for Improvement\n")
	for _, p := range pillars {
		if p.score < 75 {
			fmt.Fprintf(&b, "- **%s (%.1f):** Below the 75.0 threshold. Requires focused attention.\n", p.name, p.score)
		}
	}
	if len(s.RiskFlags) > 0 {
		fmt.Fprintf(&b, "- **Risk Flags:** %s\n", strings.Join(s.RiskFlags, ", "))
	}
	fmt.Fprintf(&b, "\n")

	// Actions
	fmt.Fprintf(&b, "## 4. Specific Actions\n")
	fmt.Fprintf(&b, "- [ ] ____________________\n")
	fmt.Fprintf(&b, "- [ ] ____________________\n")
	fmt.Fprintf(&b, "- [ ] ____________________\n\n")

	// Timeline
	fmt.Fprintf(&b, "## 5. Timeline\n")
	fmt.Fprintf(&b, "Reassessment in ______ weeks.\n\n")

	// Closing
	fmt.Fprintf(&b, "## 6. Closing Guidance\n")
	fmt.Fprintf(&b, "*[AI-generated content is a draft requiring human review.]*\n")

	return b.String()
}
