package service

import (
	"GeekTask/fifthWeek/internal/domain"
	"context"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}
