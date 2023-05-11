package service

import (
	"ZenMobileService/internal/domain"
	"ZenMobileService/internal/repository"
	"context"
)

type UsersService struct {
	repo repository.UsersRep
}

func NewUsersService(repo repository.UsersRep) *UsersService {
	return &UsersService{repo: repo}
}

func (us *UsersService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	userId, err := us.repo.CreateUser(ctx, user)
	return userId, err
}

func (us *UsersService) GetUser(ctx context.Context, userId int) (domain.User, error) {
	user, err := us.repo.GetUser(ctx, userId)
	return user, err
}
