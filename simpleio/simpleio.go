// Package simpleio contains core utils that reading in data
// optimized for both memory allocation, speed, and performance
package simpleio

import (
	"bufio"
	"bytes"
	"github.com/vertgenlab/gonomics/exception"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	gzip "github.com/klauspost/pgzip"
)

const (
	defaultBufSize = 4096
	mb             = 1024 * 1024
	gb             = 1024 * mb
	BufferSize     = mb
)

// SimpleReader implements the io.Reader interface by providing
// the Read(b []byte) method. The struct contains an embedded *bufio.Reader
// and a pointer to os.File for closeure when reading is complete.
type SimpleReader struct {
	*bufio.Reader
	Gunzip *GunzipReader
	line   []byte
	Buffer *bytes.Buffer
	close  func() error
}

// GunzipReader uncompress the input using the system's gzip. Apparently,
// the system gzip is much much faster than the go library, so I wrote some benchmarks and tests
type GunzipReader struct {
	Unzip io.Reader
	Cmd   *exec.Cmd
}

type SimpleWriter struct {
	*bufio.Writer
	Gzip   *gzip.Writer
	Buffer *bytes.Buffer
	close  func() error
}

// NewSimpleReader will process a given file and performs error handling if an error occurs.
// SimpleReader will prcoess gzipped files accordinging by performing a check on the suffix
// of the provided file.
func NewReader(filename string) *SimpleReader {
	if strings.HasPrefix(filename, "http") {
		return HttpReader(filename)
	}

	var answer SimpleReader = SimpleReader{
		line:   make([]byte, defaultBufSize),
		Buffer: &bytes.Buffer{},
	}
	switch true {
	case strings.HasSuffix(filename, ".gz"):
		var err error
		//gunzip, err := gzip.NewReader(file)
		answer.Gunzip, err = NewGunzipReader(filename)
		StdError(err)

		answer.close = answer.Gunzip.Unzip.(io.ReadCloser).Close
		answer.Reader = bufio.NewReader(answer.Gunzip)
	default:
		answer.Gunzip = nil
		file := Vim(filename)
		answer.Reader = bufio.NewReader(file)
		answer.close = file.Close
	}
	return &answer
}

func NewGunzipReader(filename string) (*GunzipReader, error) {
	cmd := exec.Command("gunzip", "-c", filename)
	stdout, err := cmd.StdoutPipe()
	StdError(err)
	err = cmd.Start()
	return &GunzipReader{Unzip: stdout, Cmd: cmd}, err
}

func NewWriter(filename string) *SimpleWriter {
	ans := SimpleWriter{}
	file := Touch(filename)

	ans.Writer = bufio.NewWriter(file)

	ans.Buffer = &bytes.Buffer{}
	if strings.HasSuffix(filename, ".gz") {
		ans.Gzip = gzip.NewWriter(ans.Writer)
		ans.Gzip.SetConcurrency(100000, 10)
	} else {
		ans.Gzip = nil
	}
	ans.close = file.Close
	return &ans
}

func Vim(filename string) *os.File {
	file, err := os.Open(filename)
	StdError(err)
	return file
}

func Touch(filename string) *os.File {
	file, err := os.Create(filename)
	StdError(err)
	return file
}

// Read reads data into p and is a method required to implement the io.Reader interface.
// It returns the number of bytes read into p.
func (reader *SimpleReader) Read(b []byte) (n int, err error) {
	if reader.Gunzip == nil {
		return reader.Read(b)
	} else {
		return reader.Gunzip.Read(b)
	}
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

func (writer *SimpleWriter) Write(p []byte) (n int, err error) {
	if writer.Gzip != nil {
		return writer.Gzip.Write(p)
	} else {
		return writer.Writer.Write(p)
	}
}

// ReadLine will return a bytes.Buffer pointing to the internal slice of bytes. Provided this function is called within a loop,
// the function will read one line at a time, and return bool to continue reading. Important to note the buffer return points to
// the internal slice belonging to the reader, meaning the slice will be overridden if the data is not copied. Please be aware the
// reader will call close on the file once the reader encounters EOF.
func ReadLine(reader *SimpleReader) (*bytes.Buffer, bool) {
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
			CatchErrThrowEOF(err)
		}
	}
	return nil, true
}

// readMore is a private helper function to deal with very long lines to
// avoid alocating too much memory upfront and only resize the size of the buffer
// only when necessary.
func readMore(reader *SimpleReader) []byte {
	_, err := reader.Buffer.Write(reader.line)
	StdError(err)
	reader.line, err = reader.ReadSlice('\n')
	if err == nil {
		return reader.line
	}
	if err == bufio.ErrBufferFull {
		_, err = reader.Buffer.Write(reader.line)
		StdError(err)
		// recursive call to read next bytes until reaching end of line character
		return readMore(reader)
	}
	StdError(err)
	return reader.line
}

func WriteLine(writer *SimpleWriter, s string) {
	writer.Buffer.WriteString(s)
	writer.Buffer.WriteByte('\n')
	io.Copy(writer, writer.Buffer)
}

