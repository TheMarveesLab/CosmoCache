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
	now := time.Now().Unix()
	return now > i.ttl
}

type cache struct {
	mu   sync.RWMutex
	data map[string]item

	cleanupInterval time.Duration
	stopCh          chan struct{}
}

func NewCache(cleanupInterval time.Duration) *cache {
	cc := &cache{
		data:            make(map[string]item, 0),
		cleanupInterval: cleanupInterval,
	}

	go cc.start()

	return cc
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

func (c *cache) start() {
	ticker := time.NewTicker(c.cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.removeExpiredItems()
		case <-c.stopCh:
			ticker.Stop()
			return
		}
	}
}

func (c *cache) Stop() {
	c.stopCh <- struct{}{}
}

func (c *cache) removeExpiredItems() {
	for k, v := range c.data {
		if v.Expired() {
			delete(c.data, k)
		}
	}
}
