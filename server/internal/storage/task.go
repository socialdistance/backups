package storage

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Command     string
	Worker_UUID uuid.UUID
	Timestamp   time.Time
}

func NewTask(command string, worker_uuid uuid.UUID, timestamp time.Time) *Task {
	id := uuid.New()

	return &Task{
		ID:          id,
		Command:     command,
		Worker_UUID: worker_uuid,
		Timestamp:   timestamp,
	}
}
