package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/calendar"
)

func TestCalendarEvents_DefaultRange(t *testing.T) {
	h := newTestHandler(t)
	h.calendarProvider = &calendar.MockCalendar{}

	req := httptest.NewRequest("GET", "/api/v1/calendar/events", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleCalendarEvents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Events []json.RawMessage `json:"events"`
		Start  string            `json:"start"`
		End    string            `json:"end"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Events == nil {
		t.Error("expected 'events' array, got nil")
	}
	if resp.Start == "" {
		t.Error("expected non-empty 'start' date string")
	}
	if resp.End == "" {
		t.Error("expected non-empty 'end' date string")
	}
}

func TestCalendarEvents_CustomRange(t *testing.T) {
	h := newTestHandler(t)
	h.calendarProvider = &calendar.MockCalendar{}

	req := httptest.NewRequest("GET", "/api/v1/calendar/events?start=2026-01-05&end=2026-01-12", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleCalendarEvents(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Events []json.RawMessage `json:"events"`
		Start  string            `json:"start"`
		End    string            `json:"end"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Start != "2026-01-05" {
		t.Errorf("expected start=2026-01-05, got %q", resp.Start)
	}
	if resp.End != "2026-01-12" {
		t.Errorf("expected end=2026-01-12, got %q", resp.End)
	}
}

func TestCalendarEvents_InvalidDate(t *testing.T) {
	h := newTestHandler(t)
	h.calendarProvider = &calendar.MockCalendar{}

	req := httptest.NewRequest("GET", "/api/v1/calendar/events?start=bad-date", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleCalendarEvents(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for invalid date, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCalendarToday(t *testing.T) {
	h := newTestHandler(t)
	h.calendarProvider = &calendar.MockCalendar{}

	req := httptest.NewRequest("GET", "/api/v1/calendar/today", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleCalendarToday(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Events []json.RawMessage `json:"events"`
		Date   string            `json:"date"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Events may be nil/empty on weekends since MockCalendar skips Sat/Sun.
	// The key assertion is that the response structure is correct (200 + date).
	if resp.Date == "" {
		t.Error("expected non-empty 'date' string")
	}
}

func TestMailSummary(t *testing.T) {
	h := newTestHandler(t)
	h.calendarProvider = &calendar.MockCalendar{}

	req := httptest.NewRequest("GET", "/api/v1/mail/summary", nil)
	req = withRoleContext(req, auth.RoleXO)
	rec := httptest.NewRecorder()

	h.handleMailSummary(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Messages    []json.RawMessage `json:"messages"`
		UnreadCount int               `json:"unreadCount"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Messages == nil {
		t.Error("expected 'messages' array, got nil")
	}
	// XO mock returns 3 messages with 2 unread
	if len(resp.Messages) == 0 {
		t.Error("expected at least one message for XO role")
	}
	// UnreadCount must be a non-negative integer (checked by type)
	if resp.UnreadCount < 0 {
		t.Errorf("expected non-negative unreadCount, got %d", resp.UnreadCount)
	}
}

func TestMailUnreadCount(t *testing.T) {
	h := newTestHandler(t)
	h.calendarProvider = &calendar.MockCalendar{}

	req := httptest.NewRequest("GET", "/api/v1/mail/unread-count", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleMailUnreadCount(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Staff mock has 1 unread message
	if resp.Count < 0 {
		t.Errorf("expected non-negative count, got %d", resp.Count)
	}
}

func TestMilToISO(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"four digit 0700", "0700", "07:00:00"},
		{"four digit 1430", "1430", "14:30:00"},
		{"three digit 530", "530", "05:30:00"},
		{"fallthrough short", "07", "07"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := milToISO(tc.input)
			if got != tc.want {
				t.Errorf("milToISO(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
