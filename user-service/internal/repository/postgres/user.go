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

func NewPostgresUserRepostiry(db *sqlx.DB, l logger.ILogger) *UserPostgresRepository {
	return &UserPostgresRepository{
		db:     db,
		logger: l,
	}
}

func (r *UserPostgresRepository) CreateUser(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	r.logger.Debug("Executing query", "query", query, "id", user.ID)
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		r.logger.Warn("Failed to create new user in repository", "username", user.Username, "email", user.Email, "err", err.Error())
		return app.Wrap(err, app.ErrNotFound, "Failed to create new user in repository")
	}

	return nil
}

func (r *UserPostgresRepository) GetByUserId(ctx context.Context, id int) (*entities.User, error) {
	user := entities.User{}
	query := "SELECT * FROM users WHERE id=$1"
	r.logger.Debug("Executing query", "query", query, "id", user.ID)
	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Failed to get user in repository", "id", user.ID)
			return nil, app.NewAppError(app.ErrNotFound, "User not found in repository", err)
		}
	}
	return &user, err
}
