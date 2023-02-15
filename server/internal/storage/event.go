package storage

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
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
	Timestamp   time.Time
	Worker_UUID uuid.UUID
}

func NewEvent(hostname, command, description string, timestamp time.Time, worker_UUID uuid.UUID) (*Event, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:          id,
		Hostname:    hostname,
		Command:     command,
		Description: description,
		Timestamp:   timestamp,
		Worker_UUID: worker_UUID,
	}, nil
}
