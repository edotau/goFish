package bash

import (
	"fmt"
	"testing"

	"github.com/commander-cli/cmd"
)

func TestCut(t *testing.T) {
	line := "werq\twerq\twerq\twerq"
	cut := Cut(line, '\t', 0, 3)
	for _, i := range cut {
		fmt.Printf("\t%s", i)
		if i != "werq" {
			t.Errorf("Error: cut function is not reading the delim correctly...\n")
		}
	}
	fmt.Printf("\n")
	if len(cut) != 2 {
		t.Errorf("Error: cut function is not reading fields corrently...\n")
	}
}

func TestSetup(t *testing.T) {
	tarfile := LatestGoResource().Stdout()
	cmd.NewCommand(fmt.Sprintf("rm %s ", tarfile), cmd.WithStandardStreams).Execute()

}
