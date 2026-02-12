package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	automation "postificus/internal/browser"
	"postificus/internal/domain"
	"time"
)

const (
	TypeSyncPlatformActivity = "task:sync_platform_activity"
)

type SyncPlatformPayload struct {
	UserID   int    `json:"user_id"`
	Platform string `json:"platform"`
}

// SyncService handles background synchronization tasks.
type SyncService struct {
	activityService *ActivityService
}

func NewSyncService(activityService *ActivityService) *SyncService {
	return &SyncService{
		activityService: activityService,
	}
}

// HandleSyncPlatformActivity is the handler for the sync task.
func (s *SyncService) HandleSyncPlatformActivity(payload []byte) error {
	ctx := context.Background()
	var p SyncPlatformPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	log.Printf("🔄 [Worker] Starting sync for User %d - Platform: %s", p.UserID, p.Platform)

	var posts []domain.UnifiedPost
	var err error

	// 1. Fetch from Platform Automation via ActivityService
	switch p.Platform {
	case "medium":
		var mediumPosts []automation.MediumPost // Explicit declaration
		mediumPosts, err = s.activityService.FetchLiveMediumActivity(ctx, p.UserID, 20)
		if err == nil {
			for _, mp := range mediumPosts {
				posts = append(posts, domain.UnifiedPost{
					Platform:    "medium",
					RemoteID:    mp.URL,
					Title:       mp.Title,
					URL:         mp.URL,
					Status:      mp.Status,
					PublishedAt: time.Now(),
				})
			}
		} else {
			err = fmt.Errorf("medium fetch failed: %w", err)
		}
	case "devto":
		var devtoPosts []automation.DevtoPost // Explicit declaration
		devtoPosts, err = s.activityService.FetchLiveDevtoActivity(ctx, p.UserID, 20)
		if err == nil {
			for _, dp := range devtoPosts {
				views := 0
				if dp.ViewsCount != nil {
					views = *dp.ViewsCount
				}
				reactions := 0
				if dp.Reactions != nil {
					reactions = *dp.Reactions
				}
				comments := 0
				if dp.Comments != nil {
					comments = *dp.Comments
				}

				var publishedAt time.Time
				if dp.UpdatedAt != "" {
					if t, err := time.Parse(time.RFC3339, dp.UpdatedAt); err == nil {
						publishedAt = t
					} else {
						log.Printf("⚠️ Failed to parse devto date '%s': %v", dp.UpdatedAt, err)
					}
				}

				posts = append(posts, domain.UnifiedPost{
					Platform:    "devto",
					RemoteID:    dp.URL,
					Title:       dp.Title,
					URL:         dp.URL,
					Status:      dp.Status,
					Views:       views,
					Reactions:   reactions,
					Comments:    comments,
					PublishedAt: publishedAt,
				})
			}
		} else {
			err = fmt.Errorf("devto fetch failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown platform: %s", p.Platform)
	}

	if err != nil {
		log.Printf("❌ [Worker] Sync failed for %s: %v", p.Platform, err)
		return err
	}

	// 2. Upsert into DB
	log.Printf("📥 [Worker] Saving %d posts for %s...", len(posts), p.Platform)
	for _, post := range posts {
		if err := s.activityService.UpsertPost(ctx, p.UserID, post); err != nil {
			log.Printf("⚠️ [Worker] Failed to save post %s: %v", post.Title, err)
			continue
		}
	}

	log.Printf("✅ [Worker] Sync complete for %s", p.Platform)
	return nil
}
