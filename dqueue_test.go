package dqueue

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	r200k []int64
	m200k sync.Mutex
	r10k  []int64
	m10k  sync.Mutex
	rc10k []int64
	mc10k sync.Mutex
)

func TestDelayPoll(t *testing.T) {
	t.Run("delayWork-10k", func(t *testing.T) {
		ch := DelayQueueInit(100)
		r10k = make([]int64, 0, 10000)
		rand.Seed(time.Now().UnixNano())
		for i := 1; i <= 10000; i++ {
			var work workOnce10k
			work.id = i
			work.tm = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(9999)+1))
			ch <- &work
		}
		for {
			m10k.Lock()
			if len(r10k) == 10000 {
				break
			}
			m10k.Unlock()
			time.Sleep(time.Second * 5)
			fmt.Printf("result %d\n", len(r10k))
		}
		var maxDif int64 = 0
		for i := 0; i < len(r10k); i++ {
			if r10k[i] > maxDif {
				maxDif = r10k[i]
			}
		}
		fmt.Printf("最大延迟 = %d纳秒 - %d微秒\n", maxDif, maxDif/int64(time.Microsecond))
	})

	t.Run("delayWork-200k", func(t *testing.T) {
		ch := DelayQueueInit(100)
		r200k = make([]int64, 0, 200000)
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < 1000; j++ { // 产生 200000 个任务
			for i := 100; i > 0; i-- {
				var work workOnce200k
				work.id = j*100 + i
				work.tm = time.Now().Add(time.Second * time.Duration(rand.Intn(99)+1))
				ch <- &work
			}
			for x := 101; x <= 200; x++ {
				var work workOnce200k
				work.id = j*100 + x
				work.tm = time.Now().Add(time.Millisecond * time.Duration(x*(rand.Intn(999)+1)))
				ch <- &work
			}
		}
		for {
			m200k.Lock()
			if len(r200k) == 200000 {
				break
			}
			m200k.Unlock()
			time.Sleep(time.Second * 5)
			fmt.Printf("result %d\n", len(r200k))
		}
		var maxDif int64 = 0
		for i := 0; i < len(r200k); i++ {
			if r200k[i] > maxDif {
				maxDif = r200k[i]
			}
		}
		fmt.Printf("最大延迟 = %d纳秒 - %d微秒\n", maxDif, maxDif/int64(time.Microsecond))
	})

	t.Run("cycle-10k", func(t *testing.T) {
		ch := DelayQueueInit(100)
		rc10k = make([]int64, 0, 10000)
		rand.Seed(time.Now().UnixNano())
		for i := 1; i <= 10000; i++ {
			var work workCycle10k
			work.id = i
			work.tickDur = time.Millisecond * time.Duration(rand.Intn(9999)+1)
			work.tick = 4 // 执行5次
			work.tm = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(9999)+1))
			ch <- &work
		}
		for {
			mc10k.Lock()
			if len(rc10k) == 50000 {
				break
			}
			mc10k.Unlock()
			time.Sleep(time.Second * 5)
			fmt.Printf("result %d\n", len(rc10k))
		}
		var maxDif int64 = 0
		for i := 0; i < len(rc10k); i++ {
			if rc10k[i] > maxDif {
				maxDif = rc10k[i]
			}
		}
		fmt.Printf("执行次数 %d\n", len(rc10k))
		//fmt.Printf("最大延迟 = %d纳秒 - %d微秒\n", maxDif, maxDif/int64(time.Microsecond))
	})
}

// 实现 IDelayJob 接口
type workOnce10k struct {
	id      int
	tickDur time.Duration
	tick    int64
	tm      time.Time
}

func (w *workOnce10k) Work() {
	// goroutine 受调度影响
	var (
		cur  = time.Now().UnixNano()
		exec = w.tm.UnixNano()
		dif  = cur - exec
	)
	m10k.Lock()
	r10k = append(r10k, dif)
	defer m10k.Unlock()
	fmt.Printf("workid - %d 已执行 当前执行时间[%d] - 应执行时间[%d], 延迟时间: %d纳秒 -  %d微秒\n", w.id, cur, exec, dif, dif/int64(time.Microsecond))
}

func (w *workOnce10k) SetTime() time.Time {
	return w.tm
}
func (w *workOnce10k) SetTick() (time.Duration, int64) {
	return w.tickDur, w.tick
}

type workCycle10k struct {
	id      int
	tickDur time.Duration
	tick    int64
	tm      time.Time
}

func (w *workCycle10k) Work() {
	// goroutine 受调度影响
	var (
		cur  = time.Now().UnixNano()
		exec = w.tm.UnixNano()
		dif  = cur - exec
	)
	mc10k.Lock()
	rc10k = append(rc10k, dif)
	defer mc10k.Unlock()
	//fmt.Printf("workid - %d 已执行 当前执行时间[%d] - 应执行时间[%d], 延迟时间: %d纳秒 -  %d微秒\n", w.id, cur, exec, dif, dif/int64(time.Microsecond))
}

func (w *workCycle10k) SetTime() time.Time {
	return w.tm
}
func (w *workCycle10k) SetTick() (time.Duration, int64) {
	return w.tickDur, w.tick
}

type workOnce200k struct {
	id      int
	tickDur time.Duration
	tick    int64
	tm      time.Time
}

func (w *workOnce200k) Work() {
	// goroutine 受调度影响
	var (
		cur  = time.Now().UnixNano()
		exec = w.tm.UnixNano()
		dif  = cur - exec
	)
	m200k.Lock()
	r200k = append(r200k, dif)
	defer m200k.Unlock()
	fmt.Printf("workid - %d 已执行 当前执行时间[%d] - 应执行时间[%d], 延迟时间: %d纳秒 -  %d微秒\n", w.id, cur, exec, dif, dif/int64(time.Microsecond))
}

// 实现 IDelayJob 接口
func (w *workOnce200k) SetTime() time.Time {
	return w.tm
}
func (w *workOnce200k) SetTick() (time.Duration, int64) {
	return w.tickDur, w.tick
}
