# go-cache

[![Build Status](https://travis-ci.org/XimingCheng/go-cache.png)](https://travis-ci.org/XimingCheng/go-cache)

go-cache is a cache system which support more than just LRU cache, it can define its own TimeToIdleSeconds and TimeToLiveSeconds and more to manage the cache data itself.

## Local Build and Test

get & install

```sh
go get github.com/XimingCheng/go-cache
```

tests

```sh
go test github.com/XimingCheng/go-cache/...
```