package service

import "context"

type MemoryCache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
	IncrementBy(ctx context.Context, key string, incrementValue int64) (int64, error)
}

type ServicesDependencies struct {
	Cache MemoryCache
}

type SignatureServicer interface {
	GenerateSignature(ctx context.Context, text, key string) (string, error)
	ParseSignature(ctx context.Context, signature, key string) (string, error)
}

type CacheServicer interface {
	IncrementValueByKey(ctx context.Context, key string, incrementValue int64) (int64, error)
	SetValueByKey(ctx context.Context, key string, value any) error
	GetValueByKey(ctx context.Context, key string) (any, error)
}

type Services struct {
	CacheService CacheServicer
	SignService  SignatureServicer
}

func NewServices(deps ServicesDependencies) *Services {
	cacheService := NewCacheService(deps.Cache)
	signService := NewSignService()
	return &Services{CacheService: cacheService, SignService: signService}
}
