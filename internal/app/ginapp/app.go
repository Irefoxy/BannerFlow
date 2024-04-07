package ginapp

import (
	"BannerFlow/internal/config"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HandlerGetter interface {
	GetHandler() http.Handler
}

type HTTPServerApp struct {
	httpConfig    *config.HTTPConfig
	server        *http.Server
	logger        *slog.Logger
	handlerGetter HandlerGetter
}

// NewHTTPServerApp creates new app to run a http server
func NewHTTPServerApp(logger *slog.Logger, cfg *config.HTTPConfig, handlerGetter HandlerGetter) *HTTPServerApp {
	return &HTTPServerApp{
		httpConfig:    cfg,
		logger:        logger,
		handlerGetter: handlerGetter,
	}
}

// MustRun Runs http server or panic on error
func (g *HTTPServerApp) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

// Run runs http server
func (g *HTTPServerApp) Run() error {
	const op = "ginapp.run"
	log := g.logger.With(op)

	r := g.handlerGetter.GetHandler()
	srv := &http.Server{
		Addr:    g.httpConfig.Address,
		Handler: r,
	}

	log.Info("Starting server...", "address", g.httpConfig.Address)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Warn(fmt.Errorf("%s %w", op, err).Error())
	}
	return nil
}

// Stop stops http server
func (g *HTTPServerApp) Stop() {
	op := "ginapp.stop"
	log := g.logger.With(op)
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), g.httpConfig.Timeout*time.Second)
	defer cancel()

	if err := g.server.Shutdown(ctx); err != nil {
		log.Warn(fmt.Errorf("%s %w", op, err).Error())
	}
}
