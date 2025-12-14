package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"postificus/internal/automation"
	"postificus/internal/database"

	"github.com/hibiken/asynq"
)

func HandlePublishTask(ctx context.Context, t *asynq.Task) error {
	// 1. Parse Payload
	var p PublishPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing publish task for platform: %s, title: %s", p.Platform, p.Title)

	// 2. Fetch Cookies/Credentials
	// Try DB first, then Env
	creds, err := fetchCredentials(ctx, p.UserID, p.Platform)
	if err != nil {
		log.Printf("Warning: Failed to fetch credentials from DB: %v. Falling back to Env.", err)
	}

	// 3. Select Platform Strategy
	switch p.Platform {
	case "medium":
		uid := getCredential(creds, "uid", "MEDIUM_UID")
		sid := getCredential(creds, "sid", "MEDIUM_SID")
		xsrf := getCredential(creds, "xsrf", "MEDIUM_XSRF")

		if uid == "" || sid == "" {
			return fmt.Errorf("medium credentials missing")
		}

		return automation.PostToMedium(uid, sid, xsrf, p.Title, p.Content)

	case "linkedin":
		liAt := getCredential(creds, "li_at", "LI_AT")
		if liAt == "" {
			return fmt.Errorf("linkedin credentials missing")
		}

		return automation.PostToLinkedIn(liAt, p.Content, p.BlogURL)

	case "devto":
		token := getCredential(creds, "token", "DEVTO_SESSION_TOKEN")
		if token == "" {
			return fmt.Errorf("devto credentials missing")
		}

		return automation.PostToDevToWithCookie(token, p.Title, p.Content, p.Tags)
	}

	return fmt.Errorf("unknown platform: %s", p.Platform)
}

// Helper to fetch credentials from DB
func fetchCredentials(ctx context.Context, userID int, platform string) (map[string]interface{}, error) {
	// We need a DB connection here. Since HandlePublishTask is in the queue package,
	// we assume database.DB is initialized globally or we need to import it.
	// We'll use postificus/internal/database

	// Note: We need to import "postificus/internal/database"
	// But circular dependency might be an issue if queue imports database and database imports queue.
	// database does NOT import queue, so it's safe.

	// However, we need to ensure database.DB is accessible.
	// In worker/main.go we init DB.

	// Let's use a raw query
	var credsJSON []byte
	query := `SELECT credentials FROM user_credentials WHERE user_id = $1 AND platform = $2`

	// We need to access the DB pool.
	// Ideally, we should pass the DB pool to the handler structure.
	// For this refactor, we'll access the global database.DB variable if exported,
	// or we might need to update the signature.
	// Checking database/db.go -> var DB *pgxpool.Pool is exported.

	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	err := database.DB.QueryRow(ctx, query, userID, platform).Scan(&credsJSON)
	if err != nil {
		return nil, err
	}

	var creds map[string]interface{}
	if err := json.Unmarshal(credsJSON, &creds); err != nil {
		return nil, err
	}

	return creds, nil
}

// Helper to get credential from map or env
func getCredential(creds map[string]interface{}, key, envVar string) string {
	if val, ok := creds[key]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			return strVal
		}
	}
	return os.Getenv(envVar)
}
