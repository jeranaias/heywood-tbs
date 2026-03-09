package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"heywood-tbs/internal/middleware"
)

// handleExportStudents exports all students as CSV.
// GET /api/v1/export/students
func (h *Handler) handleExportStudents(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot export roster data")
		return
	}

	company, _ := r.Context().Value(middleware.CompanyKey).(string)
	students := h.store.ListStudents(company, "", "", false)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=students.csv")

	cw := csv.NewWriter(w)
	cw.Write([]string{
		"ID", "Rank", "Last Name", "First Name", "Company", "Platoon", "Phase",
		"Academic", "Mil Skills", "Leadership", "Overall", "Trend", "At Risk", "Risk Flags",
	})

	for _, s := range students {
		atRisk := "No"
		if s.AtRisk {
			atRisk = "Yes"
		}
		cw.Write([]string{
			s.ID, s.Rank, s.LastName, s.FirstName, s.Company, s.Platoon, s.Phase,
			fmt.Sprintf("%.1f", s.AcademicComposite),
			fmt.Sprintf("%.1f", s.MilSkillsComposite),
			fmt.Sprintf("%.1f", s.LeadershipComposite),
			fmt.Sprintf("%.1f", s.OverallComposite),
			s.Trend, atRisk, strings.Join(s.RiskFlags, "; "),
		})
	}
	cw.Flush()
}

// handleExportAtRisk exports at-risk students as CSV.
// GET /api/v1/export/at-risk
func (h *Handler) handleExportAtRisk(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot export at-risk data")
		return
	}

	company, _ := r.Context().Value(middleware.CompanyKey).(string)
	students := h.store.AtRiskStudents(company)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=at-risk-students.csv")

	cw := csv.NewWriter(w)
	cw.Write([]string{
		"ID", "Rank", "Last Name", "First Name", "Company", "Platoon", "Phase",
		"Academic", "Mil Skills", "Leadership", "Overall", "Trend", "Risk Flags",
	})

	for _, s := range students {
		cw.Write([]string{
			s.ID, s.Rank, s.LastName, s.FirstName, s.Company, s.Platoon, s.Phase,
			fmt.Sprintf("%.1f", s.AcademicComposite),
			fmt.Sprintf("%.1f", s.MilSkillsComposite),
			fmt.Sprintf("%.1f", s.LeadershipComposite),
			fmt.Sprintf("%.1f", s.OverallComposite),
			s.Trend, strings.Join(s.RiskFlags, "; "),
		})
	}
	cw.Flush()
}

// handleExportQualRecords exports qualification records as CSV.
// GET /api/v1/export/qual-records
func (h *Handler) handleExportQualRecords(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot export qualification data")
		return
	}

	records := h.store.ListQualRecords()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=qual-records.csv")

	cw := csv.NewWriter(w)
	cw.Write([]string{
		"ID", "Instructor Name", "EDIPI", "Qual Code", "Qual Name",
		"Date Earned", "Expiration Date", "Days Until Expiration", "Status", "Renewal Status",
	})

	for _, r := range records {
		cw.Write([]string{
			r.ID, r.InstructorName, r.InstructorEDIPI, r.QualCode, r.QualName,
			r.DateEarned, r.ExpirationDate,
			fmt.Sprintf("%d", r.DaysUntilExpiration),
			r.ExpirationStatus, r.RenewalStatus,
		})
	}
	cw.Flush()
}

// handleExportCounselings exports counseling sessions as CSV.
// GET /api/v1/export/counselings
func (h *Handler) handleExportCounselings(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot export counseling data")
		return
	}

	sessions := h.store.ListCounselings("")

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=counseling-sessions.csv")

	cw := csv.NewWriter(w)
	cw.Write([]string{
		"ID", "Student ID", "Student Name", "Counselor", "Date", "Type", "Status",
		"Follow-Ups Pending", "Created At",
	})

	for _, s := range sessions {
		pending := 0
		for _, fu := range s.FollowUps {
			if fu.Status == "pending" {
				pending++
			}
		}
		cw.Write([]string{
			s.ID, s.StudentID, s.StudentName, s.CounselorName, s.Date,
			s.Type, s.Status, fmt.Sprintf("%d", pending), s.CreatedAt,
		})
	}
	cw.Flush()
}

// handleCompanyPerformanceSummary returns per-company aggregated stats.
// GET /api/v1/reports/company-summary
func (h *Handler) handleCompanyPerformanceSummary(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == "student" {
		writeError(w, http.StatusForbidden, "students cannot view company reports")
		return
	}

	// Gather students by company
	students := h.store.ListStudents("", "", "", false)
	companies := map[string]*companyStats{}

	for _, s := range students {
		co := s.Company
		if co == "" {
			co = "Unassigned"
		}
		cs, ok := companies[co]
		if !ok {
			cs = &companyStats{Name: co}
			companies[co] = cs
		}
		cs.Count++
		cs.TotalAcademic += s.AcademicComposite
		cs.TotalMilSkills += s.MilSkillsComposite
		cs.TotalLeadership += s.LeadershipComposite
		cs.TotalOverall += s.OverallComposite
		if s.AtRisk {
			cs.AtRiskCount++
		}
	}

	var result []companyStatsResponse
	for _, cs := range companies {
		n := float64(cs.Count)
		result = append(result, companyStatsResponse{
			Company:      cs.Name,
			StudentCount: cs.Count,
			AvgAcademic:  cs.TotalAcademic / n,
			AvgMilSkills: cs.TotalMilSkills / n,
			AvgLeadership: cs.TotalLeadership / n,
			AvgOverall:   cs.TotalOverall / n,
			AtRiskCount:  cs.AtRiskCount,
			AtRiskPct:    float64(cs.AtRiskCount) / n * 100,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"companies": result,
		"total":     len(students),
	})
}

type companyStats struct {
	Name           string
	Count          int
	TotalAcademic  float64
	TotalMilSkills float64
	TotalLeadership float64
	TotalOverall   float64
	AtRiskCount    int
}

type companyStatsResponse struct {
	Company       string  `json:"company"`
	StudentCount  int     `json:"studentCount"`
	AvgAcademic   float64 `json:"avgAcademic"`
	AvgMilSkills  float64 `json:"avgMilSkills"`
	AvgLeadership float64 `json:"avgLeadership"`
	AvgOverall    float64 `json:"avgOverall"`
	AtRiskCount   int     `json:"atRiskCount"`
	AtRiskPct     float64 `json:"atRiskPct"`
}
