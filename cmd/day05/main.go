package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== 1. Select 多路复用与超时 ===")
	demoSelect()

	fmt.Println("\n=== 2. Ticker 定时周期任务 ===")
	demoTicker()

	fmt.Println("\n=== 3. Worker Pool 工作池 ===")
	demoWorkerPool()
}

// --- 模块 1: Select 多路复用 ---
func demoSelect() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() { time.Sleep(1 * time.Second); ch1 <- "信号 A" }()
	go func() { time.Sleep(2 * time.Second); ch2 <- "信号 B" }()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("接收:", msg1)
		case msg2 := <-ch2:
			fmt.Println("接收:", msg2)
		case <-time.After(500 * time.Millisecond): // 如果 500ms 没收到，就报超时
			fmt.Println("🚨 响应超时！")
		}
	}
}

// --- 模块 2: Ticker 周期打点 ---
func demoTicker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop() // 最佳实践：确保退出时关闭定时器
	done := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		close(done) // 使用 close 通知停止，比发 bool 更优雅
	}()

	fmt.Println("定时器启动...")
loop:
	for {
		select {
		case t := <-ticker.C:
			fmt.Println("⏰ 打点:", t.Format("15:04:05"))
		case <-done:
			fmt.Println("🛑 定时任务停止")
			break loop // 跳出指定的 for 循环
		}
	}
}

// --- 模块 3: Worker Pool 工作池 ---
func demoWorkerPool() {
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup

	// 启动 3 个 worker
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}

	// 发送任务
	numJobs := 5
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // 发完任务必须关闭，否则 worker 的 range 不会停止

	// 监控协程：等所有 worker 干完，关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果：通过 range 自动处理关闭
	for res := range results {
		fmt.Printf("✅ 处理结果: %d\n", res)
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("  👷 Worker %d 处理任务 %d\n", id, j)
		time.Sleep(500 * time.Millisecond)
		results <- j * 2
	}
}
