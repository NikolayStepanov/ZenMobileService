package redis

import (
	"ZenMobileService/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache(cfg *config.Config) *RedisCache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Error(err)
		log.Infoln("Redis is not Connect")
	} else {
		log.Infoln(pong)
		log.Infoln("Redis is Connected")
	}
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

func (r RedisCache) IncrementBy(ctx context.Context, key string, incrementValue int64) (int64, error) {
	value, err := r.cache.IncrBy(ctx, key, incrementValue).Result()
	return value, err
}
