package task_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/WsDev69/producer-consumer/internal/domain/model"
	"github.com/WsDev69/producer-consumer/internal/domain/task"
	"github.com/WsDev69/producer-consumer/internal/domain/task/mocks"
	"github.com/WsDev69/producer-consumer/internal/monitoring"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

func TestProcessTask_Success(t *testing.T) {
	// Create mock for Task interface
	mockTask := new(mocks.Task)

	// Create a Consumer instance with the mocked Task service
	consumer := task.NewConsumer(mockTask)

	ctx := context.Background()
	taskRequest := model.TaskRequest{
		ID:    1,
		Type:  2,
		Value: 100,
	}

	// Define the mock behavior for UpdateTaskState
	mockTask.On("UpdateTaskState", mock.Anything, sqlc.UpdateTaskStateParams{
		ID:    int32(1),
		State: sqlc.TaskStateProcessing,
	}).Return(nil)

	mockTask.On("UpdateTaskState", mock.Anything, sqlc.UpdateTaskStateParams{
		ID:    int32(1),
		State: sqlc.TaskStateDone,
	}).Return(nil)

	mockTask.On("GetTask", mock.Anything, mock.Anything).Return(model.Task{
		ID:             int32(1),
		Type:           2,
		Value:          100,
		State:          model.TaskStateDone,
		CreationTime:   time.Now(),
		LastUpdateTime: time.Now(),
	}, nil)

	mockTask.On("GetTotalSumByType", mock.Anything, mock.Anything).Return(sqlc.TaskSum{
		TaskType:   2,
		TotalValue: 100,
	}, nil)

	// Call ProcessTask
	err := consumer.ProcessTask(ctx, taskRequest)

	// Assertions
	require.NoError(t, err)
	mockTask.AssertExpectations(t)

	// Check Prometheus metric has been updated
	expected := `
		# HELP task_consumer_tasks_value_sum Total sum of the 'value' field for each task type.
        # TYPE task_consumer_tasks_value_sum gauge
        task_consumer_tasks_value_sum{task_type="2"} 100
	`

	assert.NoError(t, testutil.CollectAndCompare(monitoring.TasksValueSum, strings.NewReader(expected), "task_consumer_tasks_value_sum"))
}

func TestProcessTask_UpdateTaskStateFailure(t *testing.T) {
	// Create mock for Task interface
	mockTask := new(mocks.Task)

	// Create a Consumer instance with the mocked Task service
	consumer := task.NewConsumer(mockTask)

	ctx := context.Background()
	taskRequest := model.TaskRequest{
		ID:    1,
		Type:  2,
		Value: 100,
	}

	// Simulate error during UpdateTaskState (Processing)
	mockTask.On("UpdateTaskState", mock.Anything, sqlc.UpdateTaskStateParams{
		ID:    int32(1),
		State: sqlc.TaskStateProcessing,
	}).Return(errors.New("database error"))

	// Call ProcessTask
	err := consumer.ProcessTask(ctx, taskRequest)

	// Assertions
	require.Error(t, err)
	require.EqualError(t, err, "can't update task: database error ")

	// Ensure no further updates happen
	mockTask.AssertCalled(t, "UpdateTaskState", mock.Anything, sqlc.UpdateTaskStateParams{
		ID:    int32(1),
		State: sqlc.TaskStateProcessing,
	})
	mockTask.AssertNotCalled(t, "UpdateTaskState", mock.Anything, sqlc.UpdateTaskStateParams{
		ID:    int32(1),
		State: sqlc.TaskStateDone,
	})
}
