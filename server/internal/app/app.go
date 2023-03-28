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
	Find(workerUuid uuid.UUID) (*internalstorage.Event, error)
	FindAllEvents() ([]internalstorage.Event, error)
}

type WorkerPool interface {
	Start()
	Stop()
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
	// Кеш будет обновляться каждые ~ 5 минут, воркер будет ходить в базу раз в 5 минут
	// и обновлять данные из кеша. Тем самым я гарантирую, что в кеше каждые 5 минут будут актуальные данные
	// Даже если они устареют, через 5 минут они обновяться и воркеры будут получать обновленные данные каждые 5 минут
	if !found {
		timestamp := time.Now()
		//timestamp, err := time.Parse("2006-01-02 15:04:05", time.Now().Format(time.RFC3339))
		//if err != nil {
		//	return nil, err
		//}

		event = internalstorage.NewEvent(
			address, command, hostname, timestamp, workerUuid)

		a.cache.Set(event.WorkerUuid, *event, 5*time.Minute)
		err := a.storage.CreateEvent(*event)
		if err != nil {
			a.logger.Error("[-] Failed create event: ", zap.Error(err))
			return nil, err
		}

		workerTask := internalstorage.NewTask(event.Command, event.WorkerUuid, timestamp)

		return workerTask, nil
	}

	workerTask := internalstorage.NewTask(workerEvent.Command, workerEvent.WorkerUuid, workerEvent.Timestamp)

	return workerTask, nil
}
