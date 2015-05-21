package cachetype

import (
	"testing"
)

func TestLFU(t *testing.T) {
	c, err := NewLFUCache(100)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		c.Add(i, i)
	}

	if c.Len() != 100 {
		t.Fatalf("bad len: %v", c.Len())
	}

	c.Add(10, "hahaha")
	if _, ok := c.Get(10); !ok {
		t.Fatalf("key 10 should exist")
	}

	if c.Clear(); c.Len() != 0 {
		t.Fatalf("cache clear failed!")
	}

	c, err = NewLFUCache(5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	c.Add(1, "nima")
	c.Add(2, "test")
	c.Add(3, "hahha")

	if c.Len() != 3 {
		t.Fatalf("c len() != 3")
	}

	c.Add("key", "value")
	c.Add(12, -1)
	if v, ok := c.Get(12); ok && v != -1 {
		t.Fatalf("c.Get(12) != -1")
	}
	if v, ok := c.Get(1); ok && v != "nima" {
		t.Fatalf("c.Get(1) != nima")
	}
	if v, ok := c.Get(2); ok && v != "test" {
		t.Fatalf("c.Get(2) != test")
	}
	if v, ok := c.Get(3); ok && v != "hahha" {
		t.Fatalf("c.Get(3) != hahha")
	}
	c.Add("hhhhh", "nimatest")
	for _, key := range c.Keys(true) {
		t.Logf("key %v", key)
	}
	for _, key := range c.Keys(false) {
		t.Logf("key %v", key)
	}
	if c.IsExist("key") {
		t.Fatalf("key key should not exsit")
	}

	c.Remove("hhhhh")
	if c.IsExist("hhhhh") {
		t.Fatalf("hhhhh should not exist")
	}
}
