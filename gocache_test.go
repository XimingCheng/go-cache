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

	c1, e1 := New(
		&CacheParams{"lru", "testlru1", 3, 5, false, 5})
	if e1 != nil {
		t.Fatalf("err: %v", e1)
	}
	c1.Add("key", "value")
	c1.Add("key1", "value1")
	c1.Add("key2", "value2")
	c1.Add("key3", "value3")
	c1.Add("key4", "value4")
	c1.Add("key5", "value5")
	if c1.Len() != 5 {
		t.Fatalf("err: len != 5")
	}
	time.Sleep(time.Second)
	c1.Get("key1")
	if c1.Len() != 5 {
		t.Fatalf("err: len != 5")
	}
	time.Sleep(time.Second)
	c1.Get("key2")
	time.Sleep(1500 * time.Millisecond)
	if c1.Len() != 2 {
		for _, k := range c1.Keys(false) {
			t.Logf("key is %v", k)
		}
		t.Fatalf("err: len != 2 len = %d", c1.Len())
	}
	c1.Get("key2")
	time.Sleep(1800 * time.Millisecond)
	if c1.Len() != 0 {
		t.Fatalf("err: len != 0 len = %d", c1.Len())
	}

	c2, e2 := New(
		&CacheParams{"lru", "testlru2", 3, 5, true, 5})
	if e2 != nil {
		t.Fatalf("err: %v", e2)
	}
	c2.Add("key", "value")
	c2.Add("key1", "value1")
	c2.Add("key2", "value2")
	c2.Add("key3", "value3")
	c2.Add("key4", "value4")
	c2.Add("key5", "value5")
	if c2.Len() != 5 {
		t.Fatalf("err: len != 5")
	}
	time.Sleep(6 * time.Second)
	if c2.Len() != 5 {
		t.Fatalf("err: len != 5")
	}

	c3, e3 := New(
		&CacheParams{"lru", "testlru3", 3, 5, false, 5})
	if e3 != nil {
		t.Fatalf("err: %v", e3)
	}
	c3.Add("key", "value")
	c3.Add("key1", "value1")
	c3.Add("key2", "value2")
	c3.Add("key3", "value3")
	c3.Add("key4", "value4")
	c3.Add("key5", "value5")
	c3.Remove("key5")
	time.Sleep(time.Second)
	if c3.Len() != 4 {
		t.Fatalf("err: len != 4 len = %d", c3.Len())
	}
	c3.Clear()
	if c3.Len() != 0 {
		t.Fatalf("err: len != 0 len = %d", c3.Len())
	}

	c4, e4 := New(
		&CacheParams{"lru", "testlru4", 3, 5, false, 5})
	if e4 != nil {
		t.Fatalf("err: %v", e4)
	}
	c4.Add("key", "value")
	c4.Add("key1", "value1")
	c4.Add("key2", "value2")
	c4.Add("key3", "value3")
	c4.Add("key4", "value4")
	c4.Add("key5", "value5")
	time.Sleep(time.Second)
	c4.Clear()
	time.Sleep(time.Second)
}
