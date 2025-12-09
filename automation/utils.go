package automation

import (
	"fmt"
	"math/rand"
	"os"
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
