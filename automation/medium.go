package automation

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

// Helper: Wait for "Saved" status
func waitForSaved(page *rod.Page) error {
	fmt.Print("   ⏳ Syncing: ")
	// Poll for 15 seconds
	for i := 0; i < 15; i++ {
		// Check for "Saved"
		if has, _, _ := page.HasR("span, div, p", "Saved"); has {
			fmt.Println("✅ Server Confirmed: Saved")
			return nil
		}
		// Check for "Saving..."
		if has, _, _ := page.HasR("span, div, p", "Saving..."); has {
			fmt.Print(".") // loading indicator
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("❌ Timeout")
	return fmt.Errorf("sync timeout: stuck on Saving...")
}

func PostToMedium(uid, sid, xsrf, title, content string) error {
	if Browser == nil {
		return fmt.Errorf("browser not initialized")
	}

	// CRITICAL: Initialize Stealth Page
	// stealth.MustPage(Browser) creates a page with stealth scripts pre-loaded
	page := stealth.MustPage(Browser)
	defer page.MustClose()

	// Navigate to domain to ensure cookies are set correctly
	page.MustNavigate("https://medium.com/")

	// 1. Inject Cookies
	cookies := []*proto.NetworkCookieParam{
		{Name: "uid", Value: uid, Domain: ".medium.com", Path: "/", HTTPOnly: true, Secure: true},
		{Name: "sid", Value: sid, Domain: ".medium.com", Path: "/", HTTPOnly: true, Secure: true},
		{Name: "xsrf", Value: xsrf, Domain: ".medium.com", Path: "/", HTTPOnly: true, Secure: true},
	}
	page.SetCookies(cookies)

	// 2. Navigate
	fmt.Println("Navigating to Medium Editor...")
	page.MustNavigate("https://medium.com/new-story")

	// 3. Wait for Editor (Robust Selector)
	fmt.Println("Waiting for editor...")
	// We wait for the Title field. If this times out, your Cookies/UA are still blocked.
	err := page.Timeout(15*time.Second).WaitElementsMoreThan(`[data-testid="editorTitleParagraph"]`, 0)
	if err != nil {
		page.MustScreenshot("medium_blocked.png")
		return fmt.Errorf("editor load failed: see medium_blocked.png")
	}

	// ---------------------------------------------------------
	// PART A: Title (The "Human Handshake")
	// ---------------------------------------------------------
	fmt.Println("Writing Title...")
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

	// ---------------------------------------------------------
	// PART B: Body (Hybrid Approach)
	// ---------------------------------------------------------
	fmt.Println("Writing Body...")
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
		page.MustInsertText(content)
	} else {
		time.Sleep(300 * time.Millisecond)
		// CRITICAL FIX: Use explicit Press/Type/Release for Control+V
		// Keyboard methods return error, but we ignore for brevity in this snippet
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
	fmt.Println("Initiating Publish...")

	// Click Top Publish
	page.MustElementR("button", "Publish").MustClick()
	page.MustWaitStable()

	// Click Final Publish
	// Sometimes Medium shows "Topics" selection here.
	// If you want to skip topics, just hit Publish Now.
	pubBtn := page.MustElementR("button", "Publish now")
	pubBtn.MustWaitVisible().MustClick()

	// ---------------------------------------------------------
	// PART D: Verification
	// ---------------------------------------------------------
	fmt.Println("Verifying...")
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

	fmt.Println("✅ Success! URL:", page.MustInfo().URL)
	return nil
}
