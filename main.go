package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/caarlos0/env"
	"github.com/findacocktail/backend/cmd"
	"github.com/findacocktail/backend/internal/app"
	"github.com/findacocktail/backend/internal/pkg/recipes"
)

type config struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "scrape" {
		cmd.Scrape()
		os.Exit(0)
	}

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
