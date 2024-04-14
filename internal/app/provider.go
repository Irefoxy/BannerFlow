package app

import (
	"BannerFlow/internal/app/ginapp"
	"BannerFlow/internal/auth"
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

type RunnableService interface {
	handlers.Service
	MustRun()
}

type StoppableService interface {
	handlers.Service
	Stop(ctx context.Context)
}

type SSO interface {
	handlers.Authenticator
	handlers.Authorizer
	handlers.TokenGenerator
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
	sso           SSO
}

func (p *Provider) Server() HTTPServer {
	if p.server == nil {
		p.server = ginapp.NewHTTPServerApp(p.logger, p.cfg.GinCfg, p.HandlerGetter())
	}
	return p.server
}

func (p *Provider) HandlerGetter() ginapp.HandlerGetter {
	if p.handlerGetter == nil {
		// TODO add auth
		p.handlerGetter = handlers.New(p.Service(), p.SSO(), p.SSO(), p.SSO(), p.logger)
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
		p.cache = cache.New(p.redis, p.cfg.CacheCfg)
	}
	return p.cache
}

func (p *Provider) Db() banner.Database {
	if p.db == nil {
		p.db = db.New(p.postgres)
	}
	return p.db
}

func (p *Provider) SSO() SSO {
	if p.sso == nil {
		p.sso = auth.NewAuth()
	}
	return p.sso
}
