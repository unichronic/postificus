package browser

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

var (
	mouseLogFile  *os.File
	mouseLogMutex sync.Mutex
	mouseLogOnce  sync.Once
)

func initMouseLog() {
	if os.Getenv("DEBUG_MOUSE") == "true" {
		f, err := os.OpenFile("mouse_trace.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			mouseLogFile = f
			// Write header if new file
			stat, _ := f.Stat()
			if stat.Size() == 0 {
				f.WriteString("timestamp,x,y,event\n")
			}
		}
	}
}

func LogMouse(x, y float64, event string) {
	mouseLogOnce.Do(initMouseLog)
	if mouseLogFile == nil {
		return
	}
	mouseLogMutex.Lock()
	defer mouseLogMutex.Unlock()
	timestamp := time.Now().UnixMilli()
	fmt.Fprintf(mouseLogFile, "%d,%.2f,%.2f,%s\n", timestamp, x, y, event)
}

// RandomSleep sleeps for a variable duration (base + jitter)
func RandomSleep(baseMs int, jitterMs int) {
	ms := baseMs + rand.Intn(jitterMs)
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// HumanMoveAndClick simulates a user moving the mouse to an element and clicking
// Bots teleport; Humans scroll, hover, pause, then click.
func HumanMoveAndClick(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}

	// 1. Scroll the element into view (Humans can't click what they can't see)
	el.ScrollIntoView()
	RandomSleep(200, 300)

	// 2. Hover (Moves mouse to element)
	// Log start position
	// Note: Rod doesn't expose current mouse position easily without tracking it ourselves or querying browser.
	// We will just log the target of the hover.

	// Calculate center for logging purposes (Hover does this internally)
	quads, err := el.Shape()
	if err == nil && len(quads.Quads) > 0 {
		quad := quads.Quads[0]
		cx := (quad[0] + quad[2] + quad[4] + quad[6]) / 4
		cy := (quad[1] + quad[3] + quad[5] + quad[7]) / 4
		LogMouse(cx, cy, "hover_target")
	}

	if err := el.Hover(); err != nil {
		return err
	}
	RandomSleep(150, 150) // Pause before clicking (Reaction time)

	// 3. Click
	LogMouse(0, 0, "click") // Coordinates 0,0 mean "current position" in this context or just event logging
	if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	return nil
}

// HumanType simulates human-like typing with jitter.
func HumanType(page *rod.Page, text string) {
	for i, char := range text {
		// Type the character as a keystroke
		k := input.Key(char)
		page.Keyboard.Type(k)

		// Base typing speed: 30-80ms
		latency := rand.Intn(50) + 30

		// Every 15 characters, take a "thinking pause"
		if i%15 == 0 {
			latency += rand.Intn(400)
		}

		time.Sleep(time.Duration(latency) * time.Millisecond)
	}
	// Trigger React change detection by typing a dummy character and deleting it
	page.Keyboard.MustType(input.Space)
	page.Keyboard.MustType(input.Backspace)
}

// PrepareImageUpload resolves a data URL, HTTP(S) URL, or local path into a temp file ready for SetFiles.
// It returns the file path and a cleanup function.
func PrepareImageUpload(imageRef string) (string, func(), error) {
	imageRef = strings.TrimSpace(imageRef)
	if imageRef == "" {
		return "", func() {}, nil
	}

	if strings.HasPrefix(imageRef, "data:") {
		return writeDataURLToTempFile(imageRef)
	}

	if strings.HasPrefix(imageRef, "http://") || strings.HasPrefix(imageRef, "https://") {
		return downloadImageToTempFile(imageRef)
	}

	if _, err := os.Stat(imageRef); err == nil {
		return imageRef, func() {}, nil
	}

	return "", func() {}, fmt.Errorf("unsupported image reference")
}

func writeDataURLToTempFile(dataURL string) (string, func(), error) {
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return "", func() {}, fmt.Errorf("invalid data url")
	}

	meta := parts[0]
	payload := parts[1]
	if !strings.Contains(meta, ";base64") {
		return "", func() {}, fmt.Errorf("unsupported data url encoding")
	}

	mimeType := strings.TrimPrefix(strings.SplitN(meta, ";", 2)[0], "data:")
	ext := extensionFromMime(mimeType)
	if ext == "" {
		ext = ".img"
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to decode data url: %w", err)
	}

	file, err := os.CreateTemp("", "postificus-cover-*"+ext)
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(decoded); err != nil {
		_ = os.Remove(file.Name())
		return "", func() {}, fmt.Errorf("failed to write image: %w", err)
	}

	return file.Name(), func() { _ = os.Remove(file.Name()) }, nil
}

func downloadImageToTempFile(url string) (string, func(), error) {
	resp, err := http.Get(url) //nolint:gosec // controlled by user input, used for upload
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", func() {}, fmt.Errorf("image download failed: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	ext := extensionFromMime(contentType)
	if ext == "" {
		ext = ".img"
	}

	file, err := os.CreateTemp("", "postificus-cover-*"+ext)
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		_ = os.Remove(file.Name())
		return "", func() {}, fmt.Errorf("failed to save image: %w", err)
	}

	return file.Name(), func() { _ = os.Remove(file.Name()) }, nil
}

func extensionFromMime(mimeType string) string {
	if mimeType == "" {
		return ""
	}
	if exts, err := mime.ExtensionsByType(mimeType); err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ""
}
