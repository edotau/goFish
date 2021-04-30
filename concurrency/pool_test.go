package concurrency

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// template
func TestWorkerPoolStart(t *testing.T) {
	wp := NewPool(10) // Set the maximum number of threads
	wp.SetTimeout(time.Millisecond)
	for i := 0; i < 20; i++ { // Open 20 requests
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Millisecond)
			}
			// time.Sleep(1 * time.Second)
			return nil
		})
	}

	wp.Queue.Wait()

	fmt.Println("down")
}

// Support for error return
func TestWorkerPoolError(t *testing.T) {
	wp := NewPool(10) // Set the maximum number of threads
	for i := 0; i < 20; i++ {
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				if ii == 1 {
					return errors.Cause(errors.New("my test err"))
				}
				time.Sleep(1 * time.Millisecond)
			}

			return nil
			// time.Sleep(1 * time.Second)
			// return errors.New("my test err")
		})
	}

	wp.Queue.Wait()

	fmt.Println("down")
}

// Determine whether completion (non-blocking) is placed in the workpool and wait for execution results
func TestWorkerPoolDoWait(t *testing.T) {
	wp := NewPool(5) // Set the maximum number of threads
	for i := 0; i < 10; i++ {
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))

				time.Sleep(1 * time.Millisecond)
			}

			return nil
			// time.Sleep(1 * time.Second)
			// return errors.New("my test err")
		})
	}

	wp.Queue.Wait()

	fmt.Println("down")
}
