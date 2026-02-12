package domain

import (
	"encoding/json"
	"time"
)

// UserCredential represents the stored credentials for a platform
type UserCredential struct {
	UserID      int             `json:"user_id"`
	Platform    string          `json:"platform"`
	Credentials json.RawMessage `json:"credentials"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// Draft represents a blog post draft
type Draft struct {
	ID             string    `json:"id"`
	UserID         int       `json:"user_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CoverImage     string    `json:"cover_image"`
	PublishTargets []string  `json:"publish_targets"`
	LastSavedAt    time.Time `json:"last_saved_at"`
	IsPublished    bool      `json:"is_published"`
}

// Profile represents user profile information
type Profile struct {
	UserID      int      `json:"user_id"`
	FullName    string   `json:"full_name"`
	Username    string   `json:"username"`
	Headline    string   `json:"headline"`
	Bio         string   `json:"bio"`
	Location    string   `json:"location"`
	Website     string   `json:"website"`
	PublicEmail string   `json:"public_email"`
	Skills      []string `json:"skills"`
}

// UnifiedPost represents a post from any platform (DB or Live)
type UnifiedPost struct {
	Platform       string    `json:"platform"`
	RemoteID       string    `json:"remote_id"` // ID on the platform (or URL)
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	Status         string    `json:"status"` // draft, published
	Views          int       `json:"views"`
	Reactions      int       `json:"reactions"`
	Comments       int       `json:"comments"`
	PublishedAt    time.Time `json:"published_at"`
	PublishTargets []string  `json:"publish_targets,omitempty"`
}
