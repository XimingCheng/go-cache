package gocache

import (
	"errors"
	"gocache/type"
)

type CacheManager struct {
	Type              string
	Name              string
	TimeToIdleSeconds int
	TimeToLiveSeconds int
	Capacity          int
	cacheMap          map[string]Cache
}

type Cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
	Clear()
	Len() int
	Keys(old2new bool) []interface{}
}

func New(manager *CacheManager) (cache Cache, err error) {
	if manager.cacheMap == nil {
		manager.cacheMap = make(map[string]Cache)
	}
	if c, ok := manager.cacheMap[manager.Name]; ok {
		return c, errors.New("The cache key map already exist")
	}
	switch manager.Type {
	case "lru":
		cache, err := lru.New(manager.Capacity)
		if err == nil {
			manager.cacheMap[manager.Name] = cache
			return cache, nil
		}
		return nil, err
	default:
		return nil, errors.New("No support cache type")
	}
}
