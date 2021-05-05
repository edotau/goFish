package dataflow

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// BaseActuator is the actuator interface
type BaseActuator interface {
	Exec(tasks ...Task) error
	ExecWithContext(ctx context.Context, tasks ...Task) error
}

// TimedActuator is the actuator interface within timeout method
type TimedActuator interface {
	BaseActuator
	GetTimeout() *time.Duration
	setTimeout(timeout *time.Duration)
}

// ErrorTimeOut is the error when executes tasks timeout
var ErrorTimeOut = fmt.Errorf("TimeOut")

// Task Type
type Task func() error

// Actuator is the base struct
type Actuator struct {
	timeout *time.Duration
}

// NewActuator creates an Actuator instance
func NewActuator(opt ...*Options) *Actuator {
	c := &Actuator{}
	setOptions(c, opt...)
	return c
}

// Exec is used to run tasks concurrently
func (c *Actuator) Exec(tasks ...Task) error {
	return c.ExecWithContext(context.Background(), tasks...)
}

// ExecWithContext is used to run tasks concurrently
// Return nil when tasks are all completed successfully,
// or return error when some exception happen such as timeout
func (c *Actuator) ExecWithContext(ctx context.Context, tasks ...Task) error {
	return execTasks(ctx, c, simplyRun, tasks...)
}

func simplyRun(f func()) {
	go f()
}

// GetTimeout return the timeout set before
func (c *Actuator) GetTimeout() *time.Duration {
	return c.timeout
}

// setTimeout sets the timeout
func (c *Actuator) setTimeout(timeout *time.Duration) {
	c.timeout = timeout
}

// wait waits for the notification of execution result
func wait(ctx context.Context, c TimedActuator,
	resChan chan error, cancel context.CancelFunc) error {
	if timeout := c.GetTimeout(); timeout != nil {
		return waitWithTimeout(ctx, resChan, *timeout, cancel)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-resChan:
			if err != nil {
				cancel()
				return err
			}
		}
	}
}

// waitWithTimeout waits for the notification of execution result
// when the timeout is set
func waitWithTimeout(ctx context.Context, resChan chan error,
	timeout time.Duration, cancel context.CancelFunc) error {
	for {
		select {
		case <-time.After(timeout):
			cancel()
			return ErrorTimeOut
		case <-ctx.Done():
			return nil
		case err := <-resChan:
			if err != nil {
				cancel()
				return err
			}
		}
	}
}

// execTasks uses customized function to
// execute every task, such as using the simplyRun
func execTasks(parent context.Context, c TimedActuator,
	execFunc func(f func()), tasks ...Task) error {
	size := len(tasks)
	if size == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(parent)
	resChan := make(chan error, size)
	wg := &sync.WaitGroup{}
	wg.Add(size)

	// Make sure the tasks are completed and channel is closed
	go func() {
		wg.Wait()
		cancel()
		close(resChan)
	}()

	// Sadly we can not kill a goroutine manually
	// So when an error happens, the other tasks will continue
	// But the good news is that main progress
	// will know the error immediately
	for _, task := range tasks {
		child, _ := context.WithCancel(ctx)
		f := wrapperTask(child, task, wg, resChan)
		execFunc(f)
	}

	return wait(ctx, c, resChan, cancel)
}

// Exec simply runs the tasks concurrently
// True will be returned is all tasks complete successfully
// otherwise false will be returned
func Exec(tasks ...Task) bool {
	var c int32
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))

	for _, t := range tasks {
		go func(task Task) {
			defer func() {
				if r := recover(); r != nil {
					atomic.StoreInt32(&c, 1)
					fmt.Printf("conexec panic:%v\n%s\n", r, string(debug.Stack()))
				}

				wg.Done()
			}()

			if err := task(); err != nil {
				atomic.StoreInt32(&c, 1)
			}
		}(t)
	}

	wg.Wait()
	return c == 0
}

// ExecWithError simply runs the tasks concurrently
// nil will be returned is all tasks complete successfully
// otherwise custom error will be returned
func ExecWithError(tasks ...Task) error {
	var err error
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))

	for _, t := range tasks {
		go func(task Task) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("conexec panic:%v\n%s\n", r, string(debug.Stack()))
				}

				wg.Done()
			}()

			if e := task(); e != nil {
				err = e
			}
		}(t)
	}

	wg.Wait()
	return err
}

// DurationPtr helps to make a duration ptr
func DurationPtr(t time.Duration) *time.Duration {
	return &t
}

// wrapperTask will wrapper the task in order to notice execution result
// to the main process
func wrapperTask(ctx context.Context, task Task,
	wg *sync.WaitGroup, resChan chan error) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("conexec panic:%v\n%s", r, string(debug.Stack()))
				resChan <- err
			}

			wg.Done()
		}()

		select {
		case <-ctx.Done():
			return // fast return
		case resChan <- task():
		}
	}
}

// setOptions set the options for actuator
func setOptions(c TimedActuator, options ...*Options) {
	if options == nil || len(options) == 0 || options[0] == nil {
		return
	}

	opt := options[0]
	if opt.TimeOut != nil {
		c.setTimeout(opt.TimeOut)
	}
}
