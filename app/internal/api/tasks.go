package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"heywood-tbs/internal/middleware"
)

func (h *Handler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	assignedTo := r.URL.Query().Get("assignedTo")
	if assignedTo == "" {
		assignedTo = role
	}

	// XO and staff can see all tasks
	if role == "xo" || role == "staff" {
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

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	// Only allow status, priority, assignedTo updates
	allowed := map[string]bool{"status": true, "priority": true, "assignedTo": true}
	for k := range updates {
		if !allowed[k] {
			delete(updates, k)
		}
	}

	if err := h.store.UpdateTask(id, updates); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, 404, err.Error())
		} else {
			writeError(w, 500, err.Error())
		}
		return
	}

	task, _ := h.store.GetTask(id)
	writeJSON(w, 200, task)
}

func (h *Handler) handleListMessages(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	userRole := r.URL.Query().Get("role")
	if userRole == "" {
		userRole = role
	}

	// XO and staff can see all messages
	if (role == "xo" || role == "staff") && r.URL.Query().Get("all") == "true" {
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
