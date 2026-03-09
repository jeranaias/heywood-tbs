package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/models"
)

func TestCreateTrainingEvent(t *testing.T) {
	h := newTestHandler(t)

	body := `{"title":"Patrol Base Operations","code":"TBS-PBO","phase":"Phase 2","startDate":"2026-04-01","startTime":"0600","endTime":"1800","location":"TBS Training Area 1","isGraded":true}`
	req := httptest.NewRequest("POST", "/api/v1/schedule", strings.NewReader(body))
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleCreateTrainingEvent(rec, req)

	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var event models.TrainingEvent
	json.NewDecoder(rec.Body).Decode(&event)
	if event.Title != "Patrol Base Operations" {
		t.Errorf("expected title 'Patrol Base Operations', got %s", event.Title)
	}
	if !event.IsGraded {
		t.Error("expected isGraded=true")
	}
}

func TestUpdateTrainingEvent(t *testing.T) {
	h := newTestHandler(t)

	// Get first event ID
	events := h.store.ListSchedule("")
	if len(events) == 0 {
		t.Skip("no schedule data to test with")
	}
	id := events[0].ID

	body := `{"title":"Updated Event Title","location":"New Location"}`
	req := httptest.NewRequest("PUT", "/api/v1/schedule/"+id, strings.NewReader(body))
	req.SetPathValue("id", id)
	req = withFullContext(req, auth.RoleXO, "", "")
	rec := httptest.NewRecorder()

	h.handleUpdateTrainingEvent(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var event models.TrainingEvent
	json.NewDecoder(rec.Body).Decode(&event)
	if event.Title != "Updated Event Title" {
		t.Errorf("expected updated title, got %s", event.Title)
	}
}

func TestDeleteTrainingEvent(t *testing.T) {
	h := newTestHandler(t)

	// Create an event first, then delete it
	evt := models.TrainingEvent{Title: "Temp Event", Code: "TMP-001"}
	h.store.CreateTrainingEvent(evt)

	events := h.store.ListSchedule("")
	var tempID string
	for _, e := range events {
		if e.Code == "TMP-001" {
			tempID = e.ID
			break
		}
	}
	if tempID == "" {
		t.Fatal("temp event not found")
	}

	req := httptest.NewRequest("DELETE", "/api/v1/schedule/"+tempID, nil)
	req.SetPathValue("id", tempID)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleDeleteTrainingEvent(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestScheduleRoleGating(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		name    string
		method  string
		role    string
		handler http.HandlerFunc
		body    string
		expect  int
	}{
		{"student create", "POST", auth.RoleStudent, h.handleCreateTrainingEvent, `{"title":"Test"}`, 403},
		{"spc create", "POST", auth.RoleSPC, h.handleCreateTrainingEvent, `{"title":"Test"}`, 403},
		{"student delete", "DELETE", auth.RoleStudent, h.handleDeleteTrainingEvent, "", 403},
		{"spc delete", "DELETE", auth.RoleSPC, h.handleDeleteTrainingEvent, "", 403},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, "/api/v1/schedule", strings.NewReader(tc.body))
			} else {
				req = httptest.NewRequest(tc.method, "/api/v1/schedule/EVT-001", nil)
				req.SetPathValue("id", "EVT-001")
			}
			req = withFullContext(req, tc.role, "", "")
			rec := httptest.NewRecorder()
			tc.handler(rec, req)

			if rec.Code != tc.expect {
				t.Errorf("expected %d for %s, got %d", tc.expect, tc.name, rec.Code)
			}
		})
	}
}

func TestRescheduleEvent(t *testing.T) {
	h := newTestHandler(t)

	// Create event, then update its dates (reschedule)
	evt := models.TrainingEvent{Title: "MOUT Exercise", Code: "TBS-MOUT", StartDate: "2026-04-01", EndDate: "2026-04-01"}
	h.store.CreateTrainingEvent(evt)

	events := h.store.ListSchedule("")
	var moutID string
	for _, e := range events {
		if e.Code == "TBS-MOUT" {
			moutID = e.ID
			break
		}
	}
	if moutID == "" {
		t.Fatal("MOUT event not found")
	}

	// Reschedule to April 5
	body := `{"title":"MOUT Exercise","code":"TBS-MOUT","startDate":"2026-04-05","endDate":"2026-04-05","status":"rescheduled"}`
	req := httptest.NewRequest("PUT", "/api/v1/schedule/"+moutID, strings.NewReader(body))
	req.SetPathValue("id", moutID)
	req = withFullContext(req, auth.RoleStaff, "", "")
	rec := httptest.NewRecorder()

	h.handleUpdateTrainingEvent(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.TrainingEvent
	json.NewDecoder(rec.Body).Decode(&updated)
	if updated.StartDate != "2026-04-05" {
		t.Errorf("expected rescheduled date 2026-04-05, got %s", updated.StartDate)
	}
}
