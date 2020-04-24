# 延时队列

基于最小堆实现的延迟队列  
- 线程安全
- 有任务需执行唤醒，无任务执行不唤醒
- 支持微秒级响应，(受调度器调度执行影响，平均100微秒级响应)
- 支持执行任务次数

#### 使用

`go get -u github.com/cuckooemm/dqueue`
```go
    // 实现 IDelayJob 接口
    // 到期执行回调函数
func main()  {
	var work1 work
	var work2 work
	var work3 work
        /*
        *  ex 为认为执行的绝对时间
        *  td 与 tick 成对出现 td 控制距离上次执行间隔，tick 控制次数
        ×  tick 执行次数为 0 则执行1次 为 n 则执行 n + 1 次
        ×  tick 为 -1 时持久执行 小于-1 panic
        ×  work 需实现 IDelayJob 的 Work 函数
        */
	dqueue.NewDelayWork(time.Now().Add(time.Second), 0, 0, &work1)
	dqueue.NewDelayWork(time.Now().Add(time.Second), time.Second, 4, &work2)
	dqueue.NewDelayWork(time.Now().Add(time.Second), time.Second, -1, &work3)
	select {
	}
}

type work struct {
	count int
}

func (w *work) Work() {
	w.count++
	fmt.Printf("current exec time %s, 已执行 %d 次\n", time.Now().Format("2006-01-02 15:03:04"),w.count)
}
```
详细使用看测试用例