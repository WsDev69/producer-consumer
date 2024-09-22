package task

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/WsDev69/producer-consumer/internal/domain/model"
	"github.com/WsDev69/producer-consumer/internal/monitoring"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"

	"github.com/prometheus/client_golang/prometheus"
)

type Consumer interface {
	ProcessTask(ctx context.Context, task model.TaskRequest) error
}

type consumer struct {
	taskSrv Task
}

func NewConsumer(taskSrv Task) Consumer {
	return &consumer{taskSrv: taskSrv}
}

func (c consumer) ProcessTask(ctx context.Context, task model.TaskRequest) error {
	// Task has been received. Set state to Processing
	if err := c.taskSrv.UpdateTaskState(ctx, sqlc.UpdateTaskStateParams{
		ID:    task.ID,
		State: sqlc.TaskStateProcessing,
	}); err != nil {
		return fmt.Errorf("can't update task: %w ", err)
	}

	monitoring.TasksProcessedTotal.Inc()

	// simulate work
	time.Sleep(time.Duration(task.Value) * time.Millisecond)

	if err := c.taskSrv.UpdateTaskState(ctx, sqlc.UpdateTaskStateParams{
		ID:    task.ID,
		State: sqlc.TaskStateDone,
	}); err != nil {
		return fmt.Errorf("can't update task: %w ", err)
	}

	monitoring.
		TasksValueSum.
		With(prometheus.Labels{"task_type": fmt.Sprintf("%d", task.Type)}).
		Add(float64(task.Value))

	monitoring.
		TasksValueTotal.
		With(prometheus.Labels{"task_type": fmt.Sprintf("%d", task.Type)}).
		Inc()

	t, err := c.taskSrv.GetTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("can't get task: %w", err)
	}

	totalSum, err := c.taskSrv.GetTotalSumByType(ctx, t.Type)
	if err != nil {
		return fmt.Errorf("can't get total sum: %w", err)
	}

	slog.Default().Info("Task successfully processed",
		slog.String("task", t.String()),
		slog.Int("total_sum", int(totalSum.TotalValue)),
	)

	return nil
}
