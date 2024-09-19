package postgres

import (
	"context"
	"fmt"

	"producer-consumer/internal/config"
	"producer-consumer/internal/repository"
	"producer-consumer/pkg/persistence/postgres/tx"
	"producer-consumer/pkg/persistence/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Persistence struct {
	TaskRepository repository.Task
	Conn           *tx.Conn
}

func Init(ctx context.Context, cfg *config.Postgres) (*Persistence, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	connConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = int32(cfg.MaxOpenConn)
	connConfig.MaxConnIdleTime = cfg.ConnMaxLifeTime

	conn, err := pgxpool.NewWithConfig(ctx, connConfig)
	q := sqlc.New(conn)

	p := &Persistence{}

	p.TaskRepository = NewTaskRepository(q)
	p.Conn = tx.NewConn(conn)

	return p, nil
}

func (p *Persistence) Close() error {
	// will not produce any errors
	return p.Conn.Close()
}
