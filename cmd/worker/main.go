package main

import (
	"log"
	"os"

	"postificus/internal/automation"
	"postificus/internal/database"
	"postificus/internal/queue"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, creating one if config is saved.")
	}

	// Initialize Shared Browser
	// The worker needs the browser for automation tasks
	if err := automation.InitBrowser(); err != nil {
		log.Fatal("Failed to initialize browser:", err)
	}
	defer automation.CloseBrowser()

	// Initialize Database (if needed for fetching user cookies)
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// CRITICAL: Parallelism Control
			// If you have 4GB RAM, allows max 5 browsers at once.
			// All other jobs wait in the Redis queue automatically.
			Concurrency: 5,

			// Priority Queues
			Queues: map[string]int{
				"critical": 6, // Real-time publishes
				"default":  3, // Scheduled posts
				"low":      1, // Analytics scraping
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypePublishPost, queue.HandlePublishTask)

	log.Println("🤖 Starting Background Worker...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
