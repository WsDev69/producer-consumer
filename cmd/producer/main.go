package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/WsDev69/producer-consumer/cmd/common"
	"github.com/WsDev69/producer-consumer/internal/config"
	"github.com/WsDev69/producer-consumer/internal/domain/task"
	taskgrpc "github.com/WsDev69/producer-consumer/internal/handler/grpc/gen/task"
	"github.com/WsDev69/producer-consumer/internal/monitoring"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres"
	prommethueswrap "github.com/WsDev69/producer-consumer/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	common.ShowVersion()
	common.ExposePprof("localhost:1378")
	//common.RunCPUProf()
	//common.MEMProf()

	// init context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// init config
	cfg, err := config.ReadProducer("producer")
	if err != nil {
		log.Fatalf("cannot read producer config: %v", err) //nolint:gocritic // nothing to cancel, ok to use fatal
	}

	logger := common.Logger(cfg.LogLevel, cfg.OutPut)

	// init db (postgres)
	p, err := postgres.Init(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	// init grpc client
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Default().Error("can't establish connection", slog.String("err", err.Error()))
		os.Exit(1)
	}

	client := taskgrpc.NewTaskServerClient(conn)
	taskGrpcClient := task.NewGrpcClient(client)

	// init prometheus
	reg := prometheus.NewRegistry()
	reg.MustRegister(monitoring.TasksProducedTotal)

	go func() {
		err = prommethueswrap.New(reg).Serve(":3001")
		if err != nil {
			logger.Error("can't start prometheus server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	// create task service
	taskSrv := task.NewService(p.TaskRepository, p.TaskSumsRepository, p.Conn)

	// init producer
	producerConfig := task.ProducerConfig{
		MaxBackLog: cfg.MaxBacklog,
		Config: task.Config{
			MessageRate: cfg.MessageRate,
			NumWorker:   cfg.NumWorker,
		},
	}

	r := task.NewRandom()
	producer := task.NewProducer(ctx, taskSrv, taskGrpcClient, r, producerConfig)

	logger.Info("producer started")

	// start producing
	go func() {
		err = producer.GenerateAndSendTask(ctx)
		if err != nil {
			logger.Error("producer finished with error", slog.String("err", err.Error()))

			os.Exit(1)
		}

		logger.Info("producer finished ")
		os.Exit(0)
	}()

	// wait for sigterm
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	s := <-sigCh
	logger.Info("got signal, attempting graceful shutdown", slog.String("signal", s.String()))
	if err = conn.Close(); err != nil {
		logger.Error("can't close connection", slog.String("error", err.Error()))
	}

	logger.Info("producer finished")
}
