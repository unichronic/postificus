package automation

import (
	"fmt"
	"time"

	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func PostToLinkedIn(li_at, postContent, blogURL string) error {
	if Browser == nil {
		return fmt.Errorf("browser not initialized")
	}

	// CRITICAL: Initialize Stealth Page
	page := stealth.MustPage(Browser)
	defer page.MustClose()

	// 1. Inject Cookie
	// LinkedIn only needs "li_at" for the session usually.
	cookies := []*proto.NetworkCookieParam{
		{
			Name:     "li_at",
			Value:    li_at,
			Domain:   ".linkedin.com",
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
		},
	}
	page.SetCookies(cookies)

	// 2. Navigate
	fmt.Println("Navigating to LinkedIn Feed...")
	if err := page.Navigate("https://www.linkedin.com/feed/"); err != nil {
		return fmt.Errorf("navigation failed: %v", err)
	}

	// Wait for Feed Load - Random wait to simulate reading
	fmt.Println("Waiting for feed...")
	RandomSleep(4000, 2000)

	// Check if blocked
	if page.MustInfo().URL != "https://www.linkedin.com/feed/" {
		page.MustScreenshot("linkedin_block.png")
		return fmt.Errorf("login redirected/blocked")
	}

	// 3. Open Modal
	fmt.Println("Opening Post Modal...")
	// Use the generic text search for the button, it's robust enough
	// But use human click
	err := HumanMoveAndClick(page, "button.share-box-feed-entry__trigger")
	if err != nil {
		// Fallback to text search if class fails
		btn, err := page.ElementR("span", "Start a post")
		if err == nil {
			btn.MustClick()
		} else {
			// Try one more selector
			btn = page.MustElement(".share-box-feed-entry__trigger")
			btn.MustClick()
		}
	}

	// Wait for Modal Animation
	RandomSleep(2000, 1000)

	// 4. Locate Editor (Using your Robust Selector)
	fmt.Println("Locating Editor...")

	// The Selector from your screenshot:
	const editorSelector = `div[data-test-ql-editor-contenteditable="true"]`

	// Wait for it to exist
	if err := page.Timeout(10*time.Second).WaitElementsMoreThan(editorSelector, 0); err != nil {
		page.MustScreenshot("modal_fail.png")
		return fmt.Errorf("editor not found (modal didn't load)")
	}

	// Click it humanly
	if err := HumanMoveAndClick(page, editorSelector); err != nil {
		return err
	}

	// 5. Type Content
	fmt.Println("Typing...")
	HumanType(page, postContent)

	RandomSleep(1000, 500)
	page.Keyboard.MustType(input.Enter)
	page.Keyboard.MustType(input.Enter)

	// 6. Paste Link
	fmt.Println("Pasting Link...")
	page.MustInsertText(blogURL)

	// Wait for OG Preview
	fmt.Println("Waiting for Preview...")
	time.Sleep(8 * time.Second) // Fixed wait is okay here, looks like loading time

	// 7. Post
	fmt.Println("Posting...")
	// Specific selector for the Post button inside the share actions bar
	postBtnSelector := "div.share-box_actions button.share-actions__primary-action"

	// Ensure it's enabled
	btn, err := page.Element(postBtnSelector)
	if err != nil {
		// Fallback to text
		btn = page.MustElementR("button", "Post")
	}

	disabled, _ := btn.Disabled()
	if disabled {
		return fmt.Errorf("post button disabled")
	}

	btn.MustClick()

	// 8. Verify
	fmt.Println("Verifying...")
	RandomSleep(2000, 1000)

	// Check if modal is gone
	if has, _, _ := page.Has(".share-box-modal"); !has {
		fmt.Println("✅ LinkedIn Post Published!")
		return nil
	}

	return fmt.Errorf("modal stuck open")
}
