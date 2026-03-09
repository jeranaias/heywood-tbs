package calendar

import (
	"fmt"
	"time"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/models"
)

// MockCalendar generates demo calendar events when no Outlook is connected.
// Combines TBS training schedule with simulated personal/shared events.
type MockCalendar struct{}

// GetEvents returns mock events for the given date range, filtered by role.
func (m *MockCalendar) GetEvents(role, company string, start, end time.Time) []models.CalendarEvent {
	var events []models.CalendarEvent

	// Generate events for each day in range
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dayOfWeek := d.Weekday()

		// Skip weekends
		if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
			continue
		}

		dateStr := d.Format("2006-01-02")

		// Morning PT (all roles)
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-pt-%s", dateStr),
			Title:    "Physical Training",
			Start:    dateStr + "T05:30:00",
			End:      dateStr + "T06:30:00",
			Location: "Heywood Field",
			Source:   "tbs-schedule",
			Category: "training",
			IsAllDay: false,
		})

		// Role-specific events
		switch role {
		case auth.RoleXO:
			events = append(events, m.xoEvents(dateStr, dayOfWeek)...)
		case auth.RoleStaff:
			events = append(events, m.staffEvents(dateStr, dayOfWeek)...)
		case auth.RoleSPC:
			events = append(events, m.spcEvents(dateStr, dayOfWeek, company)...)
		case auth.RoleStudent:
			events = append(events, m.studentEvents(dateStr, dayOfWeek, company)...)
		}

		// Common daily event: lunch
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-chow-%s", dateStr),
			Title:    "Chow",
			Start:    dateStr + "T11:30:00",
			End:      dateStr + "T12:30:00",
			Location: "Hawkins Hall",
			Source:   "tbs-schedule",
			Category: "admin",
			IsAllDay: false,
		})
	}

	return events
}

func (m *MockCalendar) xoEvents(date string, dow time.Weekday) []models.CalendarEvent {
	var events []models.CalendarEvent

	// Daily staff meeting
	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-staff-mtg-%s", date),
		Title:    "XO/Staff Daily Sync",
		Start:    date + "T07:00:00",
		End:      date + "T07:30:00",
		Location: "Room 201, Gruber Hall",
		Source:   "outlook",
		Category: "admin",
	})

	if dow == time.Monday {
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-weekly-%s", date),
			Title:    "Weekly Training Brief to CO",
			Start:    date + "T08:00:00",
			End:      date + "T09:00:00",
			Location: "CO Conference Room",
			Source:   "outlook",
			Category: "admin",
			Organizer: "Col Martinez",
		})
	}

	if dow == time.Wednesday {
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-readiness-%s", date),
			Title:    "Training Readiness Review",
			Start:    date + "T14:00:00",
			End:      date + "T15:00:00",
			Location: "Gruber Hall Auditorium",
			Source:   "outlook",
			Category: "admin",
		})
	}

	// Afternoon block: training observation
	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-observe-%s", date),
		Title:    "Training Observation (Alpha Co)",
		Start:    date + "T13:00:00",
		End:      date + "T15:00:00",
		Location: "Range 8 / Field",
		Source:   "tbs-schedule",
		Category: "training",
	})

	return events
}

func (m *MockCalendar) staffEvents(date string, dow time.Weekday) []models.CalendarEvent {
	var events []models.CalendarEvent

	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-instr-prep-%s", date),
		Title:    "Instructor Preparation",
		Start:    date + "T07:00:00",
		End:      date + "T08:00:00",
		Location: "Staff Office",
		Source:   "tbs-schedule",
		Category: "training",
	})

	// Morning instruction block
	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-instruct-%s", date),
		Title:    "Classroom Instruction Block",
		Start:    date + "T08:00:00",
		End:      date + "T11:30:00",
		Location: "Gruber Hall",
		Source:   "tbs-schedule",
		Category: "training",
	})

	if dow == time.Thursday {
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-counseling-%s", date),
			Title:    "Student Counseling Sessions",
			Start:    date + "T13:00:00",
			End:      date + "T16:00:00",
			Location: "Staff Office",
			Source:   "outlook",
			Category: "admin",
		})
	}

	return events
}

func (m *MockCalendar) spcEvents(date string, dow time.Weekday, company string) []models.CalendarEvent {
	var events []models.CalendarEvent
	coLabel := "Company"
	if company != "" {
		coLabel = company + " Co"
	}

	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-formation-%s", date),
		Title:    coLabel + " Morning Formation",
		Start:    date + "T06:30:00",
		End:      date + "T06:45:00",
		Location: "Company Area",
		Source:   "tbs-schedule",
		Category: "admin",
		Company:  company,
	})

	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-trng-%s", date),
		Title:    coLabel + " Training Block",
		Start:    date + "T08:00:00",
		End:      date + "T11:30:00",
		Location: "Assigned Training Area",
		Source:   "tbs-schedule",
		Category: "training",
		Company:  company,
	})

	if dow == time.Friday {
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-field-day-%s", date),
			Title:    coLabel + " Field Day / Inspection",
			Start:    date + "T13:00:00",
			End:      date + "T16:00:00",
			Location: "Barracks",
			Source:   "tbs-schedule",
			Category: "admin",
			Company:  company,
		})
	}

	return events
}

