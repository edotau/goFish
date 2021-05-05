package dataflow

import (
	"fmt"

	"github.com/edotau/goFish/api"
	"github.com/edotau/goFish/simpleio"

	"testing"
)

func TestPipeline(t *testing.T) {

	queue, err := api.NewPool(10000)
	simpleio.ErrorHandle(err)
	defer queue.Release()

	slurm := NewPooledActuator(8).WithPool(queue)

	err = slurm.Exec(
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

}
