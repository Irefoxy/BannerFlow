package app

import (
	"BannerFlow/internal/config"
	"BannerFlow/internal/repo/cache"
	"BannerFlow/internal/repo/db"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"sync"
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
	var client *redis.Client
	var postgres *pgxpool.Pool

	cfg := config.MustLoad()
	logger := setupLogger()

	logger.Info("Establishing redis connections")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer cancel()

	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		client, err = cache.NewRedisClient(ctx, cfg.RedisCfg.Address)
		return
	})
	eg.Go(func() (err error) {
		postgres, err = db.NewPostgres(ctx, cfg.PostgresCfg.DSN)
		return
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	provider := NewProvider(logger, cfg, client, postgres)

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
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), a.provider.cfg.StartTimeout)
	defer cancel()

	// TODO add errgroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		a.provider.Server().Stop(ctx)
	}()
}
