package cache

import (
	"BannerFlow/internal/config"
	"BannerFlow/internal/domain/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var testErr = errors.New("test error")

type RedisTest struct {
	suite.Suite
	redis *RedisCache
	mock  redismock.ClientMock
	ttl   time.Duration
}

func (s *RedisTest) SetupTest() {
	s.ttl = 5 * time.Minute
	var client *redis.Client
	client, s.mock = redismock.NewClientMock()
	s.redis = New(client, &config.CacheConfig{
		LocalSize: 100,
		TTL:       s.ttl,
	})
}

func (s *RedisTest) TestGetCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	banner, err := s.redis.Get(ctx, nil)
	s.Assert().ErrorIs(err, context.Canceled)
	s.Assert().Nil(banner)
}

func (s *RedisTest) TestGetError() {
	ctx := context.Background()
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     1,
	}
	key := fmt.Sprintf("feature: %d, tag: %d", options.FeatureId, options.TagId)
	s.mock.ExpectGet(key).SetErr(testErr)
	banner, err := s.redis.Get(ctx, options)
	s.Assert().ErrorIs(err, testErr)
	s.Assert().Nil(banner)
}

func (s *RedisTest) TestGetOK() {
	ctx := context.Background()
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     1,
	}
	key := fmt.Sprintf("feature: %d, tag: %d", options.FeatureId, options.TagId)
	expectedBanner := &models.UserBanner{Content: map[string]any{"title": "some"}}
	mockedResult, err := json.Marshal(&expectedBanner.Content)
	s.Require().NoError(err)
	s.mock.ExpectGet(key).SetVal(string(mockedResult))

	banner, err := s.redis.Get(ctx, options)
	s.Assert().NoError(err)
	s.Assert().Equal(*expectedBanner, *banner)
}

func (s *RedisTest) TestGetNil() {
	ctx := context.Background()
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     1,
	}
	key := fmt.Sprintf("feature: %d, tag: %d", options.FeatureId, options.TagId)
	s.mock.ExpectGet(key).RedisNil()

	banner, err := s.redis.Get(ctx, options)
	s.Assert().NoError(err)
	s.Assert().Nil(banner)
}

func (s *RedisTest) TestPutCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.redis.Put(ctx, nil, nil)
	s.Assert().ErrorIs(err, context.Canceled)
}

func (s *RedisTest) TestPutError() {
	ctx := context.Background()
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     1,
	}
	banner := &models.UserBanner{Content: map[string]any{"title": "some"}}
	key := fmt.Sprintf("feature: %d, tag: %d", options.FeatureId, options.TagId)
	mockedArg, err := json.Marshal(&banner.Content)
	s.Require().NoError(err)
	s.mock.ExpectSet(key, mockedArg, s.ttl).SetErr(testErr)
	err = s.redis.Put(ctx, options, banner)
	s.Assert().ErrorIs(err, testErr)
}

func (s *RedisTest) TestPutOK() {
	ctx := context.Background()
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     1,
	}
	banner := &models.UserBanner{Content: map[string]any{"title": "some"}}
	key := fmt.Sprintf("feature: %d, tag: %d", options.FeatureId, options.TagId)
	mockedArg, err := json.Marshal(&banner.Content)
	s.Require().NoError(err)
	s.mock.ExpectSet(key, mockedArg, s.ttl).SetVal("OK")
	err = s.redis.Put(ctx, options, banner)
	s.Assert().NoError(err)
}

func (s *RedisTest) TearDownTest() {
	s.Assert().NoError(s.mock.ExpectationsWereMet())
}

func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisTest))
}
