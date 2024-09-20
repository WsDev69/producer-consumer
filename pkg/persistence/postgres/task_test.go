package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres/test"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_CreateTask(t *testing.T) {
	// Setup PostgreSQL container
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainerWithMigration(t, ctx, "file://../../../sql/migrations/")
	defer cleanup()

	// Initialize sqlc.Queries (repository) and TaskRepository
	queries := sqlc.New(pool)
	taskRepo := postgres.NewTaskRepository(queries)

	task, err := taskRepo.CreateTask(ctx, sqlc.CreateTaskParams{
		Type:  1,
		Value: 2,
	})

	require.NoError(t, err)
	assert.NotNil(t, task)
	assert.NotEqual(t, task, sqlc.Task{}) //nolint:testifylint // we don't care about the actual valus
}

func TestTaskRepository_GetTaskByID(t *testing.T) {
	// Setup PostgreSQL container
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainerWithMigration(t, ctx, "file://../../../sql/migrations/")
	defer cleanup()

	// Initialize sqlc.Queries (repository) and TaskRepository
	queries := sqlc.New(pool)
	taskRepo := postgres.NewTaskRepository(queries)

	taskCreated, err := taskRepo.CreateTask(ctx, sqlc.CreateTaskParams{
		Type:  1,
		Value: 2,
	})
	require.NoError(t, err)

	// Test GetTaskByID
	task, err := taskRepo.GetTaskByID(ctx, taskCreated.ID)

	require.NoError(t, err)
	assert.True(t, reflect.DeepEqual(taskCreated, task))
}
