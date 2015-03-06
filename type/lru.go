package lru

import (
	"container/list"
	"errors"
	"sync"
)

// go cache basic data structure
type LRUCache struct {
	// the capacity of the cache data stored in the memory
	capacity int
	// cache data stored in the memory
	cacheData *list.List
	// the key index mapping data, used for fast searching in the cache list
	keyMap map[interface{}]*list.Element
	// the lock used in goroutines for synchronization
	lock sync.Mutex
}

// data item in the cacheData
type cacheItem struct {
	key   interface{}
	value interface{}
}

// return a new gocache with given capacity, if errors occur, return err
func New(capacity int) (cache *LRUCache, err error) {
	if capacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c := &LRUCache{
		capacity:  capacity,
		cacheData: list.New(),
		keyMap:    make(map[interface{}]*list.Element, capacity),
	}
	return c, nil
}

//
func (cache *LRUCache) Get(key interface{}) (value interface{}, ok bool) {

}

func (cache *LRUCache) Set(key interface{}, value interface{}) {

}

func (cache *LRUCache) Remove(key interface{}) {

}

func (cache *LRUCache) Clear() {

}
