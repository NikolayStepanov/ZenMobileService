package service

import "context"

type MemoryCache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
}

type ServicesDependencies struct {
	Cache MemoryCache
}

type CacheServices interface {
	IncrementValueByKey(ctx context.Context, key string, incrementValue any) (any, error)
}

type Services struct {
	CacheServices CacheServices
}

func NewServices(deps ServicesDependencies) *Services {
	cacheService := NewCacheService(deps.Cache)
	return &Services{CacheServices: cacheService}
}
