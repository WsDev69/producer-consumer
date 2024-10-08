package postgres

import (
	"context"

	"github.com/WsDev69/producer-consumer/internal/repository"
	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

type TaskSumsRepository struct {
	db *sqlc.Queries
}

func NewTaskSumsRepository(db *sqlc.Queries) repository.TaskSums {
	return &TaskSumsRepository{db: db}
}

func (t TaskSumsRepository) GetTotalValueByTaskType(ctx context.Context, taskType int32) (sqlc.TaskSum, error) {
	return t.db.GetTotalValueByTaskType(ctx, taskType)
}
