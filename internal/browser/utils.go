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
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// RandomSleep sleeps for a variable duration (base + jitter)
func RandomSleep(baseMs int, jitterMs int) {
	ms := baseMs + rand.Intn(jitterMs)
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// HumanType simulates human-like typing with jitter.
func HumanType(page *rod.Page, text string) {
	for i, char := range text {
		k := input.Key(char)
		page.Keyboard.Type(k)

		latency := rand.Intn(50) + 30
		if i%15 == 0 {
			latency += rand.Intn(400)
		}
		time.Sleep(time.Duration(latency) * time.Millisecond)
	}
	page.Keyboard.MustType(input.Space)
	page.Keyboard.MustType(input.Backspace)
}

// PrepareImageUpload resolves a data URL, HTTP(S) URL, or local path into a temp file ready for SetFiles.
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
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", func() {}, fmt.Errorf("image download failed: %s", resp.Status)
	}
	ext := extensionFromMime(resp.Header.Get("Content-Type"))
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

// HumanMoveAndClick kept for potential future use but LinkedIn removed
func HumanMoveAndClick(page *rod.Page, selector string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	el.ScrollIntoView()
	RandomSleep(200, 300)
	if err := el.Hover(); err != nil {
		return err
	}
	RandomSleep(150, 150)
	return el.Click(proto.InputMouseButtonLeft, 1)
}
