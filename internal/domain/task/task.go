package task

import (
	"context"
	"fmt"

	"github.com/WsDev69/producer-consumer/internal/domain/model"
	"github.com/WsDev69/producer-consumer/internal/repository"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres/tx"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

//go:generate mockery --name Task --output=mocks/
type Task interface {
	CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error)
	GetTask(ctx context.Context, id int32) (model.Task, error)
	UpdateTaskState(ctx context.Context, params sqlc.UpdateTaskStateParams) error
	GetUnprocessedCount(ctx context.Context) (int64, error)
	GetTotalSumByType(ctx context.Context, taskType int32) (sqlc.TaskSum, error)
}

type Config struct {
	MessageRate int
	NumWorker   int
}

type service struct {
	taskRepo    repository.Task
	taskSumRepo repository.TaskSums
	conn        *tx.Conn
}

func NewService(taskRepo repository.Task, taskSumRepo repository.TaskSums, conn *tx.Conn) Task {
	return &service{
		taskRepo:    taskRepo,
		taskSumRepo: taskSumRepo,
		conn:        conn,
	}
}

func (s service) CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error) {
	return tx.WithTx(ctx, s.conn, tx.Options{}, func(_ tx.DBTX) (sqlc.Task, error) {
		return s.taskRepo.CreateTask(ctx, params)
	})
}

func (s service) GetTask(ctx context.Context, id int32) (model.Task, error) {
	t, err := s.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return model.Task{}, fmt.Errorf("can't get task by id: %w", err)
	}

	return model.Task{
		ID:             t.ID,
		Type:           t.Type,
		Value:          t.Value,
		State:          model.GetTaskState(t.State),
		CreationTime:   t.CreationTime.Time,
		LastUpdateTime: t.LastUpdateTime.Time,
	}, err
}

func (s service) UpdateTaskState(ctx context.Context, params sqlc.UpdateTaskStateParams) error {
	return tx.WithTxExec(ctx, s.conn, tx.Options{}, func(_ tx.DBTX) error {
		return s.taskRepo.UpdateTask(ctx, params)
	})
}

func (s service) GetUnprocessedCount(ctx context.Context) (int64, error) {
	return s.taskRepo.GetUnprocessedCount(ctx)
}

func (s service) GetTotalSumByType(ctx context.Context, taskType int32) (sqlc.TaskSum, error) {
	return s.taskSumRepo.GetTotalValueByTaskType(ctx, taskType)
}
