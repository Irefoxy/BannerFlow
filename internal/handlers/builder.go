package handlers

import (
	"BannerFlow/internal/domain/models"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

const (
	tokenName = "token"
)

//go:generate mockgen -source=gin_api.go -package=mocks -destination=./mocks/mock_gin_api.go
type Service interface {
	CreateBanner(ctx context.Context, banner *models.Banner) (int, error)
	DeleteBanner(ctx context.Context, id int) error
	ListBanners(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error)
	UserGetBanners(ctx context.Context, options *models.BannerUserOptions) (*models.UserBanner, error)
	UpdateBanner(ctx context.Context, id int, banner *models.UpdateBanner) error
	ListBannerHistory(ctx context.Context, id int) ([]models.HistoryBanner, error)
	SelectBannerVersion(ctx context.Context, id, version int) error
	DeleteBannersByTagOrFeature(ctx context.Context, options *models.BannerIdentOptions) error
}

type Authenticator interface {
	Authenticate(token string) error // TODO get user?
}

type Authorizer interface {
	IsAdmin(token string) bool
}

type HandlerBuilder struct {
	srv           Service
	authenticator Authenticator
	logger        *slog.Logger
	authorizer    Authorizer
}

// New creates new handlers builder
func New(srv Service, auth Authenticator, logger *slog.Logger) *HandlerBuilder {
	return &HandlerBuilder{srv: srv, logger: logger, authenticator: auth}
}

// GetHandler initializes a default router with corresponding routes
func (b *HandlerBuilder) GetHandler() http.Handler {
	r := gin.Default()
	r.Use(b.errorMiddleware, b.authenticate)

	r.GET("/user_banner", b.handleUserGetBanner)

	postGroup := r.Group("/banner", b.authorize)
	postGroup.GET("", b.handleListBanners)
	postGroup.POST("", b.handleCreateBanner)
	postGroup.DELETE("/:id", b.handleDeleteBanner)
	postGroup.PATCH("/:id", b.handleUpdateBanner)
	postGroup.GET("/versions/:id", b.handleListBannerHistory)
	postGroup.PUT("/versions/:id/activate", b.handleSelectBannerVersion)
	postGroup.DELETE("/banners", b.handleDeleteBanner)

	return r
}

func (b *HandlerBuilder) log(c *gin.Context) {
	const op = "handlers.log"
	log := b.logger.With(op)
	for _, msg := range c.Errors.Errors() {
		log.Warn(msg)
	}
}
