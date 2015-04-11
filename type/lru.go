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

// add value into LRU cache
func (cache *LRUCache) Add(key, value interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.cacheData == nil {
		// the cache data is not set
		cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
		cache.cacheData = list.New()
	}
	if ent, ok := cache.keyMap[key]; ok {
		cache.cacheData.MoveToFront(ent)
		ent.Value.(*cacheItem).value = value
	}
	ent := &cacheItem{key, value}
	item := cache.cacheData.PushFront(ent)
	cache.keyMap[key] = item

	if cache.capacity != 0 && cache.cacheData.Len() > cache.capacity {
		cache.removeOldest()
	}
}

// get the LRU value data from the cache
func (cache *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if ent, ok := cache.keyMap[key]; ok {
		cache.cacheData.MoveToFront(ent)
		return ent.Value.(*cacheItem).value, ok
	}
	return
}

func (cache *LRUCache) Remove(key interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if ent, ok := cache.keyMap[key]; ok {
		cache.removeElement(ent)
	}
}

func (cache *LRUCache) Clear() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.cacheData = list.New()
	cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
}

func (cache *LRUCache) removeOldest() {
	ent := cache.cacheData.Back()
	if ent != nil {
		cache.removeElement(ent)
	}
}

func (cache *LRUCache) removeElement(e *list.Element) {
	cache.cacheData.Remove(e)
	kv := e.Value.(*cacheItem)
	delete(cache.keyMap, kv.key)
}
