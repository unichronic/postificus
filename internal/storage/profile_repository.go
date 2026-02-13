package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"postificus/internal/domain"

	"github.com/jackc/pgx/v5"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, userID string) (*domain.Profile, error)
	SaveProfile(ctx context.Context, profile *domain.Profile) error
}

type PostgresProfileRepository struct{}

func NewProfileRepository() *PostgresProfileRepository {
	return &PostgresProfileRepository{}
}

func (r *PostgresProfileRepository) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	query := `
		SELECT full_name, username, headline, bio, location, website, public_email, skills
		FROM user_details
		WHERE user_id = $1
	`

	var (
		fullName    string
		username    string
		headline    string
		bio         string
		location    string
		website     string
		publicEmail string
		skillsJSON  []byte
	)

	err := DB.QueryRow(ctx, query, userID).Scan(
		&fullName,
		&username,
		&headline,
		&bio,
		&location,
		&website,
		&publicEmail,
		&skillsJSON,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return &domain.Profile{UserID: userID, Skills: []string{}}, nil
		}
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}

	skills := []string{}
	if len(skillsJSON) > 0 {
		_ = json.Unmarshal(skillsJSON, &skills)
	}

	return &domain.Profile{
		UserID:      userID,
		FullName:    fullName,
		Username:    username,
		Headline:    headline,
		Bio:         bio,
		Location:    location,
		Website:     website,
		PublicEmail: publicEmail,
		Skills:      skills,
	}, nil
}

func (r *PostgresProfileRepository) SaveProfile(ctx context.Context, profile *domain.Profile) error {
	skillsJSON, err := json.Marshal(profile.Skills)
	if err != nil {
		return fmt.Errorf("failed to encode skills: %w", err)
	}

	query := `
		INSERT INTO user_details (
			user_id, full_name, username, headline, bio, location, website, public_email, skills, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			full_name = $2,
			username = $3,
			headline = $4,
			bio = $5,
			location = $6,
			website = $7,
			public_email = $8,
			skills = $9,
			updated_at = NOW()
	`

	_, err = DB.Exec(
		ctx,
		query,
		profile.UserID,
		profile.FullName,
		profile.Username,
		profile.Headline,
		profile.Bio,
		profile.Location,
		profile.Website,
		profile.PublicEmail,
		skillsJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}
	return nil
}
