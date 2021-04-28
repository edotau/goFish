package simpleio

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	gzip "github.com/klauspost/pgzip"
)

// HttpReader will fetch data from files uploaded to an internet server and stream data into an io.Reader interface.
func HttpReader(url string) *SimpleReader {
	resp, err := http.Get(url)
	ErrorHandle(err)
	var answer SimpleReader = SimpleReader{
		Buffer: &bytes.Buffer{},
		line:   make([]byte, defaultBufSize),
		close:  resp.Body.Close,
	}
	if strings.HasSuffix(url, ".gz") {
		gzipReader, err := gzip.NewReader(resp.Body)
		ErrorHandle(err)
		answer.Reader = bufio.NewReader(gzipReader)
	} else {
		answer.Reader = bufio.NewReader(resp.Body)
	}
	return &answer
}

// VimUrl is a basic function to procress a url link and print it out to stdout.
func VimUrl(url string) {
	reader := HttpReader(url)
	for i, err := ReadLine(reader); !err; i, err = ReadLine(reader) {
		fmt.Printf("%s\n", i.String())
	}
}
