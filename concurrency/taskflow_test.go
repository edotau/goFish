package concurrency

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func init() {
	DefaultOutput = ioutil.Discard
}

func Test_Register(t *testing.T) {
	testCases := []struct {
		desc  string
		task  Task
		valid bool
	}{
		{
			desc:  "good task name",
			task:  Task{Name: "my-task"},
			valid: true,
		},
		{
			desc:  "missing task name",
			task:  Task{},
			valid: false,
		},
		{
			desc:  "invalid dependency",
			task:  Task{Name: "my-task", Dependencies: Deps{RegisteredTask{}}},
			valid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			flow := New()

			_, err := flow.Register(tc.task)

			if tc.valid {
				assertNoError(t, err, "no error expected")
			} else {
				assertError(t, err, "error expected")
			}
		})
	}
}

func Test_Register_same_name(t *testing.T) {
	flow := &Taskflow{}
	task := Task{Name: "task"}
	_, err := flow.Register(task)
	requireNoError(t, err, "should be a valid task")

	_, err = flow.Register(task)

	assertError(t, err, "should not be possible to register tasks with same name twice")
}

func Test_MustRegister_panic(t *testing.T) {
	flow := New()

	act := func() { flow.MustRegister(Task{}) }

	assertPanics(t, act, "should panic because task name is empty")
}

func Test_successful(t *testing.T) {
	ctx := context.Background()
	flow := &Taskflow{}
	var executed1 int
	task1 := flow.MustRegister(Task{
		Name: "task-1",
		Command: func(*TF) {
			executed1++
		},
	})
	var executed2 int
	flow.MustRegister(Task{
		Name: "task-2",
		Command: func(*TF) {
			executed2++
		},
		Dependencies: Deps{task1},
	})
	var executed3 int
	flow.MustRegister(Task{
		Name: "task-3",
		Command: func(*TF) {
			executed3++
		},
		Dependencies: Deps{task1},
	})
	got := func() []int {
		return []int{executed1, executed2, executed3}
	}

	exitCode := flow.Run(ctx, "task-1")
	requireEqual(t, exitCode, 0, "first execution should pass")
	requireEqual(t, got(), []int{1, 0, 0}, "should execute task 1")

	exitCode = flow.Run(ctx, "task-2")
	requireEqual(t, exitCode, 0, "second execution should pass")
	requireEqual(t, got(), []int{2, 1, 0}, "should execute task 1 and 2")

	exitCode = flow.Run(ctx, "task-1", "task-2", "task-3")
	requireEqual(t, exitCode, 0, "third execution should pass")
	requireEqual(t, got(), []int{3, 2, 1}, "should execute task 1 and 2 and 3")
}

func Test_dependency_failure(t *testing.T) {
	flow := &Taskflow{}
	var executed1 int
	task1 := flow.MustRegister(Task{
		Name: "task-1",
		Command: func(tf *TF) {
			executed1++
			tf.Error("it still runs")
			executed1 += 10
			tf.FailNow()
			executed1 += 100
		},
	})
	var executed2 int
	flow.MustRegister(Task{
		Name: "task-2",
		Command: func(*TF) {
			executed2++
		},
		Dependencies: Deps{task1},
	})
	var executed3 int
	flow.MustRegister(Task{
		Name: "task-3",
		Command: func(*TF) {
			executed3++
		},
		Dependencies: Deps{task1},
	})
	got := func() []int {
		return []int{executed1, executed2, executed3}
	}

	exitCode := flow.Run(context.Background(), "task-2", "task-3")

	assertEqual(t, exitCode, 1, "should return error from first task")
	assertEqual(t, got(), []int{11, 0, 0}, "should execute task 1")
}

func Test_fail(t *testing.T) {
	flow := &Taskflow{}
	failed := false
	flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			defer func() {
				failed = tf.Failed()
			}()
			tf.Fatal("failing")
		},
	})

	exitCode := flow.Run(context.Background(), "task")

	assertEqual(t, exitCode, 1, "should return error")
	assertTrue(t, failed, "tf.Failed() should return true")
}

func Test_skip(t *testing.T) {
	flow := &Taskflow{}
	skipped := false
	flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			defer func() {
				skipped = tf.Skipped()
			}()
			tf.Skip("skipping")
		},
	})

	exitCode := flow.Run(context.Background(), "task")

	assertEqual(t, exitCode, 0, "should pass")
	assertTrue(t, skipped, "tf.Skipped() should return true")
}

func Test_task_panics(t *testing.T) {
	flow := &Taskflow{}
	flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			panic("panicked!")
		},
	})

	exitCode := flow.Run(context.Background(), "task")

	assertEqual(t, exitCode, 1, "should return error from first task")
}

func Test_cancelation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	flow := &Taskflow{}
	flow.MustRegister(Task{
		Name: "task",
	})

	exitCode := flow.Run(ctx, "task")

	assertEqual(t, exitCode, 1, "should return error canceled")
}

