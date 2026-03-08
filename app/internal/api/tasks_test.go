package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

// withRoleContext returns a request with the given role injected into context.
func withRoleContext(req *http.Request, role string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.RoleKey, role)
	return req.WithContext(ctx)
}

func TestHandleListTasks(t *testing.T) {
	h := newTestHandler(t)

	// Seed a task so there is something to list
	err := h.store.CreateTask(models.Task{
		Title:      "Review PT scores",
		AssignedTo: auth.RoleStaff,
		Priority:   "medium",
		Status:     "pending",
	})
	if err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/tasks?all=true", nil)
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleListTasks(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var tasks []models.Task
	if err := json.NewDecoder(rec.Body).Decode(&tasks); err != nil {
		t.Fatalf("failed to decode tasks: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("expected at least one task in response")
	}

	// Verify the seeded task is present
	found := false
	for _, task := range tasks {
		if task.Title == "Review PT scores" {
			found = true
			if task.Status != "pending" {
				t.Errorf("expected status 'pending', got %q", task.Status)
			}
			break
		}
	}
	if !found {
		t.Error("seeded task 'Review PT scores' not found in response")
	}
}

func TestHandleUpdateTask_ValidFields(t *testing.T) {
	h := newTestHandler(t)

	// Seed a task
	err := h.store.CreateTask(models.Task{
		Title:      "Update counseling notes",
		AssignedTo: auth.RoleStaff,
		Priority:   "low",
		Status:     "pending",
	})
	if err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	// Get the task list to find the ID
	allTasks := h.store.ListTasks("")
	if len(allTasks) == 0 {
		t.Fatal("no tasks found after seeding")
	}
	var taskID string
	for _, task := range allTasks {
		if task.Title == "Update counseling notes" {
			taskID = task.ID
			break
		}
	}
	if taskID == "" {
		t.Fatal("could not find seeded task by title")
	}

	// PATCH to update status
	body := `{"status":"in_progress"}`
	req := httptest.NewRequest("PATCH", "/api/v1/tasks/"+taskID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", taskID)
	rec := httptest.NewRecorder()

	h.handleUpdateTask(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.Task
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Status != "in_progress" {
		t.Errorf("expected status 'in_progress', got %q", updated.Status)
	}
	if updated.ID != taskID {
		t.Errorf("expected task ID %q, got %q", taskID, updated.ID)
	}
}

func TestHandleUpdateTask_TypedRequest(t *testing.T) {
	h := newTestHandler(t)

	// Seed a task
	err := h.store.CreateTask(models.Task{
		Title:      "Schedule range time",
		AssignedTo: auth.RoleSPC,
		Priority:   "medium",
		Status:     "pending",
	})
	if err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	allTasks := h.store.ListTasks("")
	var taskID string
	for _, task := range allTasks {
		if task.Title == "Schedule range time" {
			taskID = task.ID
			break
		}
	}
	if taskID == "" {
		t.Fatal("could not find seeded task by title")
	}

	// Send a typed TaskUpdateRequest with multiple fields
	body := `{"status":"completed","priority":"high","assignedTo":"xo"}`
	req := httptest.NewRequest("PATCH", "/api/v1/tasks/"+taskID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", taskID)
	rec := httptest.NewRecorder()

	h.handleUpdateTask(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.Task
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Status != "completed" {
		t.Errorf("expected status 'completed', got %q", updated.Status)
	}
	if updated.Priority != "high" {
		t.Errorf("expected priority 'high', got %q", updated.Priority)
	}
	if updated.AssignedTo != "xo" {
		t.Errorf("expected assignedTo 'xo', got %q", updated.AssignedTo)
	}
}

func TestHandleUpdateTask_NotFound(t *testing.T) {
	h := newTestHandler(t)

	body := `{"status":"completed"}`
	req := httptest.NewRequest("PATCH", "/api/v1/tasks/TSK-999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "TSK-999")
	rec := httptest.NewRecorder()

	h.handleUpdateTask(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404 for nonexistent task, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleUpdateTask_InvalidBody(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("PATCH", "/api/v1/tasks/TSK-001", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "TSK-001")
	rec := httptest.NewRecorder()

	h.handleUpdateTask(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for invalid body, got %d: %s", rec.Code, rec.Body.String())
	}
}
