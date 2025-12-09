package llm

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Summarizer struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewSummarizer(ctx context.Context) (*Summarizer, error) {
	// Check if LLM is disabled
	if os.Getenv("DISABLE_LLM") == "true" {
		fmt.Println("⚠️ LLM IS DISABLED: Using Mock Mode")
		return &Summarizer{client: nil, model: nil}, nil
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	// Use Flash for speed and cost
	model := client.GenerativeModel("gemini-2.5-flash")

	// Set temperature low (0.3) for consistent, professional output
	model.SetTemperature(0.6)

	return &Summarizer{client: client, model: model}, nil
}

func (s *Summarizer) GenerateSocialContent(ctx context.Context, blogContent string) (string, error) {
	// 1. Defend against empty or massive payloads
	if len(blogContent) > 500000 { // ~100k words
		return "", fmt.Errorf("content too large (max 500k chars)")
	}

	// Mock Mode
	if s.client == nil {
		fmt.Println("⚠️ MOCK LLM: Returning dummy content")
		return "🚀 This is a MOCK LinkedIn post because DISABLE_LLM=true.\n\nIt simulates the viral vibe without calling the API.\n\n👉 Link in comments!", nil
	}

	// 2. The Viral Prompt
	prompt := fmt.Sprintf(`
		Act as a LinkedIn ghostwriter. Analyze the following blog post and generate a "Viral Style" LinkedIn post.
		
		Constraints:
		- Hook: First line must be contrarian or a "How I" statement.
		- Format: Short lines. White space. No markdown bolding.
		- Emoji: Use sparingly (👉, 🚀).
		- Length: Under 200 words.
		- Goal: Drive clicks to the link.

		Blog Content:
		%s
	`, blogContent)

	// 3. Call API
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini api failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("empty response from llm")
	}

	// 4. Extract Text
	part := resp.Candidates[0].Content.Parts[0]
	result := fmt.Sprintf("%s", part)

	// Debug: Save to file
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("llm_output_%s.txt", timestamp)
	if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
		fmt.Printf("Warning: Failed to save LLM debug output: %v\n", err)
	} else {
		fmt.Printf("Debug: LLM output saved to %s\n", filename)
	}

	return result, nil
}
