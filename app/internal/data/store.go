package data

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"heywood-tbs/internal/models"
)

// Compile-time check that Store implements DataStore.
var _ DataStore = (*Store)(nil)

// Store holds all TBS data in memory with indexed lookups.
type Store struct {
	Students       []models.Student
	Instructors    []models.Instructor
	Schedule       []models.TrainingEvent
	Qualifications []models.Qualification
	QualRecords    []models.QualRecord
	Feedback       []models.EventFeedback
	XOSchedule     []models.XOScheduleItem

	// Exam results (immutable reference data)
	ExamResults []models.ExamResult

	// Mutable data (tasks, messages, notifications)
	Tasks         []models.Task
	Messages      []models.Message
	Notifications []models.Notification

	studentByID    map[string]*models.Student
	instructorByID map[string]*models.Instructor

	mu       sync.RWMutex // guards mutable data
	dataDir  string       // for write-through persistence
	nextID   int          // auto-increment counter for IDs
	demoMode bool         // when true, mutable data is in-memory only (no disk writes)
}

// NewStore loads all JSON data files from dataDir into memory.
// In demo mode, mutable data (tasks, messages, notifications) starts empty
// and is never written to disk — each restart is a fresh slate.
func NewStore(dataDir string) (*Store, error) {
	demoMode := os.Getenv("AUTH_MODE") != "cac"

	s := &Store{
		studentByID:    make(map[string]*models.Student),
		instructorByID: make(map[string]*models.Instructor),
		dataDir:        dataDir,
		nextID:         1,
		demoMode:       demoMode,
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
	// XO schedule is optional — not an error if missing
	_ = loadJSON(filepath.Join(dataDir, "xo-schedule.json"), &s.XOSchedule)

	// Exam results — optional
	_ = loadJSON(filepath.Join(dataDir, "exam-results.json"), &s.ExamResults)

	if demoMode {
		// Demo mode: start fresh — no loading, no disk writes
		slog.Info("demo mode: mutable data is in-memory only (tasks, messages, notifications reset on restart)")
	} else {
		// Production mode: load persisted mutable data from disk
		_ = loadJSON(filepath.Join(dataDir, "tasks.json"), &s.Tasks)
		_ = loadJSON(filepath.Join(dataDir, "messages.json"), &s.Messages)
		_ = loadJSON(filepath.Join(dataDir, "notifications.json"), &s.Notifications)

		// Set nextID based on existing mutable data
		for _, t := range s.Tasks {
			if n := parseIDNum(t.ID); n >= s.nextID {
				s.nextID = n + 1
			}
		}
		for _, m := range s.Messages {
			if n := parseIDNum(m.ID); n >= s.nextID {
				s.nextID = n + 1
			}
		}
		for _, n := range s.Notifications {
			if num := parseIDNum(n.ID); num >= s.nextID {
				s.nextID = num + 1
			}
		}
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

// TodaySchedule returns training events scheduled for the given date (YYYY-MM-DD).
func (s *Store) TodaySchedule(today string) []models.TrainingEvent {
	var result []models.TrainingEvent
	for _, evt := range s.Schedule {
		if evt.StartDate == today {
			result = append(result, evt)
		}
	}
	return result
}

// ThisWeekSchedule returns events in the same Mon-Sun week as the given date.
func (s *Store) ThisWeekSchedule(today string) []models.TrainingEvent {
	t, err := time.Parse("2006-01-02", today)
	if err != nil {
		return nil
	}
	// Find Monday of this week
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	monday := t.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)
	monStr := monday.Format("2006-01-02")
	sunStr := sunday.Format("2006-01-02")

	var result []models.TrainingEvent
	for _, evt := range s.Schedule {
		if evt.StartDate >= monStr && evt.StartDate <= sunStr {
			result = append(result, evt)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].StartDate == result[j].StartDate {
			return result[i].StartTime < result[j].StartTime
		}
		return result[i].StartDate < result[j].StartDate
	})
	return result
}

// RecentFeedback returns the last n feedback entries.
func (s *Store) RecentFeedback(n int) []models.EventFeedback {
	if n >= len(s.Feedback) {
		return s.Feedback
	}
	return s.Feedback[len(s.Feedback)-n:]
}

// ListQualifications returns all qualification reference entries.
func (s *Store) ListQualifications() []models.Qualification {
	return s.Qualifications
}

// ListQualRecords returns all instructor qualification records.
func (s *Store) ListQualRecords() []models.QualRecord {
	return s.QualRecords
}

// TotalStudentCount returns the total number of students loaded.
func (s *Store) TotalStudentCount() int {
	return len(s.Students)
}

// GetExamResults returns a student's results for a specific exam, or nil if not found.
func (s *Store) GetExamResults(studentID string, examNum int) *models.ExamResult {
	for i := range s.ExamResults {
		if s.ExamResults[i].StudentID == studentID && s.ExamResults[i].ExamNum == examNum {
			return &s.ExamResults[i]
		}
	}
	return nil
}

// XOScheduleForDate returns XO schedule items for a given date.
func (s *Store) XOScheduleForDate(date string) []models.XOScheduleItem {
	var result []models.XOScheduleItem
	for _, item := range s.XOSchedule {
		if item.Date == date {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime < result[j].StartTime
	})
	return result
}

// --- Task operations ---

func (s *Store) CreateTask(task models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.ID == "" {
		task.ID = fmt.Sprintf("TSK-%03d", s.nextID)
		s.nextID++
	}
	now := time.Now().Format(time.RFC3339)
	if task.CreatedAt == "" {
		task.CreatedAt = now
	}
	task.UpdatedAt = now
	if task.Status == "" {
		task.Status = "pending"
	}
	s.Tasks = append(s.Tasks, task)
	return s.persistJSON("tasks.json", s.Tasks)
}

func (s *Store) ListTasks(assignedTo string) []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if assignedTo == "" {
		result := make([]models.Task, len(s.Tasks))
		copy(result, s.Tasks)
		return result
	}
	var result []models.Task
	for _, t := range s.Tasks {
		if strings.EqualFold(t.AssignedTo, assignedTo) {
			result = append(result, t)
		}
	}
	return result
}

