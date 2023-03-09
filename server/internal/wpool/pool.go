package wpool

// https://github.com/swayne275/go-work/blob/main/example/main.go

import (
	"fmt"
	"sync"
)

type Pool struct {
	numWorkers int
	tasks      chan CacheTask

	start sync.Once
	stop  sync.Once

	quit chan struct{}
}

func NewPool(numWorkers int, channelSize int) (*Pool, error) {
	tasks := make(chan CacheTask, channelSize)
	quit := make(chan struct{})

	return &Pool{
		numWorkers: numWorkers,
		tasks:      tasks,
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       quit,
	}, nil
}

func (p *Pool) Start() {
	p.start.Do(func() {
		fmt.Println("[+] Starting worker pool")
		p.startWorkers()
	})
}

func (p *Pool) Stop() {
	p.stop.Do(func() {
		fmt.Println("[+] Stopping worker pool")
		close(p.quit)
	})
}

func (p *Pool) AddTask(task CacheTask) {
	select {
	case p.tasks <- task:
	case <-p.quit:
	}
}

func (p *Pool) startWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {
			fmt.Println("[+] Starting worker")

			for {
				select {
				case task, ok := <-p.tasks:
					fmt.Println("Get job:", task)
					if !ok {
						fmt.Println("fail")
						return
					}

					fmt.Println("Starting doing job")
					if err := task.Execute(); err != nil {
						task.OnFailure(err)
					}
				case <-p.quit:
					fmt.Println("[+] Stopping worker and quit channel")
					return
				}
			}
		}(i)
	}
}
