package wpool

type Task struct {
	Error    error
	TaskFunc func() error
}

func execute(t *Task) error {
	t.Error = t.TaskFunc()

	return t.Error
}
