package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"
)

func (h *Handler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	assignedTo := r.URL.Query().Get("assignedTo")
	if assignedTo == "" {
		assignedTo = role
	}

	// XO and staff can see all tasks
	if auth.IsPrivileged(role) {
		if r.URL.Query().Get("all") == "true" {
			assignedTo = ""
		}
	}

	tasks := h.store.ListTasks(assignedTo)
	writeJSON(w, 200, tasks)
}

func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, ok := h.store.GetTask(id)
	if !ok {
		writeError(w, 404, "task not found")
		return
	}
	writeJSON(w, 200, task)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		writeError(w, 400, "title is required")
		return
	}

	role := middleware.GetRole(r.Context())
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   role,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		RelatedID:   req.RelatedID,
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}

	if err := h.store.CreateTask(task); err != nil {
		writeError(w, 500, "failed to create task")
		return
	}

	// Return the most recently created task
	tasks := h.store.ListTasks("")
	if len(tasks) > 0 {
		// Broadcast task creation via SSE
		if h.sseBroker != nil {
			h.sseBroker.Broadcast("", SSEEvent{
				Type: "task",
				Data: map[string]interface{}{
					"action": "created",
					"task":   tasks[0],
				},
			})
		}
		writeJSON(w, 201, tasks[0]) // tasks sorted by created_at desc in SQL, last in slice for memory
	} else {
		writeJSON(w, 201, map[string]string{"status": "created"})
	}
}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.DeleteTask(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, 404, err.Error())
		} else {
			writeError(w, 500, err.Error())
		}
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if err := h.store.UpdateTask(id, req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, 404, err.Error())
		} else {
			writeError(w, 500, err.Error())
		}
		return
	}

	task, _ := h.store.GetTask(id)

	// Broadcast task update via SSE
	if h.sseBroker != nil && task != nil {
		h.sseBroker.Broadcast("", SSEEvent{
			Type: "task",
			Data: map[string]interface{}{
				"action": "updated",
				"task":   task,
			},
		})
	}

	writeJSON(w, 200, task)
}

func (h *Handler) handleListMessages(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	userRole := r.URL.Query().Get("role")
	if userRole == "" {
		userRole = role
	}

	// XO and staff can see all messages
	if auth.IsPrivileged(role) && r.URL.Query().Get("all") == "true" {
		userRole = ""
	}

	messages := h.store.ListMessages(userRole)
	writeJSON(w, 200, messages)
}

func (h *Handler) handleMarkMessageRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.MarkMessageRead(id); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) handleListNotifications(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	unreadOnly := r.URL.Query().Get("unread") == "true"
	notifications := h.store.ListNotifications(role, unreadOnly)
	writeJSON(w, 200, notifications)
}

func (h *Handler) handleNotificationCount(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	count := h.store.UnreadNotificationCount(role)
	writeJSON(w, 200, map[string]int{"count": count})
}

func (h *Handler) handleMarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.MarkNotificationRead(id); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
