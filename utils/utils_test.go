package utils

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestSubStrAndTrim(t *testing.T) {
	str := "/home/page/user1"
	step := "/home/page"
	if s := SubStrAndTrim(str, step); s != "user1" {
		t.Error(s)
	}
	step = "/home"
	if s := SubStrAndTrim(str, step); s != "page/user1" {
		t.Error(s)
	}
	step = "/home/page/"
	if s := SubStrAndTrim(str, step); s != "user1" {
		t.Error(s)
	}
	step = "/home/page/user1"
	if s := SubStrAndTrim(str, step); s != "" {
		t.Error(s)
	}
}

func TestStringToBytes(t *testing.T) {
	testString := "Hello, World!"
	start := time.Now()
	expectedBytes := []byte("Hello, World!")
	since := time.Since(start)
	now := time.Now()
	resultBytes := StringToBytes(testString)
	duration := time.Since(now)
	now1 := time.Now()
	resultBytes1 := StringToBytes2(testString)
	duration1 := time.Since(now1)
	fmt.Printf("[]byte() time.since:%v\nStringToBytes time.since:%v\nStringToBytes2 time.since:%v\n", since, duration, duration1)
	fmt.Println(expectedBytes, resultBytes, resultBytes1)
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Errorf("Expected: %v, but got: %v", expectedBytes, resultBytes)
	}
}
func TestTrimQuery(t *testing.T) {
	tt := "a?aaa=1&bb=2"
	t2 := "a"
	t3 := "aaa?=bda[mas]"
	t4 := "aaa?aa[sd]"
	tt1 := TrimQuery(tt)
	tt2 := TrimQuery(t2)
	tt3 := TrimQuery(t3)
	tt4 := TrimQuery(t4)
	fmt.Println(TrimQuery(tt1))
	fmt.Println(TrimQuery(tt2))
	fmt.Println(TrimQuery(tt3))
	fmt.Println(TrimQuery(tt4))
}
