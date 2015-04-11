package gocache

import (
	"gocache/type"
)

type CacheManager struct {
	Type              string
	Name              string
	TimeToIdleSeconds int
	TimeToLiveSeconds int
	Capacity          int
	Cache             Cache
}

type Cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
	Clear()
}

func New(manager *CacheManager) (cache Cache, err error) {
	switch manager.Type {
	case "lru":
		return lru.New(manager.Capacity)
	default:
		panic("No support cache type")
	}
}
