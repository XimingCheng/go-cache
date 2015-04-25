package cachetype

import (
	"container/list"
	"errors"
)

type FIFOCache struct {
	capacity  int
	cacheData *list.List
	keyMap    map[interface{}]*list.Element
}

func NewFIFOCache(capacity int) (cache *FIFOCache, err error) {
	if capacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c := &FIFOCache{
		capacity:  capacity,
		cacheData: list.New(),
		keyMap:    make(map[interface{}]*list.Element, capacity),
	}
	return c, nil
}

// add value into FIFO cache
func (cache *FIFOCache) Add(key, value interface{}) {
	if cache.cacheData == nil || cache.keyMap == nil {
		// the cache data is not set
		cache.cacheData = list.New()
		cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
	}

	if ent, ok := cache.keyMap[key]; ok { // if key value exsited
		ent.Value.(*cacheItem).value = value
		return
	}
	ele := &cacheItem{key, value}
	cache.keyMap[key] = cache.cacheData.PushBack(ele)

	if cache.capacity != 0 && cache.cacheData.Len() > cache.capacity {
		cache.removeOldest()
	}
}

func (cache *FIFOCache) removeOldest() {
	cache.removeElement(cache.cacheData.Front())
}

// get the FIFO value data from the cache
func (cache *FIFOCache) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := cache.keyMap[key]; ok {
		return ent.Value.(*cacheItem).value, ok
	}
	return nil, ok
}

func (cache *FIFOCache) Remove(key interface{}) {
	if ent, ok := cache.keyMap[key]; ok {
		cache.removeElement(ent)
	}
}

func (cache *FIFOCache) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	cache.cacheData.Remove(e)
	kv := e.Value.(*cacheItem)
	delete(cache.keyMap, kv.key)
}

func (cache *FIFOCache) Clear() {
	//golang has garbage collection
	cache.cacheData = list.New()
	cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
}

func (cache *FIFOCache) Len() int {
	return cache.cacheData.Len()
}

//iterate cache according to front to back or on the contrary
func (cache *FIFOCache) Keys(old2new bool) []interface{} {
	keys := make([]interface{}, len(cache.keyMap))
	var ent *list.Element = nil
	if !old2new {
		ent = cache.cacheData.Back()
	} else {
		ent = cache.cacheData.Front()
	}
	i := 0
	for ent != nil {
		keys[i] = ent.Value.(*cacheItem).key
		if !old2new {
			ent = ent.Prev()
		} else {
			ent = ent.Next()
		}
		i++
	}
	return keys
}
