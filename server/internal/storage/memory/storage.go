package storage

import (
	"server/internal/storage"
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func NewMemory() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.WorkerUuid]; ok {
		return storage.ErrEventExist
	}

	s.events[e.WorkerUuid] = e

	return nil
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotExist
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) UpdateCommand(e storage.Event) error {
	panic("implement me")
}

func (s *Storage) Find(workerUuid uuid.UUID) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, ok := s.events[workerUuid]; ok {
		return &event, nil
	}

	return nil, storage.ErrEventNotExist
}

func (s *Storage) FindAllEvents() ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make([]storage.Event, 0, len(s.events))

	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}
