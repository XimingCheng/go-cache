package lru

import (
	"testing"
)

func TestLRU(t *testing.T) {
	c, err := New(100)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		c.Add(i, i)
	}

	if c.Len() != 100 {
		t.Fatalf("bad len: %v", c.Len())
	}

	if _, ok := c.Get(10); ok {
		t.Fatalf("key 10 should not exist")
	}

	if v, ok := c.Get(255); !ok || v != 255 {
		t.Fatalf("key 255 failed! v %v ok %v", v, ok)
	}

	if c.Clear(); c.Len() != 0 {
		t.Fatalf("cache clear failed!")
	}

	c, err = New(10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	c.Add(1, "hahaha")
	c.Add(2, "hehehe")
	c.Add(3, "github")
	c.Add(4, "google")
	c.Add(1, "test")
	for idx, key := range c.Keys(true) {
		if idx == 0 && key != 2 {
			t.Fatalf("key 2 wrong")
		} else if idx == 1 && key != 3 {
			t.Fatalf("key 3 wrong")
		} else if idx == 2 && key != 4 {
			t.Fatalf("key 4 wrong")
		} else if idx == 3 && key != 1 {
			t.Fatalf("key 1 wrong")
		}
	}

	c.Remove(2)
	if _, ok := c.Get(2); ok {
		t.Fatalf("key 2 should not exist")
	}
	if v, ok := c.Get(1); !ok || v != "test" {
		t.Fatalf("key 1 value wrong")
	}
}
