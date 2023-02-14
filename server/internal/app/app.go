package app

import (
	"go.uber.org/zap"
	"net/http"
)

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	LogHTTP(r *http.Request, code, length int)
	Sync() error
}

type App struct {
	logger  Logger
	storage Storage
}

type Storage interface { // TODO
}

func NewApp(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

// TODO
