package wpool

import (
	"fmt"
	"time"
)

type Pool struct {
	Workers []*WorkerPool

	wCount        int
	collector     chan *Task
	runBackground chan bool
}

func NewPool(wCount int) *Pool {
	return &Pool{
		wCount:    wCount,
		collector: make(chan *Task),
	}
}

func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

func (p *Pool) RunBackground() {
	go func() {
		for {
			fmt.Println("\n[+] Waiting new task...")
			time.Sleep(time.Second * 5)
		}
	}()

	for i := 0; i < p.wCount; i++ {
		workers := NewWorkerPool(p.collector)
		p.Workers = append(p.Workers, workers)
		go workers.StartWorkerBackground()
	}
}

func (p *Pool) StopBackground() {
	for i := range p.Workers {
		p.Workers[i].StopWorkerBackGround()
	}

	p.runBackground <- true
}
