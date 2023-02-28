package wpool

import "context"

type Task struct {
	Error    error
	TaskFunc func(ctx context.Context) error
}

func (t Task) execute(ctx context.Context) error {
	err := t.TaskFunc(ctx)
	if err != nil {
		return err
	}

	return nil
}
