package automation

import (
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var Browser *rod.Browser

// InitBrowser launches a global browser instance.
func InitBrowser() error {
	path, _ := launcher.LookPath()
	// Check if running in production (Render sets PORT)
	isProduction := os.Getenv("PORT") != ""

	// 1. Create a "Clean" Launcher
	l := launcher.New().
		Bin(path).
		NoSandbox(true).
		Headless(isProduction).
		// CRITICAL: Set a real window size. WAFs flag 800x600 (default headless) immediately.
		Set("window-size", "1920,1080").
		// CRITICAL: Disable the "Chrome is being controlled by automated software" bar
		Set("disable-infobars").
		// CRITICAL: Bypass the clipboard permission prompt (headless usually blocks this)
		Set("disable-features", "ClipboardContentSetting").
		// CRITICAL: MASKING
		// This flag actively removes the "navigator.webdriver" property from the C++ level
		Set("disable-blink-features", "AutomationControlled").
		// Match your real UA exactly
		Set("disable-gpu").
		Set("no-sandbox").
		Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 2. Remove the "Enable Automation" switch that Rod adds by default
	// This is the "God Mode" fix for detection.
	l.Set("excludeSwitches", "enable-automation")

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
