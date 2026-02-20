package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"marketpulse-api/adk"
)

func main() {
	app := fiber.New()

	// Initialize Google ADK agent
	agent := adk.NewCompetitorIntelligenceAgent()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "marketpulse-api",
			"version":   "1.0.0",
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

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
