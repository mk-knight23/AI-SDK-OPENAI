package services

import (
	"context"
	"fmt"

	"github.com/mk-knight23/ai-sdk-openai/internal/models"
	"github.com/mk-knight23/ai-sdk-openai/internal/openai"
	openaiclient "github.com/sashabaranov/go-openai"
)

// ChatService handles chat completion operations
type ChatService struct {
	client *openai.Client
}

// NewChatService creates a new chat service
func NewChatService(client *openai.Client) *ChatService {
	return &ChatService{
		client: client,
	}
}

// CreateCompletion creates a non-streaming chat completion
func (s *ChatService) CreateCompletion(ctx context.Context, req models.ChatRequest) (*models.ChatResponse, error) {
	// Convert our request to OpenAI format
	openaiReq := s.convertToOpenAIRequest(req)

	// Call OpenAI
	resp, err := s.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create completion: %w", err)
	}

	// Convert response back to our format
	return s.convertFromOpenAIResponse(resp), nil
}

// CreateCompletionStream creates a streaming chat completion
func (s *ChatService) CreateCompletionStream(ctx context.Context, req models.ChatRequest) (<-chan models.StreamChunk, <-chan error) {
	chunkChan := make(chan models.StreamChunk, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		openaiReq := s.convertToOpenAIRequest(req)

		stream, err := s.client.CreateChatCompletionStream(ctx, openaiReq)
		if err != nil {
			errChan <- fmt.Errorf("failed to create stream: %w", err)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					return
				}
				errChan <- fmt.Errorf("stream error: %w", err)
				return
			}

			chunk := s.convertFromOpenAIStream(response)
			chunkChan <- chunk
		}
	}()

	return chunkChan, errChan
}

// convertToOpenAIRequest converts our request to OpenAI format
func (s *ChatService) convertToOpenAIRequest(req models.ChatRequest) openaiclient.ChatCompletionRequest {
	messages := make([]openaiclient.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openaiclient.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	openaiReq := openaiclient.ChatCompletionRequest{
		Model:       openaiclient.GPT4o,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// Set model if specified
	if req.Model != "" {
		openaiReq.Model = openaiclient.ChatCompletionModel(req.Model)
	}

	if req.Stream {
		openaiReq.Stream = true
	}

	// Convert tools
	if len(req.Tools) > 0 {
		openaiReq.Tools = make([]openaiclient.Tool, len(req.Tools))
		for i, tool := range req.Tools {
			openaiReq.Tools[i] = openaiclient.Tool{
				Type: tool.Type,
				Function: openaiclient.FunctionDefinition{
					Name:        tool.Function.Name,
					Description: tool.Function.Description,
					Parameters:  tool.Function.Parameters,
				},
			}
		}
	}

	return openaiReq
}

// convertFromOpenAIResponse converts OpenAI response to our format
func (s *ChatService) convertFromOpenAIResponse(resp openaiclient.ChatCompletionResponse) *models.ChatResponse {
	choices := make([]models.Choice, len(resp.Choices))
	for i, choice := range resp.Choices {
		choices[i] = models.Choice{
			Index: choice.Index,
			Message: models.Message{
				Role:       choice.Message.Role,
				Content:    choice.Message.Content,
				ToolCalls:  convertToolCalls(choice.Message.ToolCalls),
			},
			FinishReason: choice.FinishReason,
		}
	}

	return &models.ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: choices,
		Usage: models.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}
}

// convertFromOpenAIStream converts OpenAI stream response to our format
func (s *ChatService) convertFromOpenAIStream(resp openaiclient.ChatCompletionStreamResponse) models.StreamChunk {
	choices := make([]models.StreamChoice, len(resp.Choices))
	for i, choice := range resp.Choices {
		delta := models.MessageDelta{}
		if choice.Delta.Role != "" {
			delta.Role = choice.Delta.Role
		}
		if choice.Delta.Content != "" {
			delta.Content = choice.Delta.Content
		}
		if len(choice.Delta.ToolCalls) > 0 {
			delta.ToolCalls = convertToolCalls(choice.Delta.ToolCalls)
		}

		choices[i] = models.StreamChoice{
			Index: choice.Index,
			Delta: delta,
		}

		if choice.FinishReason != "" {
			choices[i].FinishReason = &choice.FinishReason
		}
	}

	return models.StreamChunk{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: choices,
	}
}

// convertToolCalls converts OpenAI tool calls to our format
func convertToolCalls(calls []openaiclient.ToolCall) []models.ToolCall {
	if calls == nil {
		return nil
	}

	result := make([]models.ToolCall, len(calls))
	for i, call := range calls {
		result[i] = models.ToolCall{
			ID:   call.ID,
			Type: call.Type,
			Function: models.FunctionCall{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments,
			},
		}
	}
	return result
}
