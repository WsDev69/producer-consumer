package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"producer-consumer/pkg/persistence/postgres"
	"producer-consumer/pkg/persistence/postgres/test"
	"producer-consumer/pkg/persistence/sqlc"

	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.NotEqual(t, task, sqlc.Task{})
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
	assert.NoError(t, err)

	// Test GetTaskByID
	task, err := taskRepo.GetTaskByID(ctx, taskCreated.ID)

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(taskCreated, task))
}
