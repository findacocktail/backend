package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/caarlos0/env"
	"github.com/ramonmedeiros/iba/cmd"
	"github.com/ramonmedeiros/iba/internal/app"
)

type config struct {
	Port string `env:"PORT" envDefault:"3000"`
	File string `env:"FILE" envDefault:"20241126.json"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Enabled(context.Background(), slog.LevelError)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("could not parse config", slog.Any("error", err))
		os.Exit(1)
	}

	content, err := os.ReadFile(cfg.File)
	if err != nil {
		logger.Error("could not read file", slog.Any("error", err))
		os.Exit(1)
	}

	var recipes []*cmd.Recipe
	err = json.Unmarshal(content, &recipes)
	if err != nil {
		logger.Error("could not marshal", slog.Any("error", err))
		os.Exit(1)
	}

	httpServer := app.New(cfg.Port, logger, recipes)
	httpServer.Serve()
}
