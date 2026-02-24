package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mk-knight23/ai-sdk-openai/internal/models"
	"github.com/mk-knight23/ai-sdk-openai/internal/services"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	chatService *services.ChatService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// CreateCompletion handles POST /api/chat/completions
func (h *ChatHandler) CreateCompletion(c *fiber.Ctx) error {
	var req models.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := h.validateChatRequest(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()
	resp, err := h.chatService.CreateCompletion(ctx, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create completion",
			"details": err.Error(),
		})
	}

	return c.JSON(resp)
}

// CreateCompletionStream handles POST /api/chat/completions/stream
func (h *ChatHandler) CreateCompletionStream(c *fiber.Ctx) error {
	var req models.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.Stream = true

	if err := h.validateChatRequest(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	ctx := context.Background()
	chunkChan, errChan := h.chatService.CreateCompletionStream(ctx, req)

	c.Context().SetBodyStreamWriter(func(w *fiber.StreamWriter) {
		for {
			select {
			case chunk, ok := <-chunkChan:
				if !ok {
					w.WriteString("data: [DONE]\n\n")
					return
				}
				data := string(c.App().Config().Encoder().Encode(chunk))
				w.WriteString("data: " + data + "\n\n")
			case err := <-errChan:
				if err != nil {
					w.WriteString("event: error\ndata: " + err.Error() + "\n\n")
				}
				return
			case <-c.Context().Done():
				return
			}
		}
	})

	return nil
}

// validateChatRequest validates a chat request
func (h *ChatHandler) validateChatRequest(req *models.ChatRequest) error {
	if len(req.Messages) == 0 {
		return fiber.NewError(400, "Messages are required")
	}

	for i, msg := range req.Messages {
		if msg.Role == "" {
			return fiber.NewError(400, "Message role is required")
		}
		if msg.Content == "" && len(msg.ToolCalls) == 0 {
			return fiber.NewError(400, "Message content or tool calls are required")
		}
		if msg.Role != "system" && msg.Role != "user" && msg.Role != "assistant" && msg.Role != "tool" {
			return fiber.NewError(400, "Invalid message role: " + msg.Role)
		}
		if msg.ToolCallID != "" && msg.Role != "tool" {
			return fiber.NewError(400, "ToolCallID is only valid for tool role")
		}
		if i == 0 && msg.Role != "system" && msg.Role != "user" {
			return fiber.NewError(400, "First message must be system or user")
		}
	}

	if req.MaxTokens < 0 {
		return fiber.NewError(400, "MaxTokens cannot be negative")
	}

	if req.Temperature < 0 || req.Temperature > 2 {
		return fiber.NewError(400, "Temperature must be between 0 and 2")
	}

	return nil
}
