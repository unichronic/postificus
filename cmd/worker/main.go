package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"postificus/internal/rabbitmq"
	"postificus/internal/service"
	"postificus/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 2. Init DB & Redis
	if err := storage.InitDB(); err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer storage.CloseDB()

	if err := storage.InitRedis(); err != nil {
		log.Fatal("Failed to init Redis:", err)
	}
	defer storage.CloseRedis()

	// 3. Init RabbitMQ
	rabbitAddr := os.Getenv("RABBITMQ_URL")
	if rabbitAddr == "" {
		rabbitAddr = "amqp://guest:guest@localhost:5672/"
	}

	var rabbitConn *rabbitmq.Connection
	var err error
	for i := 0; i < 10; i++ {
		rabbitConn, err = rabbitmq.NewConnection(rabbitAddr)
		if err == nil {
			break
		}
		log.Printf("Waiting for RabbitMQ... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	consumer := rabbitmq.NewConsumer(rabbitConn)

	// 4. Init Dependencies
	credsRepo := storage.NewCredentialsRepository()
	activityService := service.NewActivityService(credsRepo)
	syncWorker := service.NewSyncService(activityService)
	publishService := service.NewPublishService(credsRepo)

	// 4. Start Consumers (Parallel Workers)
	parallelism := 5
	log.Printf("🚀 Starting %d parallel workers for publishing tasks...", parallelism)

	// We simply call Consume multiple times. RabbitMQ library will multiplex them over the same connection/channel or create new channels.
	// Our `Consume` method creates a NEW channel for each call, which is thread-safe and correct.
	for i := 0; i < parallelism; i++ {
		workerID := i + 1
		go func(id int) {
			log.Printf("Starting worker %d...", id)
			err := consumer.Consume(service.TypePublishPost, publishService.HandlePublishTask)
			if err != nil {
				log.Printf("❌ Worker %d failed to start: %v", id, err)
			}
		}(workerID)
	}

	// Also start sync consumer (single worker is fine for now)
	go func() {
		log.Printf("Starting sync worker...")
		err := consumer.Consume(service.TypeSyncPlatformActivity, syncWorker.HandleSyncPlatformActivity)
		if err != nil {
			log.Printf("❌ Sync worker failed: %v", err)
		}
	}()

	// 6. Block until interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down worker...")
}
