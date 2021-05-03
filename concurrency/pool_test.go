package concurrency

import (
	"log"
	"testing"

	"github.com/pkg/errors"
)

func TestWorkerPoolStart(t *testing.T) {
	wp := NewPool(10)

	for i := 0; i < 20; i++ { // Open 20 requests
		ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ {
				t.Logf("%v->\t%v\n", ii, j)
			}
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		log.Fatal(err)
	}
	t.Log("down\n")
}
func TestWorkerPoolError(t *testing.T) {
	wp := NewPool(10) // Set the maximum number of threads
	for i := 0; i < 9; i++ {
		ii := i
		wp.Do(func() error {
			for j := 0; j < 9; j++ {
				t.Logf("%v->\t%v\n", ii, j)
				if ii == 1 {
					//t.Log(errors.New("my test err"))
					t.Log(errors.New("Worker should catch this error"))
				}
			}
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		log.Fatal(err)
	}
	t.Log("down\n")
}

// Determine whether completion (non-blocking) is placed in the workpool and wait for execution results
func TestWorkerPoolDoWait(t *testing.T) {
	wp := NewPool(5) // Set the maximum number of threads
	for i := 0; i < 10; i++ {
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ {
				t.Logf("%v->\t%v\n", ii, j)

			}

			return nil

		})
	}

	err := wp.Wait()
	if err != nil {
		t.Error(err)
	}
	t.Logf("down\n")
}

// Determine whether it is complete (non-blocking)
func TestWorkerPoolIsDone(t *testing.T) {
	wp := NewPool(5) // Set the maximum number of threads
	for i := 0; i < 10; i++ {
		//	ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ {

			}
			return nil
		})
		t.Log(wp.IsDone())
	}
	wp.Wait()
	t.Log(wp.IsDone())
	t.Log("down\n")
}
