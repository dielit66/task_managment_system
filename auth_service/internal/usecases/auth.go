package usecases

import (
	"context"
	"errors"

	"github.com/dielit66/task-management-system/internal/auth"
	"github.com/dielit66/task-management-system/internal/entities"
	app "github.com/dielit66/task-management-system/internal/errors"
	"github.com/dielit66/task-management-system/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	GetUserByUsername(context.Context, string) (*entities.User, error)
}

type AuthUseCase struct {
	repository IUserRepository
	jwtService *auth.JWTService
	logger     logger.ILogger
}

func NewAuthUseCase(repository IUserRepository, jwtService *auth.JWTService, logger logger.ILogger) *AuthUseCase {
	return &AuthUseCase{
		repository: repository,
		jwtService: jwtService,
		logger:     logger,
	}
}

type LoginResult struct {
	IsSuccessed bool
}

func (uc *AuthUseCase) LoginUser(ctx context.Context, username string, password string) (*LoginResult, error) {
	user, err := uc.repository.GetUserByUsername(ctx, username)

	if err != nil {
		var appErr *app.AppError
		if errors.As(err, &appErr) {
			if appErr.Type == app.ErrNotFound {
				uc.logger.Warn("User not found in usecase", "username", username)
				return &LoginResult{
					IsSuccessed: false,
				}, err
			}
		}
		uc.logger.Error("Failed to fetch user", "username", username, "error", err.Error())
		return nil, app.Wrap(err, app.ErrInternal, "failed to fetch user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			uc.logger.Info("Missmatch hash and password in usecase", "username", username)
			return &LoginResult{
				IsSuccessed: false,
			}, nil
		}
		uc.logger.Error("Failed to compare with hash password", "username", username, "error", err.Error())
		return nil, app.Wrap(err, app.ErrInternal, "failed to compare with hash password")
	}

	return &LoginResult{
		IsSuccessed: true,
	}, nil
}
