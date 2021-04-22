package concurrency

import (
	"runtime"
	"sync"
)

type Executor func(interface{}) interface{}

type Pipeline interface {
	Pipe(executor Executor) Pipeline
	Merge() <-chan interface{}
}

type pipeline struct {
	dataC chan interface{}
	//errC      chan error
	executors []Executor
}

func New(f func(chan interface{})) Pipeline {
	inC := make(chan interface{})

	go f(inC)

	return &pipeline{
		dataC: inC,
		//errC:      make(chan error),
		executors: []Executor{},
	}
}

func (p *pipeline) Pipe(executor Executor) Pipeline {
	p.executors = append(p.executors, executor)

	return p
}

func (p *pipeline) Merge() <-chan interface{} {
	for i := 0; i < len(p.executors); i++ {
		p.dataC = run(p.dataC, p.executors[i])
	}

	return p.dataC
}

func run(inC <-chan interface{}, f Executor) chan interface{} {
	var threads int
	if available := runtime.GOMAXPROCS(0); threads > available || threads < 1 {
		threads = available
	}
	outC := make(chan interface{})
	//errC := make(chan error)
	wg := &sync.WaitGroup{}
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			//defer close(outC)
			for v := range inC {
				res := f(v)
				//if err != nil {
				//	errC <- err
				//	continue
				//}

				outC <- res
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(outC)
	}()
	return outC
}
