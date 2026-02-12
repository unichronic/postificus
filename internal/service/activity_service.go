package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	automation "postificus/internal/browser"
	"postificus/internal/domain"
	"postificus/internal/storage"
)

const (
	DashboardCacheTTL = 10 * time.Minute
)

type ActivityService struct {
	credsRepo storage.CredentialsRepository
}

func NewActivityService(credsRepo storage.CredentialsRepository) *ActivityService {
	return &ActivityService{
		credsRepo: credsRepo,
	}
}

// UpsertPost inserts or updates a post in the unified_posts table.
func (s *ActivityService) UpsertPost(ctx context.Context, userID int, post domain.UnifiedPost) error {
	if storage.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	query := `
		INSERT INTO unified_posts (user_id, platform, remote_id, title, url, status, views, reactions, comments, published_at, last_synced_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		ON CONFLICT (user_id, platform, remote_id) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			url = EXCLUDED.url,
			status = EXCLUDED.status,
			views = EXCLUDED.views,
			reactions = EXCLUDED.reactions,
			comments = EXCLUDED.comments,
			published_at = EXCLUDED.published_at,
			last_synced_at = NOW();
	`

	_, err := storage.DB.Exec(ctx, query,
		userID,
		post.Platform,
		post.RemoteID,
		post.Title,
		post.URL,
		post.Status,
		post.Views,
		post.Reactions,
		post.Comments,
		post.PublishedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert post: %w", err)
	}

	// Invalidate Cache for this user
	if storage.RedisClient != nil {
		cacheKey := fmt.Sprintf("dashboard_activity:%d", userID)
		storage.RedisClient.Del(ctx, cacheKey)
	}

	return nil
}

// GetDashboardActivity returns the unified list of posts from local cache (Redis) or DB
func (s *ActivityService) GetDashboardActivity(ctx context.Context, userID int, limit int) ([]domain.UnifiedPost, error) {
	if storage.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	cacheKey := fmt.Sprintf("dashboard_activity:%d", userID)

	// 1. Try Cache
	if storage.RedisClient != nil {
		val, err := storage.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && val != "" {
			var cachedPosts []domain.UnifiedPost
			if err := json.Unmarshal([]byte(val), &cachedPosts); err == nil {
				// Cache Hit
				return cachedPosts, nil
			}
		}
	}

	// 2. Cache Miss - Query DB
	query := `
		SELECT
			up.platform,
			up.remote_id,
			up.title,
			up.url,
			up.status,
			up.views,
			up.reactions,
			up.comments,
			up.published_at,
			d.publish_targets
		FROM unified_posts up
		LEFT JOIN drafts d
			ON up.platform = 'postificus'
			AND d.id::text = up.remote_id
			AND d.user_id = up.user_id
		WHERE up.user_id = $1
		ORDER BY up.published_at DESC NULLS LAST, up.last_synced_at DESC
		LIMIT $2
	`

	rows, err := storage.DB.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.UnifiedPost
	for rows.Next() {
		var p domain.UnifiedPost
		var pubAt *time.Time // Handle nullable timestamp
		var publishTargets []byte

		err := rows.Scan(&p.Platform, &p.RemoteID, &p.Title, &p.URL, &p.Status, &p.Views, &p.Reactions, &p.Comments, &pubAt, &publishTargets)
		if err != nil {
			continue // Skip malformed rows
		}

		if pubAt != nil {
			p.PublishedAt = *pubAt
		}
		if len(publishTargets) > 0 {
			_ = json.Unmarshal(publishTargets, &p.PublishTargets)
		}
		posts = append(posts, p)
	}

	// 3. Set Cache (Async-ish, but blocking here for simplicity)
	if storage.RedisClient != nil && len(posts) > 0 {
		bytes, _ := json.Marshal(posts)
		// Cache for 10 minutes
		storage.RedisClient.Set(ctx, cacheKey, bytes, 10*time.Minute)
	}

	return posts, nil
}

// Live Fetch Methods

func (s *ActivityService) FetchLiveDevtoActivity(ctx context.Context, userID int, limit int) ([]automation.DevtoPost, error) {
	if limit > 50 {
		limit = 50
	}

	// Get Token
	token, err := s.getDevtoToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	posts, err := automation.FetchDevtoDashboardPosts(token, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *ActivityService) FetchLiveMediumActivity(ctx context.Context, userID int, limit int) ([]automation.MediumPost, error) {
	if limit > 50 {
		limit = 50
	}

	uid, sid, xsrf, err := s.getMediumCredentials(ctx, userID)
	if err != nil {
		return nil, err
	}

	posts, err := automation.FetchMediumPosts(uid, sid, xsrf, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// Helpers

func (s *ActivityService) getDevtoToken(ctx context.Context, userID int) (string, error) {
	creds, err := s.credsRepo.GetCredentials(ctx, userID, "devto")
	if err != nil {
		return "", err
	}
	if creds == nil {
		// Fallback to env
		if env := os.Getenv("DEVTO_SESSION_TOKEN"); env != "" {
			return env, nil
		}
		return "", errors.New("devto credentials missing")
	}

	var details map[string]interface{}
	if err := json.Unmarshal(creds.Credentials, &details); err == nil {
		if token, ok := details["remember_user_token"].(string); ok && token != "" {
			return token, nil
		}
	}
	return "", errors.New("invalid devto credentials format")
}

func (s *ActivityService) getMediumCredentials(ctx context.Context, userID int) (string, string, string, error) {
	creds, err := s.credsRepo.GetCredentials(ctx, userID, "medium")
	if err != nil {
		return "", "", "", err
	}
	if creds == nil {
		return "", "", "", errors.New("medium credentials missing")
	}

	var details map[string]interface{}
	if err := json.Unmarshal(creds.Credentials, &details); err == nil {
		uid, _ := details["uid"].(string)
		sid, _ := details["sid"].(string)
		xsrf, _ := details["xsrf"].(string)
		if uid != "" && sid != "" && xsrf != "" {
			return uid, sid, xsrf, nil
		}
	}
	return "", "", "", errors.New("incomplete medium credentials")
}
