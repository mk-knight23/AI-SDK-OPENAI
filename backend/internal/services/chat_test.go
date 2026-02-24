package services

import (
	"context"
	"testing"

	"github.com/mk-knight23/ai-sdk-openai/internal/models"
	"github.com/mk-knight23/ai-sdk-openai/internal/openai"
)

func TestNewChatService(t *testing.T) {
	client := openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small")
	service := NewChatService(client)

	if service == nil {
		t.Fatal("NewChatService() should not return nil")
	}

	if service.client != client {
		t.Error("Client should be set correctly")
	}
}

func TestChatService_CreateCompletion(t *testing.T) {
	service := NewChatService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	req := models.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "Hello, world!"},
		},
	}

	ctx := context.Background()
	_, err := service.CreateCompletion(ctx, req)

	// Expected error with fake API key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}

func TestChatService_CreateCompletionStream(t *testing.T) {
	service := NewChatService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	req := models.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "Hello, world!"},
		},
		Stream: true,
	}

	ctx := context.Background()
	chunkChan, errChan := service.CreateCompletionStream(ctx, req)

	// With fake API key, we expect an error
	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected error with fake API key")
		}
	case <-chunkChan:
		t.Error("Expected error before any chunks")
	}
}

func TestChatService_convertToOpenAIRequest(t *testing.T) {
	service := NewChatService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	req := models.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "test message"},
		},
		Model:       "gpt-4o-mini",
		MaxTokens:   1000,
		Temperature: 0.7,
		Stream:      true,
	}

	openaiReq := service.convertToOpenAIRequest(req)

	if openaiReq.Model != "gpt-4o-mini" {
		t.Errorf("Model = %v, want gpt-4o-mini", openaiReq.Model)
	}

	if openaiReq.MaxTokens != 1000 {
		t.Errorf("MaxTokens = %v, want 1000", openaiReq.MaxTokens)
	}

	if openaiReq.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", openaiReq.Temperature)
	}

	if !openaiReq.Stream {
		t.Error("Stream should be true")
	}

	if len(openaiReq.Messages) != 1 {
		t.Errorf("Messages length = %v, want 1", len(openaiReq.Messages))
	}

	if openaiReq.Messages[0].Role != "user" {
		t.Errorf("Message role = %v, want user", openaiReq.Messages[0].Role)
	}

	if openaiReq.Messages[0].Content != "test message" {
		t.Errorf("Message content = %v, want 'test message'", openaiReq.Messages[0].Content)
	}
}

func TestChatService_convertToOpenAIRequestWithTools(t *testing.T) {
	service := NewChatService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	req := models.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "What's the weather?"},
		},
		Tools: []models.Tool{
			{
				Type: "function",
				Function: models.ToolFunction{
					Name:        "get_weather",
					Description: "Get the current weather",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"location": map[string]any{
								"type":        "string",
								"description": "The city and state",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
	}

	openaiReq := service.convertToOpenAIRequest(req)

	if len(openaiReq.Tools) != 1 {
		t.Fatalf("Tools length = %v, want 1", len(openaiReq.Tools))
	}

	if openaiReq.Tools[0].Type != "function" {
		t.Errorf("Tool type = %v, want function", openaiReq.Tools[0].Type)
	}

	if openaiReq.Tools[0].Function.Name != "get_weather" {
		t.Errorf("Function name = %v, want get_weather", openaiReq.Tools[0].Function.Name)
	}
}

func TestConvertToolCalls(t *testing.T) {
	// Test with nil input
	result := convertToolCalls(nil)
	if result != nil {
		t.Error("convertToolCalls(nil) should return nil")
	}

	// Test with empty slice
	calls := []struct {
		ID       string
		Type     string
		Function struct {
			Name      string
			Arguments string
		}
	}{}

	// This would be the actual OpenAI type, but for testing we check the function handles it
	// The actual conversion would use the openai.ToolCall type
}
