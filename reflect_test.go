package gocache

import (
	"testing"
	"time"
)

func add(a, b int) int {
	time.Sleep(1 * time.Second)
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func returnTwoVals(a int) (int, int) {
	return a, a + 1
}

func TestReflectBasic(t *testing.T) {
	params := &CacheParams{
		Type:              "lru",
		Name:              "testlruReflect",
		TimeToIdleSeconds: 3,
		TimeToLiveSeconds: 5,
		Eternal:           false,
		Capacity:          5,
		ExtendParam:       nil,
	}

	params2 := &CacheParams{
		Type:              "lru",
		Name:              "testlruReflect2",
		TimeToIdleSeconds: 3,
		TimeToLiveSeconds: 5,
		Eternal:           false,
		Capacity:          5,
		ExtendParam:       nil,
	}

	err := RegsiterFunction(add, params)
	if err != nil {
		t.Fatalf("RegsiterFunction err: %v", err)
	}
	err = RegsiterFunction(returnTwoVals, params2)
	if err != nil {
		t.Fatalf("RegsiterFunction err: %v", err)
	}

	start1 := time.Now().Unix()
	outputs, e := Invoke(add, 3, 4)
	end1 := time.Now().Unix()
	if e != nil {
		t.Fatalf("Invoke err %v", e)
	}
	if len(outputs) != 1 {
		t.Fatalf("outputs number is not 1, len = %v", len(outputs))
	}
	if outputs[0] != 7 {
		t.Fatalf("outputs is not 7, is %v", outputs[0])
	}
	cost1 := end1 - start1
	start2 := time.Now().Unix()
	outputs, e = Invoke(add, 3, 4)
	end2 := time.Now().Unix()
	cost2 := end2 - start2
	t.Logf("cost1 %v, cost2 %v", cost1, cost2)
	if cost1 < cost2 {
		t.Fatalf("add Invoke in the sencod time should cost little time")
	}

	_, e1 := Invoke(sub, 3, 4)
	if e1 == nil {
		t.Fatalf("Inoke sub must be failed")
	}

	outputs2, e2 := Invoke(returnTwoVals, 5)
	if e2 != nil {
		t.Fatalf("Invoke err %v", e2)
	}
	if len(outputs2) != 2 {
		t.Fatalf("outputs number is not 2, len = %v", len(outputs2))
	}
	if outputs2[0] != 5 {
		t.Fatalf("outputs is not 5, is %v", outputs2[0])
	}
	if outputs2[1] != 6 {
		t.Fatalf("outputs is not 6, is %v", outputs2[1])
	}

	err = UnRegsiterFunction(add)
	if err != nil {
		t.Fatalf("UnRegsiterFunction err: %v", err)
	}
	err = UnRegsiterFunction(returnTwoVals)
	if err != nil {
		t.Fatalf("UnRegsiterFunction err: %v", err)
	}
	err = UnRegsiterFunction(sub)
	if err == nil {
		t.Fatalf("sub function not regsitered, UnRegsiterFunction failed")
	}
}
