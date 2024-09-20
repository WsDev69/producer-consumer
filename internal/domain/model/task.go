package model

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/WsDev69/producer-consumer/pkg/persistence/sqlc"
)

type TaskRequest struct {
	ID    int32
	Type  int32
	Value int32
}

type TaskState string

const (
	TaskStateReceived   TaskState = "received"
	TaskStateProcessing TaskState = "processing"
	TaskStateDone       TaskState = "done"
)

type Task struct {
	ID             int32
	Type           int32
	Value          int32
	State          TaskState
	CreationTime   time.Time
	LastUpdateTime time.Time
}

func (t Task) String() string {
	tStr, err := json.Marshal(t)
	if err != nil {
		slog.Default().Error("can't Marshal task ",
			slog.String("err", err.Error()),
			slog.Int("taskID", int(t.ID)),
		)

		return ""
	}

	return string(tStr)
}

func GetTaskState(state sqlc.TaskState) TaskState {
	switch state {
	case sqlc.TaskStateReceived:
		return TaskStateReceived
	case sqlc.TaskStateProcessing:
		return TaskStateProcessing
	case sqlc.TaskStateDone:
		return TaskStateDone
	}

	return ""
}
