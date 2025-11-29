package wizard

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func RunAuthFlow(startURL string, successCriteria string, cookieName string) (string, error) {
	u := launcher.New().Headless(false).NoSandbox(true).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	fmt.Printf("Browser opened. Navigating to %s\n", startURL)
	page := browser.MustPage(startURL)

	fmt.Println("Please log in manually in the opened browser window.")
	fmt.Printf("Waiting for URL to match criteria: %s\n", successCriteria)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-timeout:
			return "", errors.New("timed out waiting for login")
		case <-ticker.C:
			info, err := page.Info()
			if err != nil {
				return "", fmt.Errorf("error getting page info (did you close the browser?): %w", err)
			}

			if strings.Contains(info.URL, successCriteria) {
				fmt.Println("Login detected! Extracting cookies...")

				cookies, err := page.Cookies([]string{cookieName})
				if err != nil {
					return "", fmt.Errorf("failed to get cookies: %w", err)
				}

				for _, cookie := range cookies {
					if cookie.Name == cookieName {
						return cookie.Value, nil
					}
				}

				return "", fmt.Errorf("cookie '%s' not found after login", cookieName)
			}
		}
	}
}
