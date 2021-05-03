package concurrency

import (
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	que := NewStack()
	for i := 0; i < 10; i++ {
		que.Push(i)
	}

	go func() {
		for {
			t.Logf("%d\n", que.Pop().(int))
			//time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			t.Logf("%d\n", que.Pop().(int))
			time.Sleep(1 * time.Second)
		}
	}()
	que.Wait()
	t.Logf("down")
}

func TestClose(t *testing.T) {
	que := NewStack()
	for i := 0; i < 10; i++ {
		que.Push(i)
	}

	go func() {
		for {
			v := que.Pop()
			if v != nil {
				t.Logf("%d\n", v.(int))
				//time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			v := que.Pop()
			if v != nil {
				t.Logf("%d\n", v.(int))
				//time.Sleep(1 * time.Second)
			}
		}
	}()

	que.Close()
	que.Wait()
	t.Logf("down")
}