func Test_cancelation_during_last_task(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	flow := &Taskflow{}
	flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			cancel()
		},
	})

	exitCode := flow.Run(ctx, "task")

	assertEqual(t, exitCode, 1, "should return error canceled")
}

func Test_empty_command(t *testing.T) {
	flow := &Taskflow{}
	flow.MustRegister(Task{
		Name: "task",
	})

	exitCode := flow.Run(context.Background(), "task")

	assertEqual(t, exitCode, 0, "should pass")
}

func Test_invalid_args(t *testing.T) {
	flow := New()
	flow.MustRegister(Task{
		Name: "task",
	})

	testCases := []struct {
		desc string
		args []string
	}{
		{
			desc: "missing task name",
		},
		{
			desc: "bad flag",
			args: []string{"-badflag", "task"},
		},
		{
			desc: "bad task name",
			args: []string{"badtask"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			exitCode := flow.Run(context.Background(), tc.args...)

			assertEqual(t, exitCode, 2, "should return error bad args")
		})
	}
}

func Test_help(t *testing.T) {
	flow := New()
	a := flow.MustRegister(Task{
		Name:        "a",
		Description: "some task",
	})
	flow.DefaultTask = a
	flow.Params.SetBool("fast", false)

	exitCode := flow.Run(context.Background(), "-h")

	assertEqual(t, exitCode, 2, "should return error bad args")
}

func Test_printing(t *testing.T) {
	sb := &strings.Builder{}
	flow := &Taskflow{
		Output:  sb,
		Verbose: true,
	}
	skipped := flow.MustRegister(Task{
		Name: "skipped",
		Command: func(tf *TF) {
			tf.Skipf("Skipf %d", 0)
		},
	})
	flow.MustRegister(Task{
		Name:         "failing",
		Dependencies: Deps{skipped},
		Command: func(tf *TF) {
			tf.Log("Log", 1)
			tf.Logf("Logf %d", 2)
			tf.Error("Error", 3)
			tf.Errorf("Errorf %d", 4)
			tf.Fatalf("Fatalf %d", 5)
		},
	})
	t.Log()

	flow.Run(context.Background(), "failing")

	assertContains(t, sb.String(), "Skipf 0", "should contain proper output from \"skipped\" task")
	assertContains(t, sb.String(), `Log 1
Logf 2
Error 3
Errorf 4
Fatalf 5`, "should contain proper output from \"failing\" task")
}

func Test_concurrent_printing(t *testing.T) {
	testCases := []struct {
		verbose bool
	}{
		{verbose: false},
		{verbose: true},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("Verbose:%v", tc.verbose)
		t.Run(testName, func(t *testing.T) {
			sb := &strings.Builder{}
			flow := Taskflow{
				Verbose: tc.verbose,
				Output:  sb,
			}
			flow.MustRegister(Task{
				Name: "task",
				Command: func(tf *TF) {
					ch := make(chan struct{})
					go func() {
						defer func() { ch <- struct{}{} }()
						tf.Log("from child goroutine")
					}()
					tf.Error("from main goroutine")
					<-ch
				},
			})

			exitCode := flow.Run(context.Background(), "task")

			assertEqual(t, exitCode, CodeFailure, "should fail")
			assertContains(t, sb.String(), "from child goroutine", "should contain log from child goroutine")
			assertContains(t, sb.String(), "from main goroutine", "should contain log from main goroutine")
		})
	}
}

func Test_name(t *testing.T) {
	flow := &Taskflow{}
	taskName := "my-named-task"
	var got string
	flow.MustRegister(Task{
		Name: taskName,
		Command: func(tf *TF) {
			got = tf.Name()
		},
	})

	exitCode := flow.Run(context.Background(), taskName)

	assertEqual(t, exitCode, 0, "should pass")
	assertEqual(t, got, taskName, "should return proper Name value")
}

func Test_verbose(t *testing.T) {
	flow := &Taskflow{}
	var got bool
	flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			got = tf.Verbose()
		},
	})

	exitCode := flow.Run(context.Background(), "-v", "task")

	assertEqual(t, exitCode, 0, "should pass")
	assertTrue(t, got, "should return proper Verbose value")
}

func Test_params(t *testing.T) {
	flow := New()
	flow.Params.SetInt("x", 1)
	flow.Params["z"] = "abc"
	var got TFParams
	flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			got = tf.Params()
		},
	})

	exitCode := flow.Run(context.Background(), "y=2", "z=3", "task")

	assertEqual(t, exitCode, 0, "should pass")
	assertEqual(t, got.String("x"), "1", "x param")
	assertEqual(t, got.Int("y"), 2, "y param")
	assertEqual(t, got.Float64("z"), 3.0, "z param")
}

func Test_defaultTask(t *testing.T) {
	flow := New()
	var got TFParams
	task := flow.MustRegister(Task{
		Name: "task",
		Command: func(tf *TF) {
			got = tf.Params()
		},
	})
	flow.DefaultTask = task

	exitCode := flow.Run(context.Background(), "x=a")

	assertEqual(t, exitCode, 0, "should pass")
	assertEqual(t, got.String("x"), "a", "x param")
}