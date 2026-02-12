package controller

import (
	"net/http"
	"postificus/internal/service"

	"github.com/labstack/echo/v4"
)

// AuthController handles authentication related endpoints
type AuthController struct {
	service *service.AuthService
}

// NewAuthController creates a new instance
func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{service: service}
}

// HandleConnectPlatform handles the POST /api/connect/:platform request
func (c *AuthController) HandleConnectPlatform(ctx echo.Context) error {
	// Parse request
	var req struct {
		Platform string `json:"platform"`
		UserID   int    `json:"user_id"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Basic Validation
	if req.Platform != "medium" && req.Platform != "devto" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Unsupported platform"})
	}

	// Default User ID (MVP)
	if req.UserID == 0 {
		req.UserID = 1
	}

	// Call Service
	username, err := c.service.ConnectPlatform(ctx.Request().Context(), req.UserID, req.Platform)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Success Response
	return ctx.JSON(http.StatusOK, map[string]string{
		"status":   "connected",
		"platform": req.Platform,
		"account":  username,
		"message":  "Using account connected via " + req.Platform,
	})
}
