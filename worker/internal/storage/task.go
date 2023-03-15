package storage

import (
	"github.com/google/uuid"
	"net"
	"os"
)

type Task struct {
	WorkerUuid uuid.UUID // workerID
	Address    string
	Command    string
	Hostname   string
}

func NewTask(workerUid uuid.UUID) (*Task, error) {
	t := Task{}

	t.WorkerUuid = workerUid

	if err := t.HostnameWorker(); err != nil {
		return nil, err
	}

	if err := t.GetLocalIPWorker(); err != nil {
		return nil, err
	}

	t.Command = "test_command"

	return &t, nil
}

func (t *Task) HostnameWorker() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	t.Hostname = hostname

	return nil
}

func (t *Task) GetLocalIPWorker() error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				t.Address = ipnet.IP.String()
			}
		}
	}
	return nil
}
