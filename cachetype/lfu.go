package cachetype

import (
	"container/heap"
	"errors"
)

type dataWrapper struct {
	key       interface{}
	value     interface{}
	frequency int
	keyMapP   *map[interface{}]int
}

type dataHeap []*dataWrapper

type LFUCache struct {
	capacity  int
	cacheData *dataHeap
	keyMap    map[interface{}]int
}

func (h dataHeap) Len() int           { return len(h) }
func (h dataHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h dataHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	(*h[i].keyMapP)[h[i].key] = i
	(*h[j].keyMapP)[h[j].key] = j
}

func (h *dataHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*dataWrapper)
	(*item.keyMapP)[item.key] = n
	*h = append(*h, item)
}

func (h *dataHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewLFUCache(capacity int) (c *LFUCache, err error) {
	if capacity <= 0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c = &LFUCache{
		capacity:  capacity,
		cacheData: &dataHeap{},
		keyMap:    make(map[interface{}]int, c.capacity),
	}
	heap.Init(c.cacheData)
	return c, nil
}

func (cache *LFUCache) Add(key, value interface{}) {
	if cache.cacheData == nil {
		cache.cacheData = &dataHeap{}
		cache.keyMap = make(map[interface{}]int, cache.capacity)
		heap.Init(cache.cacheData)
	}
	if pos, ok := cache.keyMap[key]; ok {
		(*cache.cacheData)[pos].frequency++
		heap.Fix(cache.cacheData, pos)
		return
	}
	item := &dataWrapper{key, value, 1, &cache.keyMap}
	heap.Push(cache.cacheData, item)

	if cache.capacity != 0 && cache.cacheData.Len() > cache.capacity {
		d := heap.Pop(cache.cacheData)
		delete(cache.keyMap, d)
	}
}

func (cache *LFUCache) Get(key interface{}) (value interface{}, ok bool) {
	if pos, ok := cache.keyMap[key]; ok {
		(*cache.cacheData)[pos].frequency++
		return (*cache.cacheData)[pos].value, ok
	}
	return nil, ok
}

func (cache *LFUCache) Remove(key interface{}) {
	if pos, ok := cache.keyMap[key]; ok {
		heap.Remove(cache.cacheData, pos)
		delete(cache.keyMap, key)
	}
}

func (cache *LFUCache) Clear() {
	cache.cacheData = &dataHeap{}
	heap.Init(cache.cacheData)
	cache.keyMap = make(map[interface{}]int)
}

func (cache *LFUCache) Len() int {
	return cache.cacheData.Len()
}

// only the top less freq is valid
func (cache *LFUCache) Keys(old2new bool) []interface{} {
	n := len(cache.keyMap)
	keys := make([]interface{}, n)
	var i int
	if old2new {
		i = n
	} else {
		i = 0
	}
	for key := range *cache.cacheData {
		keys[i] = key
		if old2new {
			i--
		} else {
			i++
		}
	}
	return keys
}
