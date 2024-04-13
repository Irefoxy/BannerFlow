package banner

import (
	"BannerFlow/internal/config"
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Database interface {
	Add(ctx context.Context, banner *models.Banner) (int, error)
	Update(ctx context.Context, id int, banner *models.UpdateBanner) error
	List(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error)
	DeleteByIds(ctx context.Context, ids ...int) error
	DeleteByFeatureOrTag(ctx context.Context, options *models.BannerIdentOptions) error
	GetHistoryForId(ctx context.Context, id int) ([]models.HistoryBanner, error)
	SelectBannerVersion(ctx context.Context, id, version int) error
}

type Cache interface {
	Get(ctx context.Context, options *models.BannerIdentOptions) (*models.UserBanner, error)
	Put(ctx context.Context, options *models.BannerIdentOptions, banner *models.UserBanner) error
}

type Service struct {
	wg             sync.WaitGroup
	timeout        time.Duration
	logger         *slog.Logger
	db             Database
	cache          Cache
	tasksChan      chan *models.BannerIdentOptions
	activeRequests int64
}

func New(db Database, cache Cache, logger *slog.Logger, cfg *config.ServiceConfig) *Service {
	return &Service{
		timeout:   cfg.Timeout,
		logger:    logger,
		db:        db,
		cache:     cache,
		tasksChan: make(chan *models.BannerIdentOptions, 100),
	}
}

func (s *Service) Start() {
	for ident := range s.tasksChan {
		for {
			if atomic.LoadInt64(&s.activeRequests) < 200 {
				break
			}
			time.Sleep(2 * time.Second)
		}
		const op = "long deletion"
		newCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
		err := s.db.DeleteByFeatureOrTag(newCtx, ident)
		if err != nil {
			s.logger.Warn(op, fmt.Sprintf("Failed to delete: tag: %d feature: %d", ident.TagId, ident.FeatureId), err.Error())
		}
		cancel()
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
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
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

func (s *Service) ListBannerHistory(ctx context.Context, id int) ([]models.HistoryBanner, error) {
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
	const op = "banner.ListBannerHistory"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return nil, e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	banners, err := s.db.GetHistoryForId(newCtx, id)
	if err != nil {
		log.Warn(err.Error())
		return nil, e.ErrorInternal
	}
	if len(banners) == 0 {
		log.Info("no banner found")
		return nil, e.ErrorNotFound
	}
	return banners, nil
}

func (s *Service) SelectBannerVersion(ctx context.Context, id, version int) error {
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
	const op = "banner.ListBannerHistory"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.db.SelectBannerVersion(newCtx, id, version)
	if err != nil {
		log.Warn(err.Error())
		return e.ErrorInternal
	}
	return nil
}

func (s *Service) DeleteBanner(ctx context.Context, id int) error {
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
	const op = "banner.CreateBanner"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return e.ErrorInternal
	}
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.db.DeleteByIds(newCtx, id)
	if err != nil {
		log.Warn(err.Error())
		if errors.Is(err, e.ErrorNotFound) {
			return err
		}
		return e.ErrorInternal
	}
	return nil
}

func (s *Service) DeleteBannersByTagOrFeature(ctx context.Context, options *models.BannerIdentOptions) error {
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
	const op = "banner.CreateBanner"
	log := s.logger.With(op)
	if s.ctxDone(ctx, log) {
		return e.ErrorInternal
	}
	s.wg.Add(1)
	s.tasksChan <- options
	return nil
}

func (s *Service) ListBanners(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error) {
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
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
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
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
		userBanner, err = s.getBannerFromCache(newCtx, &options.BannerIdentOptions, log)
	}
	if options.UseLastRevision || err != nil {
		banner, err := s.getBannerFromDb(newCtx, &options.BannerIdentOptions, log)
		if err != nil {
			return nil, err
		}
		s.wg.Add(1)
		go s.SendBannerToCache(newCtx, &options.BannerIdentOptions, banner, log.With(op))
		userBanner = banner
	}
	return userBanner, nil
}

func (s *Service) UpdateBanner(ctx context.Context, id int, banner *models.UpdateBanner) error {
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
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

func (s *Service) getBannerFromDb(newCtx context.Context, options *models.BannerIdentOptions, log *slog.Logger) (*models.UserBanner, error) {
	const op = "banner.getBannerFromDb"
	banners, err := s.ListBanners(newCtx, &models.BannerListOptions{
		BannerIdentOptions: *options,
		Limit:              models.ZeroValue,
		Offset:             models.ZeroValue,
	})
	if err != nil {
		return nil, err
	}
	if len(banners) == 0 || !banners[0].IsActive {
		log.Info(op, "missing banner")
		return nil, e.ErrorNotFound
	}
	return &banners[0].UserBanner, nil
}

func (s *Service) getBannerFromCache(newCtx context.Context, options *models.BannerIdentOptions, log *slog.Logger) (*models.UserBanner, error) {
	const op = "banner.getBannerFromCache"
	banner, err := s.cache.Get(newCtx, options)
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

func (s *Service) SendBannerToCache(newCtx context.Context, options *models.BannerIdentOptions, banner *models.UserBanner, log *slog.Logger) {
	const op = "banner.SendBannerToCache"
	defer s.wg.Done()
	err := s.cache.Put(newCtx, options, banner)
	if err != nil {
		log.Warn(op, err.Error())
	}
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
