# 延时队列


## 性能
基于最小堆实现的毫秒级延迟队列  
- 10k 任务   延迟最大 300微秒     众数 80微妙  
- 200k 任务  延迟最大 10毫秒    众数 100微妙  

#### 使用
```go
    // 实现 IDelayJob 接口
    // 到期执行回调函数
    func (w *something) Work() {
    	
    }
    // 任务执行的绝对时间
    func (w *something) SetTime() time.Time {
    	return w.tm
    }
    // tick != 0 则 tick == n+1执行次数 tickDur 为每次间隔时间
    func (w *something) SetTick() (time.Duration, int64) {
    	return w.tickDur, w.tick
    }

    // 初始化 100 = chan容量
    ch := DelayQueueInit(100)
    // 实现IDelayQueue 接口
    var work something 
    ...
    ch <- &work
```
详细使用看测试用例