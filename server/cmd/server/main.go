package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	internalapp "server/internal/app"
	internalconfig "server/internal/config"
	internallogger "server/internal/logger"
	internalhttp "server/internal/server/http"
	internalcache "server/internal/storage/cache"
	internalstore "server/internal/storage/store"
	"syscall"
	"time"

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
	// TODO: add to config defaultExpiration and cleanupInterval
	cache := internalcache.NewCache(5*time.Minute, 10*time.Minute)

	app := internalapp.NewApp(logg, store, cache)

	httpHandler := internalhttp.NewRouter(*app, logg)
	server := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, app, httpHandler, *logg)

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
	}
	if err = server.Stop(); err != nil {
		logg.Error("[-] failed to stop http server: " + err.Error())
	}

}
