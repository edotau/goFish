package bash

import (
	"io"
	"os/exec"

	"github.com/edotau/goFish/simpleio"
)

// GunzipReader uncompress the input using the system's gzip. Apparently,
// the system gzip is much much faster than the go library,
// so I wrote some bench marks and tests
type GunzipReader struct {
	Unzip io.Reader
	Cmd   *exec.Cmd
}

func (gz GunzipReader) Read(data []byte) (int, error) {
	var err error

	var offset int
	var read_len int

	for offset = 0; offset < len(data) && err == nil; read_len, err = gz.Unzip.Read(data[offset:]) {
		offset += read_len
	}
	return offset, err
}

func (gz GunzipReader) Close() {
	gz.Unzip.(io.ReadCloser).Close()
	gz.Cmd.Wait()
}

func NewGunzipReader(filename string) (*GunzipReader, error) {
	cmd := exec.Command("gunzip", "-c", filename)
	stdout, err := cmd.StdoutPipe()
	simpleio.ErrorHandle(err)
	err = cmd.Start()
	return &GunzipReader{Unzip: stdout, Cmd: cmd}, err
}
