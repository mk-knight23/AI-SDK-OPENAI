/*
OpenAI Backend Server
AI SDK: OpenAI
Tech Stack: Go + Standard Library
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Response represents a standard API response
type Response struct {
	Status   string      `json:"status"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// AIRequest represents an AI request
type AIRequest struct {
	Prompt      string                 `json:"prompt"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Status: "healthy",
			Data: map[string]string{
				"service": "openai-api",
				"version": "1.0.0",
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Status:  "success",
			Message: fmt.Sprintf("Welcome to OpenAI API"),
			Data: map[string]interface{}{
				"version": "1.0.0",
				"endpoints": map[string]string{
					"health": "/health",
					"api":    "/api/ai",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	// AI endpoint (placeholder for now)
	http.HandleFunc("/api/ai", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var aiReq AIRequest
		if err := json.NewDecoder(r.Body).Decode(&aiReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Status:  "success",
			Message: fmt.Sprintf("Mock AI response for: %s", aiReq.Prompt),
			Data: map[string]string{
				"framework": "OpenAI",
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
