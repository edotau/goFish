package simpleio

import (
	"bufio"
	"compress/gzip"
	"os"
	"strings"
)

type SimpleWriter struct {
	*bufio.Writer
	Gzip  *gzip.Writer
	close func() error
}

func NewWriter(filename string) *SimpleWriter {
	ans := SimpleWriter{}
	file, err := os.Create(filename)
	ErrorHandle(err)
	ans.Writer = bufio.NewWriter(file)
	if strings.HasSuffix(filename, ".gz") {
		ans.Gzip = gzip.NewWriter(ans.Writer)
	} else {
		ans.Gzip = nil
		ans.close = file.Close
	}

	return &ans
}

func (writer *SimpleWriter) Write(p []byte) (n int, err error) {
	if writer.Gzip != nil {
		return writer.Gzip.Write(p)
	} else {
		return writer.Write(p)
	}
}

func (writer *SimpleWriter) WriteLine(s string) {
	if writer.Gzip == nil {
		writer.WriteString(s)
		writer.WriteByte('\n')
	} else {
		writer.Gzip.Write([]byte(s + "\n"))
	}
}

func (writer *SimpleWriter) Peek(n int) ([]byte, error) {
	return writer.Peek(n)
}

func (writer *SimpleWriter) Close() {
	if writer.Gzip != nil {
		writer.Gzip.Close()
	}
	if writer != nil {
		writer.Flush()
	}
	writer.close()
}
