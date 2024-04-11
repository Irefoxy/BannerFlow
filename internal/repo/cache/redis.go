package cache

import (
	"BannerFlow/internal/config"
	"BannerFlow/internal/services/models"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	redisCache *cache.Cache
	ttl        time.Duration
}

// TODO should be moved
// NewRedisClient new redis connection
func NewRedisClient(ctx context.Context, address string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// TODO switch client to iface for testing
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
	id := ""
	key := fmt.Sprintf("feature: %d, tag: %d", options.FeatureId, options.TagId)
	err := r.redisCache.Get(ctx, key, &id)
	if err != nil {
		return nil, err
	}
	res := models.UserBanner{}
	err = r.redisCache.Get(ctx, id, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r RedisCache) Put(ctx context.Context, banner *models.Banner) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		id := uuid.NewString()
		if err := r.redisCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   id,
			Value: *banner,
			TTL:   r.ttl,
		}); err != nil {
			return err
		}
		for _, tag := range banner.TagId {
			key := fmt.Sprintf("feature: %d, tag: %d", banner.FeatureId, tag)
			if err := r.redisCache.Set(&cache.Item{
				Ctx:   ctx,
				Key:   key,
				Value: id,
				TTL:   r.ttl,
			}); err != nil {
				return err
			}
		}
		return nil
	}
}
