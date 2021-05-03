// Package concurrency contains functions used to build concurrent processes and/or tasks in both pipelines and workflows
// very much still in development
package concurrency

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

// Goroutine is a an alias of a function designed to be programed in concurrency
type Goroutine func() error

type syncWriter struct {
	io.Writer
	mtx sync.Mutex
}

// Run calls the passed functions in a goroutine, returns a chan of errors.
func RunConcur(functions ...Goroutine) chan error {
	total := len(functions)
	errs := make(chan error, total)

	var wg sync.WaitGroup
	wg.Add(total)

	go func(errs chan error) {
		wg.Wait()
		close(errs)
	}(errs)

	for _, i := range functions {
		go func(i Goroutine, errs chan error) {
			defer wg.Done()
			errs <- i()
		}(i, errs)
	}
	return errs
}

func RaceCondition(concurrency int, tasks ...Goroutine) chan error {
	total := len(tasks)

	if concurrency <= 0 {
		concurrency = 1
	}

	if concurrency > total {
		concurrency = total
	}

	var wg sync.WaitGroup
	wg.Add(total)

	errs := make(chan error, total)
	go func(errs chan error) {
		wg.Wait()
		close(errs)
	}(errs)

	sem := make(chan struct{}, concurrency)
	defer func(sem chan<- struct{}) { close(sem) }(sem)

	for _, fn := range tasks {
		go func(fn Goroutine, sem <-chan struct{}, errs chan error) {
			defer wg.Done()
			defer func(sem <-chan struct{}) { <-sem }(sem)

			errs <- fn()
		}(fn, sem, errs)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	return errs
}

func (w *syncWriter) Write(p []byte) (int, error) {
	defer func() { w.mtx.Unlock() }()
	w.mtx.Lock()
	return w.Writer.Write(p)
}

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

// Options use to init actuator
type Options struct {
	TimeOut *time.Duration
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

// GetTimeout return the timeout set before
func (c *Actuator) GetTimeout() *time.Duration {
	return c.timeout
}

// setTimeout sets the timeout
func (c *Actuator) setTimeout(timeout *time.Duration) {
	c.timeout = timeout
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
