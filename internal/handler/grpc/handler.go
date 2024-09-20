package grpc

import (
	"github.com/WsDev69/producer-consumer/internal/domain/task"
	pb "github.com/WsDev69/producer-consumer/internal/handler/grpc/gen/task"
)

type Handler struct {
	consumer task.Consumer
}

func NewHandler(consumer task.Consumer) pb.TaskServerServer {
	return &Handler{consumer: consumer}
}
