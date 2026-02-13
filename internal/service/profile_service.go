package service

import (
	"context"
	"postificus/internal/domain"
	"postificus/internal/storage" // Changed from "postificus/internal/repository"
)

type ProfileService struct {
	profileRepo storage.ProfileRepository // Changed from repo repository.ProfileRepository
}

func NewProfileService(profileRepo storage.ProfileRepository) *ProfileService { // Changed from repository.ProfileRepository
	return &ProfileService{profileRepo: profileRepo} // Changed from repo: repo
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	return s.profileRepo.GetProfile(ctx, userID)
}

func (s *ProfileService) SaveProfile(ctx context.Context, profile *domain.Profile) error {
	return s.profileRepo.SaveProfile(ctx, profile)
}
