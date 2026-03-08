package api

import (
	"net/http"

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
