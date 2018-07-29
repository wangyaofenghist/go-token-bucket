package main

import (
	"fmt"
	"localhostTest/token-bucket/bucket"
	"runtime"
	"time"
)

var tokenBucket bucket.TokenBucket

func init() {
	//初始化 一个桶
	tokenBucket.InitBucket(time.Second, 1000)
}

func testFunc() {
	var isOk bool
	for i := 1; i < 100000; i++ {
		isOk = tokenBucket.WaitGetToken(1)
		go testFunc3(i)
		fmt.Println(tokenBucket, isOk)
	}
}
func testFunc2() {
	var isOk bool
	for i := 1; i < 100000; i++ {
		isOk = tokenBucket.WaitGetToken(1)
		go testFunc3(i)
		fmt.Println(tokenBucket, isOk)
	}
}
func testFunc3(i int) {
	time.Sleep(time.Second)
	fmt.Println(i, "this is test runtime", runtime.NumGoroutine())
}

func main() {

	num, err := tokenBucket.GetToken(10)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < num; i++ {
		println(i, "test 1")
	}
	go testFunc()
	go testFunc2()

	time.Sleep(time.Second * 100)
}
