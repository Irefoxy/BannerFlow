package app

import (
	"BannerFlow/internal/config"
	"BannerFlow/internal/repo/cache"
	"BannerFlow/internal/repo/db"
	"context"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
)

type task func() error

type App struct {
	provider *Provider
}

// NewApp  creates new main app
func NewApp() *App {
	return &App{}
}

// Run runs the app
func (a *App) Run() {
	a.provider = &Provider{}
	a.provider.cfg = config.MustLoad()
	a.provider.logger = setupLogger()

	a.provider.logger.Info("Establishing redis and postgres connections")
	ctx, cancel := context.WithTimeout(context.Background(), a.provider.cfg.InitTimeout)
	defer cancel()

	eg, _ := errgroup.WithContext(ctx)
	for _, task := range a.collectFuncsToGo(ctx) {
		eg.Go(task)
	}
	if err := eg.Wait(); err != nil {
		panic(err)
	}
	a.provider.logger.Info("connected")

	if runnable, ok := a.provider.Service().(RunnableService); ok {
		go runnable.MustRun()
		a.provider.logger.Info("Service running")
	}

	a.provider.logger.Info("starting app")
	a.provider.Server().MustRun()
}

func (a *App) collectFuncsToGo(ctx context.Context) []task {
	var result []task
	redisFn := func() (err error) {
		a.provider.redis, err = cache.NewRedisClient(ctx, a.provider.cfg.RedisCfg)
		return
	}
	postgresFn := func() (err error) {
		a.provider.postgres, err = db.NewPostgres(ctx, a.provider.cfg.PostgresCfg.DSN)
		return
	}
	result = append(result, redisFn, postgresFn)
	return result
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
	if a.provider == nil {
		a.provider.logger.Info("nothing to stop")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.provider.cfg.InitTimeout)
	defer cancel()

	a.provider.Server().Stop(ctx)

	if stoppable, ok := a.provider.Service().(StoppableService); ok {
		a.provider.logger.Info("Service is stopping")
		stoppable.Stop(ctx)
	}
}
