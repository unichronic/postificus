package browser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type MediumPost struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Status      string `json:"status,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	Claps       string `json:"claps,omitempty"`
	Responses   string `json:"responses,omitempty"`
}

// FetchMediumPosts scrapes the Medium stories page using session cookies.
func FetchMediumPosts(uid, sid, xsrf string, limit int) ([]MediumPost, error) {
	if uid == "" || sid == "" || xsrf == "" {
		return nil, fmt.Errorf("medium credentials missing")
	}

	if err := EnsureBrowser(); err != nil {
		return nil, fmt.Errorf("failed to initialize browser: %w", err)
	}

	// Use a new page on the shared browser
	page := Browser.MustPage("about:blank")
	defer page.MustClose()

	// Set Cookies
	cookies := []*proto.NetworkCookieParam{
		{
			Name:     "uid",
			Value:    uid,
			Domain:   ".medium.com",
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
		},
		{
			Name:     "sid",
			Value:    sid,
			Domain:   ".medium.com",
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
		},
		{
			Name:     "xsrf",
			Value:    xsrf,
			Domain:   ".medium.com",
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
		},
	}

	if err := page.SetCookies(cookies); err != nil {
		return nil, fmt.Errorf("failed to set cookies: %w", err)
	}

	// Navigate to Stories -> Published
	// https://medium.com/me/stories?tab=posts-published
	page.MustNavigate("https://medium.com/me/stories?tab=posts-published")
	page.MustWaitLoad()

	// Wait for list to load
	// Prefer table rows, fall back to headings if the layout changes.
	err := page.Timeout(10*time.Second).WaitElementsMoreThan("table tbody tr", 0)
	if err != nil {
		// Might be no stories or not logged in
		if isMediumLoginURL(page.MustInfo().URL) {
			return nil, fmt.Errorf("cookie expired or login failed")
		}
		// If just no stories or layout changed, continue with best-effort extraction.
	}

	return extractMediumPosts(page, limit)
}

func extractMediumPosts(page *rod.Page, limit int) ([]MediumPost, error) {
	var posts []MediumPost
	rows, err := page.Elements("table tbody tr")
	if err == nil && len(rows) > 0 {
		return extractMediumPostsFromRows(rows, limit), nil
	}

	// Fallback: grab headings if layout changes (best-effort).
	headings, err := page.Elements("h2, h3")
	if err != nil {
		return nil, err
	}
	for _, el := range headings {
		if limit > 0 && len(posts) >= limit {
			break
		}
		title := strings.TrimSpace(el.MustText())
		if title == "" {
			continue
		}
		link := findMediumHeadingLink(el)
		if link == nil {
			continue
		}
		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}
		url := normalizeMediumURL(*href)
		if url == "" {
			continue
		}
		posts = append(posts, MediumPost{
			Title:  title,
			URL:    url,
			Status: "published",
		})
	}
	return posts, nil
}

func extractMediumPostsFromRows(rows []*rod.Element, limit int) []MediumPost {
	var posts []MediumPost

	for _, row := range rows {
		if limit > 0 && len(posts) >= limit {
			break
		}

		linkEl := findMediumPostLink(row)
		if linkEl == nil {
			continue
		}

		href, err := linkEl.Attribute("href")
		if err != nil || href == nil || strings.TrimSpace(*href) == "" {
			continue
		}
		url := normalizeMediumURL(*href)
		if url == "" {
			continue
		}

		title := ""
		if h2, err := linkEl.Element("h2"); err == nil {
			title = strings.TrimSpace(h2.MustText())
		}
		if title == "" {
			if h2, err := row.Element("h2"); err == nil {
				title = strings.TrimSpace(h2.MustText())
			}
		}
		if title == "" {
			title = strings.TrimSpace(linkEl.MustText())
		}
		if title == "" {
			continue
		}

		publishedAt := extractMediumPublishedText(row)
		claps, responses := extractMediumStats(row)

		posts = append(posts, MediumPost{
			Title:       title,
			URL:         url,
			Status:      "published",
			PublishedAt: publishedAt,
			Claps:       claps,
			Responses:   responses,
		})
	}

	return posts
}

func findMediumPostLink(row *rod.Element) *rod.Element {
	anchors, err := row.Elements("a")
	if err != nil {
		return nil
	}
	for _, a := range anchors {
		href, err := a.Attribute("href")
		if err != nil || href == nil {
			continue
		}
		val := strings.TrimSpace(*href)
		if strings.Contains(val, "/@") || strings.Contains(val, "/p/") {
			return a
		}
	}
	return nil
}

func findMediumHeadingLink(el *rod.Element) *rod.Element {
	if link, err := el.Element("a"); err == nil {
		if href, err := link.Attribute("href"); err == nil && href != nil && strings.TrimSpace(*href) != "" {
			return link
		}
	}

	current := el
	for i := 0; i < 3; i++ {
		if href, err := current.Attribute("href"); err == nil && href != nil && strings.TrimSpace(*href) != "" {
			return current
		}
		parent, err := current.Parent()
		if err != nil {
			break
		}
		current = parent
	}

	return nil
}

func extractMediumPublishedText(row *rod.Element) string {
	paras, err := row.Elements("p")
	if err != nil {
		return ""
	}
	for _, p := range paras {
		text := normalizeWhitespace(p.MustText())
		if strings.HasPrefix(text, "Published ") {
			return text
		}
	}
	for _, p := range paras {
		text := normalizeWhitespace(p.MustText())
		if strings.HasPrefix(text, "Updated ") {
			return text
		}
	}
	return ""
}

func extractMediumStats(row *rod.Element) (string, string) {
	if statNodes, err := row.Elements("svg + p"); err == nil && len(statNodes) > 0 {
		claps, responses := firstTwoNumeric(statNodes)
		if claps != "" || responses != "" {
			return claps, responses
		}
	}

	paras, err := row.Elements("p")
	if err != nil {
		return "", ""
	}

	return firstTwoNumeric(paras)
}

func firstTwoNumeric(elements []*rod.Element) (string, string) {
	var counts []string
	for _, el := range elements {
		text := normalizeWhitespace(el.MustText())
		if mediumCountRegex.MatchString(text) {
			counts = append(counts, text)
		}
	}

	claps := ""
	responses := ""
	if len(counts) >= 1 {
		claps = counts[0]
	}
	if len(counts) >= 2 {
		responses = counts[1]
	}
	return claps, responses
}

func normalizeMediumURL(href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}
	if strings.HasPrefix(href, "/") {
		href = "https://medium.com" + href
	}
	if idx := strings.Index(href, "?"); idx != -1 {
		href = href[:idx]
	}
	return href
}

func isMediumLoginURL(url string) bool {
	lower := strings.ToLower(url)
	return lower == "https://medium.com/" ||
		strings.Contains(lower, "signin") ||
		strings.Contains(lower, "login") ||
		strings.Contains(lower, "ap/signin")
}

var mediumCountRegex = regexp.MustCompile(`^\d+(?:[\.,]\d+)?[KM]?$`)
