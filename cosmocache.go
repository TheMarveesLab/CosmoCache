package main

type CosmoCache interface {
	Set(key string, val string) error
	Get(key string) (string, bool)
	Del(key string)
}

type cache struct {
	data map[string]string
}

func NewCache() *cache {
	return &cache{
		data: make(map[string]string),
	}
}

func (c *cache) Set(key string, val string) error {
	c.data[key] = val
	return nil
}

func (c *cache) Get(key string) (string, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *cache) Del(key string) {
	delete(c.data, key)
}
