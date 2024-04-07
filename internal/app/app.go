package app

import (
	"BannerFlow/internal/config"
	"log/slog"
	"os"
)

type App struct {
	provider Provider
}

// NewApp  creates new main app
func NewApp() *App {
	return &App{}
}

// Run runs the app
func (a *App) Run() {
	cfg := config.MustLoad()
	logger := setupLogger()
	provider := NewProvider(logger, cfg)

	logger.Info("starting app")
	provider.Server().MustRun()
}

// setupLogger setups logger options. Some config can be added to switch different modes
func setupLogger() *slog.Logger {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}

// Stop stops the app
func (a *App) Stop() {
	a.provider.logger.Info("stopping app")
	a.provider.Server().Stop()
}
