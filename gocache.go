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
	TimeToIdleSeconds int64
	// the maximum number of seconds that an element
	// can exist in the cache whether or not is has been accessed
	TimeToLiveSeconds int64
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
	//add timer
	addTime map[interface{}]time.Time
	// timer
	timer *ConcurrenyMap
	// params pointer
	params *CacheParams
	// the lock of the current cache
	lock *sync.Mutex
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
			t, _ := NewConcurrenyMap(params.Capacity)
			gc = &GoCache{
				c:       c,
				timer:   t,
				addTime: make(map[interface{}]time.Time),
				params:  params,
				lock:    &lock,
			}
		} else {
			gc = &GoCache{
				c:       c,
				timer:   nil,
				addTime: nil,
				params:  params,
				lock:    &lock,
			}
		}
		manager.cacheMap[params.Name] = gc
		manager.paramsMap[params.Name] = params
		return gc, err
	}
	return nil, err
}

func timerRun(gc *GoCache, key interface{}) {
	if t, ok := gc.timer.Get(key); ok {
		<-t.(*time.Timer).C
		removeEle(gc, key)
	}
}

func removeEle(gc *GoCache, key interface{}) {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	gc.c.Remove(key)
	if ti, ok := gc.timer.Get(key); ok {
		t := ti.(*time.Timer)
		log.Printf("stop timer %v gc.timer %v", key, gc.timer)
		t.Stop()
		gc.timer.Delete(key)
	}
	log.Printf("delete key %v gc.addTime %v", key, gc.addTime)
	delete(gc.addTime, key)
	log.Printf("removeEle done!!!")
}

func (gc *GoCache) Add(key, value interface{}) {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	gc.c.Add(key, value)
	if !gc.params.Eternal {
		//element join time
		min := gc.params.TimeToIdleSeconds * 1000
		if min > gc.params.TimeToLiveSeconds*1000 {
			min = gc.params.TimeToLiveSeconds * 1000
		}
		gc.addTime[key] = time.Now()
		if ti, ok := gc.timer.Get(key); ok {
			t := ti.(*time.Timer)
			t.Reset(time.Duration(min) * time.Millisecond)
		} else {
			gc.timer.Add(key, time.NewTimer(time.Duration(min)*time.Millisecond))
			go timerRun(gc, key)
		}
		log.Printf("add eternal key %v", key)
	} else {
		log.Printf("Add key %v ", key)
	}
}

func (gc *GoCache) Get(key interface{}) (value interface{}, ok bool) {
	gc.lock.Lock()
	defer gc.lock.Unlock()

	value, ok = gc.c.Get(key)
	if ok && !gc.params.Eternal {
		//gc.idleTimer[key] = 0
		liveTime := gc.params.TimeToLiveSeconds*1000 - int64((time.Now().Sub(gc.addTime[key]))/1000000)
		if liveTime <= 0 {
			return nil, false
		}
		if liveTime > gc.params.TimeToIdleSeconds*1000 {
			liveTime = gc.params.TimeToIdleSeconds * 1000
		}
		if ti, ok := gc.timer.Get(key); ok {
			t := ti.(*time.Timer)
			t.Reset(time.Duration(liveTime) * time.Millisecond)
		}
		log.Printf("Get eternal key %v ", key)
	} else {
		log.Printf("Get key %v ", key)
	}
	return value, ok
}

func (gc *GoCache) Remove(key interface{}) {
	//removeEle(gc, key)
	gc.lock.Lock()
	defer gc.lock.Unlock()
	gc.c.Remove(key)
	// if !gc.params.Eternal {
	// 	if t, ok := gc.timer[key]; ok {
	// 		t.Stop()
	// 		delete(gc.timer, key)
	// 		delete(gc.addTime, key)
	// 	}
	// }
	// log.Printf("stop timer %v", key)
}

func (gc *GoCache) Clear() {
	gc.lock.Lock()
	defer gc.lock.Unlock()

	gc.c.Clear()
	// for k, _ := range gc.addTime {
	// 	delete(gc.addTime, k)
	// 	if t, ok := gc.timer[k]; ok {
	// 		t.Stop()
	// 	}
	// 	delete(gc.timer, k)
	// }
	// log.Print("Clear All")
}

func (gc *GoCache) Len() int {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	return gc.c.Len()
}

func (gc *GoCache) Keys(old2new bool) []interface{} {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	return gc.c.Keys(old2new)
}
