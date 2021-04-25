package concurrency

import (
	"errors"
	"testing"
	"time"
)

var (
	fn1     = func() error { return nil }
	fn2     = func() error { return errors.New("BOOM!") }
	timeout = time.After(2 * time.Second)
)

func TestConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			scenario: "test run",
			function: testRun,
		},
		{
			scenario: "test run limit",
			function: testaceCondition,
		},
		{
			scenario: "test run limit with concurrency value greater than passed functions",
			function: testaceConditionWithConcurrencyGreaterThanPassedFunctions,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testRun(t *testing.T) {
	var count int
	err := RunConcur(fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testaceCondition(t *testing.T) {
	var count int
	err := RaceCondition(2, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testaceConditionWithConcurrencyGreaterThanPassedFunctions(t *testing.T) {
	var count int
	err := RaceCondition(3, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}
