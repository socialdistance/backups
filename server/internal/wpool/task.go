package wpool

type Task struct {
	Error error
	TaskFunc func() error
}


