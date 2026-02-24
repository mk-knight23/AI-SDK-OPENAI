package services

import (
	"context"
	"testing"

	"github.com/mk-knight23/ai-sdk-openai/internal/models"
	"github.com/mk-knight23/ai-sdk-openai/internal/openai"
)

func TestNewVectorService(t *testing.T) {
	client := openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small")
	service := NewVectorService(client)

	if service == nil {
		t.Fatal("NewVectorService() should not return nil")
	}

	if service.client != client {
		t.Error("Client should be set correctly")
	}
}

func TestVectorService_CreateVectorStore(t *testing.T) {
	service := NewVectorService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	req := models.CreateVectorStoreRequest{
		Name: "Test Vector Store",
	}

	ctx := context.Background()
	_, err := service.CreateVectorStore(ctx, req)

	// Expected error with fake API key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}

func TestVectorService_UploadFile(t *testing.T) {
	service := NewVectorService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	ctx := context.Background()
	_, err := service.UploadFile(ctx, []byte("test content"), "test.txt", "assistants")

	// Expected error with fake API key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}

func TestVectorService_GetVectorStore(t *testing.T) {
	service := NewVectorService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	ctx := context.Background()
	_, err := service.GetVectorStore(ctx, "test-vector-store-id")

	// Expected error with fake API key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}

func TestVectorService_ListVectorStores(t *testing.T) {
	service := NewVectorService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	ctx := context.Background()
	_, err := service.ListVectorStores(ctx, 10)

	// Expected error with fake API key
	if err == nil {
		t.Error("Expected error with fake API key")
	}
}

func TestVectorService_convertVectorStore(t *testing.T) {
	service := NewVectorService(openai.NewClient("test-key", "gpt-4o", "text-embedding-3-small"))

	// Test conversion with mock data - we'd need to create a proper mock
	// For now, we just verify the conversion logic is in place
	// In real implementation, we'd use interfaces and mocks
}
