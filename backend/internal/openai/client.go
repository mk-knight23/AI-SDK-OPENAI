package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI client with our application-specific methods
type Client struct {
	client          *openai.Client
	model           string
	embeddingModel  string
}

// NewClient creates a new OpenAI client wrapper
func NewClient(apiKey, model, embeddingModel string) *Client {
	return &Client{
		client:         openai.NewClient(apiKey),
		model:          model,
		embeddingModel: embeddingModel,
	}
}

// GetModel returns the default model
func (c *Client) GetModel() string {
	return c.model
}

// GetEmbeddingModel returns the default embedding model
func (c *Client) GetEmbeddingModel() string {
	return c.embeddingModel
}

// GetClient returns the underlying OpenAI client
func (c *Client) GetClient() *openai.Client {
	return c.client
}

// CreateChatCompletion creates a non-streaming chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if req.Model == "" {
		req.Model = c.model
	}
	return c.client.CreateChatCompletion(ctx, req)
}

// CreateChatCompletionStream creates a streaming chat completion
func (c *Client) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionStream, error) {
	if req.Model == "" {
		req.Model = c.model
	}
	return c.client.CreateChatCompletionStream(ctx, req)
}

// CreateEmbedding creates an embedding for the given text
func (c *Client) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(c.embeddingModel),
	}

	resp, err := c.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return resp.Data[0].Embedding, nil
}

// CreateFile uploads a file for use with the Assistants API
func (c *Client) CreateFile(ctx context.Context, req openai.FileRequest) (openai.File, error) {
	return c.client.CreateFile(ctx, req)
}

// RetrieveFile gets a file by ID
func (c *Client) RetrieveFile(ctx context.Context, fileID string) (openai.File, error) {
	// Note: RetrieveFile method not available in current go-openai version
	// Use ListFiles and filter by ID instead
	files, err := c.client.ListFiles(ctx)
	if err != nil {
		return openai.File{}, fmt.Errorf("failed to retrieve file: %w", err)
	}
	for _, f := range files.Files {
		if f.ID == fileID {
			return f, nil
		}
	}
	return openai.File{}, fmt.Errorf("file not found: %s", fileID)
}

// DeleteFile deletes a file
func (c *Client) DeleteFile(ctx context.Context, fileID string) error {
	return c.client.DeleteFile(ctx, fileID)
}

// ListFiles lists all files
func (c *Client) ListFiles(ctx context.Context, purpose string) (openai.FilesList, error) {
	files, err := c.client.ListFiles(ctx)
	if err != nil {
		return openai.FilesList{}, err
	}
	// Filter by purpose if provided
	if purpose != "" {
		var filtered []openai.File
		for _, f := range files.Files {
			if f.Purpose == purpose {
				filtered = append(filtered, f)
			}
		}
		files.Files = filtered
	}
	return files, nil
}
