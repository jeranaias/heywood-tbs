package msgraph

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"heywood-tbs/internal/models"
)

// CalendarService reads/writes Outlook calendar events via Microsoft Graph.
type CalendarService struct {
	client *Client
}

// NewCalendarService creates a calendar service backed by the given Graph client.
func NewCalendarService(client *Client) *CalendarService {
	return &CalendarService{client: client}
}

// graphEvent is the Microsoft Graph event response shape.
type graphEvent struct {
	ID           string `json:"id"`
	Subject      string `json:"subject"`
	Start        struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
	Location struct {
		DisplayName string `json:"displayName"`
	} `json:"location"`
	BodyPreview  string `json:"bodyPreview"`
	IsAllDay     bool   `json:"isAllDay"`
	Organizer    struct {
		EmailAddress struct {
			Name string `json:"name"`
		} `json:"emailAddress"`
	} `json:"organizer"`
}

// GetEvents retrieves calendar events for a user in a date range.
// userID: UPN (user@domain.com) or user object ID.
func (s *CalendarService) GetEvents(userID string, start, end time.Time) ([]models.CalendarEvent, error) {
	if !s.client.IsConfigured() {
		return nil, nil
	}

	path := fmt.Sprintf("/users/%s/calendarView", userID)
	params := map[string]string{
		"startDateTime": start.UTC().Format("2006-01-02T15:04:05Z"),
		"endDateTime":   end.UTC().Format("2006-01-02T15:04:05Z"),
		"$select":       "id,subject,start,end,location,bodyPreview,isAllDay,organizer",
		"$orderby":      "start/dateTime",
		"$top":          "50",
	}

	body, err := s.client.Get(path, params)
	if err != nil {
		slog.Error("graph calendar query failed", "user", userID, "error", err)
		return nil, err
	}

	var resp struct {
		Value []graphEvent `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse calendar response: %w", err)
	}

	events := make([]models.CalendarEvent, 0, len(resp.Value))
	for _, ge := range resp.Value {
		events = append(events, models.CalendarEvent{
			ID:          ge.ID,
			Title:       ge.Subject,
			Start:       ge.Start.DateTime,
			End:         ge.End.DateTime,
			Location:    ge.Location.DisplayName,
			Description: ge.BodyPreview,
			Source:      "outlook",
			IsAllDay:    ge.IsAllDay,
			Organizer:   ge.Organizer.EmailAddress.Name,
		})
	}

	return events, nil
}

// GetSharedCalendar retrieves events from a shared/group calendar.
func (s *CalendarService) GetSharedCalendar(calendarID string, start, end time.Time) ([]models.CalendarEvent, error) {
	if !s.client.IsConfigured() || calendarID == "" {
		return nil, nil
	}

	path := fmt.Sprintf("/users/%s/calendarView", calendarID)
	params := map[string]string{
		"startDateTime": start.UTC().Format("2006-01-02T15:04:05Z"),
		"endDateTime":   end.UTC().Format("2006-01-02T15:04:05Z"),
		"$select":       "id,subject,start,end,location,bodyPreview,isAllDay,organizer",
		"$orderby":      "start/dateTime",
		"$top":          "100",
	}

	body, err := s.client.Get(path, params)
	if err != nil {
		slog.Error("graph shared calendar query failed", "calendar", calendarID, "error", err)
		return nil, err
	}

	var resp struct {
		Value []graphEvent `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	events := make([]models.CalendarEvent, 0, len(resp.Value))
	for _, ge := range resp.Value {
		events = append(events, models.CalendarEvent{
			ID:          ge.ID,
			Title:       ge.Subject,
			Start:       ge.Start.DateTime,
			End:         ge.End.DateTime,
			Location:    ge.Location.DisplayName,
			Description: ge.BodyPreview,
			Source:      "shared",
			IsAllDay:    ge.IsAllDay,
			Organizer:   ge.Organizer.EmailAddress.Name,
		})
	}

	return events, nil
}

// CreateEvent creates a new calendar event for a user.
func (s *CalendarService) CreateEvent(userID string, event models.CalendarEvent) (models.CalendarEvent, error) {
	if !s.client.IsConfigured() {
		return event, fmt.Errorf("Graph client not configured")
	}

	payload := map[string]interface{}{
		"subject": event.Title,
		"start": map[string]string{
			"dateTime": event.Start,
			"timeZone": "UTC",
		},
		"end": map[string]string{
			"dateTime": event.End,
			"timeZone": "UTC",
		},
		"isAllDay": event.IsAllDay,
	}
	if event.Location != "" {
		payload["location"] = map[string]string{"displayName": event.Location}
	}
	if event.Description != "" {
		payload["body"] = map[string]interface{}{
			"contentType": "text",
			"content":     event.Description,
		}
	}

	path := fmt.Sprintf("/users/%s/events", userID)
	body, err := s.client.Post(path, payload)
	if err != nil {
		return event, err
	}

	var created graphEvent
	if err := json.Unmarshal(body, &created); err != nil {
		return event, err
	}

	event.ID = created.ID
	event.Source = "outlook"
	return event, nil
}
