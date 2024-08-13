package main

import "sync"

type CosmoCache interface {
	Set(key string, val string) error
	Get(key string) (string, bool)
	Del(key string)
	Flush()
}

type cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *cache {
	return &cache{
		data: make(map[string]string),
	}
}

func (c *cache) Set(key string, val string) error {
	c.mu.Lock()
	c.data[key] = val
	c.mu.Unlock()
	return nil
}

func (c *cache) Get(key string) (string, bool) {
	c.mu.RLock()
	val, ok := c.data[key]
	c.mu.RUnlock()
	return val, ok
}

func (c *cache) Del(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *cache) Flush() {
	c.data = make(map[string]string)
}
