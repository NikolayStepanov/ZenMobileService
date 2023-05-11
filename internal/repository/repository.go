package repository

import (
	"ZenMobileService/internal/domain"
	"context"
	"github.com/jackc/pgx/v5"
)

type UsersRep interface {
	CreateUsersTable(ctx context.Context) error
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUser(ctx context.Context, userID int) (domain.User, error)
}

type Repository struct {
	UsersRep
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{UsersRep: NewUserPostgres(db)}
}
