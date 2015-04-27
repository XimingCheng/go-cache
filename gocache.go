package gocache

import (
	"errors"
	"github.com/XimingCheng/go-cache/cachetype"
	"log"
	"sync"
	"time"
)

type CacheParams struct {
	// the cache cache type name such as LRU FIFO and so on
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
	Eternal bool
	// timer cycle
	//CycleSize int
	// cache capacity
	Capacity int
}

type cacheManager struct {
	// the go cache entity
	cacheMap map[string]*GoCache
	// the go cache params entity
	paramsMap map[string]*CacheParams
	// the cache lock for each cache
	lock map[string]sync.Mutex
}

type GoCache struct {
	// cache entity
	c cache
	// idle timer
	idleTimer map[interface{}]int
	// live timer
	liveTimer map[interface{}]int
	// timer
	timer map[interface{}]*time.Ticker
	// params pointer
	params *CacheParams
	// the lock of the current cache
	Lock *sync.Mutex
}

type cache interface {
	// add key/value into cache
	Add(key, value interface{})
	// get value by key
	Get(key interface{}) (value interface{}, ok bool)
	// remove the key from the cache
	Remove(key interface{})
	// is the key exist
	IsExist(key interface{}) bool
	// clear the cache
	Clear()
	// get length of the cache
	Len() int
	// get slice of the cache keys
	Keys(old2new bool) []interface{}
}

// the global cache data map
var manager cacheManager

func New(params *CacheParams) (gc *GoCache, err error) {
	log.Print("-----------------")
	if params == nil {
		return nil, errors.New("Input cache params invalid")
	}
	if manager.cacheMap == nil {
		manager.cacheMap = make(map[string]*GoCache)
	}
	if manager.paramsMap == nil {
		manager.paramsMap = make(map[string]*CacheParams)
	}
	if manager.lock == nil {
		manager.lock = make(map[string]sync.Mutex)
	}
	if c, ok := manager.cacheMap[params.Name]; ok {
		return c, errors.New("The cache key map " + params.Name + " already exist")
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
		var lock sync.Mutex
		if !params.Eternal {
			gc = &GoCache{
				c:         c,
				idleTimer: make(map[interface{}]int),
				liveTimer: make(map[interface{}]int),
				timer:     make(map[interface{}]*time.Ticker),
				params:    params,
				Lock:      &lock,
			}
		} else {
			gc = &GoCache{
				c:         c,
				idleTimer: nil,
				liveTimer: nil,
				timer:     nil,
				params:    params,
				Lock:      &lock,
			}
		}
		manager.cacheMap[params.Name] = gc
		manager.paramsMap[params.Name] = params
		return gc, err
	}
	return nil, err
}

func timerRun(gc *GoCache, key interface{}) {
	doneChan := make(chan bool)
	for {
		select {
		case <-gc.timer[key].C:
			go func() {
				gc.Lock.Lock()
				defer gc.Lock.Unlock()

				if !gc.c.IsExist(key) {
					log.Printf("key %v not exist", key)
					removeEle(gc, key, doneChan)
				}
				gc.idleTimer[key]++
				gc.liveTimer[key]++
				//log.Printf("key %v idle %v live %v\n", key, gc.idleTimer[key], gc.liveTimer[key])

				if gc.idleTimer[key] >= 1000*gc.params.TimeToIdleSeconds {
					removeEle(gc, key, doneChan)
				} else if gc.liveTimer[key] >= 1000*gc.params.TimeToLiveSeconds {
					removeEle(gc, key, doneChan)
				}
			}()
		case <-doneChan:
			log.Printf("timer %v killed", key)
			return
		}
	}
}

func removeEle(gc *GoCache, key interface{}, doneChan chan bool) {
	log.Printf("Remove key %v idle %v live %v\n", key, gc.idleTimer[key], gc.liveTimer[key])
	gc.c.Remove(key)
	if t, ok := gc.timer[key]; ok {
		log.Printf("stop timer %v", key)
		doneChan <- true
		t.Stop()
		delete(gc.timer, key)
		delete(gc.idleTimer, key)
		delete(gc.liveTimer, key)
	} else {
		log.Printf("WARN -- KEY %v timer not exist", key)
	}
}

func (gc *GoCache) Add(key, value interface{}) {
	gc.Lock.Lock()
	defer gc.Lock.Unlock()

	gc.c.Add(key, value)
	if !gc.params.Eternal {
		gc.idleTimer[key] = 0
		gc.liveTimer[key] = 0
		log.Printf("Add key %v idle %v live %v\n", key, gc.idleTimer[key], gc.liveTimer[key])
		if t, ok := gc.timer[key]; ok {
			t.Stop()
		}
		gc.timer[key] = time.NewTicker(1 * time.Millisecond)
		go timerRun(gc, key)
	} else {
		log.Printf("Add key %v ", key)
	}
}

func (gc *GoCache) Get(key interface{}) (value interface{}, ok bool) {
	gc.Lock.Lock()
	defer gc.Lock.Unlock()

	value, ok = gc.c.Get(key)
	if ok && !gc.params.Eternal {
		gc.idleTimer[key] = 0
		log.Printf("Get key %v idle %v live %v\n", key, gc.idleTimer[key], gc.liveTimer[key])
		if t, ok := gc.timer[key]; ok {
			t.Stop()
		}
		gc.timer[key] = time.NewTicker(1 * time.Millisecond)
		go timerRun(gc, key)
	} else {
		log.Printf("Get key %v ", key)
	}
	return value, ok
}

func (gc *GoCache) Remove(key interface{}) {
	gc.Lock.Lock()
	defer gc.Lock.Unlock()

	gc.c.Remove(key)
	if !gc.params.Eternal {
		if t, ok := gc.timer[key]; ok {
			log.Printf("Del key %v idle %v live %v\n", key, gc.idleTimer[key], gc.liveTimer[key])
			t.Stop()
		}
	} else {
		log.Printf("Del key %v ", key)
	}
}

func (gc *GoCache) Clear() {
	gc.Lock.Lock()
	defer gc.Lock.Unlock()

	gc.c.Clear()
	log.Print("Clear All")
}

func (gc *GoCache) Len() int {
	gc.Lock.Lock()
	defer gc.Lock.Unlock()

	return gc.c.Len()
}

func (gc *GoCache) Keys(old2new bool) []interface{} {
	gc.Lock.Lock()
	defer gc.Lock.Unlock()

	return gc.c.Keys(old2new)
}
