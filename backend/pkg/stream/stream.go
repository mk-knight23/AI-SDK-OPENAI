package stream

import (
	"encoding/json"
	"io"
)

// StreamReader reads streaming responses from OpenAI
type StreamReader struct {
	decoder *json.Decoder
}

// NewStreamReader creates a new stream reader
func NewStreamReader(reader io.Reader) *StreamReader {
	return &StreamReader{
		decoder: json.NewDecoder(reader),
	}
}

// Next reads the next chunk from the stream
// Returns (chunk, hasMore, error)
func (sr *StreamReader) Next() (map[string]any, bool, error) {
	if !sr.decoder.More() {
		return nil, false, io.EOF
	}

	var chunk map[string]any
	if err := sr.decoder.Decode(&chunk); err != nil {
		if err == io.EOF {
			return nil, false, nil
		}
		return nil, false, err
	}

	return chunk, true, nil
}

// StreamEvent represents a server-sent event
type StreamEvent struct {
	Data  string `json:"data"`
	Event string `json:"event,omitempty"`
	ID    string `json:"id,omitempty"`
	Retry int    `json:"retry,omitempty"`
}

// StreamWriter writes streaming responses
type StreamWriter struct {
	writer io.Writer
}

// NewStreamWriter creates a new stream writer
func NewStreamWriter(writer io.Writer) *StreamWriter {
	return &StreamWriter{
		writer: writer,
	}
}

// WriteEvent writes a server-sent event
func (sw *StreamWriter) WriteEvent(event StreamEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = sw.writer.Write(data)
	if err != nil {
		return err
	}

	// Write newline
	_, err = sw.writer.Write([]byte("\n\n"))
	return err
}

// WriteData writes a data event
func (sw *StreamWriter) WriteData(data string) error {
	return sw.WriteEvent(StreamEvent{Data: data})
}

// WriteDone writes a completion signal
func (sw *StreamWriter) WriteDone() error {
	return sw.WriteData("[DONE]")
}

// ChunkBuilder helps build streaming chunks
type ChunkBuilder struct {
	deltaContent string
}

// NewChunkBuilder creates a new chunk builder
func NewChunkBuilder() *ChunkBuilder {
	return &ChunkBuilder{}
}

// Append adds content to the current chunk
func (cb *ChunkBuilder) Append(content string) {
	cb.deltaContent += content
}

// GetContent returns the accumulated content
func (cb *ChunkBuilder) GetContent() string {
	return cb.deltaContent
}

// Reset clears the accumulated content
func (cb *ChunkBuilder) Reset() {
	cb.deltaContent = ""
}

// HasContent returns true if there is accumulated content
func (cb *ChunkBuilder) HasContent() bool {
	return cb.deltaContent != ""
}
