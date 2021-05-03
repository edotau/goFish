package concurrency

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 16

// WorkPool serves incoming connections via a pool of workers
type Pool struct {
	closed     int32
	isQueTask  int32
	err        chan error
	Timeout    time.Duration
	wg         sync.WaitGroup
	JobManager chan Job
	Queue      *Stack
}

// Queue represents a single instance of the queue data structure.
type Worker struct {
	Tasks             []interface{}
	head, tail, count int
}

// Task function callbacks
type Job func() error

// NewPool workpool and set the max number of concurrencies
func NewPool(max int) *Pool {
	if max < 1 {
		max = 1
	}
	p := &Pool{
		JobManager: make(chan Job, 2*max),
		err:        make(chan error, 1),
		Queue:      NewStack(),
	}
	go p.loop(max)
	return p
}

// New constructs and returns a new Queue.
func NewWorker() *Worker {
	return &Worker{
		Tasks: make([]interface{}, minQueueLen),
	}
}

// SetTimeout Setting timeout time
func (p *Pool) SetTimeout(timeout time.Duration) {
	p.Timeout = timeout
}

// Do Add to the workpool and return immediately
func (p *Pool) Do(fn Job) {
	if p.IsClosed() { // 已关闭
		return
	}
	p.Queue.Push(fn)
	// p.task <- fn
}

// Add puts an element on the end of the queue.
func (q *Worker) Add(elem interface{}) {
	if q.count == len(q.Tasks) {
		q.resize()
	}

	q.Tasks[q.tail] = elem
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.Tasks) - 1)
	q.count++
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Worker) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.Tasks[q.head:q.tail])
	} else {
		n := copy(newBuf, q.Tasks[q.head:])
		copy(newBuf[n:], q.Tasks[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.Tasks = newBuf
}

// DoWait Add to the workpool and wait for execution to complete before returning
func (p *Pool) DoWait(task Task) { // 添加到工作池，并等待执行完成之后再返回
	if p.IsClosed() { // closed
		return
	}

	doneChan := make(chan struct{})
	p.Queue.Push(Job(func() error {
		defer close(doneChan)
		return task()
	}))
	<-doneChan
}

// Wait Waiting for the worker thread to finish executing
func (p *Pool) Wait() error {
	p.Queue.Wait()
	p.Queue.Close()
	p.waitTask()
	close(p.JobManager)
	p.wg.Wait()
	select {
	case err := <-p.err:
		return err
	default:
		return nil
	}
}

func (q *Stack) Len() int {
	return (int)(atomic.LoadInt32(&q.count))
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Worker) Peek() interface{} {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.Tasks[q.head]
}

// IsDone Determine whether it is complete (non-blocking)
func (p *Pool) IsDone() bool {
	if p == nil || p.JobManager == nil {
		return true
	}

	return p.Queue.Len() == 0 && len(p.JobManager) == 0
}

// IsClosed Has it been closed?
func (p *Pool) IsClosed() bool {
	if atomic.LoadInt32(&p.closed) == 1 {
		return true
	}
	return false
}

func (p *Pool) startQueue() {
	p.isQueTask = 1
	for {
		tmp := p.Queue.Pop()
		if p.IsClosed() {
			p.Queue.Close()
			break
		}
		if tmp != nil {
			fn := tmp.(Job)
			if fn != nil {
				p.JobManager <- fn
			}
		} else {
			break
		}

	}
	atomic.StoreInt32(&p.isQueTask, 0)
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *Worker) Remove() interface{} {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	ret := q.Tasks[q.head]
	q.Tasks[q.head] = nil
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.Tasks) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.Tasks) > minQueueLen && (q.count<<2) == len(q.Tasks) {
		q.resize()
	}
	return ret
}

func (p *Pool) waitTask() {
	for {
		runtime.Gosched()
		if p.IsDone() {
			if atomic.LoadInt32(&p.isQueTask) == 0 {
				break
			}
		}
	}
}

func (p *Pool) loop(maxWorkersCount int) {
	go p.startQueue()

	p.wg.Add(maxWorkersCount)

	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done()
			// worker 开始干活
			for wt := range p.JobManager {
				if wt == nil || atomic.LoadInt32(&p.closed) == 1 {
					continue
				}

				closed := make(chan struct{}, 1)

				if p.Timeout > 0 {
					ct, cancel := context.WithTimeout(context.Background(), p.Timeout)
					go func() {
						select {
						case <-ct.Done():
							p.err <- ct.Err()

							atomic.StoreInt32(&p.closed, 1)
							cancel()
						case <-closed:
						}
					}()
				}

				err := wt()
				close(closed)
				if err != nil {
					select {
					case p.err <- err:

						atomic.StoreInt32(&p.closed, 1)
					default:
					}
				}
			}
		}()
	}
}
