package stream

import (
	"strings"
	"testing"
)

func TestChunkBuilder(t *testing.T) {
	t.Run("new chunk builder is empty", func(t *testing.T) {
		cb := NewChunkBuilder()
		if cb.HasContent() {
			t.Error("New chunk builder should be empty")
		}
		if cb.GetContent() != "" {
			t.Error("New chunk builder should have empty content")
		}
	})

	t.Run("append and get content", func(t *testing.T) {
		cb := NewChunkBuilder()
		cb.Append("Hello")
		cb.Append(" ")
		cb.Append("World")

		content := cb.GetContent()
		if content != "Hello World" {
			t.Errorf("Content = %v, want Hello World", content)
		}
	})

	t.Run("has content after append", func(t *testing.T) {
		cb := NewChunkBuilder()
		if cb.HasContent() {
			t.Error("Should not have content initially")
		}

		cb.Append("test")
		if !cb.HasContent() {
			t.Error("Should have content after append")
		}
	})

	t.Run("reset clears content", func(t *testing.T) {
		cb := NewChunkBuilder()
		cb.Append("test")
		cb.Reset()

		if cb.HasContent() {
			t.Error("Should not have content after reset")
		}
		if cb.GetContent() != "" {
			t.Error("Content should be empty after reset")
		}
	})
}

func TestStreamWriter(t *testing.T) {
	t.Run("write data event", func(t *testing.T) {
		var sb strings.Builder
		sw := NewStreamWriter(&sb)

		err := sw.WriteData("test data")
		if err != nil {
			t.Fatalf("WriteData() error = %v", err)
		}

		output := sb.String()
		if !strings.Contains(output, "test data") {
			t.Error("Output should contain 'test data'")
		}
		if !strings.HasSuffix(output, "\n\n") {
			t.Error("Output should end with double newline")
		}
	})

	t.Run("write done signal", func(t *testing.T) {
		var sb strings.Builder
		sw := NewStreamWriter(&sb)

		err := sw.WriteDone()
		if err != nil {
			t.Fatalf("WriteDone() error = %v", err)
		}

		output := sb.String()
		if !strings.Contains(output, "[DONE]") {
			t.Error("Output should contain [DONE]")
		}
	})
}
