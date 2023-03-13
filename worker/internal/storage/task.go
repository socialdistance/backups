package storage

import (
	"os"

	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID // my workerID
	address  string
	Command  string
	Hostname string
}

func (t *Task) HostnameWorker() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	t.Hostname = hostname

	return nil
}
