package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"producer-consumer/internal/domain/model"
	"producer-consumer/internal/monitoring"

	"producer-consumer/internal/handler/grpc/gen/task"
)

func (s Handler) Process(ctx context.Context, in *task.TaskRequest) (*task.TaskResponse, error) {
	err := s.consumer.ProcessTask(ctx, model.TaskRequest{
		ID:    in.Id,
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
