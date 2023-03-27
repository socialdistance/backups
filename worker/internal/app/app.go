package app

import (
	"go.uber.org/zap"
	"os/exec"
)

type App struct {
	logger Logger
}

type Logger interface {
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
}

func NewApp(logg Logger) *App {
	return &App{
		logger: logg,
	}
}

func (a *App) ExecuteBackupScript(path string) error {
	a.logger.Info("[+] Executing backup script")

	_, err := exec.Command("/bin/sh", path).Output()
	if err != nil {
		a.logger.Error("Error execute backup script:", zap.Error(err))
		return err
	}

	return nil
}
