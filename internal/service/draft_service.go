package service

import (
	"context"
	"fmt"
	"time"

	"postificus/internal/domain"
	"postificus/internal/storage"
)

type DraftService struct {
	draftRepo storage.DraftRepository
}

func NewDraftService(draftRepo storage.DraftRepository) *DraftService {
	return &DraftService{draftRepo: draftRepo}
}

func (s *DraftService) SaveDraft(ctx context.Context, draft *domain.Draft) error {
	// 1. FAST PATH: Save to Redis
	// "SET draft:123:content '...' EX 3600" (Expire in 1 hour if inactive)
	err := storage.RedisClient.Set(ctx, fmt.Sprintf("draft:%s:content", draft.ID), draft.Content, 1*time.Hour).Err()
	if err != nil {
		// Log error but continue to DB save? Or return?
		// Handler returned error, so we will too, but strictly speaking DB save is more important.
		// Let's return error to warn client.
		return fmt.Errorf("failed to cache draft: %w", err)
	}

	// 2. Persist to DB
	if err := s.draftRepo.SaveDraft(ctx, draft); err != nil {
		return fmt.Errorf("failed to persist draft: %w", err)
	}

	// 3. Update Dashboard Cache
	if err := s.draftRepo.UpdateDashboardCache(ctx, draft); err != nil {
		return fmt.Errorf("failed to update dashboard cache: %w", err)
	}

	return nil
}

func (s *DraftService) GetDraft(ctx context.Context, id string, userID string) (*domain.Draft, error) {
	return s.draftRepo.GetDraft(ctx, id, userID)
}
