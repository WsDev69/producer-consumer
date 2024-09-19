package tx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
)

// DBTX defines an interface implemented by both *pgx.Tx and *pgx.Pool
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Options struct {
	Level IsolationLevel
}

type IsolationLevel int

const (
	ReadCommited IsolationLevel = iota
)

type Result any

type Conn struct {
	pool *pgxpool.Pool
}

func NewConn(pool *pgxpool.Pool) *Conn {
	return &Conn{pool: pool}
}

// WithTx starts a transaction and manages its lifecycle
func WithTx[T Result](ctx context.Context, conn *Conn, opts Options, fn func(tx DBTX) (T, error)) (res T, err error) {
	tx, err := GetTx(ctx, conn, opts)
	if err != nil {
		return res, err
	}

	defer func() {
		err = FinishTx(tx, ctx, err)
	}()

	return fn(tx)
}

func WithTxExec(ctx context.Context, conn *Conn, opts Options, fn func(tx DBTX) error) (err error) {
	tx, err := GetTx(ctx, conn, opts)
	if err != nil {
		return fmt.Errorf("can't start transaction: %v ", err)
	}

	defer func() {
		err = FinishTx(tx, ctx, err)
	}()

	return fn(tx)
}

func FinishTx(tx pgx.Tx, ctx context.Context, err error) error {
	if tx == nil {
		slog.Default().Error("tx is nil, but shouldn't be")
		return fmt.Errorf("tx is nil, but shouldn't be")
	}
	if p := recover(); p != nil {
		errRb := tx.Rollback(ctx)
		err = errors.Join(err, errRb)
		panic(p)
	} else if err != nil {
		errRb := tx.Rollback(ctx)
		err = errors.Join(err, errRb)
	} else {
		errC := tx.Commit(ctx)
		err = errors.Join(err, errC)

	}
	return err
}

func GetTx(ctx context.Context, conn *Conn, opts Options) (pgx.Tx, error) {
	tx, err := conn.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: ToPgxLevel(opts.Level)})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ToPgxLevel(level IsolationLevel) pgx.TxIsoLevel {
	switch level {
	case ReadCommited:
		return pgx.ReadCommitted
	}

	// default value
	return pgx.ReadCommitted
}

func (c *Conn) Close() error {
	c.pool.Close()
	return nil
}
