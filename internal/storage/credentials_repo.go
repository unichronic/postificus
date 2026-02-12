package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"postificus/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
)

// CredentialsRepository defines the interface for credential storage
type CredentialsRepository interface {
	SaveCredentials(ctx context.Context, userID int, platform string, creds map[string]string) error
	GetCredentials(ctx context.Context, userID int, platform string) (*domain.UserCredential, error)
}

// PostgresCredentialsRepository implements CredentialsRepository
type PostgresCredentialsRepository struct {
}

// NewCredentialsRepository creates a new instance
func NewCredentialsRepository() *PostgresCredentialsRepository {
	return &PostgresCredentialsRepository{}
}

// SaveCredentials upserts user credentials into the database
func (r *PostgresCredentialsRepository) SaveCredentials(ctx context.Context, userID int, platform string, creds map[string]string) error {
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	query := `
		INSERT INTO user_credentials (user_id, platform, credentials, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, platform) 
		DO UPDATE SET credentials = $3, updated_at = NOW();
	`

	_, err = DB.Exec(ctx, query, userID, platform, credsJSON)
	if err != nil {
		return fmt.Errorf("failed to execute save query: %w", err)
	}

	return nil
}

func (r *PostgresCredentialsRepository) GetCredentials(ctx context.Context, userID int, platform string) (*domain.UserCredential, error) {
	var credsJSON []byte
	var updatedAt time.Time

	// We select UpdatedAt just to fill the struct properly
	query := `SELECT credentials, updated_at FROM user_credentials WHERE user_id = $1 AND platform = $2`

	err := DB.QueryRow(ctx, query, userID, platform).Scan(&credsJSON, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found is not an error in this context, just nil
		}
		return nil, fmt.Errorf("failed to fetch credentials: %w", err)
	}

	return &domain.UserCredential{
		UserID:      userID,
		Platform:    platform,
		Credentials: json.RawMessage(credsJSON),
		UpdatedAt:   updatedAt,
	}, nil
}
