package controller

import (
	"fmt"
	"net/http"
	"time"

	"postificus/internal/domain"
	"postificus/internal/service"

	"github.com/labstack/echo/v4"
)

type DraftController struct {
	service *service.DraftService
}

func NewDraftController(service *service.DraftService) *DraftController {
	return &DraftController{service: service}
}

func (c *DraftController) UpdateDraft(ctx echo.Context) error {
	id := ctx.Param("id")
	var payload struct {
		Title          string   `json:"title"`
		Content        string   `json:"content"`
		CoverImage     string   `json:"cover_image"`
		PublishTargets []string `json:"publish_targets"`
	}
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	userID := 1 // Hardcoded for MVP

	draft := &domain.Draft{
		ID:             id,
		UserID:         userID,
		Title:          payload.Title,
		Content:        payload.Content,
		CoverImage:     payload.CoverImage,
		PublishTargets: payload.PublishTargets,
	}

	if err := c.service.SaveDraft(ctx.Request().Context(), draft); err != nil {
		fmt.Printf("Error saving draft: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save draft"})
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *DraftController) GetDraft(ctx echo.Context) error {
	id := ctx.Param("id")
	userID := 1 // Hardcoded for MVP

	draft, err := c.service.GetDraft(ctx.Request().Context(), id, userID)
	if err != nil {
		// Differentiate between not found and internal error?
		// For now simple 404
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Draft not found"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":              draft.ID,
		"title":           draft.Title,
		"content":         draft.Content,
		"cover_image":     draft.CoverImage,
		"publish_targets": draft.PublishTargets,
		"last_saved_at":   draft.LastSavedAt.UTC().Format(time.RFC3339),
		"is_published":    draft.IsPublished,
	})
}
