package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// MediumAPIClient handles direct HTTP API calls to Medium
type MediumAPIClient struct {
	UID           string
	SID           string
	XSRF          string
	Client        *http.Client
	revisionCount int // Track revision count for publishing
}

// Medium API response wrapper
type mediumResponse struct {
	Success bool            `json:"success"`
	Payload json.RawMessage `json:"payload"`
}

// Delta structures for Medium's content format
type mediumDelta struct {
	Type           int              `json:"type"`
	Index          int              `json:"index"`
	Section        *mediumSection   `json:"section,omitempty"`
	Paragraph      *mediumParagraph `json:"paragraph,omitempty"`
	VerifySameName bool             `json:"verifySameName,omitempty"`
}

type mediumSection struct {
	Name       string `json:"name"`
	StartIndex int    `json:"startIndex"`
}

type mediumParagraph struct {
	Name    string   `json:"name"`
	Type    int      `json:"type"`
	Text    string   `json:"text"`
	Markups []string `json:"markups"`
}

// NewMediumAPIClient creates a new API client
func NewMediumAPIClient(uid, sid, xsrf string) *MediumAPIClient {
	return &MediumAPIClient{
		UID:  uid,
		SID:  sid,
		XSRF: xsrf,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// stripHijackingPrefix removes Medium's XSS protection prefix
func stripHijackingPrefix(body []byte) []byte {
	prefix := []byte("])}while(1);</x>")
	return bytes.TrimPrefix(body, prefix)
}

// makeRequest performs an authenticated request to Medium
func (c *MediumAPIClient) makeRequest(method, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers (matching browser behavior)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-XSRF-Token", c.XSRF)
	req.Header.Set("X-Obvious-CID", "web")
	req.Header.Set("X-Client-Date", fmt.Sprintf("%d", time.Now().UnixMilli()))
	req.Header.Set("Origin", "https://medium.com")
	req.Header.Set("Referer", "https://medium.com/new-story")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0")

	// Set cookies
	cookieStr := fmt.Sprintf("uid=%s; sid=%s; xsrf=%s", c.UID, c.SID, c.XSRF)
	req.Header.Set("Cookie", cookieStr)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Strip Medium's XSS protection prefix
	return stripHijackingPrefix(respBody), nil
}

// CreateStory creates a new draft story and returns the post ID
func (c *MediumAPIClient) CreateStory() (string, error) {
	log.Println("📝 Creating new Medium story...")

	payload := map[string]interface{}{
		"deltas":     []interface{}{},
		"baseRev":    -1,
		"coverless":  true,
		"visibility": 0,
	}

	respBody, err := c.makeRequest("POST", "https://medium.com/new-story", payload)
	if err != nil {
		return "", fmt.Errorf("create story: %w", err)
	}

	// Parse response to extract post ID
	var resp mediumResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	// Extract post ID from payload
	// The response contains the post data, and we need to find the ID
	// Example: {"success":true,"payload":{"value":{"id":"5d18830dea37",...}}
	var payloadData map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &payloadData); err != nil {
		return "", fmt.Errorf("parse payload: %w", err)
	}

	// Navigate nested structure to get ID
	if value, ok := payloadData["value"].(map[string]interface{}); ok {
		if id, ok := value["id"].(string); ok {
			log.Printf("✅ Created story with ID: %s", id)
			return id, nil
		}
	}

	return "", fmt.Errorf("post ID not found in response")
}

// generateParagraphName generates a random 4-character name for paragraphs
func generateParagraphName() string {
	// Medium uses 4-character hex-like names (e.g., "54af", "59f6")
	// For simplicity, we'll use timestamp-based generation
	now := time.Now().UnixNano()
	return fmt.Sprintf("%04x", now%0x10000)
}

// UpdateContent updates the story content using Medium's delta format
func (c *MediumAPIClient) UpdateContent(postID, title, content string) error {
	log.Println("✍️  Updating story content...")

	// Reset revision counter
	c.revisionCount = 0

	// Build deltas for title and body
	deltas := []mediumDelta{
		// Create section
		{
			Type:  8,
			Index: 0,
			Section: &mediumSection{
				Name:       generateParagraphName(),
				StartIndex: 0,
			},
		},
		// Insert title paragraph (type 3 = title)
		{
			Type:  1,
			Index: 0,
			Paragraph: &mediumParagraph{
				Name:    generateParagraphName(),
				Type:    3,
				Text:    "",
				Markups: []string{},
			},
		},
		// Update title with actual text
		{
			Type:  3,
			Index: 0,
			Paragraph: &mediumParagraph{
				Name:    "", // Will be set below
				Type:    3,
				Text:    title,
				Markups: []string{},
			},
			VerifySameName: true,
		},
	}

	// Store the title paragraph name for the update delta
	deltas[2].Paragraph.Name = deltas[1].Paragraph.Name

	// Split content into paragraphs (simple approach: by newlines)
	paragraphs := strings.Split(strings.TrimSpace(content), "\n\n")

	idx := 1
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		paraName := generateParagraphName()

		// Insert paragraph (type 1 = body text)
		deltas = append(deltas, mediumDelta{
			Type:  1,
			Index: idx,
			Paragraph: &mediumParagraph{
				Name:    paraName,
				Type:    1,
				Text:    "",
				Markups: []string{},
			},
		})

		// Update paragraph with text
		deltas = append(deltas, mediumDelta{
			Type:  3,
			Index: idx,
			Paragraph: &mediumParagraph{
				Name:    paraName,
				Type:    1,
				Text:    para,
				Markups: []string{},
			},
			VerifySameName: true,
		})

		idx++
	}

	payload := map[string]interface{}{
		"id":      postID,
		"deltas":  deltas,
		"baseRev": -1,
	}

	url := fmt.Sprintf("https://medium.com/p/%s/deltas", postID)
	_, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("update content: %w", err)
	}

	// Update revision count based on number of deltas applied
	c.revisionCount = 0

	log.Printf("✅ Content updated (revision: %d)", c.revisionCount)
	return nil
}

