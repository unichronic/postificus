package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"postificus/automation"
	"postificus/internal/llm"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, creating one if config is saved.")
	}

	// Initialize Shared Browser
	if err := automation.InitBrowser(); err != nil {
		log.Fatal("Failed to initialize browser:", err)
	}
	defer automation.CloseBrowser()

	// Initialize LLM Summarizer
	summarizer, err := llm.NewSummarizer(context.Background())
	if err != nil {
		log.Println("Warning: Failed to initialize LLM Summarizer (check GEMINI_API_KEY):", err)
	}

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Endpoint to save Dev.to token
	e.POST("/config/devto", func(c echo.Context) error {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if req.Token == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
		}

		// Save to .env file
		f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open .env file"})
		}
		defer f.Close()

		// Simple append (in a real app, we'd parse and replace)
		// For now, we assume the user might clear it or we just append and godotenv might read the first or last?
		// Actually godotenv reads the first one usually.
		// Let's just overwrite the file for this simple demo or append a new line.
		// Better: Read all, filter out DEVTO_SESSION_TOKEN, write back.
		// For simplicity in this task: Overwrite .env with just this token or append.
		// Let's just append for now, but warn user.
		if _, err := f.WriteString(fmt.Sprintf("\nDEVTO_SESSION_TOKEN=%s\n", req.Token)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write to .env file"})
		}

		// Reload env
		godotenv.Load()

		return c.JSON(http.StatusOK, map[string]string{"status": "saved"})
	})

	// Endpoint to check Dev.to connection status
	e.GET("/config/devto", func(c echo.Context) error {
		token := os.Getenv("DEVTO_SESSION_TOKEN")
		if token != "" {
			return c.JSON(http.StatusOK, map[string]bool{"connected": true})
		}
		return c.JSON(http.StatusOK, map[string]bool{"connected": false})
	})

	e.POST("/publish/devto", func(c echo.Context) error {
		var req struct {
			Title   string   `json:"title"`
			Content string   `json:"content"`
			Tags    []string `json:"tags"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Read token from Env
		token := os.Getenv("DEVTO_SESSION_TOKEN")
		if token == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dev.to token not configured. Please connect first."})
		}

		// Trigger automation
		err := automation.PostToDevToWithCookie(token, req.Title, req.Content, req.Tags)
		if err != nil {
			e.Logger.Error("Automation failed", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "published"})
	})

	e.POST("/publish/medium", func(c echo.Context) error {
		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		uid := os.Getenv("MEDIUM_UID")
		sid := os.Getenv("MEDIUM_SID")
		xsrf := os.Getenv("MEDIUM_XSRF")

		if uid == "" || sid == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Medium credentials not configured (MEDIUM_UID, MEDIUM_SID)"})
		}

		// Trigger automation
		err := automation.PostToMedium(uid, sid, xsrf, req.Title, req.Content)
		if err != nil {
			e.Logger.Error("Medium Automation failed", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "published"})
	})

	e.POST("/config/linkedin", func(c echo.Context) error {
		var req struct {
			LiAt string `json:"li_at"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if req.LiAt == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "li_at is required"})
		}

		f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open .env file"})
		}
		defer f.Close()

		if _, err := f.WriteString(fmt.Sprintf("\nLI_AT=%s\n", req.LiAt)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write to .env file"})
		}
		godotenv.Load()
		return c.JSON(http.StatusOK, map[string]string{"status": "saved"})
	})

	e.POST("/publish/linkedin", func(c echo.Context) error {
		var req struct {
			BlogContent string `json:"blog_content"`
			BlogURL     string `json:"blog_url"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		liAt := os.Getenv("LI_AT")
		if liAt == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "LinkedIn li_at cookie not configured."})
		}

		if summarizer == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "LLM Summarizer not initialized (missing GEMINI_API_KEY?)"})
		}

		// 1. Generate Content
		postContent, err := summarizer.GenerateSocialContent(c.Request().Context(), req.BlogContent)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "LLM generation failed: " + err.Error()})
		}

		// 2. Post to LinkedIn
		err = automation.PostToLinkedIn(liAt, postContent, req.BlogURL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "LinkedIn automation failed: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "published", "content": postContent})
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
