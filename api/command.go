package api

import (
	"github.com/commander-cli/cmd"
	"github.com/edotau/goFish/simpleio"
)

// Bash struct implements the commander-cli/cmd which can be found at:
// https://github.com/commander-cli/cmd/blob/master/command.go
type Bash struct {
	cmd.Command
}

// New script calles the new command method inherited by commander-cli
func NewScript(script string) *Bash {
	sh := Bash{
		Command: *cmd.NewCommand(script, cmd.WithStandardStreams),
	}
	return &sh
}

// Run executes the command line script
func (sh *Bash) Run() {
	err := sh.Command.Execute()
	simpleio.StdError(err)
}

// Stdout returns the output to stdout
func (sh *Bash) Stdout() string {
	return sh.Command.Stdout()
}

// Stderr returns the output to stderr
func (sh *Bash) Stderr() string {
	return sh.Command.Stderr()
}
