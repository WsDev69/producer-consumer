package model

import (
	"time"
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
