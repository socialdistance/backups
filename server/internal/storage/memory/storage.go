package storage

import (
	"server/internal/storage"
	"sync"

	"github.com/gofrs/uuid"
)

type Storage struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]storage.Information
}

func NewMemory() *Storage {
	return &Storage{
		tasks: make(map[uuid.UUID]storage.Information),
	}
}
