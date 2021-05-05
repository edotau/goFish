package dataflow

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

type zeroLoadWorker int

func (w zeroLoadWorker) Run() interface{} {
	return w * 2
}

type loadWorker int

func (w loadWorker) Run() interface{} {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return w * 2
}

func TestHeapConcur(t *testing.T) {
	t.Run("Test with Preset Pool Size", func(t *testing.T) {
		max := 10
		inputChan := make(chan WorkFunction)
		wg := &sync.WaitGroup{}

		outChan := HeapProcess(inputChan, &Settings{PoolSize: 10})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(loadWorker); !ok {
					t.Error("Invalid output")
				} else {
					counter++
				}
				wg.Done()
			}
		}(t)

		// Create work and the associated order
		for work := 0; work < max; work++ {
			wg.Add(1)
			inputChan <- loadWorker(work)
		}
		close(inputChan)
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Preset Pool Size Completed")
	})
}

func TestSortedData(t *testing.T) {
	t.Run("Test if response is sorted", func(t *testing.T) {
		max := 10
		inputChan := make(chan WorkFunction)
		output := HeapProcess(inputChan, &Settings{PoolSize: 10, OutChannelBuffer: 10})
		go func() {
			for work := 0; work < max; work++ {
				inputChan <- loadWorker(work)
			}
			close(inputChan)
		}()
		var res []loadWorker
		for out := range output {
			res = append(res, out.Value.(loadWorker))
		}
		isSorted := sort.SliceIsSorted(res, func(i, j int) bool {
			return res[i] < res[j]
		})
		if !isSorted {
			t.Error("output is not sorted")
		}
		t.Log("Test if response is sorted")
	})
}

func TestSortedDataMultiple(t *testing.T) {
	for i := 0; i < 50; i++ {
		t.Run("Test if response is sorted", func(t *testing.T) {
			max := 10
			inputChan := make(chan WorkFunction)
			output := HeapProcess(inputChan, &Settings{PoolSize: 10, OutChannelBuffer: 10})
			go func() {
				for work := 0; work < max; work++ {
					inputChan <- loadWorker(work)
				}
				close(inputChan)
			}()
			var res []loadWorker
			for out := range output {
				res = append(res, out.Value.(loadWorker))
			}
			isSorted := sort.SliceIsSorted(res, func(i, j int) bool {
				return res[i] < res[j]
			})
			if !isSorted {
				t.Error("output is not sorted")
			}
			t.Log("Test if response is sorted")
		})
	}
}

func TestStreamingInput(t *testing.T) {
	t.Run("Test streaming input", func(t *testing.T) {
		inputChan := make(chan WorkFunction, 10)
		output := HeapProcess(inputChan, &Settings{PoolSize: 10, OutChannelBuffer: 10})

		ticker := time.NewTicker(100 * time.Millisecond)
		done := make(chan bool)
		wg := &sync.WaitGroup{}
		go func() {
			input := 0
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					inputChan <- zeroLoadWorker(input)
					wg.Add(1)
					input++
				default:
				}
			}
		}()

		var res []zeroLoadWorker

		go func() {
			for out := range output {
				res = append(res, out.Value.(zeroLoadWorker))
				wg.Done()
			}
		}()

		time.Sleep(1600 * time.Millisecond)
		ticker.Stop()
		done <- true
		close(inputChan)
		wg.Wait()
		isSorted := sort.SliceIsSorted(res, func(i, j int) bool {
			return res[i] < res[j]
		})
		if !isSorted {
			t.Error("output is not sorted")
		}
		t.Log("Test streaming input")
	})
}

func BenchmarkOC(b *testing.B) {
	max := 100000
	inputChan := make(chan WorkFunction)
	output := HeapProcess(inputChan, &Settings{PoolSize: 10, OutChannelBuffer: 10})
	go func() {
		for work := 0; work < max; work++ {
			inputChan <- zeroLoadWorker(work)
		}
		close(inputChan)
	}()
	for out := range output {
		_ = out
	}
}

func BenchmarkOCLoad(b *testing.B) {
	max := 10
	inputChan := make(chan WorkFunction)
	output := HeapProcess(inputChan, &Settings{PoolSize: 10, OutChannelBuffer: 10})
	go func() {
		for work := 0; work < max; work++ {
			inputChan <- loadWorker(work)
		}
		close(inputChan)
	}()
	for out := range output {
		_ = out
	}
}

func BenchmarkOC2(b *testing.B) {
	for i := 0; i < 100; i++ {
		max := 1000
		inputChan := make(chan WorkFunction)
		output := HeapProcess(inputChan, &Settings{PoolSize: 10, OutChannelBuffer: 10})
		go func() {
			for work := 0; work < max; work++ {
				inputChan <- zeroLoadWorker(work)
			}
			close(inputChan)
		}()
		for out := range output {
			_ = out
		}
	}
}
