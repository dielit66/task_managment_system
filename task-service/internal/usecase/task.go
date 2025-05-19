package usecase

import (
	"context"

	"github.com/dielit66/task-management-system/internal/entities"
	"github.com/dielit66/task-management-system/internal/logger"
)

type UserRepository interface {
	GetById(ctx context.Context, id int) (*entities.Task, error)
	GetAllByUserId(ctx context.Context, id int) ([]*entities.Task, error)
	Create(ctx context.Context, t *entities.Task) error
	Update(ctx context.Context, t *entities.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskUsecase struct {
	repository UserRepository
	logger     logger.ILogger
}

func NewTaskUsecase(r UserRepository, l logger.ILogger) *TaskUsecase {
	return &TaskUsecase{
		repository: r,
		logger:     l,
	}
}

func (uc *TaskUsecase) GetAllById(ctx context.Context, id int) ([]*entities.Task, error) {
	tasks, err := uc.repository.GetAllByUserId(ctx, id)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (uc *TaskUsecase) GetById(ctx context.Context, id int) (*entities.Task, error) {
	task, err := uc.repository.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (uc *TaskUsecase) Create(ctx context.Context, t *entities.CreateTaskDto) error {

	task := entities.Task{
		UserId: t.UserID,
		Title:  t.Title,

		Description: t.Description,
		Deadline:    t.Deadline,
	}

	err := uc.repository.Create(ctx, &task)

	if err != nil {
		return err
	}

	return nil
}

func (uc *TaskUsecase) Update(ctx context.Context, t *entities.Task) error {
	err := uc.repository.Update(ctx, t)

	if err != nil {
		return err
	}

	return nil
}

func (uc *TaskUsecase) Delete(ctx context.Context, id int) error {
	err := uc.repository.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
