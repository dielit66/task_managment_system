package repository

import (
	"context"
	"database/sql"

	"github.com/dielit66/task-management-system/internal/entities"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewPostgresUserRepostiry(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		db: db,
	}
}

func (r *UserPostgresRepository) CreateUser(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
}

func (r *UserPostgresRepository) GetByUserId(ctx context.Context, id int) (*entities.User, error) {
	return nil, nil
}
