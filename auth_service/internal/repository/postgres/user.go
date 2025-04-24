package repository

import (
	"context"
	"database/sql"

	"github.com/dielit66/task-management-system/internal/entities"
	app "github.com/dielit66/task-management-system/internal/errors"
	"github.com/dielit66/task-management-system/internal/logger"
	"github.com/jmoiron/sqlx"
)

type UserPostgresRepository struct {
	db     *sqlx.DB
	logger logger.ILogger
}

func NewPostgresUserRepostiry(db *sqlx.DB, logger logger.ILogger) *UserPostgresRepository {
	return &UserPostgresRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserPostgresRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user := entities.User{}
	query := "SELECT * FROM users WHERE username=$1"
	r.logger.Debug("Executing query", "query", query, "username", username)
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("User not found", "username", username)
			return nil, app.NewAppError(app.ErrNotFound, "User not found in repository", err)
		}
		r.logger.Error("Database error", "username", username, "error", err.Error())
		return nil, app.Wrap(err, app.ErrInternal, "failed to fetch user")
	}

	return &user, err
}
