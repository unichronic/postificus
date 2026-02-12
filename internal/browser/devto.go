package browser

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Browser logic moved to browser.go

// PostToDevToWithCookie bypasses login by injecting a valid session token.
func PostToDevToWithCookie(sessionToken, title, content, coverImage string, tags []string) error {
	log.Println("Starting PostToDevToWithCookie...")
	log.Println("Starting PostToDevToWithCookie...")
	if err := EnsureBrowser(); err != nil {
		return fmt.Errorf("browser init failed: %w", err)
	}

	// 1. USE THE SHARED BROWSER
	// Instead of launching a new browser, we just open a new TAB (Page)
	// This is instant and uses very little memory.
	log.Println("Opening new page...")
	page := Browser.MustPage("https://dev.to/404")

	// CRITICAL: Ensure we close THIS TAB when done, not the whole browser
	defer func() {
		log.Println("Closing page...")
		page.MustClose()
	}()

	// 2. Cookie Injection
	log.Println("Injecting cookies...")
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
	log.Println("Navigating to editor (https://dev.to/new)...")
	page.MustNavigate("https://dev.to/new")

	log.Println("Waiting for load...")
	page.MustWaitLoad()

	if page.MustInfo().URL == "https://dev.to/enter" {
		return fmt.Errorf("COOKIE EXPIRED: Session token invalid")
	}

	log.Println("Session valid! Writing post...")

	// 4. Fill Content
	log.Println("Waiting for stable...")
	page.MustWaitStable()

	log.Println("Inputting title...")
	page.MustElement("textarea[placeholder='New post title here...']").MustInput(title)

	if coverImage != "" {
		if err := uploadDevtoCoverImage(page, coverImage); err != nil {
			log.Printf("⚠️ Cover image upload failed: %v", err)
		}
	}

	log.Println("Inputting content...")
	page.MustElement("#article_body_markdown").MustInput(content)

	// 5. Handle Tags (THE FIX)
	// DO NOT PRESS ENTER. It submits the form prematurely.
	log.Println("Adding tags...")

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
		log.Println("❌ Failed to find tag input. Capturing screenshot...")
		page.MustScreenshot("debug_tags_missing.png")
		return fmt.Errorf("could not find tag input: %w", err)
	}

	// Input tags
	if len(tags) > 0 {
		for _, tag := range tags {
			log.Printf("Adding tag: %s", tag)
			tagInput.MustInput(tag).MustBlur()
			// Wait a bit for the tag to be processed
			time.Sleep(200 * time.Millisecond)
			// Re-focus for next tag if needed, but usually inputting appends
			tagInput.MustClick()
		}
	} else {
		// Default tag if none provided
		log.Println("Adding default tag: automation")
		tagInput.MustInput("automation").MustBlur()
	}

	// Wait for the UI to settle after blurring (Dev.to does JS processing here)
	page.MustWaitStable()

	// 6. Publish
	log.Println("Publishing...")
	// User provided selector: <button type="button" class="c-btn c-btn--primary mr-2 whitespace-nowrap">Publish</button>
	// We target the primary button with text "Publish"
	saveBtn := page.MustElementR("button.c-btn.c-btn--primary", "Publish")
	saveBtn.MustWaitVisible().MustClick()

	// 7. Robust Wait Loop
	log.Println("Waiting for save to complete...")

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
			log.Println("✅ Success: URL redirected to", info.URL)
			return nil
		}

		// Check Button Text
		if has, _, _ := page.HasR("button", "Saved"); has {
			log.Println("✅ Success: Button text changed to 'Saved'")
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

func uploadDevtoCoverImage(page *rod.Page, coverImage string) error {
	path, cleanup, err := PrepareImageUpload(coverImage)
	if err != nil {
		return err
	}
	defer cleanup()

	if path == "" {
		return nil
	}

	input, err := page.Timeout(5 * time.Second).Element(`#cover-image-input`)
	if err != nil {
		input, err = page.Timeout(5 * time.Second).Element(`[data-testid="cover-image-input"]`)
		if err != nil {
			return fmt.Errorf("cover image input not found")
		}
	}

	if err := input.SetFiles([]string{path}); err != nil {
		return fmt.Errorf("failed to set cover image: %w", err)
	}

	time.Sleep(700 * time.Millisecond)
	return nil
}
