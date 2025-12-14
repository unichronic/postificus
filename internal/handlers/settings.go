package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"postificus/internal/database"
	"time"

	"github.com/labstack/echo/v4"
)

type CredentialsRequest struct {
	Platform    string                 `json:"platform"`
	Credentials map[string]interface{} `json:"credentials"`
}

// SaveCredentials stores user credentials in the database
func SaveCredentials(c echo.Context) error {
	var req CredentialsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Platform == "" || len(req.Credentials) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "platform and credentials are required"})
	}

	// Hardcoded UserID for MVP
	userID := 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	credsJSON, err := json.Marshal(req.Credentials)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to marshal credentials"})
	}

	query := `
		INSERT INTO user_credentials (user_id, platform, credentials, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, platform) 
		DO UPDATE SET credentials = $3, updated_at = NOW()
	`

	_, err = database.DB.Exec(ctx, query, userID, req.Platform, credsJSON)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save credentials: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "saved"})
}

// GetCredentialsStatus checks if credentials exist for a platform
func GetCredentialsStatus(c echo.Context) error {
	platform := c.Param("platform")
	userID := 1 // Hardcoded for MVP

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_credentials WHERE user_id = $1 AND platform = $2)`

	err := database.DB.QueryRow(ctx, query, userID, platform).Scan(&exists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check status"})
	}

	return c.JSON(http.StatusOK, map[string]bool{"connected": exists})
}
