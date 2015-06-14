package gocache

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type httpCacheData struct {
	key   string
	value string
}

type httpKeyData struct {
	key string
}

var gc *GoCache

func RunHttpCache(port int, params *CacheParams) error {
	var err error
	gc, err = New(params)
	if err != nil {
		return err
	}

	http.HandleFunc("/add", goCacheAddHandler)
	http.HandleFunc("/remove", goCacheRemoveHandler)
	http.HandleFunc("/get", goCacheGetHandler)
	http.HandleFunc("/clear", goCacheClearHandler)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func goCacheAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "PUT" {
		io.WriteString(w, "{\"ret\":\"add must be call by protocal PUT\"}")
		return
	}
	if gc == nil {
		io.WriteString(w, "{\"ret\":\"go cache init failed\"}")
		return
	}
	var d httpCacheData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&d)
	if err != nil {
		panic(err)
	}
	gc.Add(d.key, d.value)
	io.WriteString(w, "{\"ret\":\"go cache add ok\"}")
}

func goCacheRemoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "PUT" {
		io.WriteString(w, "{\"ret\":\"remove must be call by protocal PUT\"}")
		return
	}
	if gc == nil {
		io.WriteString(w, "{\"ret\":\"go cache init failed\"}")
		return
	}
	var d httpKeyData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&d)
	if err != nil {
		panic(err)
	}
	gc.Remove(d.key)
	io.WriteString(w, "{\"ret\":\"go cache remove ok\"}")
}

// get protocal (the length of the key must not very long)
func goCacheGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		io.WriteString(w, "{\"ret\":\"get must be call by protocal GET\"}")
		return
	}
	if gc == nil {
		io.WriteString(w, "{\"ret\":\"go cache init failed\"}")
		return
	}
	key := r.URL.Query().Get("key")
	if len(key) > 0 {
		v, ok := gc.Get(key)
		if !ok {
			io.WriteString(w, "{\"ret\":\"not exsit the key\"}")
			return
		}
		msg := "{\"ret\":\"" + v.(string) + "\"}"
		io.WriteString(w, msg)
	}
}

func goCacheClearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "PUT" {
		io.WriteString(w, "{\"ret\":\"clear must be call by protocal PUT\"}")
		return
	}
	if gc == nil {
		io.WriteString(w, "{\"ret\":\"go cache init failed\"}")
		return
	}
	gc.Clear()
	io.WriteString(w, "{\"ret\":\"go cache clear ok\"}")
}
