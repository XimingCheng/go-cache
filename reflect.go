package gocache

import (
	"encoding/json"
	"errors"
	"reflect"
)

func RegsiterFunction(f interface{}, params *CacheParams) error {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return errors.New("RegsiterFunction input is not a function")
	}

	gc, err := New(params)
	if err != nil {
		return err
	}

	manager.cacheFuncMap[f] = gc
	return nil
}

func UnRegsiterFunction(f interface{}) error {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return errors.New("RegsiterFunction input is not a function")
	}

	if gc, ok := manager.cacheFuncMap[f]; ok {
		name := gc.params.Name
		gc.Clear()
		delete(manager.cacheMap, name)
		delete(manager.paramsMap, name)
		delete(manager.cacheFuncMap, f)
		return nil
	}
	return errors.New("no such function regsitered")
}

func Invoke(f interface{}, inputs ...interface{}) (outputs []interface{}, err error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, errors.New("RegsiterFunction input is not a function")
	}

	if gc, ok := manager.cacheFuncMap[f]; ok {
		inputsArgs := make([]interface{}, len(inputs))
		for idx, input := range inputs {
			inputsArgs[idx] = input
		}
		jsonInputBytes, e := json.Marshal(inputsArgs)
		if e != nil {
			return nil, e
		}
		jsonInputs := string(jsonInputBytes)
		if gc.IsExist(jsonInputs) {
			value, e := gc.Get(jsonInputs)
			if !e {
				return nil, errors.New("cache get failed")
			} else {
				return value.([]interface{}), nil
			}
		} else {
			inputsData := make([]reflect.Value, len(inputs))
			for idx, input := range inputs {
				inputsData[idx] = reflect.ValueOf(input)
			}
			outputs = make([]interface{}, t.NumOut())
			var outs []reflect.Value
			if t.IsVariadic() {
				outs = reflect.ValueOf(f).CallSlice(inputsData)
			} else {
				outs = reflect.ValueOf(f).Call(inputsData)
			}
			for idx, o := range outs {
				outputs[idx] = o.Interface()
			}
			gc.Add(jsonInputs, outputs)
			return outputs, nil
		}
	}
	return nil, errors.New("cacheManager did not exist the reg function")
}
