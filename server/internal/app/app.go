package app

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	internalstorage "server/internal/storage"
	wpool "server/internal/wpool"
)

type App struct {
	logger  Logger
	storage Storage
	cache   Cache
	wpool   WorkerPool
}

type Logger interface {
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
}

type Cache interface {
	Set(key uuid.UUID, value internalstorage.Event, duration time.Duration)
	Get(key uuid.UUID) (*internalstorage.Event, bool)
	Delete(key uuid.UUID) error
}

type Storage interface {
	CreateEvent(e internalstorage.Event) error
	DeleteEvent(id uuid.UUID) error
	UpdateCommand(e internalstorage.Event) error
	FindAllEvents() ([]internalstorage.Event, error)
}

type WorkerPool interface {
	AddTask(cacheTask wpool.CacheTask)
}

func NewApp(logger Logger, storage Storage, cache Cache, pool WorkerPool) *App {
	return &App{
		logger:  logger,
		storage: storage,
		cache:   cache,
		wpool:   pool,
	}
}

func (a *App) CommandHandlerApp(ctx context.Context, workerUuid uuid.UUID, address, command, hostname string) (*internalstorage.Task, error) {
	var event *internalstorage.Event

	workerEvent, found := a.cache.Get(workerUuid)

	timestamp := time.Now()
	event = internalstorage.NewEvent(
		address, command, hostname, timestamp, workerUuid)

	// Кеш будет обновляться каждые ~ 5 минут, воркер будет ходить в базу раз в 5 минут
	// и обновлять данные из кеша. Тем самым я гарантирую, что в кеше каждые 5 минут будут актуальные данные
	// Даже если они устареют, через 5 минут они обновяться и воркеры будут получать обновленные данные каждые 5 минут
	if !found {
		a.cache.Set(event.WorkerUuid, *event, 5*time.Minute)
		err := a.storage.CreateEvent(*event)
		if err != nil {
			a.logger.Error("[-] Failed create event: ", zap.Error(err))
			return nil, err
		}

		workerTask := internalstorage.NewTask(event.Command, event.WorkerUuid, timestamp)

		return workerTask, nil
	}

	// после того, как бекап был сделан через команду manual, нужно и в базе и кеше обновить значение обратно на cron
	if workerEvent.Command == "manual" {
		event.Command = "cron"

		err := a.cache.Delete(event.WorkerUuid)
		if err != nil {
			a.logger.Error("[-] Failed delete key from cache", zap.Error(err))
			return nil, err
		}

		a.cache.Set(event.WorkerUuid, *event, 5*time.Minute)

		err = a.storage.UpdateCommand(*event)
		if err != nil {
			a.logger.Error("[-] Failed update event in database", zap.Error(err))
			return nil, err
		}

		workerTask := internalstorage.NewTask(workerEvent.Command, event.WorkerUuid, workerEvent.Timestamp)
		return workerTask, nil
	}

	workerTask := internalstorage.NewTask(workerEvent.Command, workerEvent.WorkerUuid, workerEvent.Timestamp)

	return workerTask, nil
}
