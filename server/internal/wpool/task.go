package wpool

import (
	"fmt"
)

type CacheTask struct {
	Error error
	f     func() error
}

func NewTaskPool(f func() error) *CacheTask {
	return &CacheTask{
		f: f,
	}
}

func Execute(c CacheTask) error {
	fmt.Printf("Worker processes task\n")

	c.Error = c.f()

	return c.Error
}
