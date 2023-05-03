package redis

import (
	"ZenMobileService/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache(cfg *config.Config) *RedisCache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})
	return &RedisCache{cache: redisClient}
}

func (r RedisCache) Set(ctx context.Context, key string, value any) error {
	err := r.cache.Set(ctx, key, value, 0).Err()
	return err
}

func (r RedisCache) Get(ctx context.Context, key string) (any, error) {
	value, err := r.cache.Get(ctx, key).Result()
	return value, err
}
