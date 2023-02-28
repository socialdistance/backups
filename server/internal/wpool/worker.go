package wpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TODO: fix args in this func
func worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task, result chan<- interface{}) {
	defer wg.Done()

	for {
		select {
		case <-time.After(time.Second * 5):
			task, ok := <-tasks
			if !ok {
				return
			}
			result <- task.execute(ctx)
		case <-ctx.Done():
			fmt.Printf("Context cancelled %v", ctx.Err())
			return
		}
	}
}

type WorkerPool struct {
	wCount int
	tasks  chan Task
	result chan interface{}
}

func NewWorkerPool(wCount int) *WorkerPool {
	return &WorkerPool{
		wCount: wCount,
		tasks:  make(chan Task, wCount),
		result: make(chan interface{}, wCount),
	}
}