func (s *Store) GetTask(id string) (*models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return &s.Tasks[i], true
		}
	}
	return nil, false
}

func (s *Store) UpdateTask(id string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			if v, ok := updates["status"].(string); ok {
				s.Tasks[i].Status = v
			}
			if v, ok := updates["priority"].(string); ok {
				s.Tasks[i].Priority = v
			}
			if v, ok := updates["assignedTo"].(string); ok {
				s.Tasks[i].AssignedTo = v
			}
			s.Tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return s.persistJSON("tasks.json", s.Tasks)
		}
	}
	return fmt.Errorf("task %s not found", id)
}

// --- Message operations ---

func (s *Store) CreateMessage(msg models.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if msg.ID == "" {
		msg.ID = fmt.Sprintf("MSG-%03d", s.nextID)
		s.nextID++
	}
	if msg.CreatedAt == "" {
		msg.CreatedAt = time.Now().Format(time.RFC3339)
	}
	s.Messages = append(s.Messages, msg)
	return s.persistJSON("messages.json", s.Messages)
}

func (s *Store) ListMessages(userRole string) []models.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if userRole == "" {
		result := make([]models.Message, len(s.Messages))
		copy(result, s.Messages)
		return result
	}
	var result []models.Message
	for _, m := range s.Messages {
		if strings.EqualFold(m.To, userRole) {
			result = append(result, m)
		}
	}
	return result
}

func (s *Store) MarkMessageRead(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Messages {
		if s.Messages[i].ID == id {
			s.Messages[i].Read = true
			return s.persistJSON("messages.json", s.Messages)
		}
	}
	return fmt.Errorf("message %s not found", id)
}

// --- Notification operations ---

func (s *Store) CreateNotification(n models.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n.ID == "" {
		n.ID = fmt.Sprintf("NTF-%03d", s.nextID)
		s.nextID++
	}
	if n.CreatedAt == "" {
		n.CreatedAt = time.Now().Format(time.RFC3339)
	}
	s.Notifications = append(s.Notifications, n)
	return s.persistJSON("notifications.json", s.Notifications)
}

func (s *Store) ListNotifications(userRole string, unreadOnly bool) []models.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []models.Notification
	for _, n := range s.Notifications {
		if userRole != "" && !strings.EqualFold(n.UserRole, userRole) {
			continue
		}
		if unreadOnly && n.Read {
			continue
		}
		result = append(result, n)
	}
	// Reverse order — newest first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func (s *Store) MarkNotificationRead(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Notifications {
		if s.Notifications[i].ID == id {
			s.Notifications[i].Read = true
			return s.persistJSON("notifications.json", s.Notifications)
		}
	}
	return fmt.Errorf("notification %s not found", id)
}

func (s *Store) UnreadNotificationCount(userRole string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, n := range s.Notifications {
		if !n.Read && (userRole == "" || strings.EqualFold(n.UserRole, userRole)) {
			count++
		}
	}
	return count
}

// --- Helpers ---

func (s *Store) persistJSON(filename string, data interface{}) error {
	if s.demoMode {
		return nil // demo mode: in-memory only, no disk writes
	}
	path := filepath.Join(s.dataDir, filename)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func parseIDNum(id string) int {
	// Extract numeric suffix from IDs like "TSK-001", "MSG-042", "NTF-007"
	idx := strings.LastIndex(id, "-")
	if idx < 0 || idx+1 >= len(id) {
		return 0
	}
	n := 0
	for _, c := range id[idx+1:] {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
