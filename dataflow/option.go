package dataflow

import (
	"time"
)

// Options use to init actuator
type Options struct {
	TimeOut *time.Duration
}

/*
func TestPipeline(t *testing.T) {

	queue := api.NewPool(10000)

	defer queue.Release()

	slurm := NewJobManager(8).WithPool(queue)

	err := slurm.Exec(
		func() error {
			fmt.Println(1)

			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {

			fmt.Println(3)
			return nil
		},
	)
	simpleio.FatalErr(err)
}
*/
