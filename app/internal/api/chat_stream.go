package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

// handleStreamingChat sends the response as Server-Sent Events.
// It makes a single streaming call with tools enabled. If the model returns
// tool calls, they are executed and a follow-up is streamed without tools.
func (h *Handler) handleStreamingChat(w http.ResponseWriter, r *http.Request, messages []openai.ChatCompletionMessage, tools []openai.Tool, role, company string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, 500, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	// Build streaming request — include tools so we can detect tool calls
	// in the same stream rather than making a separate non-streaming call.
	streamReq := openai.ChatCompletionRequest{
		Model:       h.chatSvc.model,
		Messages:    messages,
		MaxTokens:   4096,
		Temperature: 0.7,
		Stream:      true,
	}
	if len(tools) > 0 {
		streamReq.Tools = tools
	}

	stream, err := h.chatSvc.client.CreateChatCompletionStream(r.Context(), streamReq)
	if err != nil {
		slog.Error("stream create error", "error", err)
		data, _ := json.Marshal(map[string]string{"error": "stream failed"})
		fmt.Fprintf(w, "data: %s\n\n", data)
		fmt.Fprintf(w, "data: [DONE]\n\n")
		flusher.Flush()
		return
	}
	defer stream.Close()

	// Accumulate tool call deltas in case the model invokes tools.
	var toolCalls []openai.ToolCall
	var finishReason openai.FinishReason

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			slog.Error("stream recv error", "error", err)
			break
		}

		if len(response.Choices) == 0 {
			continue
		}

		choice := response.Choices[0]
		finishReason = choice.FinishReason

		// Stream content tokens to the client immediately.
		if chunk := choice.Delta.Content; chunk != "" {
			data, _ := json.Marshal(map[string]string{"content": chunk})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}

		// Accumulate tool call deltas.
		if len(choice.Delta.ToolCalls) > 0 {
			toolCalls = mergeToolCallDeltas(toolCalls, choice.Delta.ToolCalls)
		}
	}

	// If the model finished with tool_calls, execute them and stream the follow-up.
	if finishReason == openai.FinishReasonToolCalls && len(toolCalls) > 0 {
		// Append the assistant message with accumulated tool calls.
		assistantMsg := openai.ChatCompletionMessage{
			Role:      "assistant",
			ToolCalls: toolCalls,
		}
		messages = append(messages, assistantMsg)

		for _, tc := range toolCalls {
			result := h.executeToolCall(tc, role, company)
			slog.Info("tool call executed", "tool", tc.Function.Name, "id", tc.ID)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}

		// Stream the final response (no tools on follow-up).
		h.streamMessages(w, r, flusher, messages)
		return
	}

	// Normal completion — stream already sent all content above.
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// streamMessages opens a streaming connection and sends SSE chunks.
func (h *Handler) streamMessages(w http.ResponseWriter, r *http.Request, flusher http.Flusher, messages []openai.ChatCompletionMessage) {
	stream, err := h.chatSvc.client.CreateChatCompletionStream(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:       h.chatSvc.model,
			Messages:    messages,
			MaxTokens:   4096,
			Temperature: 0.7,
			Stream:      true,
		},
	)
	if err != nil {
		slog.Error("stream create error", "error", err)
		data, _ := json.Marshal(map[string]string{"error": "stream failed"})
		fmt.Fprintf(w, "data: %s\n\n", data)
		fmt.Fprintf(w, "data: [DONE]\n\n")
		flusher.Flush()
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}
		if err != nil {
			slog.Error("stream recv error", "error", err)
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}

		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			if chunk != "" {
				data, _ := json.Marshal(map[string]string{"content": chunk})
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}

// mergeToolCallDeltas accumulates streaming tool call deltas into complete ToolCall objects.
// OpenAI sends tool calls as incremental chunks:
//   - First chunk for an index: has ID, function name, and partial arguments
//   - Subsequent chunks for the same index: append to arguments
func mergeToolCallDeltas(accumulated []openai.ToolCall, deltas []openai.ToolCall) []openai.ToolCall {
	for _, delta := range deltas {
		idx := 0
		if delta.Index != nil {
			idx = *delta.Index
		}

		// Grow the slice if needed.
		for len(accumulated) <= idx {
			accumulated = append(accumulated, openai.ToolCall{
				Type: openai.ToolTypeFunction,
			})
		}

		// Merge fields: ID and function name are sent once on the first chunk.
		if delta.ID != "" {
			accumulated[idx].ID = delta.ID
		}
		if delta.Type != "" {
			accumulated[idx].Type = delta.Type
		}
		if delta.Function.Name != "" {
			accumulated[idx].Function.Name = delta.Function.Name
		}
		// Arguments are sent as incremental partial JSON strings.
		accumulated[idx].Function.Arguments += delta.Function.Arguments
	}
	return accumulated
}
