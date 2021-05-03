package concurrency

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Queue queue
type Stack struct {
	sync.Mutex
	popable *sync.Cond
	Jobs    *Worker
	closed  bool
	count   int32
	once    sync.Once
}

func NewStack() *Stack {
	ch := &Stack{
		Jobs: NewWorker(),
	}
	ch.popable = sync.NewCond(&ch.Mutex)
	return ch
}

func (q *Stack) Push(v interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.Jobs.Add(v)
		atomic.AddInt32(&q.count, 1)
		q.popable.Signal()
	}
}

func (q *Stack) Wait() {
	for {
		if q.closed || q.Len() == 0 {
			break
		}
		runtime.Gosched()
	}
}

// Close Stack queue, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *Stack) Close() {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.closed = true
		atomic.StoreInt32(&q.count, 0)
		q.popable.Broadcast()
	}
}

func (q *Stack) Pop() (v interface{}) {
	c := q.popable

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.closed {
		c.Wait()
	}

	if q.closed {
		return
	}

	if q.Len() > 0 {
		buffer := q.Jobs
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
	}
	return
}
