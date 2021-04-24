// Package concurrency contains functions used to build concurrent processes and/or tasks in both pipelines and workflows
// very much still in development
package concurrency

import (
	"sync"
)

// Goroutine is a an alias of a function designed to be programed in concurrency
type Goroutine func() error

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
