package simpleio

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/biogo/hts/bgzf"
)

type BgzipReader struct {
	*bufio.Reader
	file   *os.File
	line   []byte
	Buffer *bytes.Buffer
}

func NewBgzipReader(filename string) *BgzipReader {
	var answer BgzipReader = BgzipReader{
		file:   OpenFile(filename),
		line:   make([]byte, defaultBufSize),
		Buffer: &bytes.Buffer{},
	}

	reader, err := bgzf.NewReader(answer.file, 1)
	answer.Reader = bufio.NewReader(reader)
	ErrorHandle(err)
	return &answer
}

func ReadLineBgzip(reader *BgzipReader) (*bytes.Buffer, bool) {
	var err error
	reader.line = reader.line[:0]
	reader.line, err = reader.ReadSlice('\n')
	if err == nil {
		if reader.line[len(reader.line)-1] == '\n' {
			reader.Buffer.Reset()
			_, err = reader.Buffer.Write(reader.line[:len(reader.line)-1])
			ErrorHandle(err)
			return reader.Buffer, false
		} else {
			log.Fatalf("Error: end of line did not end with an end of line character...\n")
		}
	} else {
		if err == bufio.ErrBufferFull {
			if reader.line[len(reader.line)-1] == '\n' {
				reader.Buffer.Reset()
				reader.line = reader.line[:len(reader.line)-1]
				//_, err = reader.Buffer.Write(reader.line[:len(reader.line)-1])
				//common.ExitIfError(err)
				reader.line = append(reader.line, readMoreBgzip(reader)...)
				_, err = reader.Buffer.Write(reader.line[:len(reader.line)-1])
				ErrorHandle(err)
				return reader.Buffer, false
			}
		} else {
			CatchErrThrowEOF(err)
			reader.Close()
		}
	}
	return nil, true
}

func (reader *BgzipReader) Close() {
	if reader != nil {
		err := reader.file.Close()
		ErrorHandle(err)
	}
}

func readMoreBgzip(reader *BgzipReader) []byte {
	var err error
	reader.line, err = reader.ReadBytes('\n')
	_, err = reader.Buffer.Write(reader.line[:len(reader.line)-1])
	ErrorHandle(err)
	return reader.Buffer.Bytes()
}

func ReadBgzipFile(filename string) *bytes.Buffer {
	reader := NewBgzipReader(filename)
	var ans bytes.Buffer
	for line, done := ReadLineBgzip(reader); !done; line, done = ReadLineBgzip(reader) {
		line.WriteByte('\n')
		io.Copy(&ans, line)
	}
	return &ans
}
