package task

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"producer-consumer/internal/domain/model"
	"producer-consumer/internal/handler/grpc/gen/task"
)

type Client interface {
	Process(ctx context.Context, taskRequest model.TaskRequest) error
}

type grpcClient struct {
	taskServerClient task.TaskServerClient
}

func NewGrpcClient(taskServerClient task.TaskServerClient) Client {
	return &grpcClient{taskServerClient: taskServerClient}
}

func (c *grpcClient) Process(ctx context.Context, taskRequest model.TaskRequest) error {
	_, err := c.taskServerClient.Process(ctx, &task.TaskRequest{
		Id:    taskRequest.ID,
		Type:  taskRequest.Type,
		Value: taskRequest.Value,
	}, grpc.EmptyCallOption{})
	slog.Default().Info("Got response from consumer")
	if err != nil {
		return err
	}

	return nil
}
