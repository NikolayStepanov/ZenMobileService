package redis

import (
	"ZenMobileService/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Error(err)
		log.Infoln("Redis is not connect")
	} else {
		log.Infoln(pong)
		log.Infoln("Redis is connected")
	}
	return &RedisCache{Client: redisClient}, err
}

func (r *RedisCache) Set(ctx context.Context, key string, value any) error {
	err := r.Client.Set(ctx, key, value, 0).Err()
	return err
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	value, err := r.Client.Get(ctx, key).Result()
	return value, err
}

func (r *RedisCache) IncrementBy(ctx context.Context, key string, incrementValue int64) (int64, error) {
	err := error(nil)
	valueExists := int64(0)
	value := int64(0)

	valueExists, err = r.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Error(err)
	} else if valueExists == 0 {
		err = fmt.Errorf("redis: key = %s does not exist", key)
		log.Error(err)
	} else {
		value, err = r.Client.IncrBy(ctx, key, incrementValue).Result()
	}

	return value, err
}
