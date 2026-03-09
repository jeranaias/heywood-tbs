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

func TestHandleCreateTask_Valid(t *testing.T) {
	h := newTestHandler(t)

	body := `{"title":"Review PT scores","description":"Check all at-risk students","assignedTo":"staff","priority":"high","dueDate":"2026-03-15"}`
	req := httptest.NewRequest("POST", "/api/v1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleXO)
	rec := httptest.NewRecorder()

	h.handleCreateTask(rec, req)

	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var task models.Task
	if err := json.NewDecoder(rec.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if task.Title != "Review PT scores" {
		t.Errorf("expected title 'Review PT scores', got %q", task.Title)
	}
	if task.CreatedBy != "xo" {
		t.Errorf("expected createdBy 'xo', got %q", task.CreatedBy)
	}
	if task.Priority != "high" {
		t.Errorf("expected priority 'high', got %q", task.Priority)
	}
}

func TestHandleCreateTask_MissingTitle(t *testing.T) {
	h := newTestHandler(t)

	body := `{"description":"No title provided","assignedTo":"staff"}`
	req := httptest.NewRequest("POST", "/api/v1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleStaff)
	rec := httptest.NewRecorder()

	h.handleCreateTask(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleDeleteTask_Valid(t *testing.T) {
	h := newTestHandler(t)

	// Create a task first
	h.store.CreateTask(models.Task{
		Title:      "To be deleted",
		AssignedTo: "staff",
		Priority:   "low",
		Status:     "pending",
	})

	// Find the task
	tasks := h.store.ListTasks("")
	var taskID string
	for _, task := range tasks {
		if task.Title == "To be deleted" {
			taskID = task.ID
			break
		}
	}
	if taskID == "" {
		t.Fatal("could not find task to delete")
	}

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/"+taskID, nil)
	req.SetPathValue("id", taskID)
	rec := httptest.NewRecorder()

	h.handleDeleteTask(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify task is gone
	_, ok := h.store.GetTask(taskID)
	if ok {
		t.Error("task should have been deleted")
	}
}

func TestHandleDeleteTask_NotFound(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/TSK-999", nil)
	req.SetPathValue("id", "TSK-999")
	rec := httptest.NewRecorder()

	h.handleDeleteTask(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleUpdateTask_TitleDescription(t *testing.T) {
	h := newTestHandler(t)

	// Seed a task
	h.store.CreateTask(models.Task{
		Title:       "Original title",
		Description: "Original desc",
		AssignedTo:  "staff",
		Priority:    "medium",
		Status:      "pending",
	})

	tasks := h.store.ListTasks("")
	var taskID string
	for _, task := range tasks {
		if task.Title == "Original title" {
			taskID = task.ID
			break
		}
	}
	if taskID == "" {
		t.Fatal("could not find task")
	}

	body := `{"title":"Updated title","description":"Updated desc"}`
	req := httptest.NewRequest("PATCH", "/api/v1/tasks/"+taskID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", taskID)
	rec := httptest.NewRecorder()

	h.handleUpdateTask(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.Task
	json.NewDecoder(rec.Body).Decode(&updated)

	if updated.Title != "Updated title" {
		t.Errorf("expected title 'Updated title', got %q", updated.Title)
	}
	if updated.Description != "Updated desc" {
		t.Errorf("expected description 'Updated desc', got %q", updated.Description)
	}
}
