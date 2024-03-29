package http

type TaskDTO struct {
	ID         string `json:"id" query:"id"`
	Command    string `json:"command" query:"command"`
	WorkerUuid string `json:"worker_uuid" query:"worker_uuid"`
	Timestamp  string `json:"timestamp" query:"timestamp"`
}

type WorkerTaskDTO struct {
	ID       string `json:"id" query:"id"`
	Address  string `json:"address" query:"address"`
	Command  string `json:"command" query:"command"`
	Hostname string `json:"hostname" query:"hostname"`
}

// type ResponseDTO struct {
// 	Success bool   `json:"success" from:"success" query:"success"`
// 	Error   string `json:"error" from:"error" query:"error"`
// }

//func (t *TaskDTO) GetModelTask() (*internalstorage.Task, error) {
//	time, err := time.Parse("2006-01-02 15:04:00", t.Timestamp)
//	if err != nil {
//		return nil, fmt.Errorf("error: Start exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", t.Timestamp, err)
//	}
//
//	id, err := uuid.Parse(t.ID)
//	if err != nil {
//		return nil, fmt.Errorf("ID expected to be uuid, got: %s, %w", t.ID, err)
//	}
//
//	worker_uuid, err := uuid.Parse(t.Worker_UUID)
//	if err != nil {
//		return nil, fmt.Errorf("Worker_UUID exprected to be uuid, got: %s, %w", t.Worker_UUID, err)
//	}
//
//	appTask := internalstorage.NewTask(t.Command, worker_uuid, time)
//	appTask.ID = id
//
//	return appTask, nil
//
//}
