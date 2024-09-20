package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/WsDev69/producer-consumer/internal/domain/model"
	"github.com/WsDev69/producer-consumer/internal/handler/grpc/gen/task"
	"github.com/WsDev69/producer-consumer/internal/monitoring"
)

func (s Handler) Process(ctx context.Context, in *task.TaskRequest) (*task.TaskResponse, error) {
	err := s.consumer.ProcessTask(ctx, model.TaskRequest{
		ID:    int32(in.Id), //nolint:gosec // for the current implementation, we will never overflow the int32
		Type:  in.Type,
		Value: in.Value,
	})
	if err != nil {
		slog.Default().Error(fmt.Sprintf("failed to process task: %v", err))
		return nil, err
	}

	monitoring.TasksProcessedTotal.Inc()

	return &task.TaskResponse{}, nil
}
