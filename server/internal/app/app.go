package app

import (
	internalstorage "server/internal/storage"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
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
	CreateEvent(e internalstorage.Event) (*internalstorage.Event, error)
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

func (a *App) CommandHandlerApp(ctx context.Context, worker_uuid uuid.UUID) (*internalstorage.Task, error) {
	a.logger.Info("[+] Starting command handler app")
	workerEvent, err := a.storage.Find(worker_uuid)
	if err != nil {
		// TODO: time.Now()
		timestamp, err := time.Parse("2006-01-02 15:04:05", "2022-03-14 12:00:00")
		if err != nil {
			return nil, err
		}

		// TODO: there get data from cache or database
		event := internalstorage.NewEvent(
			"hostname_test", "command_test", "description_test", timestamp, worker_uuid)

		workerEvent, err = a.storage.CreateEvent(*event)
		if err != nil {
			return nil, err
		}
	}

	workerTask := internalstorage.Task{
		ID:          workerEvent.ID,
		Command:     workerEvent.Command,
		Worker_UUID: workerEvent.Worker_UUID,
		Timestamp:   workerEvent.Timestamp,
	}

	return &workerTask, nil
}
