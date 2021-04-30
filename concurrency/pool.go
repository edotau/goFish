package concurrency

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	Queue     *Queue
	closed    int32
	isQueTask int32         // Mark whether queue retrieval is task. 标记是否队列取出任务
	errChan   chan error    // error chan
	Time      time.Duration // max timeout
	wg        sync.WaitGroup
	task      chan JobManager
}

type JobManager func() error

// NewPool initializes a new working pool with channels and qeues.
func NewPool(max int) *Pool {
	if max < 1 {
		max = 1
	}
	p := &Pool{
		task:    make(chan JobManager, 2*max),
		errChan: make(chan error, 1),
		Queue:   NewQueue(),
	}

	go p.loop(max)
	return p
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
			fn := tmp.(JobManager)
			if fn != nil {
				p.task <- fn
			}
		} else {
			break
		}

	}
	atomic.StoreInt32(&p.isQueTask, 0)
}

func (p *Pool) loop(maxWorkersCount int) {
	go p.startQueue()

	p.wg.Add(maxWorkersCount)

	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done()

			for wt := range p.task {
				if wt == nil || atomic.LoadInt32(&p.closed) == 1 {
					continue
				}

				closed := make(chan struct{}, 1)

				if p.Time > 0 {
					ct, cancel := context.WithTimeout(context.Background(), p.Time)
					go func() {
						select {
						case <-ct.Done():
							p.errChan <- ct.Err()

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
					case p.errChan <- err:

						atomic.StoreInt32(&p.closed, 1)
					default:
					}
				}
			}
		}()
	}
}

func (p *Pool) IsClosed() bool {
	if atomic.LoadInt32(&p.closed) == 1 {
		return true
	}
	return false
}

// SetTimeout Setting timeout time
func (p *Pool) SetTimeout(timeout time.Duration) { // 设置超时时间
	p.Time = timeout
}

// Do Add to the workpool and return immediately
func (p *Pool) Do(fn JobManager) {
	if p.IsClosed() {
		return
	}
	p.Queue.Push(fn)

}

// DoWait Add to the workpool and wait for execution to complete before returning
func (p *Pool) DoWait(task JobManager) {
	if p.IsClosed() {
		return
	}

	doneChan := make(chan struct{})
	p.Queue.Push(JobManager(func() error {
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
	close(p.task)
	p.wg.Wait()
	select {
	case err := <-p.errChan:
		return err
	default:
		return nil
	}
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

// IsDone Determine whether it is complete (non-blocking)
func (p *Pool) IsDone() bool {
	if p == nil || p.task == nil {
		return true
	}

	return p.Queue.Len() == 0 && len(p.task) == 0
}
