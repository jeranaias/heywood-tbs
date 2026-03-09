package api

import (
	"fmt"
	"net/http"
	"strings"

	"heywood-tbs/internal/data"
	"heywood-tbs/internal/middleware"
)

// chatPersister checks if the store supports chat history persistence.
func (h *Handler) chatPersister() (data.ChatPersister, bool) {
	cp, ok := h.store.(data.ChatPersister)
	return cp, ok
}

// handleListChatSessions returns the user's chat sessions.
func (h *Handler) handleListChatSessions(w http.ResponseWriter, r *http.Request) {
	cp, ok := h.chatPersister()
	if !ok {
		writeError(w, 501, "chat history not available with current data source")
		return
	}

	identity := middleware.GetIdentity(r.Context())
	sessions := cp.ListChatSessions(identity.ID, identity.Role)
	if sessions == nil {
		sessions = nil
	}

	writeJSON(w, 200, map[string]interface{}{"sessions": sessions})
}

// handleGetChatMessages returns messages for a specific chat session.
// Ownership is verified — users can only read their own sessions.
func (h *Handler) handleGetChatMessages(w http.ResponseWriter, r *http.Request) {
	cp, ok := h.chatPersister()
	if !ok {
		writeError(w, 501, "chat history not available")
		return
	}

	sessionID := r.PathValue("id")
	session, found := cp.GetChatSession(sessionID)
	if !found {
		writeError(w, 404, "session not found")
		return
	}

	identity := middleware.GetIdentity(r.Context())
	if session.UserID != identity.ID {
		writeError(w, 403, "access denied")
		return
	}

	messages := cp.GetChatMessages(sessionID)
	writeJSON(w, 200, map[string]interface{}{"messages": messages})
}

// handleDeleteChatSession deletes a chat session and its messages.
// Ownership is verified — users can only delete their own sessions.
func (h *Handler) handleDeleteChatSession(w http.ResponseWriter, r *http.Request) {
	cp, ok := h.chatPersister()
	if !ok {
		writeError(w, 501, "chat history not available")
		return
	}

	sessionID := r.PathValue("id")
	session, found := cp.GetChatSession(sessionID)
	if !found {
		writeError(w, 404, "session not found")
		return
	}

	identity := middleware.GetIdentity(r.Context())
	if session.UserID != identity.ID {
		writeError(w, 403, "access denied")
		return
	}

	if err := cp.DeleteChatSession(sessionID); err != nil {
		writeError(w, 500, "failed to delete session")
		return
	}

	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// handleExportChatSession exports a chat session as a markdown file.
// Ownership is verified — users can only export their own sessions.
func (h *Handler) handleExportChatSession(w http.ResponseWriter, r *http.Request) {
	cp, ok := h.chatPersister()
	if !ok {
		writeError(w, 501, "chat history not available")
		return
	}

	sessionID := r.PathValue("id")
	session, found := cp.GetChatSession(sessionID)
	if !found {
		writeError(w, 404, "session not found")
		return
	}

	identity := middleware.GetIdentity(r.Context())
	if session.UserID != identity.ID {
		writeError(w, 403, "access denied")
		return
	}

	messages := cp.GetChatMessages(sessionID)

	var sb strings.Builder
	title := session.Title
	if title == "" {
		title = "Chat with Heywood"
	}
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("Session: %s | Role: %s\n\n---\n\n", sessionID, session.UserRole))

	for _, msg := range messages {
		if msg.Role == "user" {
			sb.WriteString(fmt.Sprintf("**You:** %s\n\n", msg.Content))
		} else {
			sb.WriteString(fmt.Sprintf("**Heywood:** %s\n\n", msg.Content))
		}
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"heywood-chat-%s.md\"", sessionID))
	w.WriteHeader(200)
	w.Write([]byte(sb.String()))
}
