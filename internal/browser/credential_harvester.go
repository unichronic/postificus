package browser

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// WaitForMediumLogin launches a visible browser window and waits for the user to log in.
// It returns the uid, sid, xsrf, and username/name once detected.
func WaitForMediumLogin() (string, string, string, string, error) {
	log.Println("🚀 Launching validation browser for Medium login...")

	u := GetHardenedLauncher(false, false).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://medium.com/m/signin")

	log.Println("⏳ Waiting for user to log in...")
	fmt.Println("👉 Please log in to Medium in the opened window.")

	timeout := time.After(3 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return "", "", "", "", fmt.Errorf("login timeout: user took too long")
		case <-ticker.C:
			// 1. Check if browser is alive
			info, err := page.Info()
			if err != nil {
				return "", "", "", "", fmt.Errorf("browser closed by user")
			}

			// 2. Wait until off login page
			if info.URL == "https://medium.com/m/signin" || info.URL == "https://medium.com/ap/signin" {
				continue
			}

			// 3. Check cookies
			cookies, err := page.Browser().GetCookies()
			if err != nil {
				continue
			}

			var uid, sid, xsrf string
			for _, c := range cookies {
				switch c.Name {
				case "uid":
					uid = c.Value
				case "sid":
					sid = c.Value
				case "xsrf":
					xsrf = c.Value
				}
			}

			if uid != "" && sid != "" && xsrf != "" {
				log.Println("✅ Login detected! Fetching account info...")

				// Scraping Logic with Timeout
				username := "Medium User"

				// We fork a cleanup context for the scraping actions to ensure we don't hang
				err := rod.Try(func() {
					// Navigate to /me which redirects to /@username
					page.Timeout(10 * time.Second).MustNavigate("https://medium.com/me")
					page.Timeout(10 * time.Second).MustWaitLoad()

					// 1. Get Handle from URL
					// URL will be like https://medium.com/@myhandle
					// or https://medium.com/@myhandle/
					finalURL := page.MustInfo().URL
					if len(finalURL) > 20 && finalURL[len(finalURL)-1] == '/' {
						finalURL = finalURL[:len(finalURL)-1]
					}

					// Extract handle
					// Simple implementation: split by '@'
					parts := regexp.MustCompile(`@([\w\.\-]+)`).FindStringSubmatch(finalURL)
					handle := ""
					if len(parts) > 1 {
						handle = parts[1]
					}

					// 2. Get Name from Title
					// Title is usually "My Name – Medium"
					title := page.MustInfo().Title
					name := strings.TrimSuffix(title, " – Medium")

					if name != "" && handle != "" {
						username = fmt.Sprintf("%s (@%s)", name, handle)
					} else if handle != "" {
						username = "@" + handle
					} else if name != "" {
						username = name
					}
				})

				if err != nil {
					log.Printf("⚠️ Failed to scrape Medium username: %v", err)
				}

				return uid, sid, xsrf, username, nil
			}
		}
	}
}

// WaitForDevToLogin launches a visible browser window and waits for the user to log in to Dev.to.
// It returns the remember_user_token cookie and username once detected.
func WaitForDevToLogin() (string, string, error) {
	log.Println("🚀 Launching validation browser for Dev.to login...")

	u := GetHardenedLauncher(false, false).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://dev.to/enter")

	log.Println("⏳ Waiting for user to log in...")
	fmt.Println("👉 Please log in to Dev.to in the opened window.")

	timeout := time.After(3 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return "", "", fmt.Errorf("login timeout: user took too long")
		case <-ticker.C:
			// 1. Check browser
			_, err := page.Info()
			if err != nil {
				return "", "", fmt.Errorf("browser closed by user")
			}

			// 2. Check cookie
			cookies, err := page.Browser().GetCookies()
			if err != nil {
				continue
			}

			for _, c := range cookies {
				if c.Name == "remember_user_token" && c.Value != "" {
					log.Println("✅ Dev.to Login detected! Fetching account info...")

					username := "Dev.to User"

					// Scraping with Timeout
					err := rod.Try(func() {
						page.Timeout(10 * time.Second).MustNavigate("https://dev.to/settings")
						page.Timeout(10 * time.Second).MustWaitLoad()

						// Try explicit inputs
						// user_name is the "Name" field
						// user_username is the "Username" field

						var name, handle string

						// Use Race to find elements or timeout
						page.Race().
							Element("input[name='user[name]']").MustHandle(func(e *rod.Element) {
							name = e.MustProperty("value").String()
						}).
							Element("input[name='user[username]']").MustHandle(func(e *rod.Element) {
							handle = e.MustProperty("value").String()
						}).
							MustDo()

						if name != "" && handle != "" {
							username = fmt.Sprintf("%s (@%s)", name, handle)
						} else if name != "" {
							username = name
						} else if handle != "" {
							username = "@" + handle
						}
					})

					if err != nil {
						log.Printf("⚠️ Failed to scrape Dev.to username: %v", err)
					}

					return c.Value, username, nil
				}
			}
		}
	}
}

// GetCookies is a helper specifically to find finding cookies by name using rod's proto
func getCookies(browser *rod.Browser) ([]*proto.NetworkCookie, error) {
	return browser.GetCookies()
}
