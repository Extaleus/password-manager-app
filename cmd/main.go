package main

import (
	"log/slog"
	"os"

	"github.com/Extaleus/password-manager-app/internal/config"
	"github.com/Extaleus/password-manager-app/internal/lib/logger/sl"
	"github.com/Extaleus/password-manager-app/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	log.Info("start project", slog.String("env", cfg.Env))
	log.Debug("debug is worked")

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	err = storage.DeleteUrl("google")
	if err != nil {
		log.Error("failed delete url", sl.Err(err))
	}

	// TODO: init router: chi

	// TODO: run server:
}

func setupLogger(env string) *slog.Logger { // Argument for separation prod/dev/local logging logic
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
