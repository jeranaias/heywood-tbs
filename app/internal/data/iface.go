package data

import "heywood-tbs/internal/models"

// ChatPersister is an optional capability for stores that support chat history.
// Use type assertion to check: if cp, ok := store.(ChatPersister); ok { ... }
// The JSON-backed Store does not implement this; SQLStore does.
type ChatPersister interface {
	CreateChatSession(session models.ChatSession) error
	ListChatSessions(userID, userRole string) []models.ChatSession
	GetChatSession(id string) (*models.ChatSession, bool)
	UpdateChatSessionTitle(id, title string) error
	AddChatMessage(sessionID string, msg models.ChatMessage) error
	GetChatMessages(sessionID string) []models.ChatMessage
	DeleteChatSession(id string) error
}

// DataStore defines the contract for TBS data access.
// The JSON-backed Store implements this. Future implementations (Cosmos DB, SQL)
// can implement this interface for seamless swapping.
type DataStore interface {
	// Read operations
	ListStudents(company, phase, search string, atRiskOnly bool) []models.Student
	GetStudent(id string) (*models.Student, bool)
	StudentStats(company string) models.StudentStats
	AtRiskStudents(company string) []models.Student
	ListInstructors(company string) []models.Instructor
	GetInstructor(id string) (*models.Instructor, bool)
	QualStats() models.QualStats
	ListSchedule(phase string) []models.TrainingEvent
	TodaySchedule(today string) []models.TrainingEvent
	ThisWeekSchedule(today string) []models.TrainingEvent
	ListFeedback(eventCode string) []models.EventFeedback
	RecentFeedback(n int) []models.EventFeedback
	ListQualifications() []models.Qualification
	ListQualRecords() []models.QualRecord
	TotalStudentCount() int
	XOScheduleForDate(date string) []models.XOScheduleItem

	// Exam results
	GetExamResults(studentID string, examNum int) *models.ExamResult

	// Task operations
	CreateTask(task models.Task) error
	ListTasks(assignedTo string) []models.Task
	GetTask(id string) (*models.Task, bool)
	UpdateTask(id string, req models.TaskUpdateRequest) error

	// Message operations
	CreateMessage(msg models.Message) error
	ListMessages(userRole string) []models.Message
	MarkMessageRead(id string) error

	// Notification operations
	CreateNotification(n models.Notification) error
	ListNotifications(userRole string, unreadOnly bool) []models.Notification
	MarkNotificationRead(id string) error
	UnreadNotificationCount(userRole string) int
}
