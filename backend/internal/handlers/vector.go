package handlers

import (
	"context"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/mk-knight23/ai-sdk-openai/internal/models"
	"github.com/mk-knight23/ai-sdk-openai/internal/services"
)

// VectorHandler handles file-related HTTP requests
type VectorHandler struct {
	vectorService *services.VectorService
}

// NewVectorHandler creates a new vector handler
func NewVectorHandler(vectorService *services.VectorService) *VectorHandler {
	return &VectorHandler{
		vectorService: vectorService,
	}
}

// UploadFile handles POST /api/files
func (h *VectorHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "File is required",
			"details": err.Error(),
		})
	}

	purpose := c.FormValue("purpose", "assistants")

	// Read file content
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to open file",
			"details": err.Error(),
		})
	}
	defer fileContent.Close()

	bytes, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to read file",
			"details": err.Error(),
		})
	}

	ctx := context.Background()
	uploadedFile, err := h.vectorService.UploadFile(ctx, bytes, file.Filename, purpose)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to upload file",
			"details": err.Error(),
		})
	}

	return c.JSON(uploadedFile)
}

// GetFile handles GET /api/files/:id
func (h *VectorHandler) GetFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "File ID is required",
		})
	}

	ctx := context.Background()
	file, err := h.vectorService.GetFile(ctx, fileID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve file",
			"details": err.Error(),
		})
	}

	return c.JSON(file)
}

// ListFiles handles GET /api/files
func (h *VectorHandler) ListFiles(c *fiber.Ctx) error {
	purpose := c.Query("purpose", "assistants")

	ctx := context.Background()
	files, err := h.vectorService.ListFiles(ctx, purpose)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to list files",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"files": files,
		"count": len(files),
	})
}

// DeleteFile handles DELETE /api/files/:id
func (h *VectorHandler) DeleteFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "File ID is required",
		})
	}

	ctx := context.Background()
	err := h.vectorService.DeleteFile(ctx, fileID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete file",
			"details": err.Error(),
		})
	}

	return c.Status(204).Send(nil)
}
