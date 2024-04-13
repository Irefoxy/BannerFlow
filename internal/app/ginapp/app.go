package ginapp

import (
	"BannerFlow/internal/config"
	"BannerFlow/internal/utils"
	"context"
	"errors"
	"log/slog"
	"net/http"
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
	log := g.logger.With(utils.Text(op))
	r := g.handlerGetter.GetHandler()
	g.server = &http.Server{
		Addr:    g.httpConfig.Address,
		Handler: r,
	}

	log.Info("Starting server...", "address", g.httpConfig.Address)
	if err := g.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Warn("ListenAndServe failed", utils.Err(err))
	}
	return nil
}

// Stop stops http server
func (g *HTTPServerApp) Stop(ctx context.Context) {
	op := "ginapp.stop"
	log := g.logger.With(utils.Text(op))
	log.Info("Shutting down server...")
	if g.server == nil {
		log.Info("nothing to shutdown")
		return
	}
	if err := g.server.Shutdown(ctx); err != nil {
		log.Warn("shutdown failed", utils.Err(err))
	}
}
