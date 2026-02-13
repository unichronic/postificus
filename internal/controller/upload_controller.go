package controller

import (
	"log"
	"net/http"
	"postificus/internal/storage"

	"github.com/labstack/echo/v4"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (c *UploadController) HandleUpload(ctx echo.Context) error {
	// Source
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Validate size (e.g., max 5MB)
	if file.Size > 5*1024*1024 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "File too large (max 5MB)"})
	}

	// Upload to S3
	url, err := storage.UploadFile(ctx.Request().Context(), file)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload file"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"url": url})
}
