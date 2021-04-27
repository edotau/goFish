package concurrency

import (
	"strings"
	"testing"
)

func TestExec_success(t *testing.T) {
	sb := &strings.Builder{}
	r := Runner{
		Output: sb,
	}

	got := r.Run(Exec("go", "version"))

	assertContains(t, sb.String(), "go version go", "output should contain prefix of version report")
	assertTrue(t, got.Passed(), "task should pass")
}

func TestExec_error(t *testing.T) {
	r := Runner{}

	got := r.Run(Exec("go", "wrong"))

	assertTrue(t, got.Failed(), "task should fail")
}
