package app

import (
	"BannerFlow/internal/app/ginapp"
	"BannerFlow/internal/config"
	"BannerFlow/internal/handlers"
	"BannerFlow/internal/repo/cache"
	"BannerFlow/internal/repo/db"
	"BannerFlow/internal/services/banner"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type HTTPServer interface {
	MustRun()
	Stop(ctx context.Context)
}

type Provider struct {
	logger        *slog.Logger
	cfg           *config.Config
	redis         *redis.Client
	postgres      *pgxpool.Pool
	service       handlers.Service
	server        HTTPServer
	handlerGetter ginapp.HandlerGetter
	cache         banner.Cache
	db            banner.Database
}

func NewProvider(logger *slog.Logger, cfg *config.Config, redis *redis.Client, postgres *pgxpool.Pool) *Provider {
	return &Provider{
		logger:   logger,
		cfg:      cfg,
		redis:    redis,
		postgres: postgres,
	}
}

func (p *Provider) Server() HTTPServer {
	if p.server == nil {
		p.server = ginapp.NewHTTPServerApp(p.logger, p.cfg.GinCfg, p.handlerGetter)
	}
	return p.server
}

func (p *Provider) HandlerGetter() ginapp.HandlerGetter {
	if p.handlerGetter == nil {
		// TODO add auth
		p.handlerGetter = handlers.New(p.Service(), nil, p.logger)
	}
	return p.handlerGetter
}

func (p *Provider) Service() handlers.Service {
	if p.service == nil {
		p.service = banner.New(p.Db(), p.Cache(), p.logger, p.cfg.ServiceCfg)
	}
	return p.service
}

func (p *Provider) Cache() banner.Cache {
	if p.cache == nil {
		p.cache = cache.New(p.redis, p.cfg.RedisCfg.Cache)
	}
	return p.cache
}

func (p *Provider) Db() banner.Database {
	if p.db == nil {
		p.db = db.New(p.postgres)
	}
	return p.db
}
