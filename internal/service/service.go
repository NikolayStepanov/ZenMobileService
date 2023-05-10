package service

import (
	"ZenMobileService/internal/domain"
	"ZenMobileService/internal/repository"
	"ZenMobileService/internal/service/sign"
	"context"
)

type MemoryCache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
	IncrementBy(ctx context.Context, key string, incrementValue int64) (int64, error)
}

type ServicesDependencies struct {
	Cache MemoryCache
	Repos repository.Repository
}

type CacheServicer interface {
	IncrementValueByKey(ctx context.Context, key string, incrementValue int64) (int64, error)
	SetValueByKey(ctx context.Context, key string, value any) error
	GetValueByKey(ctx context.Context, key string) (any, error)
}

type SignatureServicer interface {
	GenerateSignature(ctx context.Context, text, key string) (string, error)
	ValidSignature(ctx context.Context, signature, text, key string) (bool, error)
}

type UsersServicer interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUser(ctx context.Context, userId int) (domain.User, error)
}

type Services struct {
	CacheService CacheServicer
	SignService  SignatureServicer
	UsersService UsersServicer
}

func NewServices(deps ServicesDependencies) *Services {
	cacheService := NewCacheService(deps.Cache)
	signService := sign.NewSignService()
	userService := NewUsersService(deps.Repos.UsersRep)
	return &Services{CacheService: cacheService, SignService: signService, UsersService: userService}
}
