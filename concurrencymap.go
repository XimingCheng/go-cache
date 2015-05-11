package gocache

import (
	"errors"
	"sync"
)

type ConcurrencyMap struct {
	items map[interface{}]interface{}
	lock  sync.RWMutex
}

func NewConcurrencyMap(capacity int) (c *ConcurrencyMap, err error) {
	if capacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c = &ConcurrencyMap{
		items: make(map[interface{}]interface{}, capacity),
	}
	return c, nil
}

func (c *ConcurrencyMap) Add(k, v interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items[k] = v
}

func (c ConcurrencyMap) Get(k interface{}) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	v, ok := c.items[k]
	return v, ok
}

func (c *ConcurrencyMap) Delete(k interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.items, k)
}

func (c ConcurrencyMap) isEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.items) == 0
}

func (c ConcurrencyMap) isExist(k interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.items[k]
	return ok
}
