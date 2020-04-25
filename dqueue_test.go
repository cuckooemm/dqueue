package dqueue

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestDelay10k(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 1; i <= 10000; i++ {
		var id = i
		var exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(9999)+1))
		NewDelayWork(exec, 0, 0, func() {
			var curTime = time.Now()
			fmt.Printf("work[%d] 应执行时间: [%s] 实际执行时间: [%s] 延迟 [%d] 纳秒\n",
				id, exec.Format(time.RFC3339Nano), curTime.Format(time.RFC3339Nano), curTime.Sub(exec).Nanoseconds())
			wg.Done()
		})
	}
	wg.Wait()
}

func TestDelay200k(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	wg.Add(200000)
	for j := 0; j < 1000; j++ { // 产生 200000 个任务
		for i := 100; i > 0; i-- {
			var id = j*1000 + i
			var exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(99)+1))
			NewDelayWork(exec, 0, 0, func() {
				var curTime = time.Now()
				fmt.Printf("work[%d] 应执行时间: [%s] 实际执行时间: [%s] 延迟 [%d] 纳秒\n",
					id, exec.Format(time.RFC3339Nano), curTime.Format(time.RFC3339Nano), curTime.Sub(exec).Nanoseconds())
				wg.Done()
			})
		}
		for i := 101; i <= 200; i++ {
			var id = j*1000 + i
			var exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(999)+1))
			NewDelayWork(exec, 0, 0, func() {
				var curTime = time.Now()
				fmt.Printf("work[%d] 应执行时间: [%s] 实际执行时间: [%s] 延迟 [%d] 纳秒\n",
					id, exec.Format(time.RFC3339Nano), curTime.Format(time.RFC3339Nano), curTime.Sub(exec).Nanoseconds())
				wg.Done()
			})
		}
	}
	wg.Wait()
}

func TestCycle(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 3; i++ {
		var id = i
		var exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(999)+1))
		NewDelayWork(exec, time.Second, -1, func() {
			var curTime = time.Now()
			fmt.Printf("持久任务  work[%d] 应执行时间: [%s] 实际执行时间: [%s] 延迟 [%d] 纳秒\n",
				id, exec.Format(time.RFC3339Nano), curTime.Format(time.RFC3339Nano), curTime.Sub(exec).Nanoseconds())
			exec = exec.Add(time.Second)
		})
	}
	// 执行30次就退出
	time.Sleep(10 * time.Second)
}
func TestCycleCount5(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	wg.Add(15)
	for i := 1; i <= 3; i++ {
		var id = i
		var exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(999)+1))
		NewDelayWork(exec, time.Second, 4, func() {
			var curTime = time.Now()
			fmt.Printf("循环任务  work[%d] 应执行时间: [%s] 实际执行时间: [%s] 延迟 [%d] 纳秒\n",
				id, exec.Format(time.RFC3339Nano), curTime.Format(time.RFC3339Nano), curTime.Sub(exec).Nanoseconds())
			exec = exec.Add(time.Second)
			wg.Done()
		})
	}
	wg.Wait()
}
