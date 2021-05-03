package concurrency

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 16

type Queue struct {
	sync.Mutex
	popable *sync.Cond
	Stack   *Jobs
	closed  bool
	count   int32
	cc      chan interface{}
	once    sync.Once
}

// Queue represents a single instance of the queue data structure.
type Jobs struct {
	Work              []interface{}
	head, tail, count int
}

func NewQueue() *Queue {
	jobs := &Queue{
		Stack: &Jobs{Work: make([]interface{}, minQueueLen)},
	}
	jobs.popable = sync.NewCond(&jobs.Mutex)
	return jobs
}

func (q *Queue) Push(v interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.Add(v)
		atomic.AddInt32(&q.count, 1)
		q.popable.Signal()
	}
}

func (q *Queue) Pop() (v interface{}) {
	work := q.popable

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.closed {
		work.Wait()
	}

	if q.closed {
		return
	}
	if q.Len() > 0 {
		work := q.Stack
		//v = work.Peek()
		work.Remove()
		atomic.AddInt32(&q.count, -1)
	}
	return
}

func (q *Queue) TryPop() (v interface{}, ok bool) {
	jobs := q.Stack

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if q.Len() > 0 {
		v = jobs.Peek()
		jobs.Remove()
		atomic.AddInt32(&q.count, -1)
		ok = true
	} else if q.closed {
		ok = true
	}
	return
}

func (q *Queue) TryPopTimeout(tm time.Duration) (v interface{}, ok bool) {
	q.once.Do(func() {
		q.cc = make(chan interface{}, 1)
	})
	go func() {
		q.popChan(&q.cc)
	}()

	ok = true
	timeout := time.After(tm)
	select {
	case v = <-q.cc:
	case <-timeout:
		if !q.closed {
			q.popable.Signal()
		}
		ok = false
	}
	return
}

func (q *Queue) popChan(v *chan interface{}) {
	//c := q.Stack

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.closed {
		q.Wait()
	}

	if q.closed {
		*v <- nil
		return
	}

	if q.Len() > 0 {
		werq := q.Stack
		//tmp := werq.Peek()
		q.Stack.Remove()
		atomic.AddInt32(&q.count, -1)
		*v <- werq
	} else {
		*v <- nil
	}
	return
}

func (q *Queue) Wait() {
	for {
		if q.closed || q.Len() == 0 {
			break
		}
		runtime.Gosched()
	}
}

func (q *Queue) Len() int {
	return (int)(atomic.LoadInt32(&q.count))
}

// Close MyQueue
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *Queue) Close() {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.closed = true
		atomic.StoreInt32(&q.count, 0)
		q.popable.Broadcast()
	}
}

// IsClose check is closed
func (q *Queue) IsClose() bool {
	return q.closed
}

// New constructs and returns a new Queue.
func NewWorkers() *Jobs {
	return &Jobs{
		Work: make([]interface{}, minQueueLen),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue) Length() int {
	return q.Stack.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if len(q.Stack.Work) > q.Stack.head {
		copy(newBuf, q.Stack.Work[q.Stack.head:q.Stack.tail])
	} else {
		n := copy(newBuf, q.Stack.Work[q.Stack.head:])
		copy(newBuf[n:], q.Stack.Work[:q.Stack.tail])
	}

	q.Stack.head = 0
	q.Stack.tail = q.Stack.count
	q.Stack.Work = newBuf
}

// Add puts an element on the end of the queue.
func (q *Queue) Add(elem interface{}) {
	if q.Stack.count == len(q.Stack.Work) {
		q.resize()
	}

	q.Stack.Work[q.Stack.tail] = elem
	// bitwise modulus
	q.Stack.tail = (q.Stack.tail + 1) & (len(q.Stack.Work) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Jobs) Peek() interface{} {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.Work[q.head]
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue) Get(i int) interface{} {
	// If indexing backwards, convert to positive index.

	if i < 0 {
		i += int(q.count)
	}
	if i < 0 || i >= int(q.count) {
		panic("queue: Get() called with index out of range")
	}
	// bitwise modulus
	return q.Stack.Work[(q.Stack.head+i)&(len(q.Stack.Work)-1)]
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *Jobs) Remove() interface{} {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	ret := q.Work[q.head]
	q.Work[q.head] = nil
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.Work) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.Work) > minQueueLen && (q.count<<2) == len(q.Work) {
		q.resize()
	}
	return ret
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Jobs) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.Work[q.head:q.tail])
	} else {
		n := copy(newBuf, q.Work[q.head:])
		copy(newBuf[n:], q.Work[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.Work = newBuf
}
