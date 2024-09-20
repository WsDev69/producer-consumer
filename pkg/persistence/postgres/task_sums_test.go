package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/WsDev69/producer-consumer/internal/repository"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres/test"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

func TestTaskSunsRepository_GetTotalValueByTaskType(t *testing.T) {
	// Setup PostgreSQL container
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainerWithMigration(t, ctx, "file://../../../sql/migrations/")
	defer cleanup()

	// Initialize sqlc.Queries (repository) and TaskRepository
	queries := sqlc.New(pool)
	taskRepo := postgres.NewTaskRepository(queries)

	taskSumsRepo := postgres.NewTaskSumsRepository(queries)

	createAndCheck(t, ctx, taskRepo, taskSumsRepo, 1)
	createAndCheck(t, ctx, taskRepo, taskSumsRepo, 2)
}

func createAndCheck(t *testing.T,
	ctx context.Context,
	taskRepo repository.Task,
	taskSums repository.TaskSums,
	multiplier int32) {
	testTypes := int32(3)
	for i := range testTypes {
		task, err := taskRepo.CreateTask(ctx, sqlc.CreateTaskParams{
			Type:  i,
			Value: i,
		})
		require.NoError(t, err)
		require.NotNil(t, task)
		assert.NotEqual(t, task, sqlc.Task{}) //nolint:testifylint // we are not interested in actual values

		err = taskRepo.UpdateTask(ctx, sqlc.UpdateTaskStateParams{
			ID:    task.ID,
			State: "done",
		})
		require.NoError(t, err)
	}

	for i := range testTypes {
		totalSum, err := taskSums.GetTotalValueByTaskType(ctx, i)
		require.NoError(t, err)
		assert.Equal(t, i, totalSum.TaskType)
		assert.Equal(t, i*multiplier, totalSum.TotalValue)
	}
}
