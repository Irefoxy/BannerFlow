package banner

import (
	"BannerFlow/internal/config"
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/services/models"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

type Database interface {
	Add(ctx context.Context, banner *models.Banner) (int, error)
	Update(ctx context.Context, id int, banner *models.Banner) error
	List(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error)
}

type Cache interface {
	Get(ctx context.Context, options *models.BannerIdentOptions) (*models.UserBanner, error)
	Put(ctx context.Context, banner *models.Banner) error
}

type Service struct {
	wg      sync.WaitGroup
	timeout time.Duration
	logger  *slog.Logger
	db      Database
	cache   Cache
}

func New(db Database, cache Cache, logger *slog.Logger, cfg *config.ServiceConfig) *Service {
	return &Service{
		timeout: cfg.Timeout,
		logger:  logger,
		db:      db,
		cache:   cache,
	}
}

func (s *Service) Stop(ctx context.Context) error {
	const op = "banner.Stop"
	s.logger.Info(op, "stopping banner service")
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		s.logger.Info(op, "failed to stop service", ctx.Err().Error())
		return ctx.Err()
	case <-done:
		s.logger.Info(op, "service stopped")
		return nil
	}
}

func (s *Service) CreateBanner(ctx context.Context, banner *models.Banner) (int, error) {
	const op = "banner.CreateBanner"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return 0, e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	id, err := s.db.Add(newCtx, banner)
	if err != nil {
		log.Warn(err.Error())
		if errors.Is(err, e.ErrorConflict) {
			return 0, err
		}
		return 0, e.ErrorInternal
	}
	return id, nil
}

func (s *Service) DeleteBanner(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ListBanners(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error) {
	const op = "banner.ListBanner"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return nil, e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	list, err := s.db.List(newCtx, options)
	if err != nil {
		log.Warn("failed to get data", err.Error())
		return nil, e.ErrorInternal
	}
	return list, nil
}

func (s *Service) UserGetBanners(ctx context.Context, options *models.BannerUserOptions) (*models.UserBanner, error) {
	const op = "banner.UserGetBanner"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return nil, e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var err error
	var userBanner *models.UserBanner
	if !options.UseLastRevision {
		userBanner, err = s.getBannerFromCache(newCtx, options, log)
	}
	if options.UseLastRevision || err != nil {
		banner, err := s.getBannerFromDb(newCtx, options, log)
		if err != nil {
			return nil, err
		}
		go s.SendBannerToCache(newCtx, banner, log.With(op))
		userBanner = &banner.UserBanner
	}
	return userBanner, nil
}

func (s *Service) getBannerFromDb(newCtx context.Context, options *models.BannerUserOptions, log *slog.Logger) (*models.Banner, error) {
	const op = "banner.getBannerFromDb"
	banners, err := s.ListBanners(newCtx, &models.BannerListOptions{
		BannerIdentOptions: options.BannerIdentOptions,
	})
	if err != nil {
		return nil, err
	}
	if len(banners) == 0 || !banners[0].IsActive {
		log.Info(op, "missing banner")
		return nil, e.ErrorNotFound
	}
	return &banners[0].Banner, nil
}

func (s *Service) getBannerFromCache(newCtx context.Context, options *models.BannerUserOptions, log *slog.Logger) (*models.UserBanner, error) {
	const op = "banner.getBannerFromCache"
	banner, err := s.cache.Get(newCtx, &options.BannerIdentOptions)
	if err != nil {
		log.Warn(op, err.Error())
		return nil, err
	}
	if banner == nil {
		log.Info(op, "no banner found in cache")
		return nil, e.ErrorNotFound
	}
	return banner, nil
}

func (s *Service) SendBannerToCache(newCtx context.Context, banner *models.Banner, log *slog.Logger) {
	const op = "banner.SendBannerToCache"
	s.wg.Add(1)
	defer s.wg.Done()
	err := s.cache.Put(newCtx, banner)
	if err != nil {
		log.Warn(op, err.Error())
	}
}

func (s *Service) UpdateBanner(ctx context.Context, id int, banner *models.Banner) error {
	const op = "banner.UpdateBanner"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.db.Update(newCtx, id, banner)
	if err != nil {
		log.Warn(op, err.Error())
		if errors.Is(err, e.ErrorNotFound) {
			return err
		}
		return e.ErrorInternal
	}
	return nil
}

func (s *Service) ctxDone(ctx context.Context, log *slog.Logger) bool {
	const op = "banner.ctxDone"
	select {
	case <-ctx.Done():
		log.Info(op, "context done")
		return true
	default:
		return false
	}
}
