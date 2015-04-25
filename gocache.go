package gocache

import (
	"errors"
	"github.com/XimingCheng/go-cache/cachetype"
	"sync"
	"time"
)

type CacheParams struct {
	Type string
	// the cache string must be unique
	Name string
	// the maximum number of seconds that an element
	// can exist in the cache without being accessed
	TimeToIdleSeconds int
	// the maximum number of seconds that an element
	// can exist in the cache whether or not is has been accessed
	TimeToLiveSeconds int
	// if set elements are allowed to exist in the cache eternally
	// and none are evicted
	Eternal  bool
	Capacity int
}

type cacheManager struct {
	cacheMap  map[string]*GoCache
	paramsMap map[string]*CacheParams
	lock      sync.Mutex
}

type GoCache struct {
	c         cache
	idleTimer map[interface{}]int
	liveTimer map[interface{}]int
	params    *CacheParams
}

type cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
	Clear()
	Len() int
	Keys(old2new bool) []interface{}
}

// the global cache data map
var manager cacheManager

func New(params *CacheParams) (gc *GoCache, err error) {
	if params == nil {
		return nil, errors.New("Input cache params invalid")
	}
	if manager.cacheMap == nil {
		manager.cacheMap = make(map[string]*GoCache)
	}
	if manager.paramsMap == nil {
		manager.paramsMap = make(map[string]*CacheParams)
	}
	if c, ok := manager.cacheMap[params.Name]; ok {
		return c, errors.New("The cache key map already exist")
	}
	var c cache
	switch params.Type {
	case "lru":
		c, err = cachetype.NewLRUCache(params.Capacity)
	case "fifo":
		c, err = cachetype.NewFIFOCache(params.Capacity)
	case "lfu":
		c, err = cachetype.NewLFUCache(params.Capacity)
	default:
		return nil, errors.New("No support cache type")
	}
	if err == nil {
		if !params.Eternal {
			gc = &GoCache{
				c:         c,
				idleTimer: make(map[interface{}]int),
				liveTimer: make(map[interface{}]int),
				params:    params,
			}
		} else {
			gc = &GoCache{
				c:         c,
				idleTimer: nil,
				liveTimer: nil,
				params:    params,
			}
		}
		manager.cacheMap[params.Name] = gc
		manager.paramsMap[params.Name] = params
		if !params.Eternal {
			// cache name chan
			nameChan := make(chan string)
			go manager.timerRun(nameChan)
			nameChan <- params.Name
		}
		return gc, err
	}
	return nil, err
}

func (manager *cacheManager) timerRun(name chan string) {
	n := <-name
	c := manager.cacheMap[n]
	idle := c.idleTimer
	live := c.liveTimer
	for {
		time.Sleep(time.Second)
		idleToDel := make([]interface{}, 0)
		liveToDel := make([]interface{}, 0)
		for k, v := range idle {
			manager.lock.Lock()
			idle[k] = v + 1
			manager.lock.Unlock()
			if idle[k] >= manager.paramsMap[n].TimeToIdleSeconds {
				idleToDel = append(idleToDel, k)
			}
		}
		for _, k := range idleToDel {
			manager.lock.Lock()
			c.Remove(k)
			manager.lock.Unlock()
		}
		for k, v := range live {
			manager.lock.Lock()
			live[k] = v + 1
			manager.lock.Unlock()
			if live[k] >= manager.paramsMap[n].TimeToLiveSeconds {
				liveToDel = append(liveToDel, k)
			}
		}
		for _, k := range liveToDel {
			manager.lock.Lock()
			c.Remove(k)
			manager.lock.Unlock()
		}
	}
}

func (gc *GoCache) Add(key, value interface{}) {
	gc.c.Add(key, value)
	if !gc.params.Eternal {
		gc.idleTimer[key] = 0
		gc.liveTimer[key] = 0
	}
}

func (gc *GoCache) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = gc.c.Get(key)
	if ok && !gc.params.Eternal {
		gc.idleTimer[key] = 0
	}
	return value, ok
}

func (gc *GoCache) Remove(key interface{}) {
	gc.c.Remove(key)
	if !gc.params.Eternal {
		delete(gc.idleTimer, key)
		delete(gc.liveTimer, key)
	}
}

func (gc *GoCache) Clear() {
	gc.c.Clear()
	if !gc.params.Eternal {
		gc.idleTimer = make(map[interface{}]int)
		gc.liveTimer = make(map[interface{}]int)
	}
}

func (gc *GoCache) Len() int {
	return gc.c.Len()
}

func (gc *GoCache) Keys(old2new bool) []interface{} {
	return gc.c.Keys(old2new)
}
