package main

import (
	"fmt"
	"github.com/cuckooemm/dqueue"
	"time"
)

func main() {
	var work1 work
	var work2 work
	var work3 work
	_ = dqueue.NewDelayWork(time.Now().Add(time.Second), 0, 0, &work1)
	_ = dqueue.NewDelayWork(time.Now().Add(time.Second), time.Second, 4, &work2)
	_ = dqueue.NewDelayWork(time.Now().Add(time.Second), time.Second, -1, &work3)
	select {}
}

type work struct {
	count int
}

func (w *work) Work() {
	w.count++
	fmt.Printf("current exec time %s, 已执行 %d 次\n", time.Now().Format("2006-01-02 15:03:04"), w.count)
}
