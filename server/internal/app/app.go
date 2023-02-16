package app

import (
	"server/internal/storage"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Sync() error
	With(fields ...zap.Field) *zap.Logger
}

type App struct {
	logger  Logger
	storage Storage
}

type Storage interface {
	CreateEvent(e storage.Event) error
	DeleteEvent(id uuid.UUID) error
	Find(id uuid.UUID) (*storage.Event, error)
	FindAllEvents() ([]storage.Event, error)
}

func NewApp(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}
