package bash

import (
	"fmt"
	"io"

	"github.com/commander-cli/cmd"

	"github.com/edotau/goFish/simpleio"
)

func LatestGoResource() *cmd.Command {
	stream := simpleio.NewReader("https://raw.githubusercontent.com/edotau/goFish/main/bash/testdata/download.sh")
	defer stream.Close()

	writer := simpleio.NewWriter("download.sh")
	io.Copy(writer.Writer, stream)
	writer.Close()

	c := cmd.NewCommand("bash download.sh; rm download.sh", cmd.WithStandardStreams)
	c.Execute()
	return c
}

func DirectorySetUp(dir string, golang string) *cmd.Command {
	untar := fmt.Sprintf("tar xf %s -C %s; rm golang-latest.tar.gz; ", golang, dir)

	stream := simpleio.NewReader("https://github.com/edotau/goFish/blob/main/bash/testdata/setupDir.sh")
	writer := simpleio.NewWriter("setupDir.sh")

	io.Copy(writer.Writer, stream)
	writer.Close()

	command := untar + "bash setupDir.sh; rm setupDir.sh"
	c := cmd.NewCommand(command, cmd.WithStandardStreams)
	c.Execute()
	return c
}
