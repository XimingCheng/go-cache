package cachetype

import (
	"errors"
)

type TWOQCache struct {
	fifoCapacity int
	lruCapacity  int
	lruCache     *LRUCache
	fifoCache    *FIFOCache
}

// return a new Two Queue cache with given capacitys(include lrucache and fifocache), if errors occur, return err
func NewTwoQCache(lruCapacity int, fifoCapacity int) (c *TWOQCache, err error) {
	if lruCapacity <= 0 || fifoCapacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c = &TWOQCache{
		fifoCapacity: fifoCapacity,
		lruCapacity:  lruCapacity,
		lruCache:     nil,
		fifoCache:    nil,
	}
	if c.lruCache, err = NewLRUCache(lruCapacity); err != nil {
		return nil, err
	}
	if c.fifoCache, err = NewFIFOCache(fifoCapacity); err != nil {
		return nil, err
	}
	return c, nil
}

func (cache *TWOQCache) Add(key, value interface{}) {
	cache.fifoCache.Add(key, value)
}

func (cache *TWOQCache) Get(key interface{}) (value interface{}, ok bool) {
	if value, ok := cache.fifoCache.Get(key); ok {
		//if fifo cache exits!
		cache.fifoCache.Remove(key)
		cache.lruCache.Add(key, value)
		return value, ok
	} else if value, ok := cache.lruCache.Get(key); ok {
		//if lru cache exits!
		return value, ok
	} else {
		//none of caches exits!
		return value, ok
	}
}

func (cache *TWOQCache) Remove(key interface{}) {
	if cache.fifoCache.IsExist(key) {
		cache.fifoCache.Remove(key)
		return
	} else if cache.lruCache.IsExist(key) {
		cache.lruCache.Remove(key)
		return
	}
}

func (cache *TWOQCache) Clear() {
	cache.fifoCache.Clear()
	cache.lruCache.Clear()
}

func (cache *TWOQCache) Len() int {
	return cache.fifoCache.Len() + cache.lruCache.Len()
}

func (cache *TWOQCache) IsExist(key interface{}) bool {
	return cache.fifoCache.IsExist(key) && cache.lruCache.IsExist(key)
}
