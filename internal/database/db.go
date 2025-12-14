package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

//go:embed schema.sql
var schemaSQL string

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	// Auto-Migrate
	if _, err := DB.Exec(ctx, schemaSQL); err != nil {
		log.Printf("Warning: Failed to apply schema: %v", err)
	} else {
		log.Println("✅ Schema applied successfully")
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("🛑 Database connection closed")
	}
}
