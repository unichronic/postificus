package database

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *pgxpool.Pool
	pgOnce     sync.Once
)

func InitPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	pgOnce.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			err = fmt.Errorf("DATABASE_URL not set")
			return
		}

		config, parseErr := pgxpool.ParseConfig(dbURL)
		if parseErr != nil {
			err = fmt.Errorf("unable to parse database URL: %w", parseErr)
			return
		}

		pgInstance, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		if pingErr := pgInstance.Ping(ctx); pingErr != nil {
			err = fmt.Errorf("unable to ping database: %w", pingErr)
			return
		}
	})

	return pgInstance, err
}

func GetPostgres() *pgxpool.Pool {
	return pgInstance
}

func ClosePostgres() {
	if pgInstance != nil {
		pgInstance.Close()
	}
}
