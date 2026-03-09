package api

import (
	"testing"
	"time"

	"heywood-tbs/internal/models"
)

func TestSSEBrokerBroadcast(t *testing.T) {
	broker := NewSSEBroker()

	ch1 := make(chan SSEEvent, 8)
	ch2 := make(chan SSEEvent, 8)

	broker.Register("staff", ch1)
	broker.Register("xo", ch2)

	// Broadcast to staff only
	broker.Broadcast("staff", SSEEvent{Type: "task", Data: map[string]interface{}{"action": "created"}})

	select {
	case evt := <-ch1:
		if evt.Type != "task" {
			t.Errorf("expected type 'task', got %s", evt.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("staff channel did not receive event")
	}

	// XO channel should be empty
	select {
	case <-ch2:
		t.Fatal("xo channel should not have received staff-targeted event")
	default:
		// expected
	}

	// Broadcast to all (role="")
	broker.Broadcast("", SSEEvent{Type: "notification", Data: map[string]interface{}{"message": "hello"}})

	select {
	case evt := <-ch1:
		if evt.Type != "notification" {
			t.Errorf("expected type 'notification', got %s", evt.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("staff channel did not receive broadcast-all event")
	}

	select {
	case evt := <-ch2:
		if evt.Type != "notification" {
			t.Errorf("expected type 'notification', got %s", evt.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("xo channel did not receive broadcast-all event")
	}
}

func TestSSEClientRegistration(t *testing.T) {
	broker := NewSSEBroker()

	if broker.ClientCount() != 0 {
		t.Fatalf("expected 0 clients, got %d", broker.ClientCount())
	}

	ch1 := make(chan SSEEvent, 4)
	ch2 := make(chan SSEEvent, 4)

	broker.Register("staff", ch1)
	broker.Register("staff", ch2)

	if broker.ClientCount() != 2 {
		t.Fatalf("expected 2 clients, got %d", broker.ClientCount())
	}

	broker.Unregister("staff", ch1)
	if broker.ClientCount() != 1 {
		t.Fatalf("expected 1 client after unregister, got %d", broker.ClientCount())
	}

	broker.Unregister("staff", ch2)
	if broker.ClientCount() != 0 {
		t.Fatalf("expected 0 clients after all unregistered, got %d", broker.ClientCount())
	}
}

func TestAtRiskMonitorTriggersNotification(t *testing.T) {
	h := newTestHandler(t)
	h.sseBroker = NewSSEBroker()

	// Subscribe to xo events
	ch := make(chan SSEEvent, 16)
	h.sseBroker.Register("xo", ch)
	defer h.sseBroker.Unregister("xo", ch)

	// Get a student and set them at-risk directly in the store
	students := h.store.ListStudents("", "", "", false)
	if len(students) == 0 {
		t.Skip("no students in test store")
	}

	// Pick a non-at-risk student
	var targetID string
	for _, s := range students {
		if !s.AtRisk {
			targetID = s.ID
			break
		}
	}
	if targetID == "" {
		t.Skip("all students are already at-risk")
	}

	// Run the monitor check with a snapshot where student was not at-risk
	knownAtRisk := make(map[string]bool)
	for _, s := range students {
		if s.AtRisk {
			knownAtRisk[s.ID] = true
		}
	}

	// Flag the student as at-risk
	atRiskTrue := true
	h.store.UpdateStudent(targetID, models.StudentUpdateRequest{AtRisk: &atRiskTrue})

	// Run the check
	h.checkAtRiskChanges(knownAtRisk)

	// Verify SSE event was sent
	select {
	case evt := <-ch:
		if evt.Type != "at-risk-alert" {
			t.Errorf("expected 'at-risk-alert', got %s", evt.Type)
		}
		data := evt.Data.(map[string]interface{})
		if data["studentId"] != targetID {
			t.Errorf("expected studentId=%s, got %v", targetID, data["studentId"])
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("did not receive at-risk-alert SSE event")
	}
}
