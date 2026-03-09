package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/data"
	"heywood-tbs/internal/models"
)

// testChatStore wraps a DataStore and adds ChatPersister capabilities for testing.
type testChatStore struct {
	data.DataStore
	sessions map[string]models.ChatSession
	messages map[string][]models.ChatMessage
}

func (s *testChatStore) CreateChatSession(sess models.ChatSession) error {
	s.sessions[sess.ID] = sess
	return nil
}

func (s *testChatStore) ListChatSessions(userID, role string) []models.ChatSession {
	var result []models.ChatSession
	for _, sess := range s.sessions {
		if sess.UserID == userID {
			result = append(result, sess)
		}
	}
	return result
}

func (s *testChatStore) GetChatSession(id string) (*models.ChatSession, bool) {
	sess, ok := s.sessions[id]
	if !ok {
		return nil, false
	}
	return &sess, true
}

func (s *testChatStore) UpdateChatSessionTitle(id, title string) error {
	if sess, ok := s.sessions[id]; ok {
		sess.Title = title
		s.sessions[id] = sess
	}
	return nil
}

func (s *testChatStore) AddChatMessage(sessionID string, msg models.ChatMessage) error {
	s.messages[sessionID] = append(s.messages[sessionID], msg)
	return nil
}

func (s *testChatStore) GetChatMessages(sessionID string) []models.ChatMessage {
	return s.messages[sessionID]
}

func (s *testChatStore) DeleteChatSession(id string) error {
	delete(s.sessions, id)
	delete(s.messages, id)
	return nil
}

func newTestHandlerWithChat(t *testing.T) (*Handler, *testChatStore) {
	t.Helper()
	h := newTestHandler(t)
	cs := &testChatStore{
		DataStore: h.store,
		sessions: map[string]models.ChatSession{
			"sess-1": {ID: "sess-1", UserID: "demo-staff", UserRole: "staff", Title: "Test Chat"},
		},
		messages: map[string][]models.ChatMessage{
			"sess-1": {
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi, how can I help?"},
			},
		},
	}
	h.store = cs
	return h, cs
}

func TestExportChatSession_Ownership(t *testing.T) {
	h, _ := newTestHandlerWithChat(t)

	t.Run("owner can export", func(t *testing.T) {
		// sess-1 is owned by "demo-staff"
		ctx := withIdentityContext(&auth.UserIdentity{
			ID: "demo-staff", Role: auth.RoleStaff, Source: "demo",
		}, "")

		req := httptest.NewRequest("GET", "/api/v1/chat/sessions/sess-1/export", nil).WithContext(ctx)
		req.SetPathValue("id", "sess-1")
		rec := httptest.NewRecorder()

		h.handleExportChatSession(rec, req)

		if rec.Code != 200 {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}

		ct := rec.Header().Get("Content-Type")
		if ct != "text/markdown; charset=utf-8" {
			t.Errorf("expected markdown content-type, got %q", ct)
		}

		cd := rec.Header().Get("Content-Disposition")
		if cd == "" {
			t.Error("expected Content-Disposition header")
		}

		body := rec.Body.String()
		if !containsCI(body, "Hello") || !containsCI(body, "how can I help") {
			t.Error("export should contain chat messages")
		}
		if !containsCI(body, "Test Chat") {
			t.Error("export should contain session title")
		}
	})

	t.Run("non-owner gets 403", func(t *testing.T) {
		ctx := withIdentityContext(&auth.UserIdentity{
			ID: "demo-xo", Role: auth.RoleXO, Source: "demo",
		}, "")

		req := httptest.NewRequest("GET", "/api/v1/chat/sessions/sess-1/export", nil).WithContext(ctx)
		req.SetPathValue("id", "sess-1")
		rec := httptest.NewRecorder()

		h.handleExportChatSession(rec, req)

		if rec.Code != 403 {
			t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
		}

		var errResp map[string]string
		json.NewDecoder(rec.Body).Decode(&errResp)
		if errResp["error"] != "access denied" {
			t.Errorf("expected 'access denied' error, got %q", errResp["error"])
		}
	})

	t.Run("missing session gets 404", func(t *testing.T) {
		ctx := withIdentityContext(&auth.UserIdentity{
			ID: "demo-staff", Role: auth.RoleStaff, Source: "demo",
		}, "")

		req := httptest.NewRequest("GET", "/api/v1/chat/sessions/nonexistent/export", nil).WithContext(ctx)
		req.SetPathValue("id", "nonexistent")
		rec := httptest.NewRecorder()

		h.handleExportChatSession(rec, req)

		if rec.Code != 404 {
			t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
		}
	})
}

func TestExportChatSession_NoPersistence(t *testing.T) {
	// Standard test handler has no ChatPersister (JSON-backed store)
	h := newTestHandler(t)

	ctx := withIdentityContext(&auth.UserIdentity{
		ID: "demo-staff", Role: auth.RoleStaff, Source: "demo",
	}, "")

	req := httptest.NewRequest("GET", "/api/v1/chat/sessions/sess-1/export", nil).WithContext(ctx)
	req.SetPathValue("id", "sess-1")
	rec := httptest.NewRecorder()

	h.handleExportChatSession(rec, req)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", rec.Code)
	}
}
