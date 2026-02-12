package browser

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

// Helper: Wait for "Saved" status
func waitForSaved(page *rod.Page) error {
	log.Print("   ⏳ Syncing: ")
	// Poll for 15 seconds
	for i := 0; i < 15; i++ {
		// Check for "Saved"
		if has, _, _ := page.HasR("span, div, p", "Saved"); has {
			log.Println("✅ Server Confirmed: Saved")
			return nil
		}
		// Check for "Saving..."
		if has, _, _ := page.HasR("span, div, p", "Saving..."); has {
			log.Print(".") // loading indicator
		}
		time.Sleep(1 * time.Second)
	}
	log.Println("❌ Timeout")
	return fmt.Errorf("sync timeout: stuck on Saving...")
}

// PostToMedium publishes to Medium using API first, falling back to browser automation
func PostToMedium(uid, sid, xsrf, title, content string) error {
	return PostToMediumWithTags(uid, sid, xsrf, title, content, []string{}, "")
}

// PostToMediumWithTags publishes to Medium with tags support
func PostToMediumWithTags(uid, sid, xsrf, title, content string, tags []string, coverImage string) error {
	if coverImage != "" {
		log.Println("🖼️ Cover image provided, skipping API and using browser automation...")
		return postToMediumBrowser(uid, sid, xsrf, title, content, tags, coverImage)
	}

	log.Println("🎯 Attempting Medium API publish...")

	// Try API-based approach first
	client := NewMediumAPIClient(uid, sid, xsrf)
	url, err := client.Publish(title, content, tags)

	if err == nil {
		log.Printf("✅ Published via API: %s", url)
		return nil
	}

	// API failed, fall back to browser automation
	log.Printf("⚠️ API failed (%v), falling back to browser automation...", err)
	return postToMediumBrowser(uid, sid, xsrf, title, content, tags, coverImage)
}

