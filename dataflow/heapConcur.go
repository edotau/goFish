package dataflow

import (
	"container/heap"
	"sync"
)

// Settings options for Process
type Settings struct {
	PoolSize         int
	OutChannelBuffer int
}

// OrderedOutput is the output channel type from Process
type OrderedOutput struct {
	Value     interface{}
	Remaining func() int
}

// WorkFunction interface
type WorkFunction interface {
	Run() interface{}
}

// Process processes work function based on input.
// It Accepts an WorkFunction read channel, work function and concurrent go routine pool size.
// It Returns an interface{} channel.
func HeapProcess(inputChan <-chan WorkFunction, options *Settings) <-chan OrderedOutput {

	outputChan := make(chan OrderedOutput, options.OutChannelBuffer)

	go func() {
		if options.PoolSize < 1 {
			// Set a minimum number of processors
			options.PoolSize = 1
		}
		processChan := make(chan *processInput, options.PoolSize)
		aggregatorChan := make(chan *processInput, options.PoolSize)

		// Go routine to print data in order
		go func() {
			var current uint64
			outputHeap := &processInputHeap{}
			defer func() {
				close(outputChan)
			}()
			remaining := func() int {
				return outputHeap.Len()
			}
			for item := range aggregatorChan {
				heap.Push(outputHeap, item)
				for top, ok := outputHeap.Peek(); ok && top.order == current; {
					outputChan <- OrderedOutput{Value: heap.Pop(outputHeap).(*processInput).value, Remaining: remaining}
					current++
				}
			}

			for outputHeap.Len() > 0 {
				outputChan <- OrderedOutput{Value: heap.Pop(outputHeap).(*processInput).value, Remaining: remaining}
			}
		}()

		poolWg := sync.WaitGroup{}
		poolWg.Add(options.PoolSize)
		// Create a goroutine pool
		for i := 0; i < options.PoolSize; i++ {
			go func(worker int) {
				defer func() {
					poolWg.Done()
				}()
				for input := range processChan {
					input.value = input.workFn.Run()
					input.workFn = nil
					aggregatorChan <- input
				}
			}(i)
		}

		go func() {
			poolWg.Wait()
			close(aggregatorChan)
		}()

		go func() {
			defer func() {
				close(processChan)
			}()
			var order uint64
			for input := range inputChan {
				processChan <- &processInput{workFn: input, order: order}
				order++
			}
		}()
	}()
	return outputChan
}

type processInput struct {
	workFn WorkFunction
	order  uint64
	value  interface{}
}

type processInputHeap []*processInput

func (h processInputHeap) Len() int {
	return len(h)
}

func (h processInputHeap) Less(i, j int) bool {
	return h[i].order < h[j].order
}

func (h processInputHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *processInputHeap) Push(x interface{}) {
	*h = append(*h, x.(*processInput))
}

func (s processInputHeap) Peek() (*processInput, bool) {
	if len(s) > 0 {
		return s[0], true
	}
	return nil, false
}

func (h *processInputHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
