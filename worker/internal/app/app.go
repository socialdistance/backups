package app

import (
	"fmt"
	"go.uber.org/zap"
	"os/exec"
)

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

func (a *App) ExecuteBackupScript(path string) error {
	a.logger.Info("[+] Executing backup script")

	cmd, err := exec.Command("/bin/sh", path).Output()
	if err != nil {
		fmt.Printf("error %s", err)
		return err
	}
	output := string(cmd)
	fmt.Println("OUTPUT", output)

	return nil
}
