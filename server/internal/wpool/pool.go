package wpool

type WorkerPool struct {
}

func NewPool(wCount int) *WorkerPool {
	return &WorkerPool{}
}
