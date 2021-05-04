package dataflow

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/edotau/goFish/api"
)

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
	api.Equal(t, Exec(getTasks(t)...), true)
	tasks, _ := getErrorTask(t)
	api.Equal(t, Exec(tasks...), false)
	api.Equal(t, Exec(), true)
	log.SetOutput(ioutil.Discard)
	api.Equal(t, Exec(getPanicTask()), false)
}

func TestExecWithError(t *testing.T) {
	api.Equal(t, ExecWithError(getTasks(t)...), nil)
	err := fmt.Errorf("TestErr")
	tasks, _ := getErrorTask(t)
	api.Equal(t, ExecWithError(tasks...), err)
	api.Equal(t, strings.Contains(ExecWithError(getPanicTask()).Error(), "panic"), true)
	api.Equal(t, ExecWithError(), nil)
}
