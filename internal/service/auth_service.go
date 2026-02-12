package service

import (
	"context"
	"encoding/json"
	"fmt"
	"postificus/internal/browser"
	"postificus/internal/storage"
)

// AuthService handles authentication business logic
type AuthService struct {
	credsRepo storage.CredentialsRepository
}

// NewAuthService creates a new instance
func NewAuthService(credsRepo storage.CredentialsRepository) *AuthService {
	return &AuthService{
		credsRepo: credsRepo,
	}
}

// ConnectPlatform handles the flow of connecting a platform account
func (s *AuthService) ConnectPlatform(ctx context.Context, userID int, platform string) (string, error) {
	var creds map[string]string
	var username string

	// 1. Trigger Login via Browser Automation
	if platform == "medium" {
		uid, sid, xsrf, uname, loginErr := browser.WaitForMediumLogin()
		if loginErr != nil {
			return "", fmt.Errorf("login failed: %w", loginErr)
		}
		creds = map[string]string{
			"uid":  uid,
			"sid":  sid,
			"xsrf": xsrf,
		}
		username = uname
	} else if platform == "devto" {
		token, uname, loginErr := browser.WaitForDevToLogin()
		if loginErr != nil {
			return "", fmt.Errorf("login failed: %w", loginErr)
		}
		creds = map[string]string{
			"remember_user_token": token,
		}
		username = uname
	} else {
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	// 2. Inject Account Name into Credentials for storage
	if username != "" {
		creds["account_name"] = username
	} else {
		creds["account_name"] = "Connected Account" // Fallback
	}

	// 3. Save Credentials via Repository
	if err := s.credsRepo.SaveCredentials(ctx, userID, platform, creds); err != nil {
		return "", fmt.Errorf("failed to save credentials: %w", err)
	}

	return username, nil
}

// ManualSaveCredentials allows saving credentials directly (non-interactive)
func (s *AuthService) ManualSaveCredentials(ctx context.Context, userID int, platform string, creds map[string]string) error {
	return s.credsRepo.SaveCredentials(ctx, userID, platform, creds)
}

// GetConnectionStatus checks if a platform is connected and returns the account name
func (s *AuthService) GetConnectionStatus(ctx context.Context, userID int, platform string) (bool, string, error) {
	cred, err := s.credsRepo.GetCredentials(ctx, userID, platform)
	if err != nil {
		return false, "", err
	}
	if cred == nil {
		return false, "", nil
	}

	var details map[string]interface{}
	if err := json.Unmarshal(cred.Credentials, &details); err != nil {
		return true, "", nil // Connected but malformed JSON? Treat as connected
	}

	accountName, _ := details["account_name"].(string)
	return true, accountName, nil
}
