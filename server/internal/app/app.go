package app

import (
	internalstorage "server/internal/storage"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
	Sync() error
}

type App struct {
	logger  Logger
	storage Storage
}

type Storage interface {
	CreateEvent(e internalstorage.Event) error
	DeleteEvent(id uuid.UUID) error
	Find(id uuid.UUID) (*internalstorage.Event, error)
	FindAllEvents() ([]internalstorage.Event, error)
}

func NewApp(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}