// postToMediumBrowser is the original browser automation implementation (preserved as fallback)
func postToMediumBrowser(uid, sid, xsrf, title, content string, tags []string, coverImage string) error {
	log.Println("Starting PostToMedium...")
	log.Println("Starting PostToMedium...")
	if err := EnsureBrowser(); err != nil {
		return fmt.Errorf("browser init failed: %w", err)
	}

	// CRITICAL: Initialize Stealth Page
	// stealth.MustPage(Browser) creates a page with stealth scripts pre-loaded
	log.Println("Creating stealth page...")
	page := stealth.MustPage(Browser)
	defer func() {
		log.Println("Closing page...")
		page.MustClose()
	}()

	// Navigate to domain to ensure cookies are set correctly
	log.Println("Navigating to medium.com to set cookies...")
	page.MustNavigate("https://medium.com/")

	// 1. Inject Cookies
	log.Println("Injecting cookies...")
	cookies := []*proto.NetworkCookieParam{
		{Name: "uid", Value: uid, Domain: ".medium.com", Path: "/", HTTPOnly: true, Secure: true},
		{Name: "sid", Value: sid, Domain: ".medium.com", Path: "/", HTTPOnly: true, Secure: true},
		{Name: "xsrf", Value: xsrf, Domain: ".medium.com", Path: "/", HTTPOnly: true, Secure: true},
	}
	page.SetCookies(cookies)

	// 2. Navigate
	log.Println("Navigating to Medium Editor (https://medium.com/new-story)...")
	page.MustNavigate("https://medium.com/new-story")

	// 3. Wait for Editor (Robust Selector)
	log.Println("Waiting for editor...")
	// We wait for the Title field. If this times out, your Cookies/UA are still blocked.
	err := page.Timeout(15*time.Second).WaitElementsMoreThan(`[data-testid="editorTitleParagraph"]`, 0)
	if err != nil {
		log.Println("❌ Editor load failed. Capturing screenshot...")
		page.MustScreenshot("medium_blocked.png")
		return fmt.Errorf("editor load failed: see medium_blocked.png")
	}

	// ---------------------------------------------------------
	// PART A: Title (The "Human Handshake")
	// ---------------------------------------------------------
	log.Println("Writing Title...")
	titleElem := page.MustElement(`[data-testid="editorTitleParagraph"]`)
	titleElem.MustClick() // Focus

	// Use Human Typing for Title
	// This generates "clean" traffic that looks real to the WAF
	HumanType(page, title)

	// SYNC CHECK 1:
	// If the title doesn't save, do not proceed. The WAF has already blocked you.
	if err := waitForSaved(page); err != nil {
		return fmt.Errorf("WAF BLOCK: Title was typed, but server rejected save")
	}

	page.Keyboard.MustType(input.Enter) // New line

	if coverImage != "" {
		if err := insertMediumInlineImage(page, coverImage); err != nil {
			log.Printf("⚠️ Cover image insert failed: %v", err)
		} else {
			time.Sleep(300 * time.Millisecond)
			page.Keyboard.MustType(input.Enter)
		}
	}

	// ---------------------------------------------------------
	// PART B: Body (Hybrid Approach)
	// ---------------------------------------------------------
	log.Println("Writing Body...")
	page.MustWaitStable()

	// For the body, typing 1000 words humanly takes too long (and might timeout execution).
	// We use the clipboard method.

	// Execute Paste
	_, err = page.Evaluate(rod.Eval(`(text) => {
        navigator.clipboard.writeText(text);
        return text;
    }`, content))

	if err != nil {
		// Fallback
		log.Println("Clipboard paste failed, using InsertText...")
		page.MustInsertText(content)
	} else {
		time.Sleep(300 * time.Millisecond)
		// CRITICAL FIX: Use explicit Press/Type/Release for Control+V
		// Keyboard methods return error, but we ignore for brevity in this snippet
		log.Println("Pasting content (Ctrl+V)...")
		page.Keyboard.Press(input.ControlLeft)
		page.Keyboard.Type(input.Key('v'))
		page.Keyboard.Release(input.ControlLeft)
	}

	// SYNC CHECK 2 (CRITICAL):
	// Large paste = Longer save time.
	// We do NOT click publish until this returns.
	if err := waitForSaved(page); err != nil {
		return fmt.Errorf("WAF BLOCK: Body content rejected")
	}

	// ---------------------------------------------------------
	// PART C: Publish
	// ---------------------------------------------------------
	// ---------------------------------------------------------
	// PART C: Publish (Pre-Publish Overlay)
	// ---------------------------------------------------------
	log.Println("Initiating Publish...")

	// Click Top Publish (usually green button)
	// User provided selector: [data-action="show-prepublish"]
	log.Println("Clicking header 'Publish' button...")

	// Try user-provided selector first, then testid, then text
	err = rod.Try(func() {
		page.MustElement(`[data-action="show-prepublish"]`).MustClick()
	})
	if err != nil {
		err = rod.Try(func() {
			page.MustElement(`[data-testid="header-publish-button"]`).MustClick()
		})
		if err != nil {
			page.MustElementR("button", "Publish").MustClick()
		}
	}
	page.MustWaitStable()

	// Wait for Overlay
	log.Println("Waiting for publish overlay...")
	// Wait specifically for the overlay dialog to be visible
	if err := page.Timeout(10 * time.Second).MustElement(`.overlay-dialog`).WaitVisible(); err != nil {
		log.Printf("⚠️ Overlay wait warning: %v", err)
	}

	// Give animation a moment to settle
	time.Sleep(1 * time.Second)

	// ---------------------------------------------------------
	// PART C.1: Tags (Topics)
	// ---------------------------------------------------------
	if len(tags) > 0 {
		log.Println("Adding tags (topics)...")

		// Use Try to safely check/wait for the input
		err := rod.Try(func() {
			// Wait for input to be interactive
			tagInput := page.Timeout(5 * time.Second).MustElement(`[data-testid="publishTopicsInput"]`)
			tagInput.MustWaitVisible()

			for _, tag := range tags {
				log.Printf("Adding tag: %s", tag)
				// Click to focus
				tagInput.MustClick()
				time.Sleep(200 * time.Millisecond)

				// Type tag
				tagInput.MustInput(tag)
				time.Sleep(500 * time.Millisecond) // Wait for dropdown suggestions

				// Press Enter to confirm tag
				page.Keyboard.MustType(input.Enter)
				time.Sleep(300 * time.Millisecond)
			}
		})

		if err != nil {
			log.Printf("⚠️ Failed to add tags: %v", err)
			// Don't error out, try to publish anyway
		}
	}

	// Click Final Publish in the modal
	// Usually "Publish now"
	log.Println("Clicking final 'Publish now' button...")

	// Wait for the modal/panel to appear
	err = rod.Try(func() {
		// User provided specific selector: [data-testid="publishConfirmButton"]
		page.Timeout(5 * time.Second).MustElement(`[data-testid="publishConfirmButton"]`).MustClick()
	})
	if err != nil {
		// Fallback to text matching
		pubBtn := page.MustElementR("button", "Publish now")
		pubBtn.MustWaitVisible().MustClick()
	}

	// ---------------------------------------------------------
	// PART D: Verification
	// ---------------------------------------------------------
	log.Println("Verifying redirect...")
	err = rod.Try(func() {
		for i := 0; i < 60; i++ {
			if page.MustInfo().URL != "https://medium.com/new-story" {
				return
			}
			time.Sleep(time.Second)
		}
		panic("timeout")
	})

	if err != nil {
		return fmt.Errorf("redirect timeout")
	}

	log.Println("✅ Success! URL:", page.MustInfo().URL)
	return nil
}

func insertMediumInlineImage(page *rod.Page, coverImage string) error {
	path, cleanup, err := PrepareImageUpload(coverImage)
	if err != nil {
		return err
	}
	defer cleanup()

	if path == "" {
		return nil
	}

	// Open inline menu (+)
	err = rod.Try(func() {
		page.MustElement(`[data-testid="editorAddButton"]`).MustClick()
	})
	if err != nil {
		err = rod.Try(func() {
			page.MustElement(`button[data-action="inline-menu"]`).MustClick()
		})
	}
	if err != nil {
		return fmt.Errorf("inline menu button not found")
	}

	// Click image option
	err = rod.Try(func() {
		page.MustElement(`button[data-action="inline-menu-image"]`).MustClick()
	})
	if err != nil {
		err = rod.Try(func() {
			page.MustElementR("button", "Add an image").MustClick()
		})
	}
	if err != nil {
		return fmt.Errorf("inline image button not found")
	}

	// Find file input and upload
	inputEl, err := page.Timeout(5 * time.Second).Element(`input[type="file"]`)
	if err != nil {
		return fmt.Errorf("image file input not found")
	}
	if err := inputEl.SetFiles([]string{path}); err != nil {
		return fmt.Errorf("failed to set image file: %w", err)
	}

	// Allow upload to complete
	time.Sleep(2 * time.Second)
	return nil
}
