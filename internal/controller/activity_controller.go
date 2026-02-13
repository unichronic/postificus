package controller

import (
	"net/http"
	"strconv"
	"time"

	"postificus/internal/service"

	"github.com/labstack/echo/v4"
)

type ActivityController struct {
	service *service.ActivityService
}

func NewActivityController(service *service.ActivityService) *ActivityController {
	return &ActivityController{service: service}
}

func (c *ActivityController) GetDevtoActivity(ctx echo.Context) error {
	userID := service.DefaultUserID()
	limit := 10
	if raw := ctx.QueryParam("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	posts, err := c.service.FetchLiveDevtoActivity(ctx.Request().Context(), userID, limit)
	if err != nil {
		// Can be improved to handle 401 specifically
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"posts":        posts,
		"count":        len(posts),
		"retrieved_at": time.Now().UTC().Format(time.RFC3339),
	})
}

func (c *ActivityController) GetMediumActivity(ctx echo.Context) error {
	userID := service.DefaultUserID()
	limit := 10
	if raw := ctx.QueryParam("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	posts, err := c.service.FetchLiveMediumActivity(ctx.Request().Context(), userID, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"posts":        posts,
		"count":        len(posts),
		"retrieved_at": time.Now().UTC().Format(time.RFC3339),
	})
}
