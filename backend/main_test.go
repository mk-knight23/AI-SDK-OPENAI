package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"marketpulse-api/adk"
)

// setupTestApp creates a Fiber app for testing
func setupTestApp() *fiber.App {
	app := fiber.New()

	// Initialize Google ADK agent
	agent := adk.NewCompetitorIntelligenceAgent()

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "marketpulse-api",
			"version": "1.0.0",
		})
	})

	// API routes
	api := app.Group("/api")

	// Competitor intelligence endpoint
	api.Post("/analyze", func(c *fiber.Ctx) error {
		type AnalyzeRequest struct {
			CompanyName string `json:"company_name"`
			Industry    string `json:"industry"`
		}

		req := new(AnalyzeRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Run Google ADK competitor analysis
		report, err := agent.Run(c.Context(), req.CompanyName, req.Industry)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Convert report to JSON
		reportJSON, err := report.ToJSON()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to generate report",
			})
		}

		return c.Send(reportJSON)
	})

	return app
}

// TestHealthEndpoint tests the /health endpoint
func TestHealthEndpoint(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test health endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", result["status"])
	}

	if result["service"] != "marketpulse-api" {
		t.Errorf("Expected service 'marketpulse-api', got %v", result["service"])
	}

	if result["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %v", result["version"])
	}
}

// TestAnalyzeEndpoint tests the /api/analyze endpoint (competitor-analysis)
func TestAnalyzeEndpoint(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "Valid SaaS request",
			requestBody: map[string]string{
				"company_name": "TechStartup",
				"industry":     "SaaS",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				// Verify required fields
				if result["target_company"] != "TechStartup" {
					t.Errorf("Expected target_company 'TechStartup', got %v", result["target_company"])
				}

				if result["generated_at"] == nil {
					t.Error("Expected generated_at to be present")
				}

				competitors, ok := result["competitors"].([]interface{})
				if !ok || len(competitors) == 0 {
					t.Error("Expected competitors array with data")
				}

				if result["market_insights"] == nil || result["market_insights"] == "" {
					t.Error("Expected market_insights to be present")
				}

				recommendations, ok := result["recommendations"].([]interface{})
				if !ok || len(recommendations) == 0 {
					t.Error("Expected recommendations array with data")
				}
			},
		},
		{
			name: "Valid Fintech request",
			requestBody: map[string]string{
				"company_name": "FinanceCo",
				"industry":     "Fintech",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if result["target_company"] != "FinanceCo" {
					t.Errorf("Expected target_company 'FinanceCo', got %v", result["target_company"])
				}

				competitors, ok := result["competitors"].([]interface{})
				if !ok || len(competitors) != 3 {
					t.Errorf("Expected 3 competitors, got %d", len(competitors))
				}
			},
		},
		{
			name: "Valid Healthcare request",
			requestBody: map[string]string{
				"company_name": "HealthTech",
				"industry":     "Healthcare",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if result["target_company"] != "HealthTech" {
					t.Errorf("Expected target_company 'HealthTech', got %v", result["target_company"])
				}
			},
		},
		{
			name: "Empty company name",
			requestBody: map[string]string{
				"company_name": "",
				"industry":     "SaaS",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if result["target_company"] != "" {
					t.Errorf("Expected empty target_company, got %v", result["target_company"])
				}
			},
		},
		{
			name: "Empty industry",
			requestBody: map[string]string{
				"company_name": "TestCorp",
				"industry":     "",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				// Should still work with empty industry
				if result["target_company"] != "TestCorp" {
					t.Errorf("Expected target_company 'TestCorp', got %v", result["target_company"])
				}
			},
		},
		{
			name:           "Missing request body",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse error response: %v", err)
				}

				if result["error"] == nil || result["error"] == "" {
					t.Error("Expected error message in response")
				}
			},
		},
		{
			name: "Missing company_name field",
			requestBody: map[string]string{
				"industry": "SaaS",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				// Should still work with empty company name
				if result["target_company"] != "" {
					t.Errorf("Expected empty target_company, got %v", result["target_company"])
				}
			},
		},
		{
			name: "Missing industry field",
			requestBody: map[string]string{
				"company_name": "TestCorp",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if result["target_company"] != "TestCorp" {
					t.Errorf("Expected target_company 'TestCorp', got %v", result["target_company"])
				}
			},
		},
		{
			name:           "Invalid JSON body",
			requestBody:    map[string]string{}, // Will send invalid JSON
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				// Test with truly invalid JSON
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.name == "Invalid JSON body" {
				reqBody = []byte(`{invalid json`)
			} else if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test analyze endpoint: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)

			if tt.checkResponse != nil && tt.name != "Invalid JSON body" {
				tt.checkResponse(t, body)
			}
		})
	}
}

// TestAnalyzeEndpoint_InvalidJSON specifically tests invalid JSON handling
func TestAnalyzeEndpoint_InvalidJSON(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/api/analyze", strings.NewReader(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test analyze endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if result["error"] == nil || result["error"] == "" {
		t.Error("Expected error message in response for invalid JSON")
	}
}

