package cache

import (
	"BannerFlow/internal/config"
	"BannerFlow/internal/domain/models"
	"context"
	"errors"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	redisCache *cache.Cache
	ttl        time.Duration
}

// NewRedisClient new redis connection
func NewRedisClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func New(rdb *redis.Client, cfg *config.CacheConfig) *RedisCache {
	redisCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(cfg.LocalSize, cfg.TTL),
	})
	return &RedisCache{
		redisCache: redisCache,
		ttl:        cfg.TTL,
	}
}

func (r RedisCache) Get(ctx context.Context, options *models.BannerIdentOptions) (*models.UserBanner, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		res, err := r.handleGet(ctx, options)
		switch {
		case errors.Is(err, cache.ErrCacheMiss):
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return res, nil
		}
	}
}

func (r RedisCache) handleGet(ctx context.Context, options *models.BannerIdentOptions) (*models.UserBanner, error) {
	adapter := RedisStorageAdapter{
		FeatureId: options.FeatureId,
		TagId:     options.TagId,
	}
	err := r.redisCache.Get(ctx, adapter.Key(), &adapter.Bytes)
	if err != nil {
		return nil, err
	}
	content, err := adapter.Content()
	if err != nil {
		return nil, err
	}
	return &models.UserBanner{Content: content}, nil
}

func (r RedisCache) Put(ctx context.Context, options *models.BannerIdentOptions, banner *models.UserBanner) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		adapter := RedisStorageAdapter{
			UserBanner: banner.Content,
			FeatureId:  options.FeatureId,
			TagId:      options.TagId,
		}
		value, err := adapter.Value()
		if err != nil {
			return err
		}
		if err := r.redisCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   adapter.Key(),
			Value: value,
			TTL:   r.ttl,
		}); err != nil {
			return err
		}
		return nil
	}
}
