package postgres

import (
	"context"
	"fmt"

	"github.com/WsDev69/producer-consumer/internal/config"
	"github.com/WsDev69/producer-consumer/internal/repository"
	"github.com/WsDev69/producer-consumer/pkg/persistence/postgres/tx"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Persistence struct {
	TaskRepository     repository.Task
	TaskSumsRepository repository.TaskSums
	Conn               *tx.Conn
}

func Init(ctx context.Context, cfg *config.Postgres) (*Persistence, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	connConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = int32(cfg.MaxOpenConn) //nolint:gosec // integer overflow conversion is not possible
	connConfig.MaxConnIdleTime = cfg.ConnMaxLifeTime

	conn, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(conn)

	p := &Persistence{}

	p.TaskRepository = NewTaskRepository(q)
	p.TaskSumsRepository = NewTaskSumsRepository(q)
	p.Conn = tx.NewConn(conn)

	return p, nil
}

func (p *Persistence) Close() error {
	// will not produce any errors
	return p.Conn.Close()
}
