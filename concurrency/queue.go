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
		Stack: New(),
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
		//q.Stack.Signal()
	}
}

func (q *Queue) Pop() (v interface{}) {
	//w := q.Stack

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for q.Len() == 0 && !q.closed {
		q.Wait()
	}

	if q.closed {
		return
	}
	if q.Len() > 0 {
		work := q.Stack
		v = work.Peek()
		work.Remove()
		atomic.AddInt32(&q.count, -1)
	}
	return
}

func (q *Queue) TryPop() (v interface{}, ok bool) {
	buffer := q.Stack

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if q.Len() > 0 {
		v = buffer.Peek()
		buffer.Remove()
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
			//	q.Stack.Signal()
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
		tmp := werq.Peek()
		q.Stack.Remove()
		atomic.AddInt32(&q.count, -1)
		*v <- tmp
	} else {
		*v <- nil
	}
	return
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
		//q.Stack.Broadcast()
	}
}

// IsClose check is closed
func (q *Queue) IsClose() bool {
	return q.closed
}

// New constructs and returns a new Queue.
func New() *Jobs {
	return &Jobs{
		Work: make([]interface{}, minQueueLen),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Jobs) Length() int {
	return q.count
}

func (q *Queue) Wait() {
	for {
		if q.closed || q.Len() == 0 {
			break
		}
		runtime.Gosched()
	}
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (w *Jobs) resize() {
	jobs := make([]interface{}, w.count<<1)

	if w.tail > w.head {
		copy(jobs, w.Work[w.head:w.tail])
	} else {
		n := copy(jobs, w.Work[w.head:])
		copy(jobs[n:], w.Work[:w.tail])
	}

	w.head = 0
	w.tail = w.count
	w.Work = jobs
}

// Add puts an element on the end of the wueue.
func (q *Queue) Add(elem interface{}) {

	if q.Stack.count == len(q.Stack.Work) {
		q.Stack.resize()
	}

	q.Stack.Work[q.Stack.tail] = elem
	// bitwise modulus
	q.Stack.tail = (q.Stack.tail + 1) & (len(q.Stack.Work) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (w *Jobs) Peek() interface{} {
	if w.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return w.Work[w.head]
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (w *Jobs) Get(i int) interface{} {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += w.count
	}
	if i < 0 || i >= w.count {
		panic("queue: Get() called with index out of range")
	}
	// bitwise modulus
	return w.Work[(w.head+i)&(len(w.Work)-1)]
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (w *Jobs) Remove() interface{} {
	if w.count <= 0 {
		panic("Jobs: Remove() called on empty queue")
	}
	ret := w.Work[w.head]
	w.Work[w.head] = nil
	// bitwise modulus
	w.head = (w.head + 1) & (len(w.Work) - 1)
	w.count--
	// Resize down if buffer 1/4 full.
	if len(w.Work) > minQueueLen && (w.count<<2) == len(w.Work) {
		w.resize()
	}
	return ret
}
