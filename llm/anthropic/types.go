package anthropic

import (
	"encoding/json"

	"github.com/giyanellow/go-llm/llm/tools"
)

type anthropicRequestBody struct {
	Model           string                 `json:"model"`
	SystemPrompt    string                 `json:"system,omitempty"`
	MaxTokens       int                    `json:"max_tokens,omitempty"`
	Input           []anthropicMessage     `json:"messages"`
	ToolsDefinition []tools.ToolDefinition `json:"tools,omitempty"`
}

type anthropicResponseBody struct {
	Model       string                 `json:"model"`
	ID          string                 `json:"id"`
	MessageType string                 `json:"type"`
	Role        string                 `json:"role"`
	Content     []anthropicContentBody `json:"content"`
	StopReason  string                 `json:"stop_reason"`
}

type anthropicContentBody struct {
	MessageType string          `json:"type"`
	Text        string          `json:"text,omitempty"`
	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Input       json.RawMessage `json:"input,omitempty"`
}

type anthropicMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type anthropicToolResultMessage struct {
	MessageType string `json:"type"`
	ToolUseID   string `json:"tool_use_id"`
	Content     string `json:"content"`
	IsError     bool   `json:"is_error,omitempty"`
}
