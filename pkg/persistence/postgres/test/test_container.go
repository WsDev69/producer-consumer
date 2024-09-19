package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/stretchr/testify/assert"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func SetupPostgresContainerWithMigration(t *testing.T, ctx context.Context, migrationPath string) (*pgxpool.Pool, func()) {
	pool, cleanup := SetupPostgresContainer(t, ctx)

	// migration
	conn, err := sql.Open("postgres", pool.Config().ConnConfig.ConnString())
	assert.NoError(t, err)

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	assert.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	assert.NoError(t, err)

	assert.NoError(t, m.Up())

	return pool, cleanup
}

func SetupPostgresContainer(t *testing.T, ctx context.Context) (*pgxpool.Pool, func()) {
	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("host=localhost port=%s user=user password=secret dbname=testdb sslmode=disable", port.Port())
		}).WithStartupTimeout(60 * time.Second),
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
		postgresC.Terminate(ctx)
	}

	return pool, cleanup
}
