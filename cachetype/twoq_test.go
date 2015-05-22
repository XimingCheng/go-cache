package cachetype

import (
	"testing"
)

func TestTwoQ(t *testing.T) {
	c, err := NewTwoQCache(20, 100)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		c.Add(i, i)
	}

	if c.Len() != 20 {
		t.Fatalf("bad len: %v", c.Len())
	}

	c.Add(10, "hahaha")
	if _, ok := c.Get(10); !ok {
		t.Fatalf("key 10 should exist")
	}

	if c.Clear(); c.Len() != 0 {
		t.Fatalf("cache clear failed!")
	}

	c, err = NewTwoQCache(4, 5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	c.Add(1, 3)
	c.Add(2, 4)
	c.Add(2, 5)
	c.Add(3, 7)
	c.Add(4, 100)
	if c.Len() != 4 {
		t.Fatalf("c len() != 4")
	}

	if v, ok := c.Get(1); ok && v != 3 {
		t.Fatalf("c.Get(1) != 3")
	}
	if v, ok := c.Get(2); ok && v != 5 {
		t.Fatalf("c.Get(2) != 5")
	}
	c.Add(101, 100)
	c.Add(102, 100)
	c.Add(103, 100)
	c.Add(104, 100)
	if c.Len() != 6 {
		t.Fatalf("bad len: %v", c.Len())
	}

	c.Remove(101)
	if c.IsExist(101) {
		t.Fatalf("101 should not exist")
	}

	c.Remove(1)
	if c.IsExist(1) {
		t.Fatalf("1 should not exist")
	}
}
