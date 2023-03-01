package wpool

import (
	"fmt"
	"time"
)

type WorkerPool struct {
	tasks chan *Task
	quit  chan struct{}
}

func NewWorkerPool(taskCh chan *Task) *WorkerPool {
	return &WorkerPool{
		tasks: taskCh,
		quit:  make(chan struct{}),
	}
}

func (wp *WorkerPool) StartWorkerBackground() {
	fmt.Println("[+] Start worker")
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case task := <-wp.tasks:
			err := execute(task)
			if err != nil {
				return
			}
		case <-wp.quit:
			fmt.Println("Test stop")
			ticker.Stop()
			return
		}
	}
}

func (wp *WorkerPool) StopWorkerBackGround() {
	fmt.Println("[+] Stopping workers in background...")

	go func() {
		
		<-wp.quit
	}()
}
