package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"postificus/internal/database"
	"postificus/internal/handlers"
	"postificus/internal/middleware"
	"postificus/internal/queue"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, creating one if config is saved.")
	}

	// Initialize Database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Initialize Redis
	if err := database.InitRedis(); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	defer database.CloseRedis()

	// Initialize LLM Summarizer (Optional for API, but used for content generation)
	// summarizer, err := llm.NewSummarizer(context.Background())
	// if err != nil {
	// 	log.Println("Warning: Failed to initialize LLM Summarizer (check GEMINI_API_KEY):", err)
	// }

	// Setup Echo
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// Rate Limiting (20 req/sec burst 50)
	rateLimiter := middleware.NewRateLimiter(rate.Limit(20), 50)
	e.Use(middleware.RateLimit(rateLimiter))

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Drafts (Auto-Save)
	e.PUT("/api/drafts/:id", handlers.UpdateDraft)

	// Config Endpoints
	e.POST("/api/settings/credentials", handlers.SaveCredentials)
	e.GET("/api/settings/credentials/:platform", handlers.GetCredentialsStatus)

	// Legacy endpoint (can be deprecated or redirected)
	e.POST("/config/devto", func(c echo.Context) error {
		// Redirect to new handler logic if needed, or keep for backward compatibility
		// For now, we'll just map it to the new handler structure manually or leave as is
		// But let's leave it as is to avoid breaking existing frontend immediately
		return c.JSON(http.StatusOK, map[string]string{"status": "saved"})
	})

	// Publish Endpoints (Enqueue tasks instead of direct execution)
	// Publish Endpoints (Enqueue tasks)
	e.POST("/publish/:platform", func(c echo.Context) error {
		platform := c.Param("platform")

		var req struct {
			Title   string   `json:"title"`
			Content string   `json:"content"`
			Tags    []string `json:"tags"`
			BlogURL string   `json:"blog_url"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		// Create Task Payload
		payload := queue.PublishPayload{
			UserID:   1, // Hardcoded for MVP
			Platform: platform,
			Title:    req.Title,
			Content:  req.Content,
			Tags:     req.Tags,
			BlogURL:  req.BlogURL,
		}

		// Enqueue Task
		client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_URL")})
		defer client.Close()

		task, err := queue.NewPublishTask(payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create task"})
		}

		info, err := client.Enqueue(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to enqueue task"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "queued",
			"task_id": info.ID,
		})
	})

	// Start server
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
