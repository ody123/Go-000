package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type metrics struct {
	success int32
	fail    int32
}

type bucket struct {
	data        metrics
	windowStart int64
}

type RollingNumber struct {
	buckets []*bucket
	size    int64
	width   int64
	tail    int64
	mux     sync.RWMutex
}

func NewRollingNumber(size, width int64) *RollingNumber {
	return &RollingNumber{
		size:    size,
		width:   width,
		buckets: make([]*bucket, size),
		tail:    0,
	}
}

func (rp *RollingNumber) getCurrent() *bucket {
	rp.mux.Lock()
	defer rp.mux.Unlock()

	current := time.Now().Unix()
	last := rp.buckets[rp.tail]
	// 超出最大时间范围清空整个数组
	if last != nil && current-last.windowStart >= rp.size*rp.width {
		rp.buckets = make([]*bucket, rp.size)
		rp.tail = 0
		last = nil
	}
	// 数组为空时直接新建桶
	if last == nil {
		bk := &bucket{
			data:        metrics{},
			windowStart: current,
		}
		rp.buckets[rp.tail] = bk
		return bk
	}
	// 数组不为空时
	for {
		// 尾桶在当前时间范围内退出循环
		if current < last.windowStart+rp.width {
			break
		}
		// 更新尾桶
		rp.tail++
		if rp.tail >= rp.size {
			rp.tail = 0
		}
		last = &bucket{
			data:        metrics{},
			windowStart: last.windowStart + rp.width,
		}
		rp.buckets[rp.tail] = last
	}
	return last
}

func (rp *RollingNumber) incrSuccess() {
	bk := rp.getCurrent()
	atomic.AddInt32(&bk.data.success, 1)
}

func (rp *RollingNumber) incrFail() {
	bk := rp.getCurrent()
	atomic.AddInt32(&bk.data.fail, 1)
}

func main() {
	// 滑动窗口计数器
	// 用于给熔断器提供数据依据
	rw := NewRollingNumber(2, 1)
	fmt.Println(time.Now().Unix())

	//test 1
	rw.incrSuccess()
	time.Sleep(time.Second * 1)
	rw.incrSuccess()
	time.Sleep(time.Second * 1)
	rw.incrSuccess()
	time.Sleep(time.Second * 1)
	fmt.Printf("%+v,%+v\n", rw.buckets[0], rw.buckets[1])
}
