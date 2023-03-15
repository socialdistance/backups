package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalapp "server/internal/app"
	internalconfig "server/internal/config"
	internallogger "server/internal/logger"
	internalhttp "server/internal/server/http"
	internalcache "server/internal/storage/cache"
	internalstore "server/internal/storage/store"
	workerpool "server/internal/wpool"

	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	logg, err := internallogger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config, err := internalconfig.LoadConfig(configFile)
	if err != nil {
		logg.Error("Failed load config", zap.Error(err))
	}

	store := internalstore.CreateStorage(ctx, *config)

	// Создаем контейнер с временем жизни по-умолчанию равным 5 минут и удалением просроченного кеша каждые 10 минут
	cache := internalcache.NewCache(time.Duration(config.Cache.DefaultExpiration)*time.Minute, time.Duration(config.Cache.CleanupInterval)*time.Minute)

	pool, err := workerpool.NewPool(config.WorkerPool.NumWorkers, config.WorkerPool.NumWorkers, logg)
	if err != nil {
		logg.Error("error making worker pool:", zap.Error(err))
		return
	}

	app := internalapp.NewApp(logg, store, cache, pool)

	pool.Start()

	httpHandler := internalhttp.NewRouter(*app, logg)
	server := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, app, httpHandler, *logg)

	doneCh := make(chan struct{})
	startCacheUpdate(store, logg, cache, *pool, doneCh)

	go func() {
		server.BuildRouters()

		if err = server.Start(); err != nil {
			logg.Info("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logg.Info("[+] app stop by signal:", zap.String("signal", s.String()))
		logg.Info("[+] workers stop by signal:", zap.String("signal", s.String()))
		pool.Stop()
		doneCh <- struct{}{}
	}
	if err = server.Stop(); err != nil {
		logg.Error("[-] failed to stop http server: " + err.Error())
	}
}

func startCacheUpdate(storage internalapp.Storage, logger internalapp.Logger, cache internalapp.Cache, pool workerpool.Pool, doneCh chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				task := workerpool.NewTaskPool(func() error {
					events, err := storage.FindAllEvents()
					if err != nil {
						logger.Error("Cant get all events for update cache", zap.Error(err))
						return err
					}

					for _, event := range events {
						cache.Set(event.Worker_UUID, event, 5*time.Minute)
					}

					logger.Info("[+] Task procceed")

					return nil
				})

				pool.AddTask(*task)
			case <-doneCh:
				logger.Info("[+] Ticker stopped")
				ticker.Stop()
			}
		}
	}()
}
