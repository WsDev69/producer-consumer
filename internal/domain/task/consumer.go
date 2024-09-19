package task

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"producer-consumer/internal/domain/model"
	"producer-consumer/internal/monitoring"

	"time"

	"producer-consumer/pkg/persistence/sqlc"
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
		ID:    int32(task.ID),
		State: sqlc.TaskStateProcessing,
	}); err != nil {
		return fmt.Errorf("can't update task: %v ", err)
	}

	// simulate work
	time.Sleep(time.Duration(task.Value) * time.Millisecond)

	if err := c.taskSrv.UpdateTaskState(ctx, sqlc.UpdateTaskStateParams{
		ID:    int32(task.ID),
		State: sqlc.TaskStateDone,
	}); err != nil {
		return fmt.Errorf("can't update task: %v ", err)
	}

	monitoring.
		TasksValueSum.
		With(prometheus.Labels{"task_type": fmt.Sprintf("%d", task.Type)}).
		Add(float64(task.Value))

	return nil
}
