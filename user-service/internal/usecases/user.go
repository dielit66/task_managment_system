package usecases

import (
	"context"

	"github.com/dielit66/task-management-system/internal/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetByUserId(ctx context.Context, id int) (*entities.User, error)
}

type UserUseCase struct {
	repository UserRepository
}

func NewUserUseCase(r UserRepository) *UserUseCase {
	return &UserUseCase{
		repository: r,
	}
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, username string, email string, password string) error {
	user := &entities.User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}

	return uc.repository.CreateUser(ctx, user)

}
