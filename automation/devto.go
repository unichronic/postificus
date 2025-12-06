package automation

import (
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

var Browser *rod.Browser

// InitBrowser launches a global browser instance.
func InitBrowser() error {
	path, _ := launcher.LookPath()
	// Check if running in production (Render sets PORT)
	isProduction := os.Getenv("PORT") != ""

	u := launcher.New().
		Bin(path).
		NoSandbox(true).
		Headless(isProduction). // Headless in prod, visible locally
		Set("disable-gpu").
		Set("disable-dev-shm-usage").
		MustLaunch()

	Browser = rod.New().ControlURL(u).MustConnect()
	return nil
}

// CloseBrowser closes the global browser instance.
func CloseBrowser() {
	if Browser != nil {
		Browser.MustClose()
	}
}

// PostToDevToWithCookie bypasses login by injecting a valid session token.
func PostToDevToWithCookie(sessionToken, title, content string) error {
	if Browser == nil {
		return fmt.Errorf("browser not initialized")
	}

	// 1. USE THE SHARED BROWSER
	// Instead of launching a new browser, we just open a new TAB (Page)
	// This is instant and uses very little memory.
	page := Browser.MustPage("https://dev.to/404")

	// CRITICAL: Ensure we close THIS TAB when done, not the whole browser
	defer page.MustClose()

	// 2. Cookie Injection
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
		return fmt.Errorf("failed to set cookies: %w", err)
	}

	// 3. Navigate
	fmt.Println("Navigating to editor...")
	page.MustNavigate("https://dev.to/new")
	page.MustWaitLoad()

	if page.MustInfo().URL == "https://dev.to/enter" {
		return fmt.Errorf("COOKIE EXPIRED: Session token invalid")
	}

	fmt.Println("Session valid! Writing post...")

	// 4. Fill Content
	page.MustWaitStable()
	page.MustElement("textarea[placeholder='New post title here...']").MustInput(title)
	page.MustElement("#article_body_markdown").MustInput(content)

	// 5. Handle Tags (THE FIX)
	// DO NOT PRESS ENTER. It submits the form prematurely.
	fmt.Println("Adding tags...")

	// Find input (robust race with wait)
	var tagInput *rod.Element
	err := rod.Try(func() {
		tagInput = page.Race().
			Element("#article_tags").MustHandle(func(e *rod.Element) {
			tagInput = e
		}).
			Element("input[placeholder*='tags']").MustHandle(func(e *rod.Element) {
			tagInput = e
		}).
			MustDo()
	})

	if err != nil {
		fmt.Println("❌ Failed to find tag input. Capturing screenshot...")
		page.MustScreenshot("debug_tags_missing.png")
		return fmt.Errorf("could not find tag input: %w", err)
	}

	// Input text and BLUR (click away) to register the tag
	tagInput.MustInput("automation").MustBlur()

	// Wait for the UI to settle after blurring (Dev.to does JS processing here)
	page.MustWaitStable()

	// 6. Save Draft
	fmt.Println("Saving draft...")
	saveBtn := page.MustElementR("button", "Save draft")
	saveBtn.MustWaitVisible().MustClick()

	// 7. Robust Wait Loop
	fmt.Println("Waiting for save to complete...")

	// We use a shorter polling interval for responsiveness
	// We wrap page.Info() to catch the "Target Closed" error if it happens
	for i := 0; i < 60; i++ { // 30 seconds max
		// Check URL
		info, err := page.Info()
		if err != nil {
			// If page crashed here, we catch it
			return fmt.Errorf("browser tab crashed during save: %w", err)
		}

		if info.URL != "https://dev.to/new" {
			fmt.Println("✅ Success: URL redirected to", info.URL)
			return nil
		}

		// Check Button Text
		if has, _, _ := page.HasR("button", "Saved"); has {
			fmt.Println("✅ Success: Button text changed to 'Saved'")
			return nil
		}

		// Check Error
		if has, el, _ := page.Has(".crayons-toast--error"); has {
			return fmt.Errorf("❌ Dev.to Error: %s", el.MustText())
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for save confirmation")
}
