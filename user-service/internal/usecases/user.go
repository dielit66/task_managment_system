package usecases

import (
	"context"
	"errors"

	"github.com/dielit66/task-management-system/internal/entities"
	app "github.com/dielit66/task-management-system/internal/errors"
	"github.com/dielit66/task-management-system/internal/logger"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetByUserId(ctx context.Context, id int) (*entities.User, error)
}

type UserUseCase struct {
	repository IUserRepository
	logger     logger.ILogger
}

func NewUserUseCase(r IUserRepository, l logger.ILogger) *UserUseCase {
	return &UserUseCase{
		repository: r,
		logger:     l,
	}
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, username string, email string, password string) error {

	hashPass, err := HashPassword(password)

	if err != nil {
		uc.logger.Warn("Failed to hash password in usecase", "err", err.Error())
		return err
	}

	user := &entities.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashPass,
	}

	return uc.repository.CreateUser(ctx, user)

}

func (uc *UserUseCase) GetUser(ctx context.Context, id int) (*entities.User, error) {

	user, err := uc.repository.GetByUserId(ctx, id)

	if err != nil {
		var appErr *app.AppError
		if errors.As(err, &appErr) {
			if appErr.Type == app.ErrNotFound {
				return nil, err
			}
		}
		return nil, app.Wrap(err, app.ErrInternal, "failed to fetch user in usecase")
	}

	return user, nil
}
