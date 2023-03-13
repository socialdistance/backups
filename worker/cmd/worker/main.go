package main

// https://blog.intelligentbee.com/2017/08/03/mysqldump-command-useful-usage-examples/

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	internalapp "worker/internal/app"
	internalconfig "worker/internal/config"
	internallogger "worker/internal/logger"
	internalhttp "worker/internal/server/http"

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

	_, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	_, err = internalconfig.LoadConfig(configFile)
	if err != nil {
		logg.Error("Failed load config", zap.Error(err))
	}

	app := internalapp.NewApp(logg)

	client := internalhttp.NewClient(*app, logg)
	client.Run()

	// httpHandler := internalhttp.NewRouter(*app, logg)
	// server := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, app, httpHandler, *logg)

	// go func() {
	// 	server.BuildRouters()

	// 	if err = server.Start(); err != nil {
	// 		logg.Info("failed to start http server: " + err.Error())
	// 		cancel()
	// 	}
	// }()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logg.Info("[+] app stop by signal:", zap.String("signal", s.String()))
	}
	// if err = server.Stop(); err != nil {
	// 	logg.Error("[-] failed to stop http server: " + err.Error())
	// }
}
