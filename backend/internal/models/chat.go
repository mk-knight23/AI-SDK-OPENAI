package models

import "time"

// Message represents a chat message
type Message struct {
	ID           string         `json:"id"`
	Role         string         `json:"role"` // "user", "assistant", "system"
	Content      string         `json:"content"`
	Timestamp    time.Time      `json:"timestamp"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	ToolCalls    []ToolCall     `json:"tool_calls,omitempty"`
	ToolCallID   string         `json:"tool_call_id,omitempty"`
}

// ToolCall represents a function/tool call requested by the model
type ToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"` // "function"
	Function FunctionCall        `json:"function"`
}

// FunctionCall represents a function call with its arguments
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatRequest represents a request to create a chat completion
type ChatRequest struct {
	Messages       []Message       `json:"messages"`
	Model          string          `json:"model,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Temperature    float64         `json:"temperature,omitempty"`
	Stream         bool            `json:"stream,omitempty"`
	Tools          []Tool          `json:"tools,omitempty"`
	ToolChoice     any             `json:"tool_choice,omitempty"` // string, map[string]any, or nil
}

// Tool represents a function that can be called by the model
type Tool struct {
	Type     string       `json:"type"` // "function"
	Function ToolFunction `json:"function"`
}

// ToolFunction describes a function that can be called
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]any         `json:"parameters"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk represents a chunk of streamed response
type StreamChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice represents a choice in a streaming response
type StreamChoice struct {
	Index        int           `json:"index"`
	Delta        MessageDelta  `json:"delta"`
	FinishReason *string       `json:"finish_reason"`
}

// MessageDelta represents the incremental message update in streaming
type MessageDelta struct {
	Role      string    `json:"role,omitempty"`
	Content   string    `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// FunctionOutput represents the result of a function call
type FunctionOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Role       string `json:"role"` // "tool"
	Name       string `json:"name,omitempty"`
	Content    string `json:"content"`
}
