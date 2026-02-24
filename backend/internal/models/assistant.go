package models

// Assistant represents an OpenAI assistant
type Assistant struct {
	ID           string            `json:"id"`
	Object       string            `json:"object"`
	CreatedAt    int64             `json:"created_at"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Model        string            `json:"model"`
	Instructions string            `json:"instructions,omitempty"`
	Tools        []AssistantTool   `json:"tools,omitempty"`
	ToolResources *ToolResources   `json:"tool_resources,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Temperature  *float64          `json:"temperature,omitempty"`
	TopP         *float64          `json:"top_p,omitempty"`
}

// AssistantTool represents a tool available to an assistant
type AssistantTool struct {
	Type     string              `json:"type"` // "file_search", "function", "code_interpreter"
	Function *ToolFunction       `json:"function,omitempty"`
}

// ToolResources contains resources for assistant tools
type ToolResources struct {
	FileSearch *FileSearchResources `json:"file_search,omitempty"`
	CodeInterpreter *CodeInterpreterResources `json:"code_interpreter,omitempty"`
}

// FileSearchResources contains vector store IDs for file search
type FileSearchResources struct {
	VectorStoreIDs []string `json:"vector_store_ids,omitempty"`
	VectorStores   []VectorStoreCreation `json:"vector_stores,omitempty"`
}

// VectorStoreCreation specifies a new vector store to create
type VectorStoreCreation struct {
	FileIDs []string `json:"file_ids"`
}

// CodeInterpreterResources contains file IDs for code interpreter
type CodeInterpreterResources struct {
	FileIDs []string `json:"file_ids,omitempty"`
}

// CreateAssistantRequest represents a request to create an assistant
type CreateAssistantRequest struct {
	Model        string            `json:"model"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Instructions string            `json:"instructions,omitempty"`
	Tools        []AssistantTool   `json:"tools,omitempty"`
	ToolResources *ToolResources   `json:"tool_resources,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Temperature  *float64          `json:"temperature,omitempty"`
	TopP         *float64          `json:"top_p,omitempty"`
}

// Thread represents a conversation thread
type Thread struct {
	ID            string            `json:"id"`
	Object        string            `json:"object"`
	CreatedAt     int64             `json:"created_at"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// CreateThreadRequest represents a request to create a thread
type CreateThreadRequest struct {
	Messages     []ThreadMessage   `json:"messages,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ThreadMessage represents a message in a thread
type ThreadMessage struct {
	Role       string                `json:"role"` // "user"
	Content    string                `json:"content"`
	Attachments []MessageAttachment  `json:"attachments,omitempty"`
	Metadata   map[string]string     `json:"metadata,omitempty"`
}

// MessageAttachment represents a file attachment
type MessageAttachment struct {
	FileID     string `json:"file_id"`
	Tools      []string `json:"tools,omitempty"` // "file_search", "code_interpreter"
}

// AssistantMessage represents a message response
type AssistantMessage struct {
	ID            string                `json:"id"`
	Object        string                `json:"object"`
	CreatedAt     int64                 `json:"created_at"`
	ThreadID      string                `json:"thread_id"`
	Role          string                `json:"role"`
	Content       []MessageContent      `json:"content"`
	Attachments   []MessageAttachment   `json:"attachments,omitempty"`
	Metadata      map[string]string     `json:"metadata,omitempty"`
	Status        string                `json:"status,omitempty"` // "in_progress", "completed"
}

// MessageContent represents content in a message
type MessageContent struct {
	Type      string                `json:"type"` // "text", "image_file"
	Text      *TextContent          `json:"text,omitempty"`
	ImageFile *ImageFileContent     `json:"image_file,omitempty"`
}

// TextContent represents text content with annotations
type TextContent struct {
	Value       string       `json:"value"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

// Annotation represents a text annotation (e.g., citations)
type Annotation struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	FileID  string `json:"file_id,omitempty"`
	StartIndex int  `json:"start_index,omitempty"`
	EndIndex   int  `json:"end_index,omitempty"`
}

// ImageFileContent represents an image file reference
type ImageFileContent struct {
	FileID string `json:"file_id"`
}

// Run represents an execution of an assistant on a thread
type Run struct {
	ID            string                 `json:"id"`
	Object        string                 `json:"object"`
	CreatedAt     int64                  `json:"created_at"`
	ThreadID      string                 `json:"thread_id"`
	AssistantID   string                 `json:"assistant_id"`
	Status        string                 `json:"status"` // "queued", "in_progress", "requires_action", "cancelling", "cancelled", "failed", "completed", "expired"
	RequiredAction *RequiredAction       `json:"required_action,omitempty"`
	LastError     *RunError             `json:"last_error,omitempty"`
	ExpiresAt     int64                  `json:"expires_at,omitempty"`
	StartedAt     int64                  `json:"started_at,omitempty"`
	CompletedAt   int64                  `json:"completed_at,omitempty"`
	CancelledAt   int64                  `json:"cancelled_at,omitempty"`
	FailedAt      int64                  `json:"failed_at,omitempty"`
	Model         string                 `json:"model"`
	Instructions  string                 `json:"instructions,omitempty"`
	Tools         []AssistantTool        `json:"tools,omitempty"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	Temperature   *float64               `json:"temperature,omitempty"`
	TopP          *float64               `json:"top_p,omitempty"`
}

// RequiredAction represents an action required to continue the run
type RequiredAction struct {
	Type       string                  `json:"type"` // "submit_tool_outputs"
	SubmitToolOutputs *SubmitToolOutputs `json:"submit_tool_outputs,omitempty"`
}

// SubmitToolOutputs represents tool outputs that need to be submitted
type SubmitToolOutputs struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

// RunError represents an error that occurred during a run
type RunError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateRunRequest represents a request to create a run
type CreateRunRequest struct {
	AssistantID   string                 `json:"assistant_id"`
	Model         string                 `json:"model,omitempty"`
	Instructions  string                 `json:"instructions,omitempty"`
	Tools         []AssistantTool        `json:"tools,omitempty"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	Temperature   *float64               `json:"temperature,omitempty"`
	TopP          *float64               `json:"top_p,omitempty"`
}

// SubmitToolOutputsRequest represents a request to submit tool outputs
type SubmitToolOutputsRequest struct {
	ToolOutputs []ToolOutput `json:"tool_outputs"`
}

// ToolOutput represents the output of a tool call
type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}
