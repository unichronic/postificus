package controller

import (
	"fmt"
	"net/http"

	"postificus/internal/domain"
	"postificus/internal/service"

	"github.com/labstack/echo/v4"
)

type SettingsController struct {
	authService    *service.AuthService
	profileService *service.ProfileService
}

func NewSettingsController(authService *service.AuthService, profileService *service.ProfileService) *SettingsController {
	return &SettingsController{
		authService:    authService,
		profileService: profileService,
	}
}

// Credentials Handling

func (c *SettingsController) SaveCredentials(ctx echo.Context) error {
	var req struct {
		Platform    string            `json:"platform"`
		Credentials map[string]string `json:"credentials"` // simplified from interface{}
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Platform == "" || len(req.Credentials) == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Platform and credentials required"})
	}

	userID := 1 // Hardcoded MVP

	// Manual save bypasses provider login
	if err := c.authService.ManualSaveCredentials(ctx.Request().Context(), userID, req.Platform, req.Credentials); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "saved"})
}

func (c *SettingsController) GetCredentialsStatus(ctx echo.Context) error {
	platform := ctx.Param("platform")
	userID := 1 // Hardcoded MVP

	connected, account, err := c.authService.GetConnectionStatus(ctx.Request().Context(), userID, platform)
	if err != nil {
		fmt.Printf("Error checking status: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check status"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"connected": connected,
		"account":   account,
	})
}

// Profile Handling

func (c *SettingsController) GetProfile(ctx echo.Context) error {
	userID := 1
	profile, err := c.profileService.GetProfile(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch profile"})
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (c *SettingsController) SaveProfile(ctx echo.Context) error {
	var profile domain.Profile
	if err := ctx.Bind(&profile); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	profile.UserID = 1 // Enforce MVP user

	if err := c.profileService.SaveProfile(ctx.Request().Context(), &profile); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "saved"})
}
