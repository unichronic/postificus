package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"postificus/internal/controller"
	"postificus/internal/middleware"
	"postificus/internal/rabbitmq"
	"postificus/internal/service"
	"postificus/internal/storage"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

func main() {
	// 1. Config & Infrastructure
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := storage.InitDB(); err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer storage.CloseDB()

	if err := storage.InitRedis(); err != nil {
		log.Fatal("Failed to init Redis:", err)
	}
	defer storage.CloseRedis()

	// 2. RabbitMQ (Producer Only)
	rabbitAddr := os.Getenv("RABBITMQ_URL")
	if rabbitAddr == "" {
		host := os.Getenv("RABBITMQ_HOST")
		port := os.Getenv("RABBITMQ_PORT")
		if host != "" {
			if port == "" {
				port = "5672"
			}
			// Default guest/guest for private service
			rabbitAddr = fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port)
		} else {
			rabbitAddr = "amqp://guest:guest@localhost:5672/"
		}
	}
	// Simple retry loop
	var rabbitConn *rabbitmq.Connection
	var err error
	for i := 0; i < 5; i++ {
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
	producer := rabbitmq.NewProducer(rabbitConn)

	// 3. DI Container (Manual)

	// Repositories
	credsRepo := storage.NewCredentialsRepository()
	draftRepo := storage.NewDraftRepository()
	profileRepo := storage.NewProfileRepository()

	// Services
	authService := service.NewAuthService(credsRepo)
	draftService := service.NewDraftService(draftRepo)
	profileService := service.NewProfileService(profileRepo)
	activityService := service.NewActivityService(credsRepo)

	// Controllers
	authController := controller.NewAuthController(authService)
	settingsController := controller.NewSettingsController(authService, profileService)
	draftController := controller.NewDraftController(draftService)
	activityController := controller.NewActivityController(activityService)
	dashboardController := controller.NewDashboardController(activityService, producer)
	publishController := controller.NewPublishController(producer)

	// 4. Server Setup
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(middleware.RateLimit(middleware.NewRateLimiter(rate.Limit(20), 50)))
	e.Use(middleware.PrometheusMiddleware)

	// 5. Routes
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth & Connect
	e.POST("/api/connect/:platform", authController.HandleConnectPlatform) // Unified

	// Settings & Profile
	e.POST("/api/settings/credentials", settingsController.SaveCredentials) // Manual override
	e.GET("/api/settings/credentials/:platform", settingsController.GetCredentialsStatus)
	e.GET("/api/profile", settingsController.GetProfile)
	e.PUT("/api/profile", settingsController.SaveProfile)

	// Drafts
	e.PUT("/api/drafts/:id", draftController.UpdateDraft)
	e.GET("/api/drafts/:id", draftController.GetDraft)

	// Dashboard & Activity
	e.GET("/api/dashboard/activity", dashboardController.GetDashboardActivity)
	e.POST("/api/dashboard/sync", dashboardController.TriggerSync)

	// Live Activity (Platform specific)
	e.GET("/api/devto/activity", activityController.GetDevtoActivity)
	e.GET("/api/medium/activity", activityController.GetMediumActivity)

	// Publishing
	e.POST("/api/publish/:platform", publishController.PublishPost)

	// 6. Start
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// 7. Cleanup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
