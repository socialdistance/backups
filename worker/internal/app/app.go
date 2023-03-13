package app

import "go.uber.org/zap"

type App struct {
	logger Logger
}

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
	Sync() error
}

func NewApp(logg Logger) *App {
	return &App{
		logger: logg,
	}
}
