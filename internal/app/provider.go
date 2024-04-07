package app

import (
	"BannerFlow/internal/app/ginapp"
	"BannerFlow/internal/config"
	"BannerFlow/internal/handlers"
	"log/slog"
)

type HTTPServer interface {
	MustRun()
	Stop()
}

type Provider struct {
	logger        *slog.Logger
	cfg           *config.Config
	server        HTTPServer
	handlerGetter ginapp.HandlerGetter
}

func NewProvider(logger *slog.Logger, cfg *config.Config) *Provider {
	return &Provider{
		logger: logger,
		cfg:    cfg,
	}
}

func (p *Provider) Server() HTTPServer {
	if p.server == nil {
		p.server = ginapp.NewHTTPServerApp(p.logger, p.cfg.GinConfig, p.handlerGetter)
	}
	return p.server
}

func (p *Provider) HandlerGetter() ginapp.HandlerGetter {
	if p.handlerGetter == nil {
		// TODO add srv and auth
		p.handlerGetter = handlers.New(nil, nil, p.logger)
	}
	return p.handlerGetter
}
