package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/mk-knight23/ai-sdk-openai/internal/config"
	"github.com/mk-knight23/ai-sdk-openai/internal/handlers"
	"github.com/mk-knight23/ai-sdk-openai/internal/openai"
	"github.com/mk-knight23/ai-sdk-openai/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize OpenAI client
	client := openai.NewClient(
		cfg.OpenAIAPIKey,
		cfg.OpenAIModel,
		cfg.OpenAIEmbeddingModel,
	)

	// Initialize services
	chatService := services.NewChatService(client)
	vectorService := services.NewVectorService(client)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(chatService)
	vectorHandler := handlers.NewVectorHandler(vectorService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "AI-SDK-OpenAI API",
		ServerHeader: "AI-SDK-OpenAI",
		BodyLimit:    50 * 1024 * 1024, // 50MB for file uploads
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "ai-sdk-openai",
			"version": "1.0.0",
			"features": fiber.Map{
				"gpt4":             true,
				"streaming":        cfg.EnableStreaming,
				"function_calling": true,
				"file_uploads":     true,
			},
		})
	})

	// API routes
	api := app.Group("/api")

	// Chat completion routes
	chat := api.Group("/chat")
	chat.Post("/completions", chatHandler.CreateCompletion)
	chat.Post("/completions/stream", chatHandler.CreateCompletionStream)

	// File routes
	files := api.Group("/files")
	files.Post("", vectorHandler.UploadFile)
	files.Get("", vectorHandler.ListFiles)
	files.Get("/:id", vectorHandler.GetFile)
	files.Delete("/:id", vectorHandler.DeleteFile)

	// Get port from environment or default to configured value
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Model: %s", cfg.OpenAIModel)
	log.Printf("Embedding Model: %s", cfg.OpenAIEmbeddingModel)
	log.Printf("Streaming: %v", cfg.EnableStreaming)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