// UpdateMetadata sets tags and other metadata
func (c *MediumAPIClient) UpdateMetadata(postID string, tags []string) error {
	log.Println("🏷️  Updating metadata...")

	payload := map[string]interface{}{
		"tags":                tags,
		"noteToCurator":       "",
		"allowCuration":       false,
		"notifyTwitter":       false,
		"pinnedPost":          false,
		"isPublishToEmail":    false,
		"isMarkedPaywallOnly": false,
	}

	url := fmt.Sprintf("https://medium.com/_/api/posts/%s/metadata", postID)
	_, err := c.makeRequest("PUT", url, payload)
	if err != nil {
		return fmt.Errorf("update metadata: %w", err)
	}

	log.Println("✅ Metadata updated")
	return nil
}

// PublishPost publishes the story
func (c *MediumAPIClient) PublishPost(postID string) (string, error) {
	log.Println("🚀 Publishing story...")

	// Use actual revision count from deltas
	latestRev := c.revisionCount
	if latestRev == 0 {
		// Fallback to -1 if no deltas were tracked
		latestRev = -1
	}

	payload := map[string]interface{}{
		"title":           "",
		"subtitle":        "",
		"metaDescription": "",
		"latestRev":       latestRev,
	}

	url := fmt.Sprintf("https://medium.com/p/%s/publish", postID)
	respBody, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("publish post: %w", err)
	}

	// Parse response to get published URL
	var resp mediumResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	// Extract URL from payload
	var payloadData map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &payloadData); err != nil {
		return "", fmt.Errorf("parse payload: %w", err)
	}

	if value, ok := payloadData["value"].(map[string]interface{}); ok {
		if mediumURL, ok := value["mediumUrl"].(string); ok {
			log.Printf("✅ Published: %s", mediumURL)
			return mediumURL, nil
		}
	}

	// Fallback: construct URL from post ID
	publishedURL := fmt.Sprintf("https://medium.com/p/%s", postID)
	log.Printf("✅ Published (inferred URL): %s", publishedURL)
	return publishedURL, nil
}

// Publish is the main method that orchestrates the entire publishing flow
func (c *MediumAPIClient) Publish(title, content string, tags []string) (string, error) {
	log.Println("🌐 Starting Medium API publish flow...")

	// Step 1: Create new story
	postID, err := c.CreateStory()
	if err != nil {
		return "", err
	}

	// Step 2: Update content
	if err := c.UpdateContent(postID, title, content); err != nil {
		return "", err
	}

	// Step 3: Update metadata (tags)
	if len(tags) > 0 {
		if err := c.UpdateMetadata(postID, tags); err != nil {
			return "", err
		}
	}

	// Step 4: Publish
	publishedURL, err := c.PublishPost(postID)
	if err != nil {
		return "", err
	}

	return publishedURL, nil
}

// ExtractPostIDFromURL extracts post ID from Medium URLs
func ExtractPostIDFromURL(url string) string {
	// Match patterns like /p/5d18830dea37/edit or /p/5d18830dea37
	re := regexp.MustCompile(`/p/([a-f0-9]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
