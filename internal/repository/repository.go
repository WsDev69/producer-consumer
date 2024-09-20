package repository

import (
	"context"

	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

type Task interface {
	GetTaskByID(ctx context.Context, id int32) (sqlc.Task, error)
	CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error)
	UpdateTask(ctx context.Context, params sqlc.UpdateTaskStateParams) error
	GetUnprocessedCount(ctx context.Context) (int64, error)
}

type TaskSums interface {
	GetTotalValueByTaskType(ctx context.Context, taskType int32) (sqlc.TaskSum, error)
}
