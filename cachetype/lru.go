package cachetype

import (
	"container/list"
	"errors"
)

// go cache basic data structure
type LRUCache struct {
	// the capacity of the cache data stored in the memory
	capacity int
	// cache data stored in the memory
	cacheData *list.List
	// the key index mapping data, used for fast searching in the cache list
	keyMap map[interface{}]*list.Element
}

// return a new gocache with given capacity, if errors occur, return err
func NewLRUCache(capacity int) (c *LRUCache, err error) {
	if capacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c = &LRUCache{
		capacity:  capacity,
		cacheData: list.New(),
		keyMap:    make(map[interface{}]*list.Element, capacity),
	}
	return c, nil
}

// add value into LRU cache
func (cache *LRUCache) Add(key, value interface{}) {
	if cache.cacheData == nil {
		// the cache data is not set
		cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
		cache.cacheData = list.New()
	}
	if ent, ok := cache.keyMap[key]; ok {
		cache.cacheData.MoveToFront(ent)
		ent.Value.(*cacheItem).value = value
		return
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
	if ent, ok := cache.keyMap[key]; ok {
		cache.cacheData.MoveToFront(ent)
		return ent.Value.(*cacheItem).value, ok
	}
	return
}

func (cache *LRUCache) Remove(key interface{}) {
	if ent, ok := cache.keyMap[key]; ok {
		cache.removeElement(ent)
	}
}

func (cache *LRUCache) Clear() {
	cache.cacheData = list.New()
	cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
}

func (cache *LRUCache) Len() int {
	return cache.cacheData.Len()
}

// Keys returns a slice of the keys in the cache
// old2new true from oldest to newest
func (cache *LRUCache) Keys(old2new bool) []interface{} {
	keys := make([]interface{}, len(cache.keyMap))
	var ent *list.Element = nil
	if old2new {
		ent = cache.cacheData.Back()
	} else {
		ent = cache.cacheData.Front()
	}
	i := 0
	for ent != nil {
		keys[i] = ent.Value.(*cacheItem).key
		if old2new {
			ent = ent.Prev()
		} else {
			ent = ent.Next()
		}
		i++
	}
	return keys
}

func (cache *LRUCache) removeOldest() {
	ent := cache.cacheData.Back()
	cache.removeElement(ent)
}

func (cache *LRUCache) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	cache.cacheData.Remove(e)
	kv := e.Value.(*cacheItem)
	delete(cache.keyMap, kv.key)
}
