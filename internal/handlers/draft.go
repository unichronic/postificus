package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"postificus/internal/database"

	"github.com/labstack/echo/v4"
)

func UpdateDraft(c echo.Context) error {
	id := c.Param("id")
	var payload struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&payload); err != nil {
		return err
	}

	// 1. FAST PATH: Save to Redis
	ctx := context.Background()
	// "SET draft:123:content '...' EX 3600" (Expire in 1 hour if inactive)
	err := database.RedisClient.Set(ctx, fmt.Sprintf("draft:%s:content", id), payload.Content, 1*time.Hour).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to cache draft"})
	}

	// 2. TRIGGER PERSISTENCE (Optional / Debounced)
	// In a real production app, we might push a job to Asynq here to persist to DB eventually,
	// or rely on a background ticker. For now, we just cache it.

	return c.NoContent(http.StatusOK)
}
