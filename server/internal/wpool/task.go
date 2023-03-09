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

func (c *CacheTask) Execute() error {
	fmt.Printf("Worker processes task\n")

	c.Error = c.f()

	return c.Error
}

func (c *CacheTask) OnFailure(error) {
	panic("implement me")
}
