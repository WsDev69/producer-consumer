package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/WsDev69/producer-consumer/internal/config"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres/test"
)

func TestInit(t *testing.T) {
	// Setup PostgreSQL container
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainerWithMigration(t, ctx, "file://../../../sql/migrations/")
	defer cleanup()

	cfg := &config.Postgres{
		Host:            pool.Config().ConnConfig.Host,
		Port:            int(pool.Config().ConnConfig.Port),
		User:            "user",
		Password:        "secret",
		DB:              "testdb",
		MaxOpenConn:     10,
		ConnMaxLifeTime: 30,
	}

	// Call Init to create a new Persistence instance
	persistence, err := postgres.Init(ctx, cfg)
	require.NoError(t, err)
	assert.NotNil(t, persistence)
	assert.NotNil(t, persistence.TaskRepository)

	// Clean up
	err = persistence.Close()
	assert.NoError(t, err)
}
