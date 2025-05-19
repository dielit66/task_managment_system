package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/dielit66/task-management-system/internal/entities"
	"github.com/dielit66/task-management-system/internal/logger"
	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db     *sqlx.DB
	logger logger.ILogger
}

func NewTaskRepository(db *sqlx.DB, l logger.ILogger) *TaskRepository {
	return &TaskRepository{
		db:     db,
		logger: l,
	}
}

func (r *TaskRepository) GetById(ctx context.Context, id int) (*entities.Task, error) {
	query := "SELECT * FROM tasks  WHERE id = $1"
	task := entities.Task{}
	err := r.db.Get(&task, query, id)
	if err != nil {
		return nil, err
	}

	return &task, nil

}

func (r *TaskRepository) GetAllByUserId(ctx context.Context, id int) ([]*entities.Task, error) {
	query := "SELECT * FROM tasks WHERE user_id = $1"
	var tasks []*entities.Task
	err := r.db.SelectContext(ctx, &tasks, query, id)
	if err != nil {
		return nil, err
	}

	return tasks, nil

}

func (r *TaskRepository) Create(ctx context.Context, t *entities.Task) error {
	query := "INSERT INTO tasks (title, description, deadline, user_id, status_id) VALUES ($1,$2,$3,$4,$5)"
	fmt.Println(t.UserId)
	row := r.db.QueryRowContext(ctx, query, t.Title, t.Description, t.Deadline, t.UserId, 1)

	if row.Err() != nil {
		r.logger.Error("Error while inserting new task in repository", "title", t.Title, "desc", t.Description, "deadline", t.Deadline, "user_id", t.UserId)
		fmt.Println("qqqqqqqqq")
		return row.Err()
	}

	return nil
}

func (r *TaskRepository) Update(ctx context.Context, t *entities.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, deadline = $3 WHERE id = $4"
	result, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Deadline, t.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %d", t.ID)
	}

	log.Printf("Updated %d rows", rowsAffected)

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE from tasks WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	r.logger.Info("Deleted %d rows", rowsAffected)

	return err

}
