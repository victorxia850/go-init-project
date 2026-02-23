package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value int64
}

func (c *Counter) Increment() int64 {
	// AddInt64 会返回修改后的新值，我们可以利用这个值来做日志判断
	return atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func main() {
	var counter Counter
	var upCount, downCount int64 // 专门记录执行次数
	var wg sync.WaitGroup

	// 点赞
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
			atomic.AddInt64(&upCount, 1) // 记录成功次数
			fmt.Printf("****** 点赞任务 %d 完成 ******\n", i)
		}()
	}

	// 取消
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Decrement()
			atomic.AddInt64(&downCount, 1) // 记录成功次数
			fmt.Printf("====== 取消任务 %d 完成 ======\n", i)

		}()
	}

	wg.Wait()
	fmt.Printf("统计结果 -> 点赞执行: %d次, 取消执行: %d次, 最终值: %d\n",
		atomic.LoadInt64(&upCount),
		atomic.LoadInt64(&downCount),
		counter.Value())
}
