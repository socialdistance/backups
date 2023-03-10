package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.zap.Debug(message, fields...)
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	l.zap.Info(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	l.zap.Error(message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.zap.Fatal(message, fields...)
}

func (l *Logger) With(fields ...zap.Field) *zap.Logger {
	return l.zap.With(fields...)
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func NewLogger() (*Logger, error) {
	cfg := zap.NewDevelopmentConfig()

	cfg.OutputPaths = []string{"stderr"}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Can't build logger %s", err)
	}

	logger.Info("[+] logger construction succeeded")

	return &Logger{
		zap: logger,
	}, nil
}
