package bash

import (
	"bufio"
	"bytes"

	"io"
	"log"
	"os/exec"

	"github.com/edotau/goFish/simpleio"
	"github.com/klauspost/pgzip"
)

const (
	defaultBufSize = 4096
)

// GunzipReader uncompress the input using the system's gzip. Apparently,
// the system gzip is much much faster than the go library,
// so I wrote some bench marks and tests
type GunzipReader struct {
	*bufio.Reader
	Unzip  io.Reader
	Cmd    *exec.Cmd
	line   []byte
	Buffer *bytes.Buffer
}

type GzipWriter struct {
	*bufio.Writer
	Gzip   *pgzip.Writer
	Buffer *bytes.Buffer
	close  func() error
}

func NewPzip(filename string, data []string) {
	ans := GzipWriter{}
	file := simpleio.Touch(filename)
	ans.Writer = bufio.NewWriter(file)
	ans.Gzip = pgzip.NewWriter(ans.Writer)
	ans.Gzip.SetConcurrency(100000, 10)
	for i := 0; i < len(data); i++ {
		ans.Gzip.Write([]byte(data[i] + "\n"))
	}
	file.Close()
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
func ReadLine(reader *GunzipReader) (*bytes.Buffer, bool) {
	var err error
	reader.line, err = reader.ReadSlice('\n')
	reader.Buffer.Reset()
	if err == nil {
		if reader.line[len(reader.line)-1] == '\n' {
			return BytesToBuffer(reader), false
		} else {
			log.Fatalf("Error: end of line did not end with an end of line character...\n")
		}
	} else {
		if err == bufio.ErrBufferFull {
			reader.line = readMore(reader)
			return BytesToBuffer(reader), false
		} else {
			simpleio.CatchErrThrowEOF(err)
		}
	}
	return nil, true
}

// BytesToBuffer will parse []byte and return a pointer to the same underlying bytes.Buffer
func BytesToBuffer(reader *GunzipReader) *bytes.Buffer {
	_, err := reader.Buffer.Write(reader.line[:len(reader.line)-1])
	simpleio.ErrorHandle(err)
	return reader.Buffer
}

// readMore is a private helper function to deal with very long lines to
// avoid alocating too much memory upfront and only resize the size of the buffer
// only when necessary.
func readMore(reader *GunzipReader) []byte {
	_, err := reader.Buffer.Write(reader.line)
	simpleio.ErrorHandle(err)
	reader.line, err = reader.ReadSlice('\n')
	if err == nil {
		return reader.line
	}
	if err == bufio.ErrBufferFull {
		_, err = reader.Buffer.Write(reader.line)
		simpleio.ErrorHandle(err)
		// recursive call to read next bytes until reaching end of line character
		return readMore(reader)
	}
	simpleio.ErrorHandle(err)
	return reader.line
}

func NewGunzipReader(filename string) *GunzipReader {
	var answer GunzipReader = GunzipReader{
		line:   make([]byte, defaultBufSize),
		Buffer: &bytes.Buffer{},
	}
	cmd := exec.Command("gunzip", "-c", filename)
	stdout, err := cmd.StdoutPipe()
	simpleio.ErrorHandle(err)
	err = cmd.Start()
	simpleio.ErrorHandle(err)
	answer.Reader = bufio.NewReader(GunzipReader{Unzip: stdout, Cmd: cmd})
	return &answer
}
