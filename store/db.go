package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamichhaneBibek/dev-ops/config"
	_ "github.com/lib/pq"
)

func NewPostgresDb(conf *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.DatabaseUrl())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
