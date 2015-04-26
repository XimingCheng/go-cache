package gocache

import (
	"testing"
	"time"
)

func TestBasicGoCache(t *testing.T) {
	c, e := New(
		&CacheParams{"lru", "testlru", 1, 2, false, 5})
	if e != nil {
		t.Fatalf("err: %v", e)
	}
	c.Add(1, "2")
	c.Add("ahahah", "ok")
	time.Sleep(3 * time.Second)
	if c.Len() != 0 {
		t.Fatalf("err: len != 0")
	}

	c, e = New(
		&CacheParams{"lru", "testlru1", 3, 5, false, 5})
	if e != nil {
		t.Fatalf("err: %v", e)
	}
	c.Add("key", "value")
	c.Add("key1", "value1")
	c.Add("key2", "value2")
	c.Add("key3", "value3")
	c.Add("key4", "value4")
	c.Add("key5", "value5")
	if c.Len() != 5 {
		t.Fatalf("err: len != 5")
	}
	time.Sleep(time.Second)
	c.Get("key1")
	if c.Len() != 5 {
		t.Fatalf("err: len != 5")
	}
	time.Sleep(time.Second)
	c.Get("key2")
	time.Sleep(1500 * time.Millisecond)
	if c.Len() != 2 {
		for _, k := range c.Keys(false) {
			t.Logf("key is %v", k)
		}
		t.Fatalf("err: len != 2 len = %d", c.Len())
	}
	c.Get("key2")
	time.Sleep(1800 * time.Millisecond)
	if c.Len() != 0 {
		t.Fatalf("err: len != 0 len = %d", c.Len())
	}

	c, e = New(
		&CacheParams{"lru", "testlru2", 3, 5, true, 5})
	if e != nil {
		t.Fatalf("err: %v", e)
	}
	c.Add("key", "value")
	c.Add("key1", "value1")
	c.Add("key2", "value2")
	c.Add("key3", "value3")
	c.Add("key4", "value4")
	c.Add("key5", "value5")
	if c.Len() != 5 {
		t.Fatalf("err: len != 5")
	}
	time.Sleep(6 * time.Second)
	if c.Len() != 5 {
		t.Fatalf("err: len != 5")
	}

	c, e = New(
		&CacheParams{"lru", "testlru3", 3, 5, false, 5})
	if e != nil {
		t.Fatalf("err: %v", e)
	}
	c.Add("key", "value")
	c.Add("key1", "value1")
	c.Add("key2", "value2")
	c.Add("key3", "value3")
	c.Add("key4", "value4")
	c.Add("key5", "value5")
	c.Remove("key5")
	time.Sleep(time.Second)
	c.Stop()
	time.Sleep(4 * time.Second)
	if c.Len() != 4 {
		t.Fatalf("err: len != 4 len = %d", c.Len())
	}

	c, e = New(
		&CacheParams{"lru", "testlru4", 3, 5, false, 5})
	if e != nil {
		t.Fatalf("err: %v", e)
	}
	c.Add("key", "value")
	c.Add("key1", "value1")
	c.Add("key2", "value2")
	c.Add("key3", "value3")
	c.Add("key4", "value4")
	c.Add("key5", "value5")
	time.Sleep(time.Second)
	c.Clear()
}
