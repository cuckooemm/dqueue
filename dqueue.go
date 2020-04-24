package dqueue

import (
	"errors"
	"time"
)

const (
	maxTimeDuration = 1<<63 - 1
)

var (
	ch    chan delayJob
	qPoll queuePoll
)

type IDelayJob interface {
	// 执行的回调函数
	Work()
}

/*
 *  ex 为认为执行的绝对时间
 *  td 与 tick 成对出现 td 控制距离上次执行间隔，tick 控制次数
 ×  tick 执行次数为 0 则执行1次 为 n 则执行 n + 1 次
 ×  tick 为 -1 时持久执行 小于-1 返回error
 ×  work 需实现 IDelayJob 的 Work 函数
*/
func NewDelayWork(ex time.Time, td time.Duration, tick int64, work IDelayJob) error {
	if tick < -1 {
		return errors.New("Tick parameter input error")
	}
	ch <- delayJob{tm: ex.UnixNano(), tickDur: int64(td), tick: tick, work: work}
	return nil
}

type delayJob struct {
	tickDur int64
	tick    int64
	tm      int64
	work    IDelayJob
}

type queuePoll struct {
	size     int
	maxSize  int
	workPoll []delayJob
}

func init() {
	ch = make(chan delayJob, 100)
	go start(ch)
}

//func DelayQueueInit(chanSize int) chan<- IDelayJob {
//	if chanSize < 1 {
//		chanSize = 100
//	}
//	var ch = make(chan IDelayJob, chanSize)
//	go start(ch)
//	return ch
//}

func start(ch <-chan delayJob) {
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
		tk.Reset(time.Duration(job.tm - curTime))
		select {
		case <-tk.C:
			go job.work.Work()
			if job.tick > 0 {
				job.tm += job.tickDur
				job.tick--
				qPoll.addJob(job)
			} else if job.tick == -1 {
				job.tm += job.tickDur
				qPoll.addJob(job)
			}
			qPoll.deleteTop()
		case q := <-ch:
			qPoll.addJob(q)
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
	if p.maxSize > 32 && p.size < p.maxSize>>2 {
		p.maxSize >>= 1
		p.workPoll = p.workPoll[:p.maxSize:p.maxSize]
	}
	p.workPoll[0] = p.workPoll[p.size]
	p.size--
	p.sink(0)
}

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
