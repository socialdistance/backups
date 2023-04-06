package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() (*Logger, error) {
	cfg := zap.NewDevelopmentConfig()

	cfg.OutputPaths = []string{"./logs/log.log", "stderr"}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Can't build logger %s", err)
	}

	logger.Info("[+] logger construction succeeded")

	return &Logger{
		logger,
	}, nil
}
