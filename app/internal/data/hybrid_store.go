package data

import "heywood-tbs/internal/models"

// HybridStore wraps a read-only reference source (Excel, SharePoint, etc.)
// with a mutable JSON store for tasks, messages, and notifications.
// Reference data (students, instructors, schedule, quals) comes from the
// reference store; mutable operations fall through to the JSON store.
type HybridStore struct {
	ref     DataStore // read-only reference data
	mutable *Store    // JSON store for tasks, messages, notifications
}

// NewHybridStore creates a hybrid store with ref for reads, mutable for writes.
func NewHybridStore(ref DataStore, mutable *Store) *HybridStore {
	return &HybridStore{ref: ref, mutable: mutable}
}

// ---- Reference data: delegates to ref store ----

func (h *HybridStore) ListStudents(company, phase, search string, atRiskOnly bool) []models.Student {
	return h.ref.ListStudents(company, phase, search, atRiskOnly)
}
func (h *HybridStore) GetStudent(id string) (*models.Student, bool) {
	return h.ref.GetStudent(id)
}
func (h *HybridStore) UpdateStudent(id string, req models.StudentUpdateRequest) error {
	return h.ref.UpdateStudent(id, req)
}
func (h *HybridStore) CreateStudentNote(note models.StudentNote) error {
	return h.mutable.CreateStudentNote(note)
}
func (h *HybridStore) ListStudentNotes(studentID string) []models.StudentNote {
	return h.mutable.ListStudentNotes(studentID)
}
func (h *HybridStore) StudentStats(company string) models.StudentStats {
	return h.ref.StudentStats(company)
}
func (h *HybridStore) AtRiskStudents(company string) []models.Student {
	return h.ref.AtRiskStudents(company)
}
func (h *HybridStore) ListInstructors(company string) []models.Instructor {
	return h.ref.ListInstructors(company)
}
func (h *HybridStore) GetInstructor(id string) (*models.Instructor, bool) {
	return h.ref.GetInstructor(id)
}
func (h *HybridStore) QualStats() models.QualStats {
	return h.ref.QualStats()
}
func (h *HybridStore) ListSchedule(phase string) []models.TrainingEvent {
	return h.ref.ListSchedule(phase)
}
func (h *HybridStore) CreateTrainingEvent(event models.TrainingEvent) error {
	return h.mutable.CreateTrainingEvent(event)
}
func (h *HybridStore) UpdateTrainingEvent(id string, event models.TrainingEvent) error {
	return h.mutable.UpdateTrainingEvent(id, event)
}
func (h *HybridStore) DeleteTrainingEvent(id string) error {
	return h.mutable.DeleteTrainingEvent(id)
}
func (h *HybridStore) TodaySchedule(today string) []models.TrainingEvent {
	return h.ref.TodaySchedule(today)
}
func (h *HybridStore) ThisWeekSchedule(today string) []models.TrainingEvent {
	return h.ref.ThisWeekSchedule(today)
}
func (h *HybridStore) ListFeedback(eventCode string) []models.EventFeedback {
	return h.ref.ListFeedback(eventCode)
}
func (h *HybridStore) RecentFeedback(n int) []models.EventFeedback {
	return h.ref.RecentFeedback(n)
}
func (h *HybridStore) ListQualifications() []models.Qualification {
	return h.ref.ListQualifications()
}
func (h *HybridStore) ListQualRecords() []models.QualRecord {
	return h.ref.ListQualRecords()
}
func (h *HybridStore) TotalStudentCount() int {
	return h.ref.TotalStudentCount()
}
func (h *HybridStore) XOScheduleForDate(date string) []models.XOScheduleItem {
	return h.ref.XOScheduleForDate(date)
}
func (h *HybridStore) GetExamResults(studentID string, examNum int) *models.ExamResult {
	return h.ref.GetExamResults(studentID, examNum)
}

// ---- Counseling: delegates to mutable store ----

func (h *HybridStore) CreateCounseling(session models.CounselingSession) error {
	return h.mutable.CreateCounseling(session)
}
func (h *HybridStore) ListCounselings(studentID string) []models.CounselingSession {
	return h.mutable.ListCounselings(studentID)
}
func (h *HybridStore) GetCounseling(id string) (*models.CounselingSession, bool) {
	return h.mutable.GetCounseling(id)
}
func (h *HybridStore) UpdateCounseling(id string, session models.CounselingSession) error {
	return h.mutable.UpdateCounseling(id, session)
}

// ---- Mutable data: delegates to JSON store ----

func (h *HybridStore) CreateTask(task models.Task) error {
	return h.mutable.CreateTask(task)
}
func (h *HybridStore) ListTasks(assignedTo string) []models.Task {
	return h.mutable.ListTasks(assignedTo)
}
func (h *HybridStore) GetTask(id string) (*models.Task, bool) {
	return h.mutable.GetTask(id)
}
func (h *HybridStore) UpdateTask(id string, req models.TaskUpdateRequest) error {
	return h.mutable.UpdateTask(id, req)
}
func (h *HybridStore) DeleteTask(id string) error {
	return h.mutable.DeleteTask(id)
}
func (h *HybridStore) CreateMessage(msg models.Message) error {
	return h.mutable.CreateMessage(msg)
}
func (h *HybridStore) ListMessages(userRole string) []models.Message {
	return h.mutable.ListMessages(userRole)
}
func (h *HybridStore) MarkMessageRead(id string) error {
	return h.mutable.MarkMessageRead(id)
}
func (h *HybridStore) CreateNotification(n models.Notification) error {
	return h.mutable.CreateNotification(n)
}
func (h *HybridStore) ListNotifications(userRole string, unreadOnly bool) []models.Notification {
	return h.mutable.ListNotifications(userRole, unreadOnly)
}
func (h *HybridStore) MarkNotificationRead(id string) error {
	return h.mutable.MarkNotificationRead(id)
}
func (h *HybridStore) UnreadNotificationCount(userRole string) int {
	return h.mutable.UnreadNotificationCount(userRole)
}
