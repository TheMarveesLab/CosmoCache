package main

import (
	"sync"
	"testing"
	"time"
)

func TestCacheGet(t *testing.T) {
	cc := NewCache(5 * time.Millisecond)

	if val1, ok := cc.Get("key1"); ok {
		t.Errorf("should not find any value: %v", val1)
	}
	if val2, ok := cc.Get("key2"); ok {
		t.Errorf("should not find any value: %v", val2)
	}

	cc.data["key1"] = item{content: "value for key1"}
	cc.data["key2"] = item{content: "value for key2"}

	val1, ok := cc.Get("key1")
	if !ok {
		t.Error("could not find value for key1")
	}
	if val1 != "value for key1" {
		t.Error("should return `value for key1`")
	}
	val2, ok := cc.Get("key2")
	if !ok {
		t.Error("could not find value for key2")
	}
	if val2 != "value for key2" {
		t.Error("should return `value for key2`")
	}
}

func TestCacheSet(t *testing.T) {
	cc := NewCache(5 * time.Millisecond)

	cc.Set("key1", "value for key1", time.Duration(NotExpires))

	i, ok := cc.data["key1"]
	if !ok {
		t.Error("should get key `key1`")
	}

	if i.content != "value for key1" {
		t.Error("should return `value for key1`")
	}

	if i.ttl != NotExpires {
		t.Error("should return NotExpires")
	}

	cc.Set("key2", "value for key2", time.Duration(5*time.Minute))

	i, ok = cc.data["key2"]
	if !ok {
		t.Error("should get key `key2`")
	}

	if i.content != "value for key2" {
		t.Error("should return `value for key2`")
	}

	if i.ttl <= 0 {
		t.Error("should return ttl > zero")
	}
}

func TestCacheDel(t *testing.T) {
	cc := NewCache(5 * time.Millisecond)
	cc.data["key1"] = item{content: "val for key1"}
	cc.data["key2"] = item{content: "val for key2"}
	cc.data["key3"] = item{content: "val for key3"}

	cc.Del("key2")

	if _, ok := cc.data["key2"]; ok {
		t.Error("should not find key2")
	}

	if len(cc.data) != 2 {
		t.Error("should have 2 items")
	}
}

func TestCacheFlush(t *testing.T) {
	cc := NewCache(5 * time.Millisecond)
	cc.data["key1"] = item{content: "val for key1"}
	cc.data["key2"] = item{content: "val for key2"}
	cc.data["key3"] = item{content: "val for key3"}

	cc.Flush()

	if len(cc.data) != 0 {
		t.Error("should have no items")
	}
}

func TestCacheRaceConditrion(t *testing.T) {
	cc := NewCache(5 * time.Millisecond)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cc.Set("key1", "value for key1", time.Duration(NotExpires))
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = cc.Get("key1")
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cc.Del("key1")
		}()
	}
	wg.Wait()
}

func TestItemExpired(t *testing.T) {
	exp := time.Duration(2 * time.Second)
	i := item{
		content: "a very important text",
		ttl:     time.Now().Add(exp).Unix(),
	}

	if i.Expired() {
		t.Error("should not be expired at this time")
	}

	time.Sleep(3 * time.Second)

	if !i.Expired() {
		t.Error("should be expired at this time")
	}
}

func TestRemoveExpiredItems(t *testing.T) {
	cc := NewCache(5 * time.Millisecond)
	cc.Set("key1", "val 2 sec", 2*time.Second)
	cc.Set("key2", "val 3 sec", 3*time.Second)
	cc.Set("key3", "val 5 sec", 5*time.Second)

	time.Sleep(4 * time.Second)

	if len(cc.data) != 1 {
		t.Errorf("shoudl have 1 item and got %d", len(cc.data))
	}

	time.Sleep(2 * time.Second)

	if len(cc.data) != 0 {
		t.Errorf("shoudl have 0 item and got %d", len(cc.data))
	}
}
