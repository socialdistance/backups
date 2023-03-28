package storage

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID         uuid.UUID
	Command    string
	WorkerUuid uuid.UUID
	Timestamp  time.Time
}

func NewTask(command string, workerUuid uuid.UUID, timestamp time.Time) *Task {
	id := uuid.New()

	return &Task{
		ID:         id,
		Command:    command,
		WorkerUuid: workerUuid,
		Timestamp:  timestamp,
	}
}
