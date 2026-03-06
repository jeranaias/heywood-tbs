package data

import "heywood-tbs/internal/models"

// DataStore defines the contract for TBS data access.
// The JSON-backed Store implements this. Future implementations (Cosmos DB, SQL)
// can implement this interface for seamless swapping.
type DataStore interface {
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
}
