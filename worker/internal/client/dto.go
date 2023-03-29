package client

type TaskDTO struct {
	ID         string `json:"id" query:"id"`
	Command    string `json:"command" query:"command"`
	WorkerUuid string `json:"worker_uuid" query:"worker_uuid"`
	Timestamp  string `json:"timestamp" query:"timestamp"`
}

type WorkerTaskDTO struct {
	ID string `json:"id" query:"id"`
}
