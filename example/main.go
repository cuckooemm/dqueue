package main

import (
	"fmt"
	"github.com/cuckooemm/dqueue"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 3; i++ {
		var id = i
		/*
		   *  ex 为任务执行的绝对时间
		   *  td 与 tick 成对出现 td 控制距离上次执行间隔，tick 控制次数
		   ×  tick 执行次数为 0 则执行1次 为 n 则执行 n + 1 次
		   ×  tick 为 -1 时持久执行 小于-1 panic
		   ×  work func() 闭包
		*/
		dqueue.NewDelayWork(time.Now().Add(time.Millisecond*time.Duration(rand.Intn(999)+1)), time.Second, -1, func() {
			fmt.Printf("持久任务  work[%d] 执行时间: [%s]\n", id, time.Now().Format(time.RFC3339Nano))
		})
	}
	select {}
}