func (m *MockCalendar) studentEvents(date string, dow time.Weekday, company string) []models.CalendarEvent {
	var events []models.CalendarEvent
	coLabel := "Company"
	if company != "" {
		coLabel = company + " Co"
	}

	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-stu-form-%s", date),
		Title:    coLabel + " Morning Formation",
		Start:    date + "T06:30:00",
		End:      date + "T06:45:00",
		Location: "Company Area",
		Source:   "tbs-schedule",
		Category: "admin",
		Company:  company,
	})

	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-stu-class-%s", date),
		Title:    "Classroom Instruction",
		Start:    date + "T08:00:00",
		End:      date + "T11:30:00",
		Location: "Gruber Hall",
		Source:   "tbs-schedule",
		Category: "training",
	})

	events = append(events, models.CalendarEvent{
		ID:       fmt.Sprintf("mock-stu-prac-%s", date),
		Title:    "Practical Application",
		Start:    date + "T13:00:00",
		End:      date + "T16:00:00",
		Location: "Training Area",
		Source:   "tbs-schedule",
		Category: "training",
	})

	if dow == time.Tuesday || dow == time.Thursday {
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-study-%s", date),
			Title:    "Study Group / Mentoring",
			Start:    date + "T18:00:00",
			End:      date + "T19:30:00",
			Location: "Library",
			Source:   "outlook",
			Category: "personal",
		})
	}

	if dow == time.Wednesday {
		events = append(events, models.CalendarEvent{
			ID:       fmt.Sprintf("mock-exam-prep-%s", date),
			Title:    "Exam Review Session",
			Start:    date + "T17:00:00",
			End:      date + "T18:30:00",
			Location: "Gruber Hall Room 103",
			Source:   "outlook",
			Category: "training",
		})
	}

	return events
}

// GetMailSummary returns mock email summaries for demo mode.
func (m *MockCalendar) GetMailSummary(role string) []models.MailSummary {
	now := time.Now()

	switch role {
	case auth.RoleXO:
		return []models.MailSummary{
			{ID: "mail-1", Subject: "Training Schedule Update — Week 12", From: "S-3 Operations", Preview: "Updated training matrix attached. Key changes to Range 8 allocation...", Received: now.Add(-2 * time.Hour).Format(time.RFC3339), IsRead: false, HasAttach: true},
			{ID: "mail-2", Subject: "RE: At-Risk Student Board Results", From: "Capt Rodriguez", Preview: "Concur with recommendations. Lt Chen and Lt Park should receive additional...", Received: now.Add(-5 * time.Hour).Format(time.RFC3339), IsRead: false},
			{ID: "mail-3", Subject: "Quarterly Training Brief — Draft", From: "GySgt Williams", Preview: "Sir, please review the attached draft brief for CO approval...", Received: now.Add(-24 * time.Hour).Format(time.RFC3339), IsRead: true, HasAttach: true},
		}
	case auth.RoleStaff:
		return []models.MailSummary{
			{ID: "mail-4", Subject: "Instructor Qual Renewal — Deadline 15 Mar", From: "Training Admin", Preview: "The following instructors have qualifications expiring within 30 days...", Received: now.Add(-1 * time.Hour).Format(time.RFC3339), IsRead: false},
			{ID: "mail-5", Subject: "Range Request Approved", From: "Range Control", Preview: "Your range request for 12-14 Mar has been approved. Time block: 0800-1600...", Received: now.Add(-8 * time.Hour).Format(time.RFC3339), IsRead: true},
		}
	case auth.RoleSPC:
		return []models.MailSummary{
			{ID: "mail-6", Subject: "Alpha Co — Weekly Student Status Report", From: "1stSgt Davis", Preview: "SPC, updated student tracker attached. Two new at-risk nominations...", Received: now.Add(-3 * time.Hour).Format(time.RFC3339), IsRead: false, HasAttach: true},
		}
	case auth.RoleStudent:
		return []models.MailSummary{
			{ID: "mail-7", Subject: "Land Navigation Practical — Packing List", From: "Capt Thompson", Preview: "Ensure the following items are in your pack for Thursday's land nav...", Received: now.Add(-6 * time.Hour).Format(time.RFC3339), IsRead: false},
			{ID: "mail-8", Subject: "Study Group — Thursday 1800", From: "Lt Martinez", Preview: "Hey, we're meeting at the library Thursday at 1800 to prep for Exam 2...", Received: now.Add(-10 * time.Hour).Format(time.RFC3339), IsRead: true},
		}
	default:
		return nil
	}
}

// CreateEvent adds a mock event (in demo mode, just returns the event with a generated ID).
func (m *MockCalendar) CreateEvent(event models.CalendarEvent) (models.CalendarEvent, error) {
	event.ID = fmt.Sprintf("demo-%d", time.Now().UnixMilli())
	event.Source = "outlook"
	return event, nil
}

// SendMail is a no-op in demo mode.
func (m *MockCalendar) SendMail(role string, to []string, subject, body string) error {
	return nil
}

// ReplyToMail is a no-op in demo mode.
func (m *MockCalendar) ReplyToMail(role string, messageID, body string) error {
	return nil
}

// RespondToEvent is a no-op in demo mode.
func (m *MockCalendar) RespondToEvent(role string, eventID, response string) error {
	return nil
}
