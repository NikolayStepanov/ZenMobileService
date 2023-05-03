package service

import "context"

type CacheService struct {
	cache MemoryCache
}

func NewCacheService(cache MemoryCache) *CacheService {
	return &CacheService{cache: cache}
}

func (c CacheService) IncrementValueByKey(ctx context.Context, key string, incrementValue any) (any, error) {
	//TODO implement me
	panic("implement me")
}
