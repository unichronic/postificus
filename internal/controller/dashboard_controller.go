package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"postificus/internal/rabbitmq"
	"postificus/internal/service"

	"github.com/labstack/echo/v4"
)

type DashboardController struct {
	activityService *service.ActivityService
	producer        *rabbitmq.Producer
}

func NewDashboardController(activityService *service.ActivityService, producer *rabbitmq.Producer) *DashboardController {
	return &DashboardController{
		activityService: activityService,
		producer:        producer,
	}
}

// GetDashboardActivity returns the unified list of posts from the local database.
func (c *DashboardController) GetDashboardActivity(ctx echo.Context) error {
	userID := service.DefaultUserID()

	limit := 20
	if raw := ctx.QueryParam("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	posts, err := c.activityService.GetDashboardActivity(ctx.Request().Context(), userID, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"posts": posts,
		"count": len(posts),
	})
}

// TriggerSync enqueues a background task to sync activity for a specific platform (or all).
func (c *DashboardController) TriggerSync(ctx echo.Context) error {
	userID := service.DefaultUserID()

	var req struct {
		Platform string `json:"platform"` // "medium", "devto", or "all"
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	platforms := []string{}
	if req.Platform == "all" || req.Platform == "" {
		platforms = []string{"medium", "devto"}
	} else {
		platforms = []string{req.Platform}
	}

	enqueued := 0
	for _, p := range platforms {
		payload := service.SyncPlatformPayload{
			UserID:   userID,
			Platform: p,
		}

		bytes, err := json.Marshal(payload)
		if err != nil {
			continue // Skip invalid
		}

		// Publish to RabbitMQ
		err = c.producer.Publish(service.TypeSyncPlatformActivity, bytes)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to enqueue %s: %v", p, err)})
		}
		fmt.Printf("Enqueued task for platform %s\n", p)
		enqueued++
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":  fmt.Sprintf("Triggered sync for %d platforms", enqueued),
		"enqueued": enqueued,
	})
}
