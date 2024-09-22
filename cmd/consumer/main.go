package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/WsDev69/producer-consumer/cmd/common"
	"github.com/WsDev69/producer-consumer/internal/config"
	"github.com/WsDev69/producer-consumer/internal/domain/task"
	handler "github.com/WsDev69/producer-consumer/internal/handler/grpc"
	"github.com/WsDev69/producer-consumer/internal/midlleware"
	"github.com/WsDev69/producer-consumer/internal/monitoring"
	"github.com/WsDev69/producer-consumer/internal/server/grpc"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres"
	prommethueswrap "github.com/WsDev69/producer-consumer/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/ratelimit"
	pb "google.golang.org/grpc"
)

func main() {
	common.ShowVersion()
	common.ExposePprof("localhost:1377")
	//common.RunCPUProf()
	//common.MEMProf()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ReadConsumer("consumer")
	if err != nil {
		log.Fatalf("can't read config: %v", err) //nolint:gocritic // nothing to cancel, we can call fatal
	}

	logger := common.Logger(cfg.LogLevel, cfg.OutPut)

	p, err := postgres.Init(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	taskSrv := task.NewService(p.TaskRepository, p.TaskSumsRepository, p.Conn)
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

	rl := ratelimit.New(cfg.MessageRate)
	wg := &sync.WaitGroup{}
	grpcServer.Serve(ctx, wg, pb.UnaryInterceptor(midlleware.UnaryServerInterceptor(rl)))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	s := <-sigCh
	logger.Info("got signal, attempting graceful shutdown", slog.String("signal", s.String()))
	cancel()
	wg.Done()
}
