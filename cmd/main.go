package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/LamichhaneBibek/dev-ops/apiserver"
	"github.com/LamichhaneBibek/dev-ops/config"
	"github.com/LamichhaneBibek/dev-ops/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := config.New()
	if err != nil {
		return err
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	db, err := store.NewPostgresDb(config)
	if err != nil {
		return err
	}

	jwtManager := apiserver.NewJWTManager(config)
	dataStore := store.New(db)
	server := apiserver.New(config, logger, dataStore, jwtManager)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
