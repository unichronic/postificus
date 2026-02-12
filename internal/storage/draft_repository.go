package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"postificus/internal/domain"
)

type DraftRepository interface {
	SaveDraft(ctx context.Context, draft *domain.Draft) error
	GetDraft(ctx context.Context, id string, userID int) (*domain.Draft, error)
	UpdateDashboardCache(ctx context.Context, draft *domain.Draft) error
}

type PostgresDraftRepository struct{}

func NewDraftRepository() *PostgresDraftRepository {
	return &PostgresDraftRepository{}
}

func (r *PostgresDraftRepository) SaveDraft(ctx context.Context, draft *domain.Draft) error {
	publishTargetsJSON, err := json.Marshal(draft.PublishTargets)
	if err != nil {
		return fmt.Errorf("failed to marshal publish targets: %w", err)
	}

	query := `
		INSERT INTO drafts (id, user_id, title, content, cover_image, publish_targets, last_saved_at, is_published)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), FALSE)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			cover_image = EXCLUDED.cover_image,
			publish_targets = EXCLUDED.publish_targets,
			last_saved_at = NOW();
	`

	_, err = DB.Exec(ctx, query, draft.ID, draft.UserID, draft.Title, draft.Content, draft.CoverImage, publishTargetsJSON)
	if err != nil {
		return fmt.Errorf("failed to execute save query: %w", err)
	}
	return nil
}

func (r *PostgresDraftRepository) UpdateDashboardCache(ctx context.Context, draft *domain.Draft) error {
	unifiedQuery := `
        INSERT INTO unified_posts (user_id, platform, remote_id, title, url, status, published_at, last_synced_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        ON CONFLICT (user_id, platform, remote_id)
        DO UPDATE SET
            title = EXCLUDED.title,
            url = EXCLUDED.url,
            status = EXCLUDED.status,
            published_at = EXCLUDED.published_at,
            last_synced_at = NOW();
    `
	// Platform is "postificus" for local drafts
	_, err := DB.Exec(ctx, unifiedQuery, draft.UserID, "postificus", draft.ID, draft.Title, fmt.Sprintf("/editor?draft=%s", draft.ID), "draft")
	if err != nil {
		return fmt.Errorf("failed to update dashboard cache: %w", err)
	}
	return nil
}

func (r *PostgresDraftRepository) GetDraft(ctx context.Context, id string, userID int) (*domain.Draft, error) {
	query := `
		SELECT title, content, cover_image, publish_targets, last_saved_at, is_published
		FROM drafts
		WHERE id = $1 AND user_id = $2
	`

	var (
		title          string
		content        string
		coverImage     string
		publishTargets []byte
		lastSavedAt    time.Time
		isPublished    bool
	)

	err := DB.QueryRow(ctx, query, id, userID).Scan(&title, &content, &coverImage, &publishTargets, &lastSavedAt, &isPublished)
	if err != nil {
		return nil, err
	}

	targets := []string{}
	if len(publishTargets) > 0 {
		_ = json.Unmarshal(publishTargets, &targets)
	}

	return &domain.Draft{
		ID:             id,
		UserID:         userID,
		Title:          title,
		Content:        content,
		CoverImage:     coverImage,
		PublishTargets: targets,
		LastSavedAt:    lastSavedAt,
		IsPublished:    isPublished,
	}, nil
}
