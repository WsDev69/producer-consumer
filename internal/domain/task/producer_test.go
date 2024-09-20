package task_test

import (
	"context"
	"testing"
	"time"

	"github.com/WsDev69/producer-consumer/internal/domain/task"
	"github.com/WsDev69/producer-consumer/internal/domain/task/mocks"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateAndSendTask(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockTask := new(mocks.Task)
	mockClient := new(mocks.Client)
	mockRandom := new(mocks.Random)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configure mock expectations
	mockTask.On("GetUnprocessedCount", mock.Anything).Return(int64(50), nil)
	mockTask.On("CreateTask", mock.Anything, mock.Anything).Return(sqlc.Task{
		ID:    1,
		Type:  1,
		Value: 1,
	}, nil)

	mockClient.On("Process", mock.Anything, mock.Anything).Return(nil)

	mockRandom.On("Int63n", mock.Anything).Return(int64(3))

	// Create a producer configuration
	producerConfig := task.ProducerConfig{
		MaxBackLog: 100,
		Config: task.Config{
			NumWorker:   5,
			MessageRate: 100,
		},
	}

	// Create serviceProducer
	producer := task.NewProducer(ctx, mockTask, mockClient, mockRandom, producerConfig)

	// Run the test in a separate goroutine (non-blocking)
	go func() {
		err := producer.GenerateAndSendTask(ctx)
		assert.NoError(t, err)
	}()

	// Give some time for the producer to run before canceling the context
	time.Sleep(1 * time.Second)
	cancel()

	// Assert that the mocks were called as expected
	mockTask.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}
