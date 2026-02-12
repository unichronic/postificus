package browser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type DevtoPost struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Status     string `json:"status,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	Reactions  *int   `json:"reactions,omitempty"`
	Comments   *int   `json:"comments,omitempty"`
	Views      string `json:"views,omitempty"`
	ViewsCount *int   `json:"views_count,omitempty"`
}

// FetchDevtoDashboardPosts scrapes the Dev.to dashboard using a session cookie.
// It returns a best-effort list of the user's posts with basic metadata.
func FetchDevtoDashboardPosts(sessionToken string, limit int) ([]DevtoPost, error) {
	if sessionToken == "" {
		return nil, fmt.Errorf("devto session token missing")
	}

	if err := EnsureBrowser(); err != nil {
		return nil, fmt.Errorf("failed to initialize browser: %w", err)
	}

	page := Browser.MustPage("https://dev.to/404")
	defer page.MustClose()

	cookie := &proto.NetworkCookieParam{
		Name:     "remember_user_token",
		Value:    sessionToken,
		Domain:   "dev.to",
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: proto.NetworkCookieSameSiteLax,
	}

	if err := page.SetCookies([]*proto.NetworkCookieParam{cookie}); err != nil {
		return nil, fmt.Errorf("failed to set cookie: %w", err)
	}

	page.MustNavigate("https://dev.to/dashboard")
	page.MustWaitLoad()
	page.MustWaitStable()

	if page.MustInfo().URL == "https://dev.to/enter" {
		return nil, fmt.Errorf("cookie expired: session token invalid")
	}

	posts, err := extractDevtoDashboardPosts(page)
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(posts) > limit {
		posts = posts[:limit]
	}

	return posts, nil
}

func extractDevtoDashboardPosts(page *rod.Page) ([]DevtoPost, error) {
	containers, err := findDevtoPostContainers(page)
	if err != nil {
		return nil, err
	}

	posts := make([]DevtoPost, 0, len(containers))
	seen := map[string]struct{}{}

	for _, el := range containers {
		title := firstText(el, []string{
			".dashboard-story__title a",
			".crayons-story__title a",
			"h3 a",
			"h2 a",
			"a",
		})

		href := firstAttr(el, []string{
			".dashboard-story__title a",
			".crayons-story__title a",
			"h3 a",
			"h2 a",
			"a",
		}, "href")

		if title == "" || href == "" {
			continue
		}

		if !isLikelyDevtoPostURL(href) {
			continue
		}

		url := normalizeDevtoURL(href)
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}

		reactions, comments, viewsLabel, viewsCount := findDevtoStats(el)

		posts = append(posts, DevtoPost{
			Title:      title,
			URL:        url,
			Status:     findDevtoStatus(el),
			UpdatedAt:  findDevtoTimestamp(el),
			Reactions:  reactions,
			Comments:   comments,
			Views:      viewsLabel,
			ViewsCount: viewsCount,
		})
	}

	return posts, nil
}

func findDevtoPostContainers(page *rod.Page) ([]*rod.Element, error) {
	selectors := []string{
		`.dashboard-story`,
		`.js-dashboard-story`,
		`.crayons-card .dashboard-story`,
		`.crayons-card .js-dashboard-story`,
		`[data-testid="dashboard-article"]`,
		`[data-testid="dashboard-post"]`,
		`.dashboard__article`,
		`.crayons-story`,
		`.article-card`,
		`.crayons-card`,
	}

	for _, sel := range selectors {
		els, err := page.Elements(sel)
		if err == nil && len(els) > 0 {
			return els, nil
		}
	}

	// Fallback: grab any article cards and hope for the best.
	els, err := page.Elements("article")
	if err != nil {
		return nil, err
	}
	return els, nil
}

func firstText(el *rod.Element, selectors []string) string {
	for _, sel := range selectors {
		child, err := el.Element(sel)
		if err != nil {
			continue
		}
		text := strings.TrimSpace(child.MustText())
		if text != "" {
			return text
		}
	}
	return ""
}

func firstAttr(el *rod.Element, selectors []string, attr string) string {
	for _, sel := range selectors {
		child, err := el.Element(sel)
		if err != nil {
			continue
		}
		val, err := child.Attribute(attr)
		if err == nil && val != nil && strings.TrimSpace(*val) != "" {
			return strings.TrimSpace(*val)
		}
	}
	return ""
}

func findDevtoStatus(el *rod.Element) string {
	if strong, err := el.Element(".js-dashboard-story-details strong"); err == nil {
		label := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(strong.MustText(), ":")))
		switch label {
		case "published":
			return "published"
		case "draft":
			return "draft"
		case "scheduled":
			return "scheduled"
		case "archived":
			return "archived"
		case "unlisted":
			return "unlisted"
		}
	}

	labels, _ := el.Elements(".crayons-pill, .crayons-tag, .crayons-badge, .dashboard-article__status, .article-status")
	for _, label := range labels {
		text := strings.ToLower(strings.TrimSpace(label.MustText()))
		switch {
		case strings.Contains(text, "draft"):
			return "draft"
		case strings.Contains(text, "published"):
			return "published"
		case strings.Contains(text, "archived"):
			return "archived"
		case strings.Contains(text, "scheduled"):
			return "scheduled"
		case strings.Contains(text, "unlisted"):
			return "unlisted"
		}
	}
	return ""
}

func findDevtoTimestamp(el *rod.Element) string {
	if t, err := el.Element(".js-dashboard-story-details time"); err == nil {
		if dt, err := t.Attribute("datetime"); err == nil && dt != nil && strings.TrimSpace(*dt) != "" {
			return strings.TrimSpace(*dt)
		}
		text := strings.TrimSpace(t.MustText())
		if text != "" {
			return text
		}
	}

	if t, err := el.Element("time"); err == nil {
		if dt, err := t.Attribute("datetime"); err == nil && dt != nil && strings.TrimSpace(*dt) != "" {
			return strings.TrimSpace(*dt)
		}
		text := strings.TrimSpace(t.MustText())
		if text != "" {
			return text
		}
	}
	return ""
}

func findDevtoStats(el *rod.Element) (*int, *int, string, *int) {
	container := el
	if analytics, err := el.Element("[data-analytics-pageviews]"); err == nil {
		container = analytics
	}

	var reactions *int
	if span, err := container.Element(`span[title="Reactions"]`); err == nil {
		text := normalizeWhitespace(span.MustText())
		reactions = parseCount(text)
	}

	var comments *int
	if span, err := container.Element(`span[title="Comments"]`); err == nil {
		if countEl, err := span.Element(".spec__comments-count"); err == nil {
			text := normalizeWhitespace(countEl.MustText())
			comments = parseCount(text)
		} else {
			text := normalizeWhitespace(span.MustText())
			comments = parseCount(text)
		}
	}

	var viewsLabel string
	var viewsCount *int
	if span, err := container.Element(`span[title="Views"]`); err == nil {
		text := normalizeWhitespace(span.MustText())
		viewsLabel = text
		if strings.Contains(text, "<") {
			viewsCount = nil
		} else {
			viewsCount = parseCount(text)
		}
	}

	return reactions, comments, viewsLabel, viewsCount
}

func parseCount(text string) *int {
	clean := strings.ReplaceAll(text, ",", "")
	match := countRegex.FindString(clean)
	if match == "" {
		return nil
	}
	val, err := strconv.Atoi(match)
	if err != nil {
		return nil
	}
	return &val
}

func normalizeWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func normalizeDevtoURL(href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	if strings.HasPrefix(href, "/") {
		return "https://dev.to" + href
	}
	return "https://dev.to/" + href
}

func isLikelyDevtoPostURL(href string) bool {
	lower := strings.ToLower(strings.TrimSpace(href))
	if lower == "" {
		return false
	}
	if strings.Contains(lower, "dashboard") ||
		strings.Contains(lower, "settings") ||
		strings.Contains(lower, "notifications") ||
		strings.Contains(lower, "signout") ||
		strings.Contains(lower, "enter") {
		return false
	}

	trimmed := strings.TrimPrefix(lower, "https://dev.to")
	trimmed = strings.TrimPrefix(trimmed, "http://dev.to")
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return false
	}

	// Expect at least username/slug
	return strings.Count(trimmed, "/") >= 1
}

var countRegex = regexp.MustCompile(`\d+`)
