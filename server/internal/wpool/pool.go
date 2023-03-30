package wpool

// https://github.com/swayne275/go-work/blob/main/example/main.go

import (
	"sync"

	"go.uber.org/zap"
)

type Logger interface {
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
}

type Pool struct {
	numWorkers int
	tasks      chan CacheTask

	start sync.Once
	stop  sync.Once

	quit chan struct{}

	logger Logger
}

func NewPool(numWorkers int, channelSize int, logger Logger) (*Pool, error) {
	tasks := make(chan CacheTask, channelSize)
	quit := make(chan struct{})

	return &Pool{
		numWorkers: numWorkers,
		tasks:      tasks,
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       quit,
		logger:     logger,
	}, nil
}

func (p *Pool) Start() {
	p.start.Do(func() {
		p.logger.Info("[+] Starting worker pool")
		p.startWorkers()
	})
}

func (p *Pool) Stop() {
	p.stop.Do(func() {
		p.logger.Info("[+] Stopping worker pool")
		p.quit <- struct{}{}
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
			p.logger.Info("[+] Starting workers")

			for {
				select {
				case task, ok := <-p.tasks:
					if !ok {
						return
					}

					err := Execute(task)
					if err != nil {
						return
					}
				case <-p.quit:
					p.logger.Info("[+] Stopping worker and quit channel")
					return
				}
			}
		}(i)
	}
}
