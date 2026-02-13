package controller

import (
	"net/http"

	"postificus/internal/rabbitmq"
	"postificus/internal/service"

	"github.com/labstack/echo/v4"
)

type PublishController struct {
	producer *rabbitmq.Producer
}

func NewPublishController(producer *rabbitmq.Producer) *PublishController {
	return &PublishController{producer: producer}
}

func (c *PublishController) PublishPost(ctx echo.Context) error {
	platform := ctx.Param("platform")

	var req struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		CoverImage string   `json:"cover_image"`
		Tags       []string `json:"tags"`
		BlogURL    string   `json:"blog_url"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Create Task Payload
	payload := service.PublishPayload{
		UserID:     service.DefaultUserID(),
		Platform:   platform,
		Title:      req.Title,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		Tags:       req.Tags,
		BlogURL:    req.BlogURL,
	}

	bytes, err := service.NewPublishPayload(payload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create task"})
	}

	// Publish to RabbitMQ
	if err := c.producer.Publish(service.TypePublishPost, bytes); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to enqueue task"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":  "queued",
		"message": "Task submitted to RabbitMQ",
	})
}
