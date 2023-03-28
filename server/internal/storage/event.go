package storage

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrEventExist    = errors.New("Event from worker already exist")
	ErrEventNotExist = errors.New("Event from worker not exist")
)

// struct event from workers
type Event struct {
	ID         uuid.UUID
	Address    string
	Command    string
	Hostname   string
	WorkerUuid uuid.UUID
	Timestamp  time.Time
}

func NewEvent(address, command, hostname string, timestamp time.Time, workerUuid uuid.UUID) *Event {
	id := uuid.New()

	return &Event{
		ID:         id,
		Address:    address,
		Command:    command,
		Hostname:   hostname,
		WorkerUuid: workerUuid,
		Timestamp:  timestamp,
	}
}
