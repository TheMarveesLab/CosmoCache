package main

import (
	"sync"
	"testing"
)

func TestCacheGet(t *testing.T) {
	cc := NewCache()

	if val1, ok := cc.Get("key1"); ok {
		t.Errorf("should not find any value: %v", val1)
	}
	if val2, ok := cc.Get("key2"); ok {
		t.Errorf("should not find any value: %v", val2)
	}

	cc.data["key1"] = "value for key1"
	cc.data["key2"] = "value for key2"

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
	cc := NewCache()

	cc.Set("key1", "value for key1")

	val, ok := cc.data["key1"]
	if !ok {
		t.Error("should get key `key1`")
	}

	if val != "value for key1" {
		t.Error("should return `value for key1`")
	}
}

func TestCacheDel(t *testing.T) {
	cc := NewCache()
	cc.data["key1"] = "val for key1"
	cc.data["key2"] = "val for key2"
	cc.data["key3"] = "val for key3"

	cc.Del("key2")

	if _, ok := cc.data["key2"]; ok {
		t.Error("should not find key2")
	}

	if len(cc.data) != 2 {
		t.Error("should have 2 items")
	}
}

func TestCacheFlush(t *testing.T) {
	cc := NewCache()
	cc.data["key1"] = "val for key1"
	cc.data["key2"] = "val for key2"
	cc.data["key3"] = "val for key3"

	cc.Flush()

	if len(cc.data) != 0 {
		t.Error("should have no items")
	}
}

func TestCacheRaceConditrion(t *testing.T) {
	cc := NewCache()

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cc.Set("key1", "value for key1")
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
