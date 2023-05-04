package service

import "context"

type CacheService struct {
	cache MemoryCache
}

func NewCacheService(cache MemoryCache) *CacheService {
	return &CacheService{cache: cache}
}

func (c *CacheService) IncrementValueByKey(ctx context.Context, key string, incrementValue int64) (int64, error) {
	value, err := c.cache.IncrementBy(ctx, key, incrementValue)
	return value, err
}

func (c *CacheService) SetValueByKey(ctx context.Context, key string, value any) error {
	err := c.cache.Set(ctx, key, value)
	return err
}

func (c *CacheService) GetValueByKey(ctx context.Context, key string) (any, error) {
	value, err := c.cache.Get(ctx, key)
	return value, err
}
