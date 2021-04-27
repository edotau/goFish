package concurrency

func Example() {
	flow := New()

	task1 := flow.MustRegister(Task{
		Name:        "task-1",
		Description: "Print Go version",
		Command:     Exec("go", "version"),
	})

	task2 := flow.MustRegister(Task{
		Name: "task-2",
		Command: func(tf *TF) {
			tf.Skip("skipping")
		},
	})

	task3 := flow.MustRegister(Task{
		Name: "task-3",
		Command: func(tf *TF) {
			tf.Error("hello from", tf.Name())
			tf.Log("this will be printed")
		},
	})

	flow.MustRegister(Task{
		Name:         "all",
		Dependencies: Deps{task1, task2, task3},
	})

	flow.Main()
}
