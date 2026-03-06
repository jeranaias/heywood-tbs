package calendar

import (
	"log/slog"
	"time"

	"heywood-tbs/internal/models"
	"heywood-tbs/internal/msgraph"
)

// CalendarProvider abstracts calendar data retrieval.
// Implementations: MockCalendar (demo), OutlookCalendar (production).
type CalendarProvider interface {
	GetEvents(role, company string, start, end time.Time) []models.CalendarEvent
	GetMailSummary(role string) []models.MailSummary
	CreateEvent(event models.CalendarEvent) (models.CalendarEvent, error)
}

// OutlookCalendar connects to Microsoft Graph API for real calendar/mail data.
type OutlookCalendar struct {
	calSvc          *msgraph.CalendarService
	mailSvc         *msgraph.MailService
	masterCalID     string // shared master calendar ID (optional)
	userLookup      func(role, company string) string // resolves role → Graph user ID
}

// NewOutlookCalendar creates a production calendar provider backed by Microsoft Graph.
// graphClient: authenticated Graph API client
// masterCalID: optional shared calendar ID for battalion-wide events
// userLookup: function that maps (role, company) → user principal name for Graph queries
func NewOutlookCalendar(graphClient *msgraph.Client, masterCalID string, userLookup func(role, company string) string) *OutlookCalendar {
	return &OutlookCalendar{
		calSvc:      msgraph.NewCalendarService(graphClient),
		mailSvc:     msgraph.NewMailService(graphClient),
		masterCalID: masterCalID,
		userLookup:  userLookup,
	}
}

// GetEvents retrieves calendar events from Outlook via Microsoft Graph.
// Merges personal calendar + shared master calendar (if configured).
func (o *OutlookCalendar) GetEvents(role, company string, start, end time.Time) []models.CalendarEvent {
	var all []models.CalendarEvent

	// Get personal calendar events
	userID := o.resolveUser(role, company)
	if userID != "" {
		events, err := o.calSvc.GetEvents(userID, start, end)
		if err != nil {
			slog.Error("Outlook calendar query failed", "role", role, "user", userID, "error", err)
		} else {
			all = append(all, events...)
		}
	}

	// Get shared/master calendar events (if configured)
	if o.masterCalID != "" {
		shared, err := o.calSvc.GetSharedCalendar(o.masterCalID, start, end)
		if err != nil {
			slog.Error("shared calendar query failed", "calendarID", o.masterCalID, "error", err)
		} else {
			all = append(all, shared...)
		}
	}

	return all
}

// GetMailSummary retrieves recent emails from Outlook via Microsoft Graph.
func (o *OutlookCalendar) GetMailSummary(role string) []models.MailSummary {
	userID := o.resolveUser(role, "")
	if userID == "" {
		return nil
	}

	mails, err := o.mailSvc.GetMailSummary(userID, false)
	if err != nil {
		slog.Error("Outlook mail query failed", "role", role, "user", userID, "error", err)
		return nil
	}

	return mails
}

// CreateEvent creates a calendar event in Outlook via Microsoft Graph.
func (o *OutlookCalendar) CreateEvent(event models.CalendarEvent) (models.CalendarEvent, error) {
	// Use a default user for event creation; in production this would be
	// the authenticated user's UPN from the CAC certificate.
	userID := o.resolveUser("staff", "")
	if userID == "" {
		return event, nil
	}

	return o.calSvc.CreateEvent(userID, event)
}

func (o *OutlookCalendar) resolveUser(role, company string) string {
	if o.userLookup != nil {
		return o.userLookup(role, company)
	}
	return ""
}
