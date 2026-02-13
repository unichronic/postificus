package browser

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var Browser *rod.Browser
var browserOnce sync.Once
var browserInitErr error

// GetHardenedLauncher returns a launcher with anti-detection flags pre-configured
func GetHardenedLauncher(headless bool, useFirefox bool) *launcher.Launcher {
	l := launcher.New()

	if useFirefox {
		// FIREFOX CONFIGURATION
		l.Bin("/usr/bin/firefox")
		// Firefox requires specific flags for remote debugging, but Rod handles basic setup.
		// We set user agent and window size.
		// Crucial for Firefox: Use a unique temporary profile to avoid "Already Running" errors
		// Rod's default 'user-data-dir' is for Chrome; Firefox needs '-profile'
		// SNAP FIX: Snap Firefox cannot access /tmp, so we must use a local dir in HOME.
		homeDir, _ := os.UserHomeDir()
		profileBase := filepath.Join(homeDir, ".postificus_profiles")
		_ = os.MkdirAll(profileBase, 0755)

		profileDir, _ := os.MkdirTemp(profileBase, "ff-profile-")
		l.Set("profile", profileDir)

		l.Set("no-remote") // Essential for running parallel instances
		// l.Set("new-instance")
		if headless {
			l.Set("MOZ_HEADLESS", "1")
		}

		// Use fixed port to avoid stdout parsing issues with Snap
		l.RemoteDebuggingPort(9222)

		// Note: 'stealth' lib and many Chrome flags don't apply to Firefox.
		// We rely on Firefox's inherent differences to bypass Chrome-targeted bot detection.
	} else {
		// CHROME CONFIGURATION
		path, _ := launcher.LookPath()
		if bin := os.Getenv("BROWSER_BIN"); bin != "" {
			path = bin
		}
		l.Bin(path)
		l.NoSandbox(true)
		l.Headless(headless)

		// CRITICAL ANTI-DETECTION FLAGS (Chrome Only)
		l.Set("window-size", "1920,1080")
		l.Set("disable-infobars")
		l.Set("disable-features", "ClipboardContentSetting")
		l.Set("disable-blink-features", "AutomationControlled") // Removing navigator.webdriver
		l.Set("excludeSwitches", "enable-automation")           // Hides info bar
		l.Set("disable-gpu")
		l.Set("no-sandbox")
	}

	// Shared Config
	l.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0")

	return l
}

// InitBrowser launches a global browser instance.
func InitBrowser() error {
	// Check if running in production (Render sets PORT, Docker sets APP_ENV)
	isProduction := os.Getenv("PORT") != "" || os.Getenv("APP_ENV") == "production"

	headless := isProduction
	// Allow manual override via env var
	if val := os.Getenv("BROWSER_HEADLESS"); val != "" {
		headless = (val == "true" || val == "1")
	}

	// 1. Create a "Clean" Launcher using shared logic (Default to Chrome for background tasks)
	l := GetHardenedLauncher(headless, false)

	// 3. Launch
	u := l.MustLaunch()

	Browser = rod.New().ControlURL(u).MustConnect()
	return nil
}

// CloseBrowser closes the global browser instance.
func CloseBrowser() {
	if Browser != nil {
		Browser.MustClose()
	}
}

// EnsureBrowser initializes the shared browser exactly once.
func EnsureBrowser() error {
	browserOnce.Do(func() {
		browserInitErr = InitBrowser()
	})
	return browserInitErr
}
