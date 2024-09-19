package main

import (
	"context"

	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"producer-consumer/cmd/common"
	"producer-consumer/internal/config"
	"producer-consumer/internal/domain/task"
	handler "producer-consumer/internal/handler/grpc"
	"producer-consumer/internal/midlleware"
	"producer-consumer/internal/monitoring"
	"producer-consumer/internal/server/grpc"
	"producer-consumer/pkg/persistence/postgres"
	prommethueswrap "producer-consumer/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/ratelimit"
	pb "google.golang.org/grpc"
)

func main() {
	common.ShowVersion()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Read("consumer")
	if err != nil {
		logger.Error("can't read config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	p, err := postgres.Init(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	taskSrv := task.NewService(p.TaskRepository, p.Conn)
	taskConsumer := task.NewConsumer(taskSrv)

	// init prometheus
	reg := prometheus.NewRegistry()
	reg.MustRegister(monitoring.TasksProcessedTotal)
	reg.MustRegister(monitoring.TasksDoneTotal)
	reg.MustRegister(monitoring.TasksValueTotal)
	reg.MustRegister(monitoring.TasksValueSum)

	go func() {
		err = prommethueswrap.New(reg).Serve(":3002")
		if err != nil {
			logger.Error("can't start prometheus server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	h := handler.NewHandler(taskConsumer)
	grpcServer := grpc.NewServer(cfg.GRPC, h)

	rl := ratelimit.New(10)
	wg := &sync.WaitGroup{}
	grpcServer.Serve(ctx, wg, pb.UnaryInterceptor(midlleware.UnaryServerInterceptor(rl)))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	s := <-sigCh
	logger.Info("got signal, attempting graceful shutdown", slog.String("signal", s.String()))
	cancel()
	wg.Done()
}
