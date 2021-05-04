package dataflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	//println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

func TestNewWorker(t *testing.T) {
	pool := make(chan *Worker)
	worker := newWorker(pool)
	worker.start()
	assert.NotNil(t, worker)

	worker = <-pool
	assert.NotNil(t, worker, "Worker should register itself to the pool")

	called := false
	done := make(chan bool)

	job := func() {
		called = true
		done <- true
	}

	worker.jobChannel <- job
	<-done
	assert.Equal(t, true, called)
}

func TestNewPool(t *testing.T) {
	pool := NewPool(1000, 10000)
	defer pool.Release()

	iterations := 1000000
	pool.WaitCount(iterations)
	var counter uint64 = 0

	for i := 0; i < iterations; i++ {
		arg := uint64(1)

		job := func() {
			defer pool.JobDone()
			atomic.AddUint64(&counter, arg)
			assert.Equal(t, uint64(1), arg)
		}

		pool.JobQueue <- job
	}

	pool.WaitAll()

	counterFinal := atomic.LoadUint64(&counter)
	assert.Equal(t, uint64(iterations), counterFinal)
}

func TestRelease(t *testing.T) {
	grNum := runtime.NumGoroutine()
	pool := NewPool(5, 10)
	defer func() {
		pool.Release()

		// give some time for all goroutines to quit
		assert.Equal(t, grNum, runtime.NumGoroutine(), "All goroutines should be released after Release() call")
	}()

	pool.WaitCount(1000)

	for i := 0; i < 1000; i++ {
		job := func() {
			defer pool.JobDone()
		}

		pool.JobQueue <- job
	}

	pool.WaitAll()
}

func BenchmarkPool(b *testing.B) {
	b.ReportAllocs()
	// Testing with just 1 goroutine
	// to benchmark the non-parallel part of the code

	pool := NewPool(2, 10)
	defer pool.Release()
	log.SetOutput(ioutil.Discard)

	runtime.GOMAXPROCS(2)
	//log.Printf("%d\n", runtime.NumCPU())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pool.JobQueue <- func() {
			//b.Logf("%s", s)
		}
	}

}

func testTimeout(t *testing.T, c TimeManager) {
	st := time.Now().UnixNano()
	err := c.Exec(getTasks(t)...)

	Equal(t, err, ErrorTimeOut)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testError(t *testing.T, c TimeManager) {
	st := time.Now().UnixNano()
	tasks, te := getErrorTask(t)
	err := c.Exec(tasks...)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testManyError(t *testing.T, c TimeManager) {
	tasks := make([]Task, 0)
	tmp, te := getErrorTask(t)
	tasks = append(tasks, tmp...)

	for i := 0; i < 100; i++ {
		tmp, _ = getErrorTask(t)
		tasks = append(tasks, tmp...)
	}

	st := time.Now().UnixNano()
	err := c.Exec(tasks...)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testNormal(t *testing.T, c TimeManager) {
	Equal(t, c.Exec(getTasks(t)...), nil)
}

func testPanic(t *testing.T, c TimeManager) {
	NotEqual(t, c.Exec(getPanicTask()), nil)
}

func testEmpty(t *testing.T, c TimeManager) {
	Equal(t, c.Exec(), nil)
}

func getTasks(t *testing.T) []Task {
	return []Task{
		func() error {
			t.Logf("%d\n", 1)
			time.Sleep(time.Millisecond * 100)
			return nil
		},
		func() error {
			t.Logf("%d\n", 2)
			return nil
		},
		func() error {
			time.Sleep(time.Millisecond * 200)
			t.Logf("%d\n", 3)
			return nil
		},
	}
}

func getErrorTask(t *testing.T) ([]Task, error) {
	te := errors.New("TestErr")

	tasks := getTasks(t)
	tasks = append(tasks,
		func() error {
			t.Logf("%d\n", 4)
			return te
		},
		func() error {
			time.Sleep(time.Millisecond * 300)
			t.Logf("%d\n", 5)
			return te
		},
		func() error {
			time.Sleep(time.Second)
			t.Logf("%d\n", 6)
			return te
		})

	return tasks, te
}

func getPanicTask() Task {
	return func() error {
		var i *int64
		num := *i + 1
		fmt.Println(num)
		return nil
	}
}
