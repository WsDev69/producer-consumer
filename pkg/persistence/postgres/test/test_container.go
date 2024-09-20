package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func SetupPostgresContainerWithMigration(t *testing.T,
	ctx context.Context,
	migrationPath string) (*pgxpool.Pool, func()) { // nolint:gocritic
	pool, cleanup := SetupPostgresContainer(t, ctx)

	// migration
	conn, err := sql.Open("postgres", pool.Config().ConnConfig.ConnString())
	require.NoError(t, err)

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	require.NoError(t, err)

	require.NoError(t, m.Up())

	return pool, cleanup
}

func SetupPostgresContainer(t *testing.T, ctx context.Context) (*pgxpool.Pool, func()) { //nolint:gocritic / test container
	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(_ string, port nat.Port) string {
			return fmt.Sprintf("host=localhost port=%s user=user password=secret dbname=testdb sslmode=disable", port.Port())
		}).WithStartupTimeout(60 * time.Second), //nolint:mnd / it's a helpful method that represents 60 seconds, aka 1 minute
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	// Get container host and port
	host, err := postgresC.Host(ctx)
	assert.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	// Build connection URL
	dsn := fmt.Sprintf("postgres://user:secret@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Create pgxpool.Pool
	pool, err := pgxpool.New(context.Background(), dsn)
	assert.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		pool.Close()
		_ = postgresC.Terminate(ctx)
	}

	return pool, cleanup
}
