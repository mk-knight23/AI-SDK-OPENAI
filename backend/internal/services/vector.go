package services

import (
	"context"
	"fmt"

	"github.com/mk-knight23/ai-sdk-openai/internal/models"
	"github.com/mk-knight23/ai-sdk-openai/internal/openai"
	"github.com/sashabaranov/go-openai"
)

// VectorService handles file operations for RAG
type VectorService struct {
	client *openai.Client
}

// NewVectorService creates a new vector service
func NewVectorService(client *openai.Client) *VectorService {
	return &VectorService{
		client: client,
	}
}

// UploadFile uploads a file
func (s *VectorService) UploadFile(ctx context.Context, fileData []byte, filename, purpose string) (*models.File, error) {
	req := openai.FileRequest{
		FileName: filename,
		Bytes:    fileData,
		Purpose:  openai.Purpose(purpose),
	}

	file, err := s.client.CreateFile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return s.convertFile(file), nil
}

// GetFile retrieves a file by ID
func (s *VectorService) GetFile(ctx context.Context, fileID string) (*models.File, error) {
	file, err := s.client.RetrieveFile(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file: %w", err)
	}

	return s.convertFile(file), nil
}

// DeleteFile deletes a file
func (s *VectorService) DeleteFile(ctx context.Context, fileID string) error {
	return s.client.DeleteFile(ctx, fileID)
}

// ListFiles lists all files
func (s *VectorService) ListFiles(ctx context.Context, purpose string) ([]models.File, error) {
	list, err := s.client.ListFiles(ctx, purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]models.File, len(list.Files))
	for i, f := range list.Files {
		files[i] = *s.convertFile(&f)
	}

	return files, nil
}

// convertFile converts OpenAI file to our format
func (s *VectorService) convertFile(f *openai.File) *models.File {
	return &models.File{
		ID:        f.ID,
		Object:    f.Object,
		Bytes:     f.Bytes,
		CreatedAt: f.CreatedAt,
		Filename:  f.FileName,
		Purpose:   string(f.Purpose),
	}
}
