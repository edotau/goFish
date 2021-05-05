package dataflow

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/edotau/goFish/api"
	"github.com/edotau/goFish/simpleio"
)

var (
	// ErrorUsingActuator is the error when goroutine pool has exception
	ErrorUsingActuator = fmt.Errorf("ErrorUsingActuator")
)

// GoroutinePool is the base routine pool interface
// User can use custom goroutine pool by implementing this interface
type GoroutinePool interface {
	Submit(f func()) error
	Release()
}

// PooledActuator is a actuator which has a worker pool
type PooledActuator struct {
	timeout *time.Duration

	workerNum int
	pool      GoroutinePool
	Wg        sync.WaitGroup
	initOnce  sync.Once
}

// NewPooledActuator creates an PooledActuator instance
func NewPooledActuator(workerNum int, opt ...*Options) *PooledActuator {
	c := &PooledActuator{
		workerNum: workerNum,
	}
	setOptions(c, opt...)
	return c
}

// WithPool will support for using custom goroutine pool
func (c *PooledActuator) WithPool(pool GoroutinePool) *PooledActuator {
	newActuator := c.clone()
	newActuator.pool = pool
	return newActuator
}

// PooledActuator is a actuator which has a worker pool
type JobManager struct {
	timeout *time.Duration

	workerNum int
	pool      GoroutinePool

	initOnce sync.Once
}

// NewPooledActuator creates an PooledActuator instance
func NewJobManager(workerNum, poolSize int, opt ...*Options) *JobManager {
	queue, err := api.NewPool(poolSize)
	simpleio.ErrorHandle(err)
	c := &JobManager{
		workerNum: workerNum,
		pool:      queue,
	}
	setOptions(c, opt...)
	return c
}

// Exec is used to run tasks concurrently
func (c *JobManager) Exec(tasks ...Task) error {
	return c.ExecWithContext(context.Background(), tasks...)
}

// ExecWithContext uses goroutine pool to run tasks concurrently
// Return nil when tasks are all completed successfully,
// or return error when some exception happen such as timeout
func (c *JobManager) ExecWithContext(ctx context.Context, tasks ...Task) error {
	// ensure the actuator can init correctly
	c.initOnce.Do(func() {
		c.initPooledActuator()
	})

	if c.workerNum == -1 {
		return ErrorUsingActuator
	}

	return execTasks(ctx, c, c.runWithPool, tasks...)
}

// GetTimeout return the timeout set before
func (c *JobManager) GetTimeout() *time.Duration {
	return c.timeout
}

// setTimeout sets the timeout
func (c *JobManager) setTimeout(timeout *time.Duration) {
	c.timeout = timeout
}

// clone will clone this PooledActuator without goroutine pool
func (c *JobManager) clone() *JobManager {
	return &JobManager{
		timeout:   c.timeout,
		workerNum: c.workerNum,
		initOnce:  sync.Once{},
	}
}

// runWithPool used the goroutine pool to execute the tasks
func (c *JobManager) runWithPool(f func()) {
	err := c.pool.Submit(f)
	if err != nil {
		log.Printf("submit task err:%s\n", err.Error())
	}
}

// initPooledActuator init the pooled actuator once while the runtime
// If the workerNum is zero or negative,
// default worker num will be used
func (c *JobManager) initPooledActuator() {
	if c.pool != nil {
		// just pass
		c.workerNum = 1
		return
	}

	if c.workerNum <= 0 {
		c.workerNum = runtime.NumCPU() << 1
	}

	var err error
	c.pool, err = api.NewPool(c.workerNum)

	if err != nil {
		c.workerNum = -1
		fmt.Println("initPooledActuator err")
	}
}

// Exec is used to run tasks concurrently
func (c *PooledActuator) Exec(tasks ...Task) error {
	return c.ExecWithContext(context.Background(), tasks...)
}

// ExecWithContext uses goroutine pool to run tasks concurrently
// Return nil when tasks are all completed successfully,
// or return error when some exception happen such as timeout
func (c *PooledActuator) ExecWithContext(ctx context.Context, tasks ...Task) error {
	// ensure the actuator can init correctly
	c.initOnce.Do(func() {
		c.initPooledActuator()
	})

	if c.workerNum == -1 {
		return ErrorUsingActuator
	}

	return execTasks(ctx, c, c.runWithPool, tasks...)
}

// GetTimeout return the timeout set before
func (c *PooledActuator) GetTimeout() *time.Duration {
	return c.timeout
}

// Release will release the pool
func (c *PooledActuator) Release() {
	if c.pool != nil {
		c.pool.Release()
	}
}

// initPooledActuator init the pooled actuator once while the runtime
// If the workerNum is zero or negative,
// default worker num will be used
func (c *PooledActuator) initPooledActuator() {
	if c.pool != nil {
		// just pass
		c.workerNum = 1
		return
	}

	if c.workerNum <= 0 {
		c.workerNum = runtime.NumCPU() << 1
	}

	var err error
	c.pool, err = api.NewPool(c.workerNum)

	if err != nil {
		c.workerNum = -1
		fmt.Println("initPooledActuator err")
	}
}

// runWithPool used the goroutine pool to execute the tasks
func (c *PooledActuator) runWithPool(f func()) {
	err := c.pool.Submit(f)
	if err != nil {
		log.Printf("submit task err:%s\n", err.Error())
	}
}

// setTimeout sets the timeout
func (c *PooledActuator) setTimeout(timeout *time.Duration) {
	c.timeout = timeout
}

// clone will clone this PooledActuator without goroutine pool
func (c *PooledActuator) clone() *PooledActuator {
	return &PooledActuator{
		timeout:   c.timeout,
		workerNum: c.workerNum,
		initOnce:  sync.Once{},
	}
}
