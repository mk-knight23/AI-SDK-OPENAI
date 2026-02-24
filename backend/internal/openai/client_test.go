package openai

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	model := "gpt-4o"
	embeddingModel := "text-embedding-3-small"

	client := NewClient(apiKey, model, embeddingModel)

	if client == nil {
		t.Fatal("NewClient() should not return nil")
	}

	if client.model != model {
		t.Errorf("model = %v, want %v", client.model, model)
	}

	if client.embeddingModel != embeddingModel {
		t.Errorf("embeddingModel = %v, want %v", client.embeddingModel, embeddingModel)
	}
}

func TestClient_GetModel(t *testing.T) {
	model := "gpt-4o-mini"
	client := NewClient("test-key", model, "text-embedding-3-small")

	if client.GetModel() != model {
		t.Errorf("GetModel() = %v, want %v", client.GetModel(), model)
	}
}

func TestClient_GetEmbeddingModel(t *testing.T) {
	embeddingModel := "text-embedding-3-large"
	client := NewClient("test-key", "gpt-4o", embeddingModel)

	if client.GetEmbeddingModel() != embeddingModel {
		t.Errorf("GetEmbeddingModel() = %v, want %v", client.GetEmbeddingModel(), embeddingModel)
	}
}

func TestClient_CreateChatCompletion(t *testing.T) {
	// This test verifies the client method exists and has correct signature
	// Actual API calls should be tested with integration tests
	client := NewClient("test-key", "gpt-4o", "text-embedding-3-small")

	ctx := context.Background()
	req := map[string]any{
		"messages": []map[string]any{
			{"role": "user", "content": "test"},
		},
	}

	// We expect this to fail with an invalid API key, but we're testing
	// that the method signature and call structure is correct
	_, err := client.CreateChatCompletion(ctx, req)

	// The error is expected with a fake key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}

func TestClient_CreateEmbedding(t *testing.T) {
	client := NewClient("test-key", "gpt-4o", "text-embedding-3-small")

	ctx := context.Background()
	text := "test text for embedding"

	_, err := client.CreateEmbedding(ctx, text)

	// Expected error with fake API key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}
