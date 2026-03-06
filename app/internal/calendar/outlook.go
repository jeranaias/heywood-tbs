package calendar

import (
	"time"

	"heywood-tbs/internal/models"
)

// CalendarProvider abstracts calendar data retrieval.
// Implementations: MockCalendar (demo), OutlookCalendar (production).
type CalendarProvider interface {
	GetEvents(role, company string, start, end time.Time) []models.CalendarEvent
	GetMailSummary(role string) []models.MailSummary
	CreateEvent(event models.CalendarEvent) (models.CalendarEvent, error)
}

// OutlookCalendar connects to Microsoft Graph API for real calendar/mail data.
// This is a stub — full implementation requires msgraph-sdk-go dependency.
type OutlookCalendar struct {
	TenantID string
	ClientID string
	Cloud    string // "commercial", "gcc-high", "dod"
}

// GetEvents retrieves calendar events from Outlook via Microsoft Graph.
// TODO: implement with msgraph-sdk-go when dependency is added.
func (o *OutlookCalendar) GetEvents(role, company string, start, end time.Time) []models.CalendarEvent {
	// Stub: returns empty until Graph SDK is integrated
	return nil
}

// GetMailSummary retrieves unread/recent mail from Outlook via Microsoft Graph.
// TODO: implement with msgraph-sdk-go when dependency is added.
func (o *OutlookCalendar) GetMailSummary(role string) []models.MailSummary {
	// Stub: returns empty until Graph SDK is integrated
	return nil
}

// CreateEvent creates a calendar event in Outlook via Microsoft Graph.
// TODO: implement with msgraph-sdk-go when dependency is added.
func (o *OutlookCalendar) CreateEvent(event models.CalendarEvent) (models.CalendarEvent, error) {
	return event, nil
}

// GraphEndpoint returns the Microsoft Graph base URL for the configured cloud.
func (o *OutlookCalendar) GraphEndpoint() string {
	switch o.Cloud {
	case "gcc-high":
		return "https://graph.microsoft.us"
	case "dod":
		return "https://dod-graph.microsoft.us"
	default:
		return "https://graph.microsoft.com"
	}
}
