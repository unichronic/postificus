package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"postificus/internal/breaker"
	"postificus/internal/browser"
	"postificus/internal/metrics"
	"postificus/internal/storage"
)

const (
	TypePublishPost = "publish:post"
)

type PublishPayload struct {
	UserID     int      `json:"user_id"`
	Platform   string   `json:"platform"` // "medium", "linkedin", "devto"
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	CoverImage string   `json:"cover_image,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	BlogURL    string   `json:"blog_url,omitempty"` // For LinkedIn
}

// NewPublishPayload helper
func NewPublishPayload(p PublishPayload) ([]byte, error) {
	return json.Marshal(p)
}

// PublishService handles the execution of publishing tasks via browser automation.
type PublishService struct {
	credsRepo storage.CredentialsRepository
	breakers  map[string]*breaker.CircuitBreaker
}

func NewPublishService(credsRepo storage.CredentialsRepository) *PublishService {
	breakers := make(map[string]*breaker.CircuitBreaker)
	// Initialize breakers for known platforms
	breakers["medium"] = breaker.NewCircuitBreakerWithName("medium", 3, 1*time.Minute)
	breakers["devto"] = breaker.NewCircuitBreakerWithName("devto", 3, 1*time.Minute)
	breakers["linkedin"] = breaker.NewCircuitBreakerWithName("linkedin", 3, 1*time.Minute)

	return &PublishService{
		credsRepo: credsRepo,
		breakers:  breakers,
	}
}

// HandlePublishTask executes the actual publishing logic (Consumer Handler)
func (s *PublishService) HandlePublishTask(payload []byte) error {
	ctx := context.Background()

	// 1. Parse Payload
	var p PublishPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	log.Printf("Processing publish task for platform: %s, title: %s", p.Platform, p.Title)

	// 2. Fetch Credentials
	credsMap, err := s.fetchCredentials(ctx, p.UserID, p.Platform)
	if err != nil {
		log.Printf("Warning: Failed to fetch credentials from DB: %v. Falling back to Env.", err)
		credsMap = make(map[string]string)
	}

	// 3. Execute Automation based on Platform
	var url string

	switch p.Platform {
	case "medium":
		cb, ok := s.breakers["medium"]
		if !ok {
			cb = breaker.NewCircuitBreakerWithName("medium", 3, 1*time.Minute)
		}

		err = cb.Execute(func() error {
			uid := getCredential(credsMap, "uid", "MEDIUM_UID")
			sid := getCredential(credsMap, "sid", "MEDIUM_SID")
			xsrf := getCredential(credsMap, "xsrf", "MEDIUM_XSRF")

			if uid == "" || sid == "" {
				return fmt.Errorf("medium credentials missing")
			}
			return browser.PostToMediumWithTags(uid, sid, xsrf, p.Title, p.Content, p.Tags, p.CoverImage)
		})
		url = "https://medium.com/me/stories/public" // Placeholder

	case "devto":
		cb, ok := s.breakers["devto"]
		if !ok {
			cb = breaker.NewCircuitBreakerWithName("devto", 3, 1*time.Minute)
		}

		err = cb.Execute(func() error {
			token := getCredential(credsMap, "remember_user_token", "DEVTO_SESSION_TOKEN")
			if token == "" {
				token = getCredential(credsMap, "token", "DEVTO_SESSION_TOKEN") // Legacy fallback
			}
			if token == "" {
				return fmt.Errorf("devto credentials missing")
			}
			return browser.PostToDevToWithCookie(token, p.Title, p.Content, p.CoverImage, p.Tags)
		})
		url = "https://dev.to/dashboard" // Placeholder

	case "linkedin":
		cb, ok := s.breakers["linkedin"]
		if !ok {
			cb = breaker.NewCircuitBreakerWithName("linkedin", 3, 1*time.Minute)
		}

		err = cb.Execute(func() error {
			liAt := getCredential(credsMap, "li_at", "LI_AT")
			if liAt == "" {
				return fmt.Errorf("linkedin credentials missing")
			}
			return browser.PostToLinkedIn(liAt, p.Content, p.BlogURL)
		})
		url = "https://www.linkedin.com/feed/" // Placeholder

	default:
		return fmt.Errorf("unsupported platform: %s", p.Platform)
	}

	start := time.Now()
	if err != nil {
		metrics.PostPublishTotal.WithLabelValues(p.Platform, "error").Inc()
		return fmt.Errorf("publish failed: %w", err)
	}

	metrics.PostPublishTotal.WithLabelValues(p.Platform, "success").Inc()
	metrics.PostPublishDuration.WithLabelValues(p.Platform).Observe(time.Since(start).Seconds())

	log.Printf("✅ Published to %s: %s", p.Platform, url)
	return nil
}

// Helper to fetch credentials from DB and unmarshal them
func (s *PublishService) fetchCredentials(ctx context.Context, userID int, platform string) (map[string]string, error) {
	cred, err := s.credsRepo.GetCredentials(ctx, userID, platform)
	if err != nil {
		return nil, err
	}
	if cred == nil {
		return nil, fmt.Errorf("credentials not found")
	}

	var credsMap map[string]string
	if err := json.Unmarshal(cred.Credentials, &credsMap); err != nil {
		return nil, err
	}
	return credsMap, nil
}

// Helper to get credential from map or env
func getCredential(creds map[string]string, key, envVar string) string {
	if val, ok := creds[key]; ok && val != "" {
		return val
	}
	return os.Getenv(envVar)
}
