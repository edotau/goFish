package dataflow

import (
	"runtime"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

func TestPooledTimeOut(t *testing.T) {
	timeout := time.Millisecond * 50
	opt := &Options{TimeOut: &timeout}

	c := NewPooledActuator(5, opt)
	testTimeout(t, c)
	c = NewPooledActuator(-1, opt)
	testTimeout(t, c)
}

func TestPooledError(t *testing.T) {
	timeout := time.Second
	opt := &Options{TimeOut: &timeout}

	c := NewPooledActuator(5, opt)
	testError(t, c)
}

func TestPooledNormal(t *testing.T) {
	c := NewPooledActuator(5)
	testNormal(t, c)

	timeout := time.Minute
	opt := &Options{TimeOut: &timeout}
	c = NewPooledActuator(5, opt)
	testNormal(t, c)

	c.Release()
	c = &PooledActuator{}
	testNormal(t, c)
}

func TestPooledEmpty(t *testing.T) {
	c := NewPooledActuator(5)
	testEmpty(t, c)
}

func TestPooledPanic(t *testing.T) {
	c := NewPooledActuator(5)
	testPanic(t, c)
}

func TestWithPool(t *testing.T) {
	pool, _ := ants.NewPool(5)
	c := NewPooledActuator(5).WithPool(pool)
	testNormal(t, c)
	testError(t, c)
}

func BenchmarkTaskPool(b *testing.B) {
	b.ReportAllocs()

	// Testing with just 1 goroutine
	// to benchmark the non-parallel part of the code
	pool, _ := ants.NewPool(10)

	defer pool.Release()
	//log.SetOutput(ioutil.Discard)

	runtime.GOMAXPROCS(2)
	//log.Printf("%d\n", runtime.NumCPU())
	b.ResetTimer()
	c := &PooledActuator{}
	for n := 0; n < b.N; n++ {

		c = NewPooledActuator(2).WithPool(pool)
		c.Release()

	}
}

/*
func BenchmarkPool(b *testing.B) {
	b.ReportAllocs()
	// Testing with just 1 goroutine
	// to benchmark the non-parallel part of the code

	pool := NewPool(2, 10)
	defer pool.Release()
	log.SetOutput(ioutil.Discard)

	runtime.GOMAXPROCS(6)
	//log.Printf("%d\n", runtime.NumCPU())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pool.JobQueue <- func() {
			//b.Logf("%s", s)
		}
	}

}

func BenchmarkTaskPool(b *testing.B) {
	b.ReportAllocs()

	// Testing with just 1 goroutine
	// to benchmark the non-parallel part of the code
	pool, _ := ants.NewPool(10)

	defer pool.Release()
	//log.SetOutput(ioutil.Discard)

	runtime.GOMAXPROCS(6)
	//log.Printf("%d\n", runtime.NumCPU())
	b.ResetTimer()

	for n := 0; n < b.N; n++ {

		c := NewPooledActuator(2).WithPool(pool)

		err := c.Exec(
			func() error {
				fmt.Println(1)
				time.Sleep(time.Second * 2)
				return nil
			},
		)
		if err == nil {
			c.Release()
		} else {
			simpleio.ErrorHandle(err)
		}
		//

	}
}
*/
