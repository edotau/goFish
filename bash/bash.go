// bash package implements simple bash commands that are sometimes useful when writing go scripts
package bash

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

// Cut mimics the basic bash command to cut a column by any delim and returns the fields indices you specify
// indices start at 0
func Cut(line string, d byte, f ...int) []string {
	columns := strings.Split(line, string(d))
	var ans []string
	for _, i := range f {
		ans = append(ans, columns[i])
	}
	return ans
}

func GetColumnCount(line string, delim byte) int {
	return len(strings.Split(line, string(delim)))
}

// Mkdir will create all parent directories of a provided path
func Mkdir(path string, perm os.FileMode) {
	//0755
	err := os.MkdirAll("/Users/temp", perm)
	simpleio.StdError(err)
}

// Package exec runs external commands. It wraps os.exec to
// allow using full command string as arguments,
// and provides functions
// of providing (stdin, stdout, stderr) channel for
// (stdin, stdout, stderr) pipe.
//
// Attention, this package is experimental.

type Cmd struct {
	*exec.Cmd
}

/*
// Command returns the Cmd struct to execute the command.
// No need to split the path and arguments.
// The error may be caused by unclosed quote.
func Command(name string) (*Cmd, error) {
	path, argv, err := parseCommandName(name)
	if err != nil {
		return nil, err
	}
	return &Cmd{exec.Command(path, argv...)}, nil
}*/

// StdinChannel returns a channel that will be connected to
// the command's standard error when the command starts.
func (c *Cmd) StdinChannel() (chan string, error) {
	ch := make(chan string, 100)
	pipe, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(pipe)
	go func() {
		for {
			select {
			case str := <-ch:
				writer.WriteString(str)
			}
		}
	}()
	return ch, nil
}

// parseCommandName split the full command into path and arguments.
func parseCommandName(name string) (string, []string, error) {
	if len(strings.Trim(name, " ")) == 0 {
		return "", nil, errors.New("no command given")
	}

	var (
		quoted    bool = false
		quotation rune
		tmp       []rune   = make([]rune, 0)
		argv      []string = make([]string, 0)
	)
	for _, b := range name {
		switch b {
		case ' ':
			if quoted {
				tmp = append(tmp, b)
			} else {
				if len(strings.Trim(string(tmp), " ")) > 0 {
					argv = append(argv, string(tmp))
				}
				tmp = make([]rune, 0)
			}
		case '"':
			if quoted {
				if quotation == '"' {
					quoted, quotation = false, '_'
					argv = append(argv, string(tmp))
					tmp = make([]rune, 0)
				} else {
					tmp = append(tmp, b)
				}
			} else {
				quoted, quotation = true, '"'
			}
		case '\'':
			if quoted {
				if quotation == '\'' {
					quoted, quotation = false, '_'
					argv = append(argv, string(tmp))
					tmp = make([]rune, 0)
				} else {
					tmp = append(tmp, b)
				}
			} else {
				quoted, quotation = true, '\''
			}
		default:
			tmp = append(tmp, b)
		}
	}
	if len(strings.Trim(string(tmp), " ")) > 0 {
		argv = append(argv, string(tmp))
	}

	path := argv[0]
	var arg []string
	if len(argv) > 1 {
		arg = argv[1:]
	}

	if quoted {
		return path, arg, errors.New("unclosed quote")
	}
	return path, arg, nil
}
