# go-cache

[![Build Status](https://travis-ci.org/XimingCheng/go-cache.png)](https://travis-ci.org/XimingCheng/go-cache)

go-cache is a cache system which support more than just LRU cache, it can define its own TimeToIdleSeconds and TimeToLiveSeconds and more to manage the cache data itself.
User or Developer can use the cache library to speed the database retrieval or query.

## Local Build and Test

get & install

```sh
go get github.com/XimingCheng/go-cache
```

tests

```sh
go test github.com/XimingCheng/go-cache/...
```

## Features

* Support LRU/LFU/FIFO/TwoQueue cache type
* Support use-defined cache parameters
* Goroutine cache key management
* Golang function invoke with reflection by gocache

## Example

```go
func add(a, b int) int {
    // simulate the database query time cost
    time.Sleep(1 * time.Second)
    return a + b
}

func Test() {
    // user defined cache parameters
    // user can choose its cache type/timer type/size
    params := &CacheParams{
        Type:              "lru",
        Name:              "testlruReflect",
        TimeToIdleSeconds: 3,
        TimeToLiveSeconds: 5,
        Eternal:           false,
        Capacity:          5,
        ExtendParam:       nil,
    }

    // user can regsiter his own function
    err := RegsiterFunction(add, params)
    if err != nil {
        t.Fatalf("RegsiterFunction err: %v", err)
    }

    // get the invoke start and end time (first time)
    start1 := time.Now().Unix()
    // outputs -> 7
    outputs, _ := Invoke(add, 3, 4)
    end1 := time.Now().Unix()
    cost1 := end1 - start1

    // get the invoke start and end time (second time)
    start2 := time.Now().Unix()
    // outputs -> 7
    outputs, _ = Invoke(add, 3, 4)
    end2 := time.Now().Unix()
    cost2 := end2 - start2

    // cost1 > cost2, second time is faster than the first time
    fmt.Printf("cost1 %v, cost2 %v", cost1, cost2)

    UnRegsiterFunction(add)
}
```

