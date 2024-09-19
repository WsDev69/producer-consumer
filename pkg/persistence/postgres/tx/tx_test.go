package tx_test

import (
	"context"
	"fmt"
	"testing"

	"producer-consumer/pkg/persistence/postgres/test"
	"producer-consumer/pkg/persistence/postgres/tx"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestWithTx_Commit(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainer(t, ctx)
	defer cleanup()

	conn := tx.NewConn(pool)

	// Create table
	_, err := pool.Exec(ctx, `CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)

	// Run transaction
	ok, err := tx.WithTx(ctx, conn, tx.Options{Level: tx.ReadCommited}, func(tx tx.DBTX) (bool, error) {
		_, err := tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
		return true, err
	})
	assert.NoError(t, err)
	assert.True(t, ok)

	// Verify the row was committed
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestWithTx_Rollback(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainer(t, ctx)
	defer cleanup()

	conn := tx.NewConn(pool)

	// Create table
	_, err := pool.Exec(ctx, `CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)

	// Run transaction that will fail and rollback
	ok, err := tx.WithTx(ctx, conn, tx.Options{Level: tx.ReadCommited}, func(tx tx.DBTX) (bool, error) {
		_, err := tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Bob")
		if err != nil {
			return false, err
		}

		// Simulate an error to trigger rollback
		return false, fmt.Errorf("trigger rollback")
	})
	assert.Error(t, err)
	assert.False(t, ok)

	// Verify the row was not committed
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestWithTxExec_Commit(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainer(t, ctx)
	defer cleanup()

	conn := tx.NewConn(pool)

	// Create table
	_, err := pool.Exec(ctx, `CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)

	// Run transaction
	err = tx.WithTxExec(ctx, conn, tx.Options{Level: tx.ReadCommited}, func(tx tx.DBTX) error {
		_, err := tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
		return err
	})
	assert.NoError(t, err)

	// Verify the row was committed
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestWithTxExec_Rollback(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := test.SetupPostgresContainer(t, ctx)
	defer cleanup()

	conn := tx.NewConn(pool)

	// Create table
	_, err := pool.Exec(ctx, `CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)`)
	assert.NoError(t, err)

	// Run transaction that will fail and rollback
	err = tx.WithTxExec(ctx, conn, tx.Options{Level: tx.ReadCommited}, func(tx tx.DBTX) error {
		_, err := tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Bob")
		if err != nil {
			return err
		}

		// Simulate an error to trigger rollback
		return fmt.Errorf("trigger rollback")
	})
	assert.Error(t, err)

	// Verify the row was not committed
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
