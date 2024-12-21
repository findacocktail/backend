package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/caarlos0/env"
	"github.com/ramonmedeiros/iba/internal/app"
	"github.com/ramonmedeiros/iba/internal/pkg/recipes"
)

type config struct {
	Port string `env:"PORT" envDefault:"3000"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Enabled(context.Background(), slog.LevelError)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("could not parse config", slog.Any("error", err))
		os.Exit(1)
	}

	recipesService, err := recipes.New(logger)
	if err != nil {
		logger.Error("could not index recipes", slog.Any("error", err))
		os.Exit(1)
	}

	httpServer := app.New(cfg.Port, logger, recipesService)
	httpServer.Serve()
}
