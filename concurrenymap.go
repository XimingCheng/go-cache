package gocache

import (
	"errors"
	"sync"
)

type ConcurrenyMap struct {
	items map[interface{}]interface{}
	lock  sync.RWMutex
}

func NewConcurrenyMap(capacity int) (c *ConcurrenyMap, err error) {
	if capacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c = &ConcurrenyMap{
		items: make(map[interface{}]interface{}, capacity),
	}
	return c, nil
}

func (c *ConcurrenyMap) Add(k, v interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items[k] = v
}

func (c ConcurrenyMap) Get(k interface{}) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	v, ok := c.items[k]
	return v, ok
}
func (c *ConcurrenyMap) Delete(k interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.items, k)
}

func (c ConcurrenyMap) isEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.items) == 0
}

func (c ConcurrenyMap) isExist(k interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.items[k]
	return ok
}
