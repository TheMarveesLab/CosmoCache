package main

import "testing"

func TestCache(t *testing.T) {

	c := NewCache()

	if val1, ok := c.Get("key1"); ok {
		t.Errorf("should not find any value: %v", val1)
	}
	if val2, ok := c.Get("key2"); ok {
		t.Errorf("should not find any value: %v", val2)
	}

	if err := c.Set("key1", "val 1"); err != nil {
		t.Errorf("should not return error: %v", err)
	}
	if err := c.Set("key2", "val 2"); err != nil {
		t.Errorf("should not return error: %v", err)
	}

	val1, ok := c.Get("key1")
	if !ok {
		t.Errorf("could not find value for key1")
	}
	if val1 != "val 1" {
		t.Errorf("should return `val 1`")
	}
	val2, ok := c.Get("key2")
	if !ok {
		t.Errorf("could not find value for key2")
	}
	if val2 != "val 2" {
		t.Errorf("should return `val 2`")
	}

}
