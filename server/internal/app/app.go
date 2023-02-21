package app

import (
	internalstorage "server/internal/storage"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type App struct {
	logger  Logger
	storage Storage
	cache   Cache
}

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
	Sync() error
}

type Cache interface {
	Set(key uuid.UUID, value internalstorage.Event, duration time.Duration)
	Get(key uuid.UUID) (*internalstorage.Event, bool)
	Delete(key uuid.UUID) error
}

type Storage interface {
	CreateEvent(e internalstorage.Event) error
	DeleteEvent(id uuid.UUID) error
	Find(worker_UUID uuid.UUID) (*internalstorage.Event, error)
	FindAllEvents() ([]internalstorage.Event, error)
}

func NewApp(logger Logger, storage Storage, cache Cache) *App {
	return &App{
		logger:  logger,
		storage: storage,
		cache:   cache,
	}
}

func (a *App) CommandHandlerApp(ctx context.Context, worker_uuid uuid.UUID) (*internalstorage.Task, error) {
	a.logger.Info("[+] Starting command handler app")
	// 1. check data in cache
	// 2. check data in database/memorystorage?

	var event *internalstorage.Event

	workerEvent, found := a.cache.Get(worker_uuid)
	if !found {
		timestamp, err := time.Parse("2006-01-02 15:04:05", "2022-03-14 12:00:00")
		if err != nil {
			return nil, err
		}

		// TODO: there set data to cache and database
		event = internalstorage.NewEvent(
			"hostname_test", "command_test", "description_test", timestamp, worker_uuid)

		a.cache.Set(event.Worker_UUID, *event, 5*time.Minute)
		err = a.storage.CreateEvent(*event)
		if err != nil {
			a.logger.Error("Failed create event", zap.Error(err))
			return nil, err
		}

		workerTask := internalstorage.NewTask(event.Command, event.Worker_UUID, timestamp)

		return workerTask, nil
	}

	workerTask := internalstorage.NewTask(workerEvent.Command, workerEvent.Worker_UUID, workerEvent.Timestamp)

	return workerTask, nil
}
