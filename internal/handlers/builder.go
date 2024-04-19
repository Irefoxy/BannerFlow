package handlers

import (
	"BannerFlow/internal/domain/models"
	"BannerFlow/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

const (
	tokenName = "token"
)

//go:generate mockgen -source=builder.go -package=mocks -destination=./mocks/mock_handlers.go
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
	Authenticate(token string) error
}

type Authorizer interface {
	IsAdmin(token string) bool
}

type TokenGenerator interface {
	GenerateToken(isAdmin bool) (string, error)
}

type HandlerBuilder struct {
	srv           Service
	authenticator Authenticator
	logger        *slog.Logger
	authorizer    Authorizer
	generator     TokenGenerator
}

// New creates new handlers builder
func New(srv Service, auth Authenticator, authorizer Authorizer, generator TokenGenerator, logger *slog.Logger) *HandlerBuilder {
	return &HandlerBuilder{srv: srv, logger: logger, authenticator: auth, authorizer: authorizer, generator: generator}
}

// GetHandler initializes a default router with corresponding routes
func (b *HandlerBuilder) GetHandler() http.Handler {
	r := gin.Default()
	r.Use(b.errorMiddleware)

	r.GET("/get_token/*admin", b.handleTokenGeneration)

	authenticateGroup := r.Group("/", b.authenticate)
	authenticateGroup.GET("/user_banner", b.handleUserGetBanner)

	adminGroup := authenticateGroup.Group("/banner", b.authorize)
	adminGroup.GET("", b.handleListBanners)
	adminGroup.POST("", b.handleCreateBanner)
	adminGroup.DELETE("/:id", b.handleDeleteBanner)
	adminGroup.PATCH("/:id", b.handleUpdateBanner)
	adminGroup.GET("/versions/:id", b.handleListBannerHistory)
	adminGroup.PUT("/versions/:id/activate", b.handleSelectBannerVersion)
	adminGroup.DELETE("/del", b.handleDeleteBannerByTagOrFeature)

	return r
}

func (b *HandlerBuilder) log(c *gin.Context) {
	const op = "handlers.log"
	log := b.logger.With(utils.Text(op))
	for _, msg := range c.Errors.Errors() {
		log.Warn(msg)
	}
}
