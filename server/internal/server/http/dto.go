package http

import (
	"fmt"
	internalstorage "server/internal/storage"
	"time"

	"github.com/google/uuid"
)

type TaskDTO struct {
	ID          string `json:"id"`
	Command     string `json:"command"`
	Worker_UUID string `json:"worker_uuid"`
	Timestamp   string `json:"timestamp"`
}

type ResponseDTO struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (t *TaskDTO) GetModelTask() (*internalstorage.Task, error) {
	time, err := time.Parse("2006-01-02 15:04:00", "2020-10-20 12:30:00")
	if err != nil {
		return nil, fmt.Errorf("error: Start exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", t.Timestamp, err)
	}

	id, err := uuid.Parse(t.ID)
	if err != nil {
		return nil, fmt.Errorf("ID exprected to be uuid, got: %s, %w", t.ID, err)
	}

	worker_uuid, err := uuid.Parse(t.Worker_UUID)
	if err != nil {
		return nil, fmt.Errorf("Worker_UUID exprected to be uuid, got: %s, %w", t.Worker_UUID, err)
	}

	appTask := internalstorage.NewTask(t.Command, worker_uuid, time)
	appTask.ID = id

	return appTask, nil

}
