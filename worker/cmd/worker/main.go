package main

// https://blog.intelligentbee.com/2017/08/03/mysqldump-command-useful-usage-examples/

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"os/signal"
	"syscall"
	internalhttp "worker/internal/client"

	"go.uber.org/zap"
	internalapp "worker/internal/app"
	internalconfig "worker/internal/config"
	internallogger "worker/internal/logger"
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
		panic(err)
	}

	app := internalapp.NewApp(logg)

	workerUuid := uuid.New()
	worker := internalhttp.NewClient(*app, logg, config.HTTP.TargetUrl, config.File.FileNameBackup, workerUuid)

	if err = worker.Run(ctx); err != nil {
		logg.Error("Failed start client", zap.Error(err))
	}

	select {
	case <-ctx.Done():
		logg.Info("[+] app stop by signal")
	}
}
