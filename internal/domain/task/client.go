package task

import (
	"context"
	"google.golang.org/grpc"

	"github.com/WsDev69/producer-consumer/internal/domain/model"
	"github.com/WsDev69/producer-consumer/internal/handler/grpc/gen/task"
)

//go:generate mockery --name Client --output=mocks/
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
		Id:    int64(taskRequest.ID),
		Type:  taskRequest.Type,
		Value: taskRequest.Value,
	}, grpc.EmptyCallOption{})
	if err != nil {
		return err
	}

	return nil
}
