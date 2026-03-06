package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"heywood-tbs/internal/calendar"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

// calendarProvider is the active calendar backend (mock or Outlook).
var calendarProvider calendar.CalendarProvider = &calendar.MockCalendar{}

// InitCalendar sets the calendar provider. Call from main.go if Outlook is configured.
func InitCalendar(provider calendar.CalendarProvider) {
	calendarProvider = provider
}

// handleCalendarEvents returns merged events (TBS schedule + calendar) for a date range.
func (h *Handler) handleCalendarEvents(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	company := middleware.GetCompany(r.Context())

	// Parse date range from query params
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			writeError(w, 400, "invalid start date format (use YYYY-MM-DD)")
			return
		}
	} else {
		// Default to start of current week
		now := time.Now()
		start = now.AddDate(0, 0, -int(now.Weekday()))
	}

	if endStr != "" {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			writeError(w, 400, "invalid end date format (use YYYY-MM-DD)")
			return
		}
	} else {
		// Default to 7 days from start
		end = start.AddDate(0, 0, 7)
	}

	events := calendarProvider.GetEvents(role, company, start, end)

	// Also merge TBS training schedule from the data store
	scheduleEvents := h.scheduleToCalendarEvents(role, company, start, end)
	events = append(events, scheduleEvents...)

	writeJSON(w, 200, map[string]interface{}{
		"events": events,
		"start":  start.Format("2006-01-02"),
		"end":    end.Format("2006-01-02"),
	})
}

// handleCalendarToday returns today's events for the current user.
func (h *Handler) handleCalendarToday(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	company := middleware.GetCompany(r.Context())

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := today.AddDate(0, 0, 1)

	events := calendarProvider.GetEvents(role, company, today, endOfDay)
	scheduleEvents := h.scheduleToCalendarEvents(role, company, today, endOfDay)
	events = append(events, scheduleEvents...)

	writeJSON(w, 200, map[string]interface{}{
		"events": events,
		"date":   today.Format("2006-01-02"),
	})
}

// handleMailSummary returns unread mail count and recent subjects.
func (h *Handler) handleMailSummary(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	mails := calendarProvider.GetMailSummary(role)

	unread := 0
	for _, m := range mails {
		if !m.IsRead {
			unread++
		}
	}

	writeJSON(w, 200, map[string]interface{}{
		"messages":    mails,
		"unreadCount": unread,
	})
}

// handleMailUnreadCount returns just the unread mail badge number.
func (h *Handler) handleMailUnreadCount(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	mails := calendarProvider.GetMailSummary(role)
	unread := 0
	for _, m := range mails {
		if !m.IsRead {
			unread++
		}
	}

	writeJSON(w, 200, map[string]interface{}{"count": unread})
}

// handleCreateCalendarEvent creates a new calendar event.
func (h *Handler) handleCreateCalendarEvent(w http.ResponseWriter, r *http.Request) {
	var event models.CalendarEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid event data"})
		return
	}

	created, err := calendarProvider.CreateEvent(event)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 201, created)
}

// scheduleToCalendarEvents converts TBS training schedule to calendar events.
func (h *Handler) scheduleToCalendarEvents(role, company string, start, end time.Time) []models.CalendarEvent {
	today := start.Format("2006-01-02")
	allEvents := h.store.ListSchedule("")

	var events []models.CalendarEvent
	for _, e := range allEvents {
		// Filter by date range
		if e.StartDate < today || e.StartDate > end.Format("2006-01-02") {
			continue
		}

		// Role-based filtering
		if role == "student" || role == "spc" {
			if company != "" && e.Company != "" && !strings.EqualFold(e.Company, company) && e.Company != "all" {
				continue
			}
		}

		events = append(events, models.CalendarEvent{
			ID:          "sched-" + e.ID,
			Title:       e.Title,
			Start:       e.StartDate + "T" + milToISO(e.StartTime),
			End:         e.EndDate + "T" + milToISO(e.EndTime),
			Location:    e.Location,
			Description: e.Notes,
			Source:      "tbs-schedule",
			Category:    e.Category,
			Company:     e.Company,
			IsAllDay:    false,
		})
	}

	return events
}

// milToISO converts military time "0700" to ISO "07:00:00".
func milToISO(mil string) string {
	if len(mil) == 4 {
		return mil[:2] + ":" + mil[2:] + ":00"
	}
	if len(mil) == 3 {
		return "0" + mil[:1] + ":" + mil[1:] + ":00"
	}
	return mil
}
