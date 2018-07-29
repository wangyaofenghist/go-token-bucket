package bucket

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

//实现一个令牌通
type TokenBucket struct {
	cap    int
	len    int
	lock   sync.Mutex
	ticker *time.Ticker
	ch     chan int
}

//初始化bucket 一旦初始化 不允许更改
func (T *TokenBucket) InitBucket(t time.Duration, cap int) {
	if T.cap != 0 {
		return
	}
	*T = TokenBucket{
		cap:    cap,               //容量
		len:    cap,               //长度
		lock:   sync.Mutex{},      //锁
		ticker: time.NewTicker(t), //定时器
		ch:     make(chan int),    //channel 阻塞 获取每次减少的令牌数
	}
	//启动一个协程对桶进行维护
	go T.bucketDeamon(t)
}

//bucket 守护协程 用于定时填充/充实 token 长度
func (T *TokenBucket) bucketDeamon(t time.Duration) {
	for runt := range T.ticker.C {
		fmt.Println(runt)
		T.lock.Lock()
		//恢复桶的长度与容量相同
		T.len = T.cap

		T.lock.Unlock()
		select {
		case num := <-T.ch:
			T.len -= num
		case <-time.After(time.Millisecond * 100):
			continue

		}
	}
}

//获取token 用于 获取Token 不支持队列获取等方式
func (T *TokenBucket) GetToken(num int) (tokenNum int, err error) {
	tokenNum = 0
	T.lock.Lock()
	if num > T.cap {
		err = errors.New("get num exceed max num: " + strconv.Itoa(T.cap))
		return
	}
	if T.len >= num {
		tokenNum = num
		T.len -= num
	} else {
		err = errors.New("bucker len is less: " + strconv.Itoa(T.len))
	}
	T.lock.Unlock()
	return
}

//等待获取token 直到得到足够的token
func (T *TokenBucket) WaitGetToken(num int) (ok bool) {
	ok = false
	_, err := T.GetToken(num)
	if err != nil {
		T.ch <- num
		ok = true
	}
	return
}

//关闭token 池填充 销毁 bucket
func (T *TokenBucket) Destory() {
	T.ticker.Stop()
}
