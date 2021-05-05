package dataflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/edotau/goFish/api"
	. "github.com/go-playground/assert/v2"
)

func testTimeout(t *testing.T, c TimedActuator) {
	st := time.Now().UnixNano()
	err := c.Exec(getTasks(t)...)

	Equal(t, err, ErrorTimeOut)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testError(t *testing.T, c TimedActuator) {
	st := time.Now().UnixNano()
	tasks, te := getErrorTask(t)
	err := c.Exec(tasks...)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testManyError(t *testing.T, c TimedActuator) {
	tasks := make([]Task, 0)
	tmp, te := getErrorTask(t)
	tasks = append(tasks, tmp...)

	for i := 0; i < 100; i++ {
		tmp, _ = getErrorTask(t)
		tasks = append(tasks, tmp...)
	}

	st := time.Now().UnixNano()
	err := c.Exec(tasks...)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testNormal(t *testing.T, c TimedActuator) {
	Equal(t, c.Exec(getTasks(t)...), nil)
}

func testPanic(t *testing.T, c TimedActuator) {
	NotEqual(t, c.Exec(getPanicTask()), nil)
}

func testEmpty(t *testing.T, c TimedActuator) {
	Equal(t, c.Exec(), nil)
}

func getTasks(t *testing.T) []Task {
	return []Task{
		func() error {
			t.Logf("%d\n", 1)
			time.Sleep(time.Millisecond * 100)
			return nil
		},
		func() error {
			t.Logf("%d\n", 2)
			return nil
		},
		func() error {
			time.Sleep(time.Millisecond * 200)
			t.Logf("%d\n", 3)
			return nil
		},
	}
}

func getErrorTask(t *testing.T) ([]Task, error) {
	te := errors.New("TestErr")

	tasks := getTasks(t)
	tasks = append(tasks,
		func() error {
			t.Logf("%d\n", 4)
			return te
		},
		func() error {
			time.Sleep(time.Millisecond * 300)
			t.Logf("%d\n", 5)
			return te
		},
		func() error {
			time.Sleep(time.Second)
			t.Logf("%d\n", 6)
			return te
		})

	return tasks, te
}

func getPanicTask() Task {
	return func() error {
		var i *int64
		num := *i + 1
		fmt.Println(num)
		return nil
	}
}

func TestTimeOut(t *testing.T) {
	timeout := time.Millisecond * 50
	opt := &Options{TimeOut: &timeout}
	c := NewActuator(opt)
	testTimeout(t, c)
}

func TestError(t *testing.T) {
	timeout := time.Second
	opt := &Options{TimeOut: &timeout}
	c := NewActuator(opt)
	testError(t, c)
}

func TestNormal(t *testing.T) {
	c := NewActuator()
	testNormal(t, c)

	timeout := time.Minute
	opt := &Options{TimeOut: &timeout}
	c = NewActuator(opt)
	testNormal(t, c)
}

func TestEmpty(t *testing.T) {
	c := NewActuator()
	testEmpty(t, c)
}

func TestPanic(t *testing.T) {
	c := NewActuator()
	testPanic(t, c)
}

func TestManyError(t *testing.T) {
	timeout := time.Second
	opt := &Options{TimeOut: &timeout}
	c := NewActuator(opt)
	testManyError(t, c)
}

func TestDurationPtr(t *testing.T) {
	timeout := time.Minute
	api.Equal(t, timeout, *DurationPtr(timeout))
}

func TestExec(t *testing.T) {
	Equal(t, Exec(getTasks(t)...), true)
	tasks, _ := getErrorTask(t)
	Equal(t, Exec(tasks...), false)
	Equal(t, Exec(), true)
}

func TestExecWithError(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	api.Equal(t, ExecWithError(getTasks(t)...), nil)
	err := fmt.Errorf("TestErr")
	tasks, _ := getErrorTask(t)
	api.Equal(t, ExecWithError(tasks...), err)
	api.NotEqual(t, strings.Contains(ExecWithError(getPanicTask()).Error(), "panic"), false)
	api.Equal(t, ExecWithError(), nil)
}
