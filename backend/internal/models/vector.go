package models

// VectorStore represents an OpenAI vector store
type VectorStore struct {
	ID           string            `json:"id"`
	Object       string            `json:"object"`
	CreatedAt    int64             `json:"created_at"`
	Name         string            `json:"name,omitempty"`
	UsageBytes   int64             `json:"usage_bytes"`
	FileCounts   FileCounts        `json:"file_counts"`
	Status       string            `json:"status"` // "active", "in_progress", "completed"
	ExpiresAfter *ExpiresAfter     `json:"expires_after,omitempty"`
	ExpiresAt    *int64            `json:"expires_at,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// FileCounts represents file count statistics
type FileCounts struct {
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Total      int `json:"total"`
}

// ExpiresAfter represents expiration policy
type ExpiresAfter struct {
	Anchor string `json:"anchor"` // "last_active_at"
	Days   int    `json:"days"`
}

// CreateVectorStoreRequest represents a request to create a vector store
type CreateVectorStoreRequest struct {
	Name         string            `json:"name,omitempty"`
	FileIDs      []string          `json:"file_ids,omitempty"`
	ExpiresAfter *ExpiresAfter     `json:"expires_after,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// VectorStoreFile represents a file in a vector store
type VectorStoreFile struct {
	ID          string    `json:"id"`
	Object      string    `json:"object"`
	CreatedAt   int64     `json:"created_at"`
	VectorStoreID string  `json:"vector_store_id"`
	Status      string    `json:"status"` // "in_progress", "completed", "cancelled", "failed"
	LastError   *FileError `json:"last_error,omitempty"`
}

// FileError represents a file error
type FileError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// File represents an uploaded file
type File struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Bytes     int64     `json:"bytes"`
	CreatedAt int64     `json:"created_at"`
	Filename  string    `json:"filename"`
	Purpose   string    `json:"purpose"` // "assistants", "assistants_output"
}

// CreateFileRequest represents a request to upload a file
type CreateFileRequest struct {
	File     []byte `json:"file"`
	Filename string `json:"filename"`
	Purpose  string `json:"purpose"`
}

// RAGQueryRequest represents a RAG (Retrieval Augmented Generation) query
type RAGQueryRequest struct {
	Query        string   `json:"query"`
	VectorStoreID string  `json:"vector_store_id"`
	TopK         int      `json:"top_k,omitempty"`
}

// RAGQueryResponse represents the response to a RAG query
type RAGQueryResponse struct {
	Query         string    `json:"query"`
	Context       string    `json:"context"`
	SourceFiles   []string  `json:"source_files"`
	Answer        string    `json:"answer,omitempty"`
}
