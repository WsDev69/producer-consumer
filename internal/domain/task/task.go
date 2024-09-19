package task

import (
	"context"
	"producer-consumer/internal/repository"
	"producer-consumer/pkg/persistence/postgres/tx"
	"producer-consumer/pkg/persistence/sqlc"
)

type Task interface {
	CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error)
	GetTask(ctx context.Context, id int32) (sqlc.Task, error)
	UpdateTaskState(ctx context.Context, params sqlc.UpdateTaskStateParams) error
	GetUnprocessedCount(ctx context.Context) (int64, error)
}

type Config struct {
	MessageRate int
	NumWorker   int
}

type service struct {
	taskRepo repository.Task
	conn     *tx.Conn
}

func NewService(taskRepo repository.Task, conn *tx.Conn) Task {
	return &service{taskRepo: taskRepo,
		conn: conn,
	}
}

func (s service) CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error) {
	return tx.WithTx(ctx, s.conn, tx.Options{}, func(tx tx.DBTX) (sqlc.Task, error) {
		return s.taskRepo.CreateTask(ctx, params)
	})
}

func (s service) GetTask(ctx context.Context, id int32) (sqlc.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s service) UpdateTaskState(ctx context.Context, params sqlc.UpdateTaskStateParams) error {
	return tx.WithTxExec(ctx, s.conn, tx.Options{}, func(tx tx.DBTX) error {
		return s.taskRepo.UpdateTask(ctx, params)
	})
}

func (s service) GetUnprocessedCount(ctx context.Context) (int64, error) {
	return s.taskRepo.GetUnprocessedCount(ctx)
}
