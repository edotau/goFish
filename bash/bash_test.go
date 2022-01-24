package bash

import (
	"testing"
)

func TestCut(t *testing.T) {
	line := "werq\twerq\twerq\twerq"
	cut := Cut(line, '\t', 0, 3)
	for _, i := range cut {
		t.Logf("\t%s", i)
		if i != "werq" {
			t.Errorf("Error: cut function is not reading the delim correctly...\n")
		}
	}
	t.Log("\n")
	if len(cut) != 2 {
		t.Errorf("Error: cut function is not reading fields corrently...\n")
	}
}

/*
func TestSetup(t *testing.T) {
	LatestGoResource().Stdout()
	rm := api.NewScript("rm go*.tar.gz")
	rm.Run()
}*/
