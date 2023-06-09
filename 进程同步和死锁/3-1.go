package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
3.1 编程模拟生产者/消费者问题
	以生产者/消费者模型为根据，编写一个图形界面程序，创建 n 个线程，使用 windows 信号量机制，模拟生产者和消费者的工作流程。
*/

/*
	创建5个生产者和5个消费者线程，它们共享一个大小为10的仓库。
	生产者随机生产物品并将它们放入仓库中，消费者从仓库中取出物品并消费它们。
	通过使用信号量机制，程序保证了仓库中同时最多只能有10个物品，并且在仓库已满或已空时，生产者和消费者会被阻塞等待。
*/

var (
	wg    sync.WaitGroup   // 控制程序结束
	mutex sync.Mutex       // 互斥锁，保证同时只有一个线程对仓库进行操作
	full  chan int         // 信号量，表示仓库中已有的物品数量
	empty chan int         // 信号量，表示仓库中空余的物品数量
	q     = make([]int, 0) // 存放物品的队列
)

func producer(id int, numItems int) {
	defer wg.Done()

	for i := 0; i < numItems; i++ {
		item := rand.Intn(1000)
		<-empty             // 将空槽个数减一
		mutex.Lock()        // 进入临界区
		q = append(q, item) // 将生产的数据放到缓存区
		fmt.Printf("Producer %d produced item %d\n", id, item)
		mutex.Unlock()                                                // 离开临界区
		full <- 1                                                     // 将满槽个数加一
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // 生产物品需要一定时间
	}
}

func consumer(id int, numItems int) {
	defer wg.Done()

	for i := 0; i < numItems; i++ {
		<-full       // 将满槽个数减一
		mutex.Lock() // 进入临界区
		item := q[0] // 从缓冲区取出数据
		q = q[1:]
		fmt.Printf("Consumer %d consumed item %d\n", id, item)
		mutex.Unlock()                                                // 离开临界区
		empty <- 1                                                    // 将空槽个数加一
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // 消费物品需要一定时间
	}
}

func main() {
	var maxItems, numItems, numProducers, numConsumers int

	fmt.Print("请输入仓库最大容量(例如 10)：")
	fmt.Scanln(&maxItems)

	fmt.Print("请输入每个生产者和消费者需要生产/消费的物品数目(例如 20)：")
	fmt.Scanln(&numItems)

	fmt.Print("请输入生产者线程数(例如 5)：")
	fmt.Scanln(&numProducers)

	fmt.Print("请输入消费者线程数(例如 5)：")
	fmt.Scanln(&numConsumers)

	full = make(chan int, maxItems)  // 初始化信号量，表示仓库中已有的物品数量
	empty = make(chan int, maxItems) // 初始化信号量，表示仓库中空余的物品数量

	for i := 0; i < maxItems; i++ {
		empty <- 1 // 将空槽个数初始化为maxItems
	}

	for i := 0; i < numProducers; i++ {
		wg.Add(1)
		go producer(i, numItems)
	}

	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go consumer(i, numItems)
	}

	wg.Wait() // 等待所有线程结束
}
