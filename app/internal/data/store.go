package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"heywood-tbs/internal/models"
)

// Store holds all TBS data in memory with indexed lookups.
type Store struct {
	Students       []models.Student
	Instructors    []models.Instructor
	Schedule       []models.TrainingEvent
	Qualifications []models.Qualification
	QualRecords    []models.QualRecord
	Feedback       []models.EventFeedback

	studentByID    map[string]*models.Student
	instructorByID map[string]*models.Instructor
}

// NewStore loads all JSON data files from dataDir into memory.
func NewStore(dataDir string) (*Store, error) {
	s := &Store{
		studentByID:    make(map[string]*models.Student),
		instructorByID: make(map[string]*models.Instructor),
	}

	if err := loadJSON(filepath.Join(dataDir, "students.json"), &s.Students); err != nil {
		return nil, fmt.Errorf("load students: %w", err)
	}
	if err := loadJSON(filepath.Join(dataDir, "instructors.json"), &s.Instructors); err != nil {
		return nil, fmt.Errorf("load instructors: %w", err)
	}
	if err := loadJSON(filepath.Join(dataDir, "schedule.json"), &s.Schedule); err != nil {
		return nil, fmt.Errorf("load schedule: %w", err)
	}
	if err := loadJSON(filepath.Join(dataDir, "qualifications.json"), &s.Qualifications); err != nil {
		return nil, fmt.Errorf("load qualifications: %w", err)
	}
	if err := loadJSON(filepath.Join(dataDir, "qual-records.json"), &s.QualRecords); err != nil {
		return nil, fmt.Errorf("load qual-records: %w", err)
	}
	if err := loadJSON(filepath.Join(dataDir, "feedback.json"), &s.Feedback); err != nil {
		return nil, fmt.Errorf("load feedback: %w", err)
	}

	// Build indexes
	for i := range s.Students {
		s.studentByID[s.Students[i].ID] = &s.Students[i]
	}
	for i := range s.Instructors {
		s.instructorByID[s.Instructors[i].ID] = &s.Instructors[i]
	}

	return s, nil
}

func loadJSON(path string, dest interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// ListStudents returns students filtered by optional criteria.
func (s *Store) ListStudents(company, phase, search string, atRiskOnly bool) []models.Student {
	var result []models.Student
	for _, st := range s.Students {
		if company != "" && !strings.EqualFold(st.Company, company) {
			continue
		}
		if phase != "" && !strings.EqualFold(st.Phase, phase) {
			continue
		}
		if atRiskOnly && !st.AtRisk {
			continue
		}
		if search != "" {
			q := strings.ToLower(search)
			name := strings.ToLower(st.LastName + " " + st.FirstName)
			if !strings.Contains(name, q) && !strings.Contains(strings.ToLower(st.ID), q) {
				continue
			}
		}
		result = append(result, st)
	}
	return result
}

// GetStudent returns a single student by ID.
func (s *Store) GetStudent(id string) (*models.Student, bool) {
	st, ok := s.studentByID[id]
	return st, ok
}

// StudentStats computes aggregate KPIs for students, optionally filtered by company.
func (s *Store) StudentStats(company string) models.StudentStats {
	stats := models.StudentStats{
		ByPhase:         make(map[string]int),
		ByStandingThird: make(map[string]int),
	}
	var totalComposite float64
	for _, st := range s.Students {
		if company != "" && !strings.EqualFold(st.Company, company) {
			continue
		}
		if st.Status != "Active" && st.Status != "" {
			// Include all for now; could filter to Active only
		}
		stats.ActiveStudents++
		totalComposite += st.OverallComposite
		if st.AtRisk {
			stats.AtRiskCount++
		}
		stats.ByPhase[st.Phase]++
		if st.ClassStandingThird != "" {
			stats.ByStandingThird[st.ClassStandingThird]++
		}
	}
	if stats.ActiveStudents > 0 {
		stats.AvgComposite = totalComposite / float64(stats.ActiveStudents)
		stats.AtRiskPercent = float64(stats.AtRiskCount) / float64(stats.ActiveStudents) * 100
	}
	return stats
}

// AtRiskStudents returns only students flagged as at-risk.
func (s *Store) AtRiskStudents(company string) []models.Student {
	return s.ListStudents(company, "", "", true)
}

// GetInstructor returns a single instructor by ID.
func (s *Store) GetInstructor(id string) (*models.Instructor, bool) {
	inst, ok := s.instructorByID[id]
	return inst, ok
}

// ListInstructors returns instructors, optionally filtered by company.
func (s *Store) ListInstructors(company string) []models.Instructor {
	if company == "" {
		return s.Instructors
	}
	var result []models.Instructor
	for _, inst := range s.Instructors {
		if strings.EqualFold(inst.Company, company) {
			result = append(result, inst)
		}
	}
	return result
}

// QualStats computes qualification KPI aggregates.
func (s *Store) QualStats() models.QualStats {
	stats := models.QualStats{
		TotalRecords: len(s.QualRecords),
	}

	for _, qr := range s.QualRecords {
		switch {
		case strings.Contains(strings.ToLower(qr.ExpirationStatus), "expired"):
			stats.ExpiredCount++
		case strings.Contains(strings.ToLower(qr.ExpirationStatus), "critical"):
			stats.Expiring30++
		case strings.Contains(strings.ToLower(qr.ExpirationStatus), "warning"):
			stats.Expiring60++
		case strings.Contains(strings.ToLower(qr.ExpirationStatus), "caution"):
			stats.Expiring90++
		default:
			stats.CurrentCount++
		}
	}

	// Compute coverage gaps: for each qualification that has a minimum requirement,
	// count how many current instructors hold it.
	qualMinimums := make(map[string]int)    // qualCode -> minimumPerEvent
	qualNames := make(map[string]string)    // qualCode -> name
	qualCurrent := make(map[string]int)     // qualCode -> count of current holders

	for _, q := range s.Qualifications {
		if q.MinimumPerEvent > 0 {
			qualMinimums[q.Code] = q.MinimumPerEvent
			qualNames[q.Code] = q.Name
		}
	}

	for _, qr := range s.QualRecords {
		status := strings.ToLower(qr.ExpirationStatus)
		if strings.Contains(status, "current") || strings.Contains(status, "caution") {
			qualCurrent[qr.QualCode]++
		}
	}

	for code, required := range qualMinimums {
		current := qualCurrent[code]
		if current < required {
			stats.CoverageGaps = append(stats.CoverageGaps, models.CoverageGap{
				QualCode:       code,
				QualName:       qualNames[code],
				QualifiedCount: current,
				RequiredCount:  required,
				Gap:            required - current,
			})
		}
	}

	return stats
}

// ListSchedule returns training events, optionally filtered by phase.
func (s *Store) ListSchedule(phase string) []models.TrainingEvent {
	if phase == "" {
		return s.Schedule
	}
	var result []models.TrainingEvent
	for _, evt := range s.Schedule {
		if strings.EqualFold(evt.Phase, phase) {
			result = append(result, evt)
		}
	}
	return result
}

// ListFeedback returns feedback, optionally filtered by event code.
func (s *Store) ListFeedback(eventCode string) []models.EventFeedback {
	if eventCode == "" {
		return s.Feedback
	}
	var result []models.EventFeedback
	for _, fb := range s.Feedback {
		if strings.EqualFold(fb.EventCode, eventCode) {
			result = append(result, fb)
		}
	}
	return result
}
