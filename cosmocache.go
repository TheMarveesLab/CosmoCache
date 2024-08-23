package main

import (
	"sync"
	"time"
)

const (
	NotExpires = int64(-1)
)

type CosmoCache interface {
	Set(key string, val string, ttl time.Duration) error
	Get(key string) (string, bool)
	Del(key string)
	Flush()
}

type item struct {
	content string
	ttl     int64
}

func (i item) Expired() bool {
	if i.ttl == NotExpires {
		return false
	}
	return time.Now().Unix() < i.ttl
}

type cache struct {
	mu   sync.RWMutex
	data map[string]item
}

func NewCache() *cache {
	return &cache{
		data: make(map[string]item, 0),
	}
}

func (c *cache) Set(key string, val string, ttl time.Duration) error {
	exp := NotExpires
	if ttl > 0 {
		exp = time.Now().Add(ttl).Unix()
	}
	c.mu.Lock()
	c.data[key] = item{
		content: val,
		ttl:     exp,
	}
	c.mu.Unlock()
	return nil
}

func (c *cache) Get(key string) (string, bool) {
	c.mu.RLock()
	i, ok := c.data[key]
	c.mu.RUnlock()
	return i.content, ok
}

func (c *cache) Del(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *cache) Flush() {
	c.data = make(map[string]item)
}
