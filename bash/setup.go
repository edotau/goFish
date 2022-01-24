package bash

import (
	"fmt"
	"io"

	"github.com/edotau/goFish/api"
	"github.com/edotau/goFish/simpleio"
)

func LatestGoResource() *api.Bash {
	stream := simpleio.NewReader("https://raw.githubusercontent.com/edotau/goFish/main/bash/testdata/download.sh")
	defer stream.Close()

	writer := simpleio.NewWriter("download.sh")
	io.Copy(writer.Writer, stream)
	writer.Close()

	sh := api.NewScript("bash download.sh; rm download.sh")
	sh.Run()
	return sh
}

func DirectorySetUp(dir string, golang string) *api.Bash {
	untar := fmt.Sprintf("tar xf %s -C %s; rm golang-latest.tar.gz; ", golang, dir)

	stream := simpleio.NewReader("https://github.com/edotau/goFish/blob/main/bash/testdata/setupDir.sh")
	writer := simpleio.NewWriter("setupDir.sh")

	io.Copy(writer.Writer, stream)
	writer.Close()

	command := untar + "bash setupDir.sh; rm setupDir.sh"
	sh := api.NewScript(command)
	sh.Run()
	return sh
}
