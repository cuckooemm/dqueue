package cdelay_queue

import (
	"time"
)

const (
	maxTimeDuration = 1<<63 - 1
)

var (
	qPoll queuePoll
)

type IDelayJob interface {
	// 执行的回调函数
	Work()
	// 设置任务执行的绝对时间
	SetTime() time.Time
	// 设置任务的执行次数与每次执行与上次执行的时间差 0 == 1次 n == n+1次
	SetTick() (td time.Duration, n int64)
}

type delayJob struct {
	tickDur time.Duration
	tick    int64
	tm      int64
	work    IDelayJob
}

type queuePoll struct {
	size     int
	maxSize  int
	workPoll []delayJob
}

func DelayPollStart(chanSize int) chan<- IDelayJob {
	if chanSize < 1 {
		chanSize = 100
	}
	var ch = make(chan IDelayJob, chanSize)
	go start(ch)
	return ch
}

func start(ch <-chan IDelayJob) {
	var (
		tk  = time.NewTimer(maxTimeDuration)
		job delayJob
	)
	// 初始化堆空间
	qPoll = queuePoll{
		size:     -1,
		maxSize:  -1,
		workPoll: make([]delayJob, 0, 16),
	}

	for {
		job = qPoll.getJob()
		curTime := time.Now().UnixNano()
		if !tk.Stop() && len(tk.C) > 0 {
			<-tk.C
		}
		//if job.tm - curTime < 99999 {
		//	job.work.Work()
		//	if job.tick > 0 {
		//		job.tm += int64(job.tickDur)
		//		job.tick--
		//		qPoll.addJob(job)
		//	}
		//	qPoll.deleteTop()
		//	continue
		//}
		tk.Reset(time.Duration(job.tm - curTime))
		select {
		case <-tk.C:
			go job.work.Work()
			if job.tick > 0 {
				job.tm += int64(job.tickDur)
				job.tick--
				qPoll.addJob(job)
			}
			qPoll.deleteTop()

		case q := <-ch:
			var dq = delayJob{
				tm:   q.SetTime().UnixNano(),
				work: q,
			}
			dq.tickDur, dq.tick = q.SetTick()
			qPoll.addJob(dq)
		}
	}
}
func (p *queuePoll) addJob(q delayJob) {
	p.size++
	if p.size > p.maxSize {
		p.workPoll = append(p.workPoll, q)
		p.maxSize++
	} else {
		p.workPoll[p.size] = q
	}
	p.upFloat(p.size)
}

func (p *queuePoll) getJob() delayJob {
	if p.size < 0 {
		return delayJob{tm: maxTimeDuration}
	}
	return p.workPoll[0]
}

func (p *queuePoll) deleteTop() {
	// 当前使用数据小于总数据1/4释放空间
	if p.maxSize > 16 && p.size < p.maxSize>>2 {
		p.maxSize >>= 1
		p.workPoll = p.workPoll[:p.maxSize]
	}
	p.workPoll[0] = p.workPoll[p.size]
	p.size--
	p.sink(0)
}

// 上浮
func (p *queuePoll) upFloat(i int) {
	if p.size == 0 || i == 0 {
		return
	}
	var idx = (i - 1) >> 1
	if p.workPoll[idx].tm > p.workPoll[i].tm {
		p.workPoll[idx], p.workPoll[i] = p.workPoll[i], p.workPoll[idx]
		p.upFloat(idx)
	}
	return
}

// 下沉
func (p *queuePoll) sink(i int) {
	var (
		l   = (i << 1) + 1
		r   = l + 1
		mid = i
	)
	if l <= p.size && p.workPoll[mid].tm > p.workPoll[l].tm {
		mid = l
	}
	if l <= p.size && p.workPoll[mid].tm > p.workPoll[r].tm {
		mid = r
	}
	if mid != i {
		p.workPoll[mid], p.workPoll[i] = p.workPoll[i], p.workPoll[mid]
		p.sink(mid)
	}
}
