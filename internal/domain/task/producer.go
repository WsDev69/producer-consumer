package task

import (
	"context"
	"github.com/WsDev69/producer-consumer/internal/monitoring"
	"log/slog"
	"time"

	"github.com/WsDev69/producer-consumer/internal/domain/model"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

const millisecondInSeconds = 1000

type Producer interface {
	GenerateAndSendTask(ctx context.Context) error
}

type ProducerConfig struct {
	MaxBackLog int64
	Config
}

type serviceProducer struct {
	taskSrv Task
	client  Client
	intRand Random

	taskChan     chan model.Task
	producerTick <-chan time.Time

	maxBackLog int64
}

// NewProducer returns a work to produce tasks
func NewProducer(ctx context.Context,
	taskSrv Task,
	client Client,
	intRandom Random,
	config ProducerConfig) Producer {
	s := &serviceProducer{
		taskSrv: taskSrv,

		client:  client,
		intRand: intRandom,

		maxBackLog: config.MaxBackLog,
	}

	s.taskChan = make(chan model.Task, config.NumWorker)
	for i := range config.NumWorker {
		go s.worker(ctx, i+1, s.taskChan)
	}

	s.producerTick = time.Tick(time.Duration(millisecondInSeconds/config.MessageRate) * time.Millisecond)

	return s
}

func (s serviceProducer) GenerateAndSendTask(ctx context.Context) error {
	defer func() {
		close(s.taskChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.producerTick:
			unprocessedCount, err := s.taskSrv.GetUnprocessedCount(ctx)
			if err != nil {
				return err
			}
			if unprocessedCount >= s.maxBackLog {
				slog.Default().Info("MaxBacklog reached")
				return nil
			}

			t, err := s.taskSrv.CreateTask(ctx, s.getTask())
			if err != nil {
				return err
			}

			s.taskChan <- model.Task{
				ID:    t.ID,
				Type:  t.Type,
				Value: t.Value,
			}
		}
	}
}

// getTask generate a random task
func (s serviceProducer) getTask() sqlc.CreateTaskParams {
	return sqlc.CreateTaskParams{
		Type:  int32(s.intRand.Int63n(10)), //nolint:gosec,mnd // no magic numbers
		Value: int32(s.intRand.Int63n(64)), //nolint:gosec,mnd // no magic numbers
	}
}

// Worker function that consumes tasks from the channel
func (s serviceProducer) worker(ctx context.Context, id int, tasks <-chan model.Task) {
	log := slog.Default()
	for {
		select {
		case <-ctx.Done():
			log.Warn("context canceled")
			return
		case t, ok := <-tasks:
			if !ok {
				log.Warn("channel closed")
				return
			}

			log.Debug("Worker processing task", slog.Int("workerID", id), slog.Int("taskID", int(t.ID)))
			if err := s.client.Process(ctx, model.TaskRequest{
				ID:    t.ID,
				Type:  t.Type,
				Value: t.Value,
			}); err != nil {
				log.Error("Failed to process task", slog.String("err", err.Error()))
			}
			monitoring.
				TasksProducedTotal.
				Inc()
			log.Info("task successfully processed", slog.Int("workerID", id), slog.Int("taskID", int(t.ID)))
		}
	}
}
