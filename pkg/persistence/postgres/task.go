package postgres

import (
	"context"

	"producer-consumer/internal/repository"
	"producer-consumer/pkg/persistence/sqlc"
)

type TaskRepository struct {
	db *sqlc.Queries
}

func NewTaskRepository(db *sqlc.Queries) repository.Task {
	return &TaskRepository{db: db}
}

func (t *TaskRepository) GetTaskByID(ctx context.Context, id int32) (sqlc.Task, error) {
	task, err := t.db.GetTaskByID(ctx, id)
	if err != nil {
		return sqlc.Task{}, err
	}

	return task, nil
}

func (t *TaskRepository) CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error) {
	task, err := t.db.CreateTask(ctx, params)
	if err != nil {
		return sqlc.Task{}, err
	}

	return task, nil
}

func (t *TaskRepository) UpdateTask(ctx context.Context, params sqlc.UpdateTaskStateParams) error {
	err := t.db.UpdateTaskState(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskRepository) GetUnprocessedCount(ctx context.Context) (int64, error) {
	return t.db.GetUnprocessedCount(ctx)
}