// BytesToBuffer will parse []byte and return a pointer to the same underlying bytes.Buffer
func BytesToBuffer(reader *SimpleReader) *bytes.Buffer {
	var err error
	_, err = reader.Buffer.Write(bytes.TrimSpace(reader.line))

	// if reader.line[len(reader.line)-2] == '\r' {
	// 	reader.Buffer.Write(reader.line[:len(reader.line)-2])
	// } else {
	// 	_, err = reader.Buffer.Write(reader.line[:len(reader.line)-1])
	// }
	exception.PanicOnErr(err)
	return reader.Buffer
}

// StdError will simply print and handle errors returned.
func StdError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CatchErrThrowEOF will silently handles and throws the EOF error and will log and exit any other errors.
func CatchErrThrowEOF(err error) {
	if err == io.EOF {
		return
	} else {
		StdError(err)
	}
}

// Close closes the File, rendering it unusable for I/O. On files that support SetDeadline,
// any pending I/O operations will be canceled and return immediately with an error.
// Close will return an error if it has already been called.
func (reader *SimpleReader) Close() {
	if reader != nil {
		StdError(reader.close())
	}
}

func (writer *SimpleWriter) Close() {

	if writer.Gzip != nil {
		StdError(writer.Gzip.Close())

	}
	if writer != nil {
		writer.Writer.Flush()
	}
	StdError(writer.close())
}

func Rm(filename string) {
	err := os.Remove(filename)
	StdError(err)
}

func ReadFromFile(filename string) []string {
	reader := NewReader(filename)
	var ans []string
	for i, err := ReadLine(reader); !err; i, err = ReadLine(reader) {
		if !strings.HasPrefix(i.String(), "#") {
			ans = append(ans, i.String())
		}
	}
	reader.Close()
	return ans
}

func (reader *SimpleReader) ToString() string {
	var ans string
	for i, err := ReadLine(reader); !err; i, err = ReadLine(reader) {
		ans += i.String() + "\n"
	}
	reader.Close()
	return ans
}

func WriteToFile(filename string, data []string) {
	writer := NewWriter(filename)
	for i := 0; i < len(data); i++ {
		writer.Write([]byte(data[i] + "\n"))
		//writer.Buffer.WriteByte('\n')
		//io.Copy(writer, writer.Buffer)
	}
	writer.Close()
}

func GetBuffer(bufferPool sync.Pool) bytes.Buffer {
	return bufferPool.Get().(bytes.Buffer)
}

func PutBuffer(buf bytes.Buffer, bufferPool sync.Pool) {
	buf.Reset()
	bufferPool.Put(buf)
}

// SimplyRun uses a new goroutine to run the function
func SimplyRun(f func()) {
	go f()
}

/*
func SimpleBufioPoolTest(filename string) {
	reader := NewSimpleReader(filename)
	var simplePool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	writer := fileio.EasyCreate("testdata/simplePoolTest.vcf.gz")

	var line []byte
	var err error
	var errW error
	var curr *bytes.Buffer
	for {
		curr = simplePool.Get().(*bytes.Buffer)

		line, err = reader.ReadSlice('\n')

		if err == nil {
			//line.b[len(line.b)-1] == '\n'
			curr.Write(line[:len(line)-1])

			//_, errW = fmt.Fprintf(writer, "%s\n", line)
			//fmt.Printf("%s\n", )
			simpleLogErr(errW)

			simplePool.Put(line)

		} else if err == io.EOF {
			CatchErrThrowEOF(err)
			break
		} else {
			simpleLogErr(err)
		}

	}
	writer.Close()
}*/

// BufferPool implements a pool of bytes.Buffers in the form of a bounded
// channel.
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new BufferPool bounded to the given size.
func NewBufferPool(size int) *BufferPool {
	var bp BufferPool
	bp.pool.New = func() interface{} {
		return new(bytes.Buffer)
	}
	return &bp
}

// Get gets a Buffer from the BufferPool, or creates a new one if none are
// available in the pool.
func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

// Put returns the given Buffer to the BufferPool.
func (bp *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.pool.Put(b)
}

/*
func client() {
    for {
        var b *bytes.Buffer
        // Grab a buffer if available; allocate if not.
        select {
        case b = <-freeList:
            // Got one; nothing more to do.
        default:
            // None free, so allocate a new one.
            b = new(Buffer)
        }
        load(b)              // Read next message from the net.
        serverChan <- b      // Send to server.
    }
}

func server() {
    for {
        b := <-serverChan    // Wait for work.
        process(b)
        // Reuse buffer if there's room.
        select {
        case freeList <- b:
            // Buffer on free list; nothing more to do.
        default:
            // Free list full, just carry on.
        }
    }
}*/
/*
func routineReader(reader *SimpleReader, data chan []byte) {
	for curr, done := ReadLine(reader); !done; curr, done = ReadLine(reader) {
		data <- curr
	}
	close(data)
}

func ConcurrentReader(filename string) <- chan []byte {
	simpleReader := NewSimpleReader(filename)
	ans := make(chan []byte, 100000)
	go routineReader(simpleReader, ans)
	return ans
}*/
