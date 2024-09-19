package grpc

import (
	"producer-consumer/internal/domain/task"
	pb "producer-consumer/internal/handler/grpc/gen/task"
)

type Handler struct {
	consumer task.Consumer
}

func NewHandler(consumer task.Consumer) pb.TaskServerServer {
	return &Handler{consumer: consumer}
}
