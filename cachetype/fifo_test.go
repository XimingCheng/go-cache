package cachetype

import (
	"testing"
)

func TestFIFO(t *testing.T) {
	c, err := NewFIFOCache(100)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		c.Add(i+1, i)
	}

	if c.Len() != 100 {
		t.Fatalf("bad len: %v", c.Len())
	}

	if _, ok := c.Get(20); ok {
		t.Fatalf("key 10 should not exist")
	}

	if v, ok := c.Get(256); !ok || v != 255 {
		t.Fatalf("key 256 failed! v %v ok %v", v, ok)
	}

	if c.Clear(); c.Len() != 0 {
		t.Fatalf("cache clear failed!")
	}

	c, err = NewFIFOCache(5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	c.Add(1, "first")
	c.Add(2, "second")
	c.Add(3, "third")
	c.Add(4, "fourth")
	c.Add(5, "fifth")
	c.Add(6, "sixth")
	c.Add(7, "seventh")
	c.Add(3, "third_2")

	for idx, key := range c.Keys(true) {
		if idx == 0 && key != 3 {
			t.Fatalf("key 3 wrong")
		} else if idx == 1 && key != 4 {
			t.Fatalf("key 4 wrong")
		} else if idx == 2 && key != 5 {
			t.Fatalf("key 5 wrong")
		} else if idx == 3 && key != 6 {
			t.Fatalf("key 6 wrong")
		}
	}

	c.Remove(4)
	if _, ok := c.Get(4); ok {
		t.Fatalf("key 4 should not exist")
	}
	if v, ok := c.Get(3); !ok || v != "third_2" {
		t.Fatalf("key 3 value wrong")
	}
}
