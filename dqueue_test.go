package dqueue

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDelayPoll(t *testing.T) {
	t.Run("delayWork-10k", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		for i := 1; i <= 10000; i++ {
			var work work
			work.id = i
			work.a = 2
			work.b = 3
			work.exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(9999)+1)).UnixNano()
			_ = NewDelayWork(time.Unix(0, work.exec), 0, 0, &work)
		}
		time.Sleep(time.Second * 15)
	})

	t.Run("delayWork-200k", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		for j := 0; j < 1000; j++ { // 产生 200000 个任务
			for i := 100; i > 0; i-- {
				var work work
				work.id = j*100 + i
				work.a = 2
				work.b = 3
				work.exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(99)+1)).UnixNano()
				_ = NewDelayWork(time.Unix(0, work.exec), 0, 0, &work)
			}
			for x := 101; x <= 200; x++ {
				var work work
				work.id = j*100 + x
				work.a = 2
				work.b = 3
				work.exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(999)+1)).UnixNano()
				_ = NewDelayWork(time.Unix(0, work.exec), 0, 0, &work)
			}
		}
		time.Sleep(time.Minute)
	})

	t.Run("cycle", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		for i := 1; i <= 3; i++ {
			var work work
			work.id = i
			work.a = 2
			work.b = 3
			work.continu = true
			work.dif = time.Second
			work.exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(9999)+1)).UnixNano()
			_ = NewDelayWork(time.Unix(0, work.exec), work.dif, -1, &work)
		}
		time.Sleep(time.Second * 30)
	})

	t.Run("cycle-count-5", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		for i := 1; i <= 3; i++ {
			var work work
			work.id = i
			work.a = 2
			work.b = 3
			work.dif = time.Second
			work.exec = time.Now().Add(time.Millisecond * time.Duration(rand.Intn(9999)+1)).UnixNano()
			_ = NewDelayWork(time.Unix(0, work.exec), work.dif, 4, &work)
		}
		time.Sleep(time.Minute)
	})
}

// 实现 IDelayJob 接口
type work struct {
	id      int
	continu bool
	exec    int64
	dif     time.Duration
	a, b    int
}

func (w *work) Work() {
	// goroutine 受调度影响
	var cur = time.Now().UnixNano()
	fmt.Printf("workid - [%d] 已执行 当前执行时间[%d] - 应执行时间[%d] = 延迟[%d] %s a +b = %d\n", w.id, cur, w.exec, w.exec-cur, b2s(w.continu), w.a+w.b)
	w.exec += int64(w.dif)
}

func b2s(b bool) string {
	if b {
		return " - 持久任务 - "
	}
	return ""
}
