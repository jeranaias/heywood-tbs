package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/models"

	openai "github.com/sashabaranov/go-openai"
)

func (h *Handler) handleChat(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		writeError(w, 400, "message is required")
		return
	}

	role := middleware.GetRole(r.Context())
	company := middleware.GetCompany(r.Context())
	studentID := middleware.GetStudentID(r.Context())

	// Resolve or create chat session for persistence
	sessionID := strings.TrimSpace(req.SessionID)
	cp, hasPersistence := h.chatPersister()
	if hasPersistence && sessionID == "" {
		sessionID = fmt.Sprintf("chat-%d", time.Now().UnixNano())
		identity := middleware.GetIdentity(r.Context())
		cp.CreateChatSession(models.ChatSession{
			ID: sessionID, UserID: identity.ID, UserRole: role, Company: company,
		})
	}

	// Build system prompt and context based on role
	systemPrompt, userContext := h.buildChatContext(role, company, studentID, req.Message)

	// If no OpenAI client, use mock responses
	if h.chatSvc == nil {
		response := h.mockResponse(role, company, studentID, req.Message)
		if hasPersistence && sessionID != "" {
			cp.AddChatMessage(sessionID, models.ChatMessage{Role: "user", Content: req.Message})
			cp.AddChatMessage(sessionID, models.ChatMessage{Role: "assistant", Content: response})
		}
		writeJSON(w, 200, models.ChatResponse{Response: response, SessionID: sessionID})
		return
	}

	// Build messages for OpenAI
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Add history (limited to last 20 messages)
	historyStart := 0
	if len(req.History) > 20 {
		historyStart = len(req.History) - 20
	}
	for _, msg := range req.History[historyStart:] {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add current message with injected context
	userMessage := req.Message
	if userContext != "" {
		userMessage = req.Message + "\n\n---\n[Relevant data context]\n" + userContext
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: userMessage,
	})

	// Tools are available for all roles
	tools := ai.HeywoodTools

	// Streaming mode
	if req.Stream {
		h.handleStreamingChat(w, r, messages, tools, role, company)
		return
	}

	// Non-streaming mode
	chatReq := openai.ChatCompletionRequest{
		Model:       h.chatSvc.model,
		Messages:    messages,
		MaxTokens:   4096,
		Temperature: 0.7,
	}
	if len(tools) > 0 {
		chatReq.Tools = tools
	}

	resp, err := h.chatSvc.client.CreateChatCompletion(r.Context(), chatReq)
	if err != nil {
		slog.Error("openai error", "error", err)
		response := h.mockResponse(role, company, studentID, req.Message)
		writeJSON(w, 200, models.ChatResponse{Response: response + "\n\n*(Note: AI service temporarily unavailable — showing cached response)*"})
		return
	}

	if len(resp.Choices) == 0 {
		writeError(w, 500, "no response from AI")
		return
	}

	// Handle tool calls — execute tools and re-send for final answer
	choice := resp.Choices[0]
	if choice.FinishReason == openai.FinishReasonToolCalls && len(choice.Message.ToolCalls) > 0 {
		messages = append(messages, choice.Message)
		for _, tc := range choice.Message.ToolCalls {
			result := h.executeToolCall(tc, role, company)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}
		// Re-send without tools to get final response
		resp2, err := h.chatSvc.client.CreateChatCompletion(r.Context(), openai.ChatCompletionRequest{
			Model:       h.chatSvc.model,
			Messages:    messages,
			MaxTokens:   4096,
			Temperature: 0.7,
		})
		if err != nil {
			slog.Error("openai tool follow-up error", "error", err)
			writeError(w, 500, "AI follow-up failed")
			return
		}
		if len(resp2.Choices) > 0 {
			content := resp2.Choices[0].Message.Content
			if hasPersistence && sessionID != "" {
				cp.AddChatMessage(sessionID, models.ChatMessage{Role: "user", Content: req.Message})
				cp.AddChatMessage(sessionID, models.ChatMessage{Role: "assistant", Content: content})
			}
			writeJSON(w, 200, models.ChatResponse{Response: content, SessionID: sessionID})
			return
		}
	}

	responseContent := choice.Message.Content
	if hasPersistence && sessionID != "" {
		cp.AddChatMessage(sessionID, models.ChatMessage{Role: "user", Content: req.Message})
		cp.AddChatMessage(sessionID, models.ChatMessage{Role: "assistant", Content: responseContent})
	}
	writeJSON(w, 200, models.ChatResponse{
		Response:  responseContent,
		SessionID: sessionID,
	})
}
