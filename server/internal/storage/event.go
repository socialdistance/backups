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
	Address     string
	Command     string
	Hostname    string
	Worker_UUID uuid.UUID
	Timestamp   time.Time
}

func NewEvent(address, command, hostname string, timestamp time.Time, worker_UUID uuid.UUID) *Event {
	id := uuid.New()

	return &Event{
		ID:          id,
		Address:     address,
		Command:     command,
		Hostname:    hostname,
		Timestamp:   timestamp,
		Worker_UUID: worker_UUID,
	}
}
