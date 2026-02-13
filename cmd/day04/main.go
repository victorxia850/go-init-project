package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Amount float64
}

// 1. 生产者：只在边界节点打日志
func orderProducer(orderChan chan<- Order, num int) {
	fmt.Printf("🚀 [生产] 开始创建 %d 个订单...\n", num)
	for i := 1; i <= num; i++ {
		orderChan <- Order{ID: i, Amount: rand.Float64() * 100}
		time.Sleep(10 * time.Millisecond)
	}
	close(orderChan)
	fmt.Println("✅ [生产] 订单全部发送完毕并关闭通道。")
}

// 2. 处理器：重点在于展示哪个工人干了哪个活
func orderProcessor(id int, in <-chan Order, out chan<- Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range in {
		// 模拟耗时，让不同工人交替出现
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

		fmt.Printf("  👷 [工人%d] 已处理单据: #%d\n", id, order.ID)
		out <- order
	}
	fmt.Printf("  🔌 [工人%d] 任务领完，退出。\n", id)
}

// 3. 收集器：作为流水线终点，负责汇总
func orderResultCollector(out <-chan Order, done chan<- bool) {
	count := 0
	for range out {
		count++
		// 只有每处理 5 个才打一个进度，或者保持每条一行（如果总数不多）
		if count%5 == 0 {
			fmt.Printf("    📊 [进度] 已收集 %d 个处理结果...\n", count)
		}
	}
	fmt.Printf("🏁 [总结] 流程结束，累计成功处理: %d\n", count)
	done <- true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	orderChan := make(chan Order, 10)
	resultChan := make(chan Order, 10)
	done := make(chan bool)
	var wg sync.WaitGroup

	// 流程开始
	go orderProducer(orderChan, 20)

	// 启动 3 个工人并行工作
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go orderProcessor(i, orderChan, resultChan, &wg)
	}

	// 监控工人们，完工了就关掉结果通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 启动结果收集
	go orderResultCollector(resultChan, done)

	<-done
}
