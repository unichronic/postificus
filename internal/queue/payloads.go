package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypePublishPost = "publish:post"
)

type PublishPayload struct {
	UserID   int      `json:"user_id"`
	Platform string   `json:"platform"` // "medium", "linkedin"
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags,omitempty"`
	BlogURL  string   `json:"blog_url,omitempty"` // For LinkedIn
}

// NewPublishTask creates a new task for publishing a post.
func NewPublishTask(p PublishPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePublishPost, payload), nil
}