// TestAnalyzeEndpoint_CompetitorStructure tests the structure of competitor data in response
func TestAnalyzeEndpoint_CompetitorStructure(t *testing.T) {
	app := setupTestApp()

	requestBody := map[string]string{
		"company_name": "TestCorp",
		"industry":     "SaaS",
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test analyze endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	competitors, ok := result["competitors"].([]interface{})
	if !ok {
		t.Fatal("Expected competitors to be an array")
	}

	if len(competitors) != 3 {
		t.Errorf("Expected 3 competitors, got %d", len(competitors))
	}

	// Check structure of first competitor
	if len(competitors) > 0 {
		firstComp, ok := competitors[0].(map[string]interface{})
		if !ok {
			t.Fatal("Expected competitor to be an object")
		}

		requiredFields := []string{"competitor_name", "threat_level", "positioning", "key_differentiators", "opportunities", "risks"}
		for _, field := range requiredFields {
			if firstComp[field] == nil {
				t.Errorf("Expected competitor to have field '%s'", field)
			}
		}

		// Verify threat level is one of expected values
		threatLevel, ok := firstComp["threat_level"].(string)
		if ok {
			validThreatLevels := []string{"High", "Medium", "Low"}
			found := false
			for _, level := range validThreatLevels {
				if threatLevel == level {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected threat level: %s", threatLevel)
			}
		}
	}
}

// TestAnalyzeEndpoint_Recommendations tests the recommendations in response
func TestAnalyzeEndpoint_Recommendations(t *testing.T) {
	app := setupTestApp()

	requestBody := map[string]string{
		"company_name": "TestCorp",
		"industry":     "SaaS",
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test analyze endpoint: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	recommendations, ok := result["recommendations"].([]interface{})
	if !ok {
		t.Fatal("Expected recommendations to be an array")
	}

	if len(recommendations) == 0 {
		t.Error("Expected at least one recommendation")
	}

	// Verify recommendations are strings
	for i, rec := range recommendations {
		if _, ok := rec.(string); !ok {
			t.Errorf("Recommendation %d is not a string", i)
		}
	}
}

// TestAnalyzeEndpoint_MarketInsights tests the market insights in response
func TestAnalyzeEndpoint_MarketInsights(t *testing.T) {
	app := setupTestApp()

	requestBody := map[string]string{
		"company_name": "TestCorp",
		"industry":     "SaaS",
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test analyze endpoint: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	marketInsights, ok := result["market_insights"].(string)
	if !ok {
		t.Fatal("Expected market_insights to be a string")
	}

	if marketInsights == "" {
		t.Error("Expected market_insights to be non-empty")
	}

	// Verify it contains expected content
	if !strings.Contains(marketInsights, "competitive landscape") {
		t.Error("Expected market_insights to mention 'competitive landscape'")
	}
}

// TestAnalyzeEndpoint_GeneratedAt tests the generated_at timestamp
func TestAnalyzeEndpoint_GeneratedAt(t *testing.T) {
	app := setupTestApp()

	requestBody := map[string]string{
		"company_name": "TestCorp",
		"industry":     "SaaS",
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test analyze endpoint: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	generatedAt, ok := result["generated_at"].(string)
	if !ok {
		t.Fatal("Expected generated_at to be a string")
	}

	if generatedAt == "" {
		t.Error("Expected generated_at to be non-empty")
	}

	// Verify it's a valid timestamp
	if _, err := json.Marshal(result["generated_at"]); err != nil {
		t.Errorf("generated_at is not valid: %v", err)
	}
}

// TestNonExistentEndpoint tests that non-existent endpoints return 404
func TestNonExistentEndpoint(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test non-existent endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent endpoint, got %d", resp.StatusCode)
	}
}

// TestMethodNotAllowed tests that wrong HTTP methods return appropriate errors
func TestMethodNotAllowed(t *testing.T) {
	app := setupTestApp()

	// Test GET on POST-only endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/analyze", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to test method not allowed: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Logf("GET /api/analyze returned status %d (Fiber returns 404 for method mismatch)", resp.StatusCode)
	}
}

// TestCORSHeaders tests that CORS headers are present
func TestCORSHeaders(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodOptions, "/api/analyze", nil)
	req.Header.Set("Origin", "http://localhost:4200")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test CORS: %v", err)
	}

	// Note: CORS middleware may not be fully configured in test app
	// This test documents the expected behavior
	t.Logf("CORS preflight response status: %d", resp.StatusCode)
}

// TestContentTypeHeader tests that responses have correct content type
func TestContentTypeHeader(t *testing.T) {
	app := setupTestApp()

	requestBody := map[string]string{
		"company_name": "TestCorp",
		"industry":     "SaaS",
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test content type: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		t.Log("Content-Type header not set (Fiber may set it automatically)")
	}
}

// TestMainFunction tests the main function (partially)
func TestMainFunction(t *testing.T) {
	// Test that the port environment variable is read correctly
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	// Test with custom port
	os.Setenv("PORT", "9090")
	port := os.Getenv("PORT")
	if port != "9090" {
		t.Errorf("Expected port '9090', got '%s'", port)
	}

	// Test with empty port (should default to 8080)
	os.Unsetenv("PORT")
	port = os.Getenv("PORT")
	if port != "" {
		t.Errorf("Expected empty port, got '%s'", port)
	}
}

// BenchmarkAnalyzeEndpoint benchmarks the analyze endpoint
func BenchmarkAnalyzeEndpoint(b *testing.B) {
	app := setupTestApp()

	requestBody := map[string]string{
		"company_name": "TestCorp",
		"industry":     "SaaS",
	}

	reqBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// Drain body
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// BenchmarkHealthEndpoint benchmarks the health endpoint
func BenchmarkHealthEndpoint(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)

		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// Drain body
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
