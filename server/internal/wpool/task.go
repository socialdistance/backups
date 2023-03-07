package wpool

type CacheTask struct {
}

func (c *CacheTask) Execute() error {
	panic("implement me")
}

func (c *CacheTask) OnFailure(error) {
	panic("implement me")
}
