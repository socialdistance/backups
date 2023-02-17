package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventExist    = errors.New("Event from worker already exist")
	ErrEventNotExist = errors.New("Event from worker not exist")
)

// struct event from workers
type Event struct {
	ID          uuid.UUID
	Hostname    string
	Command     string
	Description string
	Worker_UUID uuid.UUID
	Timestamp   time.Time
}

func NewEvent(hostname, command, description string, timestamp time.Time, worker_UUID uuid.UUID) *Event {
	id := uuid.New()

	return &Event{
		ID:          id,
		Hostname:    hostname,
		Command:     command,
		Description: description,
		Timestamp:   timestamp,
		Worker_UUID: worker_UUID,
	}
}
