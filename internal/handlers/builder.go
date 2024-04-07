package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

//go:generate mockgen -source=gin_api.go -package=mocks -destination=./mocks/mock_gin_api.go

type Service interface {
}

type Authenticator interface {
}

type HandlerBuilder struct {
	srv    Service
	auth   Authenticator
	logger *slog.Logger
}

// New creates new handlers builder
func New(srv Service, auth Authenticator, logger *slog.Logger) *HandlerBuilder {
	return &HandlerBuilder{srv: srv, logger: logger, auth: auth}
}

// GetHandler initializes a default router with corresponding routes
func (b *HandlerBuilder) GetHandler() http.Handler {
	r := gin.Default()
	r.Use(b.errorMiddleware, b.authenticate)

	r.GET("/user_banner", b.handlerUserGetBanner)

	postGroup := r.Group("/banner", b.authorize)
	postGroup.POST("", b.handlerCreateBanner)
	postGroup.GET("", b.handlerListBanners)
	postGroup.DELETE("/*id", b.handlerDeleteBanner)
	postGroup.PATCH("/*id", b.handlerUpdateBanner)

	return r
}
