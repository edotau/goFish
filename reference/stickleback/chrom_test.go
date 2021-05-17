package stickleback

import(
    "testing"
)

func TestChromStdout(t *testing.T) {
    for _, chr := range Chr {
        GetChrom(chr)
    }
}
